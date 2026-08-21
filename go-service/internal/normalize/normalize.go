package normalize

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var ErrNoPrice = errors.New("normalize: no parsable price")

var (
	priceRe    = regexp.MustCompile(`(\d{1,3}(?:,\d{3})*(?:\.\d{1,2})?|\d+(?:\.\d{1,2})?)`)
	nonAlnumRe = regexp.MustCompile(`[^a-z0-9]+`)
	packRe     = regexp.MustCompile(`(?i)\b(\d+)\s*[- ]?\s*(ct|count|pack|pk|pc|pcs|piece|sheets?|ream|box|case|each|ea)\b`)
	// "Set of 400", "Pack of 12", "Case of 6" -- the dominant form on
	// education-supply sites, where bulk classpacks are the whole catalogue.
	setOfRe = regexp.MustCompile(`(?i)\b(?:set|pack|box|case|classpack)\s+of\s+(\d+)\b`)
)

// PriceCents turns whatever the scraper handed back into integer cents.
// Handles "$12.99", "12,499.00 USD", "USD 8.5", " $1,299 ".
// Money never goes near a float in this codebase past this function.
func PriceCents(raw string) (int64, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, ErrNoPrice
	}
	m := priceRe.FindString(s)
	if m == "" {
		return 0, ErrNoPrice
	}
	m = strings.ReplaceAll(m, ",", "")
	f, err := strconv.ParseFloat(m, 64)
	if err != nil {
		return 0, ErrNoPrice
	}
	if f <= 0 {
		return 0, ErrNoPrice
	}
	// +0.5 so 12.99*100 = 1298.9999... rounds to 1299 instead of truncating.
	return int64(f*100 + 0.5), nil
}

// CanonicalKey is the crude product-matching key: same key means we treat two
// retailer listings as the same physical item. Deliberately simple. When you
// have time, upgrade to brand + normalized model + pack size with a fuzzy
// fallback, but do not let matching block the pipeline this week.
func CanonicalKey(brand, name string, unitCount float64, unitType string) string {
	parts := []string{
		slug(brand),
		slug(name),
	}
	if unitCount > 0 {
		parts = append(parts, strconv.FormatFloat(unitCount, 'f', -1, 64)+slug(unitType))
	}
	out := strings.Trim(strings.Join(parts, "-"), "-")
	out = strings.ReplaceAll(out, "--", "-")
	return out
}

func slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonAlnumRe.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// Pack pulls "24 ct", "500 sheets", "pack of 12" out of a product title so
// price-per-unit is comparable across retailers that sell different pack sizes.
// This is the whole point of the calculator: a district buying a case of 10
// reams needs cost per sheet, not cost per SKU.
func Pack(title string) (count float64, unitType string) {
	// "Set of N" wins when present: a title like "Value Pack - Set of 400"
	// contains both forms and the explicit count is the right one.
	if m := setOfRe.FindStringSubmatch(title); len(m) == 2 {
		if n, err := strconv.ParseFloat(m[1], 64); err == nil && n > 0 {
			return n, "each"
		}
	}
	m := packRe.FindStringSubmatch(title)
	if len(m) != 3 {
		return 0, "each"
	}
	n, err := strconv.ParseFloat(m[1], 64)
	if err != nil || n <= 0 {
		return 0, "each"
	}
	return n, canonUnit(m[2])
}

func canonUnit(u string) string {
	switch strings.ToLower(u) {
	case "ct", "count", "pc", "pcs", "piece":
		return "each"
	case "pk", "pack":
		return "pack"
	case "sheet", "sheets":
		return "sheet"
	case "ea":
		return "each"
	default:
		return strings.ToLower(u)
	}
}

// PricePerUnitCents is the derived figure the frontend actually compares on.
func PricePerUnitCents(priceCents int64, unitCount float64) float64 {
	if unitCount <= 0 {
		return float64(priceCents)
	}
	return float64(priceCents) / unitCount
}

// Title cleanup: collapse whitespace, strip the retailer's marketing suffix.
func Title(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	for _, cut := range []string{" | ", " - Buy ", " :: "} {
		if i := strings.Index(s, cut); i > 20 {
			s = s[:i]
		}
	}
	return strings.TrimSpace(s)
}
