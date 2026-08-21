package pipeline

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/savern02/into-the-scrape-verse/go-service/internal/brightdata"
	"github.com/savern02/into-the-scrape-verse/go-service/internal/normalize"
	"github.com/savern02/into-the-scrape-verse/go-service/internal/store"
)

// Job is one snapshot of scraped rows waiting to be processed.
type Job struct {
	CollectorID string
	SnapshotID  string
	ObjectKey   string // where the raw payload was preserved
	Payload     []byte // raw JSON array from Bright Data
	ObservedAt  time.Time
}

// RawRow is the loose shape coming off the scraper. Everything is a string
// because that is what extraction gives you and because assuming otherwise is
// how pipelines break silently.
type RawRow struct {
	Title       string `json:"title"`
	Brand       string `json:"brand"`
	Price       string `json:"price"`
	ListPrice   string `json:"list_price"`
	Currency    string `json:"currency"`
	SKU         string `json:"sku"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Availability  string `json:"availability"`
	EduPricing  string `json:"edu_pricing"`
	MinBulkQty  int    `json:"min_bulk_qty"`
}

type Pipeline struct {
	store    *store.Store
	bd       *brightdata.Client
	jobs     chan Job
	wg       sync.WaitGroup
	workers  int
	autoHeal    bool
	autoApprove bool
	brokenAt    float64
	healTO   time.Duration
	log      *slog.Logger
}

func New(st *store.Store, bd *brightdata.Client, workers int, autoHeal, autoApprove bool, brokenAt float64, healTO time.Duration, log *slog.Logger) *Pipeline {
	return &Pipeline{
		store:    st,
		bd:       bd,
		jobs:     make(chan Job, 128),
		workers:  workers,
		autoHeal:    autoHeal,
		autoApprove: autoApprove,
		brokenAt:    brokenAt,
		healTO:   healTO,
		log:      log,
	}
}

// Start spins up the worker pool. Ingestion is concurrent because a discovery
// scraper can hand you thousands of rows and you do not want the webhook
// handler blocking while you chew through them.
func (p *Pipeline) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func(id int) {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-p.jobs:
					if !ok {
						return
					}
					if err := p.process(ctx, job); err != nil {
						p.log.Error("process failed",
							"worker", id, "collector", job.CollectorID, "err", err)
					}
				}
			}
		}(i)
	}
}

// Submit is non-blocking-ish: it returns false if the queue is full so the
// webhook can answer 503 rather than hanging Bright Data's retry logic.
func (p *Pipeline) Submit(j Job) bool {
	select {
	case p.jobs <- j:
		return true
	default:
		return false
	}
}

func (p *Pipeline) Stop() {
	close(p.jobs)
	p.wg.Wait()
}

func (p *Pipeline) process(ctx context.Context, job Job) error {
	col, err := p.store.GetCollector(ctx, job.CollectorID)
	if err != nil {
		return err
	}

	ingestID, err := p.store.StartIngest(ctx, job.CollectorID, job.SnapshotID, job.ObjectKey)
	if err != nil {
		return err
	}

	var rows []RawRow
	if err := json.Unmarshal(job.Payload, &rows); err != nil {
		_ = p.store.FinishIngest(ctx, ingestID, 0, err)
		return err
	}

	observedAt := job.ObservedAt
	if observedAt.IsZero() {
		observedAt = time.Now().UTC()
	}

	// Count empties per field as we go. These counts are the entire basis for
	// deciding whether the scraper is still working.
	nulls := map[string]int{"title": 0, "price": 0, "url": 0}
	written, dropped := 0, 0

	for _, r := range rows {
		if r.Title == "" {
			nulls["title"]++
		}
		if r.URL == "" {
			nulls["url"]++
		}

		priceCents, perr := normalize.PriceCents(r.Price)
		if perr != nil {
			nulls["price"]++
		}

		// A row without a price or a URL is not an observation, it's a gap.
		if perr != nil || r.URL == "" {
			dropped++
			continue
		}

		listCents, _ := normalize.PriceCents(r.ListPrice)
		count, unitType := normalize.Pack(r.Title)
		title := normalize.Title(r.Title)
		key := normalize.CanonicalKey(r.Brand, title, count, unitType)

		inStock := parseStock(r.Availability)
		currency := r.Currency
		if currency == "" {
			currency = "USD"
		}

		rec := store.Record{
			CollectorID:       job.CollectorID,
			RetailerID:        col.RetailerID,
			Name:              title,
			Brand:             r.Brand,
			Category:          r.Category,
			UnitType:          unitType,
			UnitCount:         count,
			SKU:               r.SKU,
			URL:               r.URL,
			EduEligible:       r.EduPricing != "",
			MinBulkQty:        r.MinBulkQty,
			PriceCents:        priceCents,
			ListPriceCents:    listCents,
			Currency:          currency,
			InStock:           inStock,
			QtyBreak:          1,
			PricePerUnitCents: normalize.PricePerUnitCents(priceCents, count),
			ObservedAt:        observedAt,
		}

		if err := p.store.UpsertRecord(ctx, ingestID, key, rec); err != nil {
			p.log.Warn("upsert failed", "url", r.URL, "err", err)
			dropped++
			continue
		}
		written++
	}

	if err := p.store.FinishIngest(ctx, ingestID, written, nil); err != nil {
		p.log.Warn("finish ingest", "err", err)
	}

	p.checkHealth(ctx, col, ingestID, len(rows), written, dropped, nulls)
	return nil
}

// checkHealth is the self-healing loop. Extraction coming back empty where it
// used to come back with a value is the signal; heal is the response.
func (p *Pipeline) checkHealth(ctx context.Context, col store.Collector, ingestID int64, total, written, dropped int, nulls map[string]int) {
	rates := map[string]float64{}
	if total > 0 {
		for f, n := range nulls {
			rates[f] = float64(n) / float64(total)
		}
	}

	verdict := "healthy"
	switch {
	case total == 0:
		verdict = "broken"
	case rates["price"] >= p.brokenAt || rates["url"] >= p.brokenAt || rates["title"] >= p.brokenAt:
		verdict = "broken"
	case dropped > 0 && float64(dropped)/float64(total) > 0.1:
		verdict = "degraded"
	}

	h := store.Health{
		CollectorID:  col.ID,
		RawIngestID:  ingestID,
		RowsReturned: total,
		RowsDropped:  dropped,
		NullRates:    rates,
		Verdict:      verdict,
	}

	if verdict == "broken" {
		desc := brightdata.Describe(rates, col.FieldSpec)
		h.HealCommand = joinArgs(brightdata.HealCommand(col.ID, desc))

		if p.autoHeal {
			p.log.Warn("collector broken, healing", "collector", col.ID, "rates", rates)
			res, err := brightdata.Heal(ctx, col.ID, desc, p.autoApprove, p.healTO)
			h.HealOutput = res.Raw
			h.HealStatus = res.Status
			h.Healed = res.Committed()

			switch {
			case err != nil:
				p.log.Error("heal failed", "collector", col.ID, "err", err)
			case res.Committed():
				p.log.Info("healed and committed", "collector", col.ID)
			case res.Status == "awaiting_approval":
				// Not a failure. The fix is built and waiting on a human.
				p.log.Warn("heal awaiting approval -- review preview, then approve",
					"collector", col.ID,
					"view_url", res.ViewURL,
					"next_step", res.NextStep)
			default:
				p.log.Warn("heal ended in unexpected state",
					"collector", col.ID, "status", res.Status)
			}
		} else {
			p.log.Warn("collector broken, auto-heal off",
				"collector", col.ID, "command", h.HealCommand)
		}
	}

	if err := p.store.RecordHealth(ctx, h); err != nil {
		p.log.Warn("record health", "err", err)
	}
}

func parseStock(s string) *bool {
	if s == "" {
		return nil
	}
	v := false
	switch {
	case containsAny(s, "in stock", "instock", "available", "in_stock", "true"):
		v = true
	case containsAny(s, "out of stock", "outofstock", "unavailable", "backorder", "false"):
		v = false
	default:
		return nil
	}
	return &v
}

func containsAny(s string, subs ...string) bool {
	ls := strings.ToLower(s)
	for _, sub := range subs {
		if strings.Contains(ls, sub) {
			return true
		}
	}
	return false
}

func joinArgs(args []string) string { return strings.Join(args, " ") }
