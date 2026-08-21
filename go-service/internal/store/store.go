package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct{ pool *pgxpool.Pool }

// NewPool connects to Supabase.
//
// IMPORTANT: use the connection string from Supabase -> Connect -> "Session
// pooler" or "Transaction pooler". The direct-connection host is IPv6-only and
// most cloud VMs will not reach it.
//
// If you use the TRANSACTION pooler (port 6543), pgbouncer does not support
// prepared statements, so pgx must fall back to the simple protocol or every
// query fails with "prepared statement already exists". That's what the
// DefaultQueryExecMode line below is doing.
func NewPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	cfg.MaxConns = 8 // Supabase free tier caps you around 60 total; leave room
	cfg.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("open pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return pool, nil
}

func New(pool *pgxpool.Pool) *Store { return &Store{pool: pool} }

// ---------------------------------------------------------------
// Raw layer
// ---------------------------------------------------------------

func (s *Store) StartIngest(ctx context.Context, collectorID, snapshotID, objectKey string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `
		insert into raw_ingests (collector_id, snapshot_id, object_key, status)
		values ($1, $2, $3, 'processing')
		on conflict (collector_id, snapshot_id) do update
		    set status = 'processing', error = null
		returning id`,
		collectorID, snapshotID, objectKey,
	).Scan(&id)
	return id, err
}

func (s *Store) FinishIngest(ctx context.Context, id int64, rowCount int, ingestErr error) error {
	status, msg := "done", (*string)(nil)
	if ingestErr != nil {
		status = "failed"
		e := ingestErr.Error()
		msg = &e
	}
	_, err := s.pool.Exec(ctx, `
		update raw_ingests
		   set status = $2, row_count = $3, error = $4, processed_at = now()
		 where id = $1`,
		id, status, rowCount, msg)
	return err
}

// ---------------------------------------------------------------
// Canonical layer
// ---------------------------------------------------------------

type Record struct {
	CollectorID string
	RetailerID  int64

	Name      string
	Brand     string
	Category  string
	UnitType  string
	UnitCount float64

	SKU         string
	URL         string
	EduEligible bool
	MinBulkQty  int

	PriceCents        int64
	ListPriceCents    int64
	Currency          string
	InStock           *bool
	QtyBreak          int
	PricePerUnitCents float64

	ObservedAt time.Time
}

// UpsertRecord writes one normalized observation across products/offers/
// price_observations in a single transaction. Products and offers are
// idempotent; price_observations is append-only.
func (s *Store) UpsertRecord(ctx context.Context, rawIngestID int64, canonicalKey string, r Record) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op after Commit

	var productID int64
	err = tx.QueryRow(ctx, `
		insert into products (canonical_key, name, brand, category, unit_type, unit_count)
		values ($1, $2, nullif($3,''), nullif($4,''), nullif($5,''), nullif($6,0))
		on conflict (canonical_key) do update
		    set name = excluded.name
		returning id`,
		canonicalKey, r.Name, r.Brand, r.Category, r.UnitType, r.UnitCount,
	).Scan(&productID)
	if err != nil {
		return fmt.Errorf("upsert product: %w", err)
	}

	var offerID int64
	err = tx.QueryRow(ctx, `
		insert into offers (product_id, retailer_id, sku, url, edu_eligible, min_bulk_qty)
		values ($1, $2, nullif($3,''), $4, $5, nullif($6,0))
		on conflict (retailer_id, url) do update
		    set product_id    = excluded.product_id,
		        sku           = coalesce(excluded.sku, offers.sku),
		        edu_eligible  = excluded.edu_eligible
		returning id`,
		productID, r.RetailerID, r.SKU, r.URL, r.EduEligible, r.MinBulkQty,
	).Scan(&offerID)
	if err != nil {
		return fmt.Errorf("upsert offer: %w", err)
	}

	_, err = tx.Exec(ctx, `
		insert into price_observations
		    (offer_id, raw_ingest_id, price_cents, list_price_cents, currency,
		     in_stock, qty_break, price_per_unit_cents, observed_at)
		values ($1, $2, $3, nullif($4,0), $5, $6, $7, $8, $9)
		on conflict (offer_id, qty_break, observed_at) do nothing`,
		offerID, rawIngestID, r.PriceCents, r.ListPriceCents, r.Currency,
		r.InStock, r.QtyBreak, r.PricePerUnitCents, r.ObservedAt,
	)
	if err != nil {
		return fmt.Errorf("insert observation: %w", err)
	}

	return tx.Commit(ctx)
}

// ---------------------------------------------------------------
// Collector metadata + health
// ---------------------------------------------------------------

type Collector struct {
	ID         string
	RetailerID int64
	Type       string
	TargetURL  string
	FieldSpec  map[string]string // field name -> plain-language description
}

func (s *Store) GetCollector(ctx context.Context, id string) (Collector, error) {
	var c Collector
	var specRaw []byte
	err := s.pool.QueryRow(ctx, `
		select id, retailer_id, scraper_type, coalesce(target_url,''), field_spec
		  from collectors where id = $1`, id,
	).Scan(&c.ID, &c.RetailerID, &c.Type, &c.TargetURL, &specRaw)
	if err != nil {
		return c, fmt.Errorf("collector %s: %w", id, err)
	}
	if len(specRaw) > 0 {
		_ = json.Unmarshal(specRaw, &c.FieldSpec)
	}
	return c, nil
}

type Health struct {
	CollectorID   string
	RawIngestID   int64
	RowsReturned  int
	RowsDropped   int
	NullRates     map[string]float64
	Verdict       string
	HealCommand   string
	HealOutput    string
	HealStatus    string // awaiting_approval | done | failed | rejected
	Healed        bool   // true only when the fix is actually committed
}

func (s *Store) RecordHealth(ctx context.Context, h Health) error {
	rates, _ := json.Marshal(h.NullRates)
	var healedAt *time.Time
	if h.Healed {
		now := time.Now()
		healedAt = &now
	}
	_, err := s.pool.Exec(ctx, `
		insert into collector_health
		    (collector_id, raw_ingest_id, rows_returned, rows_dropped,
		     field_null_rates, verdict, heal_command, heal_output, heal_status, healed_at)
		values ($1, nullif($2,0), $3, $4, $5, $6, nullif($7,''), nullif($8,''), nullif($9,''), $10)`,
		h.CollectorID, h.RawIngestID, h.RowsReturned, h.RowsDropped,
		rates, h.Verdict, h.HealCommand, h.HealOutput, h.HealStatus, healedAt)
	return err
}

func (s *Store) Close() { s.pool.Close() }
