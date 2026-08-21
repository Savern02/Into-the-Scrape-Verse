package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Supabase pooler connection string. Use the SESSION pooler (port 5432)
	// or the TRANSACTION pooler (6543). See store.NewPool for the pgx caveat.
	DatabaseURL string

	// Bright Data API token, from your teammate who owns the scraping side.
	BrightDataToken string

	// Shared secret you register with the Bright Data webhook so random
	// POSTs to your VM can't inject rows into your database.
	WebhookSecret string

	Addr    string
	Workers int

	// When true the pipeline shells out to `bdata scraper heal` on a broken
	// collector. When false it only records the command it *would* have run.
	// Leave false until you've watched it work once by hand.
	AutoHeal bool

	// Skip the human approval gate. `bdata scraper heal` normally stops and
	// waits for you to review preview_result. True = fully autonomous.
	AutoApprove bool

	// Fraction of rows with an empty price above which a collector is broken.
	BrokenNullRate float64

	HealTimeout time.Duration

	// Cloudflare R2. Leave R2Bucket empty to fall back to local disk so the
	// pipeline is never blocked waiting on someone else's dashboard access.
	R2AccountID string
	R2AccessKey string
	R2SecretKey string
	R2Bucket    string
	LocalRawDir string
}

func Load() Config {
	loadDotEnv(".env")

	c := Config{
		DatabaseURL:     req("DATABASE_URL"),
		BrightDataToken: os.Getenv("BRIGHTDATA_TOKEN"),
		WebhookSecret:   os.Getenv("WEBHOOK_SECRET"),
		Addr:            def("ADDR", ":8080"),
		Workers:         intDef("WORKERS", 4),
		AutoHeal:        os.Getenv("AUTO_HEAL") == "true",
		AutoApprove:     os.Getenv("AUTO_APPROVE") == "true",
		BrokenNullRate:  floatDef("BROKEN_NULL_RATE", 0.5),
		HealTimeout:     30 * time.Minute,
		R2AccountID:     os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKey:     os.Getenv("R2_ACCESS_KEY_ID"),
		R2SecretKey:     os.Getenv("R2_SECRET_ACCESS_KEY"),
		R2Bucket:        os.Getenv("R2_BUCKET"),
		LocalRawDir:     def("LOCAL_RAW_DIR", "./raw"),
	}
	return c
}

func req(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("config: %s is required", k)
	}
	return v
}

func def(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func intDef(k string, d int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return d
}

func floatDef(k string, d float64) float64 {
	if v := os.Getenv(k); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return d
}
