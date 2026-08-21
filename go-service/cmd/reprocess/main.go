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
	"os/signal"
	"syscall"
	"time"

	"github.com/savern02/into-the-scrape-verse/go-service/internal/config"
	"github.com/savern02/into-the-scrape-verse/go-service/internal/storage"
	"github.com/savern02/into-the-scrape-verse/go-service/internal/store"
)

// reprocess replays preserved raw payloads through the ingest pipeline.
//
// This is the payoff of the raw-first architecture. When the normalization
// logic changes -- a new pack-size pattern, a different unit conversion, a
// bug fix in price parsing -- every observation ever collected can be
// re-derived from object storage. No scraping, no credits, no waiting on a
// retailer's rate limits.
//
//	go run ./cmd/reprocess -list              # what is available to replay
//	go run ./cmd/reprocess -all               # replay everything
//	go run ./cmd/reprocess -collector c_x     # one collector's history
//
// Snapshot IDs are reused deliberately: StartIngest upserts on
// (collector_id, snapshot_id), so replaying updates the existing ingest
// record rather than inflating history with duplicates.
func main() {
	var (
		collectorID = flag.String("collector", "", "limit to one collector")
		all         = flag.Bool("all", false, "replay every preserved snapshot")
		list        = flag.Bool("list", false, "show what could be replayed, then exit")
		limit       = flag.Int("limit", 20, "maximum snapshots to consider")
		webhook     = flag.String("webhook", "http://localhost:8080/webhook/brightdata", "ingest endpoint")
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

	raw, err := newStorage(ctx, cfg, log)
	if err != nil {
		log.Error("storage", "err", err)
		os.Exit(1)
	}

	ingests, err := st.ListIngests(ctx, *collectorID, *limit)
	if err != nil {
		log.Error("list ingests", "err", err)
		os.Exit(1)
	}
	if len(ingests) == 0 {
		fmt.Println("No preserved snapshots found.")
		return
	}

	if *list || !*all {
		fmt.Printf("%d preserved snapshot(s):\n\n", len(ingests))
		for _, in := range ingests {
			fmt.Printf("  %-22s %-28s %4d rows  %s\n",
				in.CollectorID, in.SnapshotID, in.RowCount, in.ObjectKey)
		}
		if !*list {
			fmt.Println("\nPass -all to replay these through the pipeline.")
		}
		return
	}

	var ok, failed int
	for _, in := range ingests {
		if err := replay(ctx, raw, *webhook, cfg.WebhookSecret, in); err != nil {
			log.Error("replay failed",
				"snapshot", in.SnapshotID, "key", in.ObjectKey, "err", err)
			failed++
			continue
		}
		log.Info("replayed",
			"collector", in.CollectorID, "snapshot", in.SnapshotID, "rows", in.RowCount)
		ok++
	}

	fmt.Printf("\nReplayed %d snapshot(s), %d failed. No scraping credits spent.\n", ok, failed)
}

func replay(ctx context.Context, raw storage.Store, webhook, secret string, in store.Ingest) error {
	payload, err := raw.Get(ctx, in.ObjectKey)
	if err != nil {
		return fmt.Errorf("fetch %s: %w", in.ObjectKey, err)
	}

	// Guard against a truncated or corrupted object before it reaches the
	// pipeline: a bad replay would look identical to a broken scraper.
	var rows []json.RawMessage
	if err := json.Unmarshal(payload, &rows); err != nil {
		return fmt.Errorf("stored payload is not a JSON array: %w", err)
	}

	body, err := json.Marshal(map[string]any{
		"collector_id": in.CollectorID,
		"snapshot_id":  in.SnapshotID,
		"data":         json.RawMessage(payload),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if secret != "" {
		req.Header.Set("X-Webhook-Secret", secret)
	}

	resp, err := (&http.Client{Timeout: 2 * time.Minute}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		return fmt.Errorf("ingest returned %s: %s", resp.Status, buf.String())
	}
	return nil
}

func newStorage(ctx context.Context, cfg config.Config, log *slog.Logger) (storage.Store, error) {
	if cfg.R2Bucket == "" {
		log.Warn("R2 not configured, reading raw payloads from local disk", "dir", cfg.LocalRawDir)
		return &storage.Local{Dir: cfg.LocalRawDir}, nil
	}
	return storage.NewR2(ctx, storage.R2Config{
		AccountID: cfg.R2AccountID,
		AccessKey: cfg.R2AccessKey,
		SecretKey: cfg.R2SecretKey,
		Bucket:    cfg.R2Bucket,
	})
}
