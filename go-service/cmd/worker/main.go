package main

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/savern02/into-the-scrape-verse/go-service/internal/brightdata"
	"github.com/savern02/into-the-scrape-verse/go-service/internal/config"
	"github.com/savern02/into-the-scrape-verse/go-service/internal/pipeline"
	"github.com/savern02/into-the-scrape-verse/go-service/internal/storage"
	"github.com/savern02/into-the-scrape-verse/go-service/internal/store"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(log)

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := store.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("database", "err", err)
		os.Exit(1)
	}
	st := store.New(pool)
	defer st.Close()

	raw, err := newStorage(ctx, cfg)
	if err != nil {
		log.Error("storage", "err", err)
		os.Exit(1)
	}

	bd := brightdata.New(cfg.BrightDataToken)
	p := pipeline.New(st, bd, cfg.Workers, cfg.AutoHeal, cfg.AutoApprove, cfg.BrokenNullRate, cfg.HealTimeout, log)
	p.Start(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	mux.HandleFunc("POST /webhook/brightdata", webhookHandler(p, bd, raw, cfg.WebhookSecret, log))

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Info("listening", "addr", cfg.Addr, "workers", cfg.Workers, "auto_heal", cfg.AutoHeal)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	p.Stop()
	log.Info("stopped")
}

// webhookPayload is what Bright Data POSTs when a run finishes. Shapes vary by
// scraper type -- log one real delivery and adjust before you rely on it.
type webhookPayload struct {
	CollectorID string `json:"collector_id"`
	SnapshotID  string `json:"snapshot_id"`
	Status      string `json:"status"`
	// Either the rows arrive inline, or a URL to fetch them from.
	URL  string          `json:"url"`
	Data json.RawMessage `json:"data"`
}

func newStorage(ctx context.Context, cfg config.Config) (storage.Store, error) {
	if cfg.R2Bucket == "" {
		slog.Warn("R2 not configured, preserving raw payloads on local disk",
			"dir", cfg.LocalRawDir)
		return &storage.Local{Dir: cfg.LocalRawDir}, nil
	}
	return storage.NewR2(ctx, storage.R2Config{
		AccountID: cfg.R2AccountID,
		AccessKey: cfg.R2AccessKey,
		SecretKey: cfg.R2SecretKey,
		Bucket:    cfg.R2Bucket,
	})
}

func webhookHandler(p *pipeline.Pipeline, bd *brightdata.Client, raw storage.Store, secret string, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Anyone who finds your VM's IP can POST here otherwise. Register the
		// same value as a header on the Bright Data webhook config.
		if secret != "" {
			got := r.Header.Get("X-Webhook-Secret")
			if subtle.ConstantTimeCompare([]byte(got), []byte(secret)) != 1 {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		body, err := io.ReadAll(io.LimitReader(r.Body, 64<<20)) // 64 MiB ceiling
		if err != nil {
			http.Error(w, "read body", http.StatusBadRequest)
			return
		}

		var wp webhookPayload
		if err := json.Unmarshal(body, &wp); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		if wp.CollectorID == "" {
			http.Error(w, "missing collector_id", http.StatusBadRequest)
			return
		}

		payload := []byte(wp.Data)
		if len(payload) == 0 && wp.URL != "" {
			fetchCtx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
			defer cancel()
			payload, err = bd.FetchSnapshot(fetchCtx, wp.URL)
			if err != nil {
				log.Error("fetch snapshot", "collector", wp.CollectorID, "err", err)
				http.Error(w, "fetch failed", http.StatusBadGateway)
				return
			}
		}
		if len(payload) == 0 {
			http.Error(w, "no data", http.StatusBadRequest)
			return
		}

		// Raw-first: preserve the bytes BEFORE any parsing. If this write
		// fails we refuse the delivery rather than processing data we cannot
		// reproduce later -- Bright Data will retry.
		now := time.Now().UTC()
		objectKey := storage.Key(wp.CollectorID, wp.SnapshotID, now)
		if err := raw.Put(r.Context(), objectKey, payload); err != nil {
			log.Error("preserve raw", "key", objectKey, "err", err)
			http.Error(w, "storage failed", http.StatusServiceUnavailable)
			return
		}

		job := pipeline.Job{
			CollectorID: wp.CollectorID,
			SnapshotID:  wp.SnapshotID,
			ObjectKey:   objectKey,
			Payload:     payload,
			ObservedAt:  now,
		}

		if !p.Submit(job) {
			// Backpressure: tell Bright Data to retry rather than dropping rows.
			http.Error(w, "queue full", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"accepted":true}`))
	}
}
