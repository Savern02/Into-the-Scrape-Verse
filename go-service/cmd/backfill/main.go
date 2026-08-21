package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/savern02/into-the-scrape-verse/go-service/internal/config"
	"github.com/savern02/into-the-scrape-verse/go-service/internal/storage"
)

// backfill uploads raw payloads that landed on local disk into object storage.
//
// Snapshots ingested before R2 was configured were preserved locally, so their
// object_key in Postgres points at a bucket location that was never written.
// This walks the local directory and puts each file at its recorded key,
// making the whole history replayable from storage rather than only the part
// collected after the bucket existed.
//
//	go run ./cmd/backfill -dry-run
//	go run ./cmd/backfill
func main() {
	var (
		dryRun = flag.Bool("dry-run", false, "list what would be uploaded, upload nothing")
		dir    = flag.String("dir", "", "local raw directory (default: LOCAL_RAW_DIR from config)")
	)
	flag.Parse()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()

	if cfg.R2Bucket == "" {
		fmt.Println("R2_BUCKET is empty -- nothing to back fill into. Configure R2 first.")
		os.Exit(1)
	}

	localDir := cfg.LocalRawDir
	if *dir != "" {
		localDir = *dir
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	r2, err := storage.NewR2(ctx, storage.R2Config{
		AccountID: cfg.R2AccountID,
		AccessKey: cfg.R2AccessKey,
		SecretKey: cfg.R2SecretKey,
		Bucket:    cfg.R2Bucket,
	})
	if err != nil {
		log.Error("r2", "err", err)
		os.Exit(1)
	}

	var uploaded, skipped, failed int

	err = filepath.WalkDir(localDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}

		// The local writer joins LOCAL_RAW_DIR with the storage key, and the
		// key already begins with "raw/". So ./raw/raw/c_x/date/snap.json maps
		// back to the key raw/c_x/date/snap.json.
		rel, relErr := filepath.Rel(localDir, path)
		if relErr != nil {
			return relErr
		}
		key := filepath.ToSlash(rel)

		if *dryRun {
			fmt.Printf("  would upload %s -> %s\n", path, key)
			skipped++
			return nil
		}

		body, readErr := os.ReadFile(path)
		if readErr != nil {
			log.Error("read", "path", path, "err", readErr)
			failed++
			return nil // one bad file should not abort the whole backfill
		}

		if putErr := r2.Put(ctx, key, body); putErr != nil {
			log.Error("upload", "key", key, "err", putErr)
			failed++
			return nil
		}

		log.Info("uploaded", "key", key, "bytes", len(body))
		uploaded++
		return nil
	})

	if err != nil {
		log.Error("walk", "dir", localDir, "err", err)
		os.Exit(1)
	}

	if *dryRun {
		fmt.Printf("\n%d file(s) would be uploaded from %s to bucket %s.\n",
			skipped, localDir, cfg.R2Bucket)
		return
	}
	fmt.Printf("\nUploaded %d, failed %d. Run `go run ./cmd/reprocess -all` to verify.\n",
		uploaded, failed)
}
