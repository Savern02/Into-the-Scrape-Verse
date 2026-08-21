package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/savern02/into-the-scrape-verse/go-service/internal/config"
	"github.com/savern02/into-the-scrape-verse/go-service/internal/store"
)

// collect runs every active collector and pushes the results into the ingest
// webhook. Run it from cron, or with -every to keep it resident.
//
//	go run ./cmd/collect                 # one pass over every active collector
//	go run ./cmd/collect -collector c_x  # just one
//	go run ./cmd/collect -every 6h       # resident scheduler
//	go run ./cmd/collect -dry-run        # print the plan, scrape nothing
//
// Each run costs Bright Data credits (1 per page load), so -dry-run first
// when you are changing anything.
func main() {
	var (
		collectorID = flag.String("collector", "", "run only this collector (default: all active)")
		every       = flag.Duration("every", 0, "repeat on this interval instead of exiting")
		webhook     = flag.String("webhook", "http://localhost:8080/webhook/brightdata", "ingest endpoint")
		inputFile   = flag.String("input-file", "", "file of URLs, one per line, instead of the collector's target_url")
		dryRun      = flag.Bool("dry-run", false, "print what would run without scraping")
		timeout     = flag.Duration("timeout", 20*time.Minute, "per-collector timeout")
	)
	flag.Parse()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := store.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()
	st := store.New(pool)

	run := func() {
		collectors, err := st.ActiveCollectors(ctx)
		if err != nil {
			log.Error("list collectors", "err", err)
			return
		}

		for _, c := range collectors {
			if *collectorID != "" && c.ID != *collectorID {
				continue
			}
			if c.TargetURL == "" && *inputFile == "" {
				log.Warn("skipping: no target_url and no -input-file", "collector", c.ID)
				continue
			}
			if err := collectOne(ctx, log, c, *webhook, cfg.WebhookSecret, *inputFile, *dryRun, *timeout); err != nil {
				// One bad collector must not stop the others. A partial run
				// beats no run when three of four retailers are healthy.
				log.Error("collect failed", "collector", c.ID, "err", err)
			}
		}
	}

	run()
	if *every <= 0 {
		return
	}

	log.Info("scheduler running", "interval", every.String())
	ticker := time.NewTicker(*every)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info("scheduler stopped")
			return
		case <-ticker.C:
			run()
		}
	}
}

func collectOne(
	ctx context.Context,
	log *slog.Logger,
	c store.Collector,
	webhook, secret, inputFile string,
	dryRun bool,
	timeout time.Duration,
) error {
	args := []string{
		"--yes", "--package", "@brightdata/cli",
		"brightdata", "scraper", "run", c.ID,
	}
	if inputFile != "" {
		args = append(args, "--input-file", inputFile)
	} else {
		args = append(args, c.TargetURL)
	}
	args = append(args, "--json")

	if dryRun {
		log.Info("dry run", "collector", c.ID, "command", "npx "+fmt.Sprint(args))
		return nil
	}

	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	log.Info("scraping", "collector", c.ID, "type", c.Type)
	start := time.Now()

	cmd := exec.CommandContext(runCtx, "npx", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("scraper run: %w: %s", err, truncate(stderr.String(), 400))
	}

	// The CLI prints progress on stderr and data on stdout, so stdout should
	// be a clean JSON array. Verify before shipping it downstream.
	var rows []json.RawMessage
	if err := json.Unmarshal(stdout.Bytes(), &rows); err != nil {
		return fmt.Errorf("scraper output was not a JSON array: %w: %s",
			err, truncate(stdout.String(), 300))
	}
	if len(rows) == 0 {
		// Not an error worth aborting on: the ingest side is what decides a
		// collector is broken, and it needs to see the empty result to say so.
		log.Warn("scraper returned zero rows", "collector", c.ID)
	}

	snapshotID := fmt.Sprintf("sched_%s", time.Now().UTC().Format("20060102T150405"))
	body, err := json.Marshal(map[string]any{
		"collector_id": c.ID,
		"snapshot_id":  snapshotID,
		"data":         json.RawMessage(stdout.Bytes()),
	})
	if err != nil {
		return err
	}

	if err := postWebhook(runCtx, webhook, secret, body); err != nil {
		return err
	}

	log.Info("collected",
		"collector", c.ID,
		"rows", len(rows),
		"snapshot", snapshotID,
		"took", time.Since(start).Round(time.Second).String())
	return nil
}

func postWebhook(ctx context.Context, url, secret string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if secret != "" {
		req.Header.Set("X-Webhook-Secret", secret)
	}

	resp, err := (&http.Client{Timeout: 5 * time.Minute}).Do(req)
	if err != nil {
		return fmt.Errorf("post to ingest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		return fmt.Errorf("ingest returned %s: %s", resp.Status, truncate(buf.String(), 200))
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
