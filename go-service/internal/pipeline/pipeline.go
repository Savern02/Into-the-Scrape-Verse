package pipeline

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
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
//
// The AI names fields from your create/heal description, and those names drift
// between scrapers -- a discovery scraper returns product_title/product_url
// while a PDP scraper returns title/url. Rather than rewrite this struct every
// time, accept both and resolve in Normalize().
type RawRow struct {
	Title        string `json:"title"`
	ProductTitle string `json:"product_title"`
	Name         string `json:"name"`

	Brand string `json:"brand"`

	// Scraper Studio returns money either as a plain string ("$12.99") or as a
	// nested object ({"value":99.99,"currency":"USD","symbol":"$"}). A single
	// field of the wrong type kills the whole json.Unmarshal, so both shapes
	// have to be accepted here -- see Money below.
	Price      Money  `json:"price"`
	ListPrice  Money  `json:"list_price"`
	Currency   string `json:"currency"`
	SKU        string `json:"sku"`
	ItemNumber string `json:"item_number"`

	URL            string `json:"url"`
	ProductURL     string `json:"product_url"`
	ProductPageURL string `json:"product_page_url"`

	Category     string `json:"category"`
	Availability string `json:"availability"`
	InStock      *bool  `json:"in_stock"`
	EduPricing   string `json:"edu_pricing"`
	MinBulkQty   int    `json:"min_bulk_qty"`
}

// Money accepts a bare string, a bare number, or an object with a value field.
// It always presents itself downstream as a string for normalize.PriceCents.
type Money struct {
	Text     string
	Currency string
}

func (m *Money) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	if s == "" || s == "null" {
		return nil
	}

	// Bare string: "$12.99"
	var asStr string
	if err := json.Unmarshal(b, &asStr); err == nil {
		m.Text = asStr
		return nil
	}

	// Bare number: 12.99
	var asNum float64
	if err := json.Unmarshal(b, &asNum); err == nil {
		m.Text = strconv.FormatFloat(asNum, 'f', -1, 64)
		return nil
	}

	// Object: {"value":99.99,"currency":"USD","symbol":"$"}
	var obj struct {
		Value    json.RawMessage `json:"value"`
		Amount   json.RawMessage `json:"amount"`
		Currency string          `json:"currency"`
	}
	if err := json.Unmarshal(b, &obj); err != nil {
		// Unknown shape. Swallow it rather than failing the whole batch --
		// the row will be dropped and counted as a null, which is exactly
		// what the health check is for.
		return nil
	}
	m.Currency = obj.Currency

	raw := obj.Value
	if len(raw) == 0 {
		raw = obj.Amount
	}
	if len(raw) == 0 {
		return nil
	}
	var n float64
	if err := json.Unmarshal(raw, &n); err == nil {
		m.Text = strconv.FormatFloat(n, 'f', -1, 64)
		return nil
	}
	var t string
	if err := json.Unmarshal(raw, &t); err == nil {
		m.Text = t
	}
	return nil
}

func (m Money) String() string { return m.Text }

// Normalize collapses the alternate field names down to the canonical ones and
// undoes the duplication we see when a selector matches a product tile twice
// (many sites render desktop and mobile variants of the same markup).
func (r *RawRow) Normalize() {
	if r.Currency == "" {
		r.Currency = firstNonEmpty(r.Price.Currency, r.ListPrice.Currency)
	}
	r.Title = firstNonEmpty(r.Title, r.ProductTitle, r.Name)
	r.URL = firstNonEmpty(r.URL, r.ProductURL, r.ProductPageURL)
	r.SKU = firstNonEmpty(r.SKU, r.ItemNumber)

	r.Title = dedupeRepeat(r.Title)
	r.Brand = dedupeRepeat(r.Brand)
	r.SKU = dedupeRepeat(r.SKU)

	// "Item # ALLSILK" -> "ALLSILK"
	if i := strings.Index(strings.ToLower(r.SKU), "item #"); i >= 0 {
		r.SKU = strings.TrimSpace(r.SKU[i+len("item #"):])
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// dedupeRepeat turns "Crayola Crayola" into "Crayola" and an 8x repeated title
// back into one. Works by finding the shortest prefix that, repeated, rebuilds
// the whole string. Leaves genuinely non-repeating text alone.
func dedupeRepeat(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if s == "" {
		return s
	}
	words := strings.Fields(s)
	n := len(words)
	for size := 1; size <= n/2; size++ {
		if n%size != 0 {
			continue
		}
		unit := words[:size]
		match := true
		for i := size; i < n && match; i += size {
			for j := 0; j < size; j++ {
				if !strings.EqualFold(words[i+j], unit[j]) {
					match = false
					break
				}
			}
		}
		if match {
			return strings.Join(unit, " ")
		}
	}
	return s
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
		r.Normalize()

		if r.Title == "" {
			nulls["title"]++
		}
		if r.URL == "" {
			nulls["url"]++
		}

		priceCents, perr := normalize.PriceCents(r.Price.String())
		if perr != nil {
			nulls["price"]++
		}

		// A row without a price or a URL is not an observation, it's a gap.
		if perr != nil || r.URL == "" {
			dropped++
			continue
		}

		listCents, _ := normalize.PriceCents(r.ListPrice.String())
		count, unitType := normalize.Pack(r.Title)
		title := normalize.Title(r.Title)
		key := normalize.CanonicalKey(r.Brand, title, count, unitType)

		inStock := r.InStock
		if inStock == nil {
			inStock = parseStock(r.Availability)
		}
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
