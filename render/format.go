package render

import (
	"fmt"
	"time"

	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

// formatNumber renders a number using locale-aware grouping + decimal
// separator. Mirrors browser Intl.NumberFormat with no special style.
func formatNumber(tag language.Tag, v any) string {
	p := message.NewPrinter(tag)
	return p.Sprintf("%v", number.Decimal(v))
}

// formatCurrency renders an amount with the locale-specific currency
// presentation. currencyCode is an ISO 4217 code like "USD"/"CAD"/"EUR".
// Mirrors Intl.NumberFormat { style: 'currency' }.
func formatCurrency(tag language.Tag, v float64, currencyCode string) string {
	unit, err := currency.ParseISO(currencyCode)
	if err != nil {
		// Be lenient: fall back to the plain number rather than failing
		// the whole render. Templates run at request time and we can't
		// surface errors gracefully past the FuncMap boundary.
		return formatNumber(tag, v)
	}
	p := message.NewPrinter(tag)
	return p.Sprintf("%v", currency.Symbol(unit.Amount(v)))
}

// formatPercent renders a 0..1 ratio as a percent ("20%" / "20 %").
func formatPercent(tag language.Tag, v float64) string {
	p := message.NewPrinter(tag)
	return p.Sprintf("%v", number.Percent(v))
}

// byteUnits maps the cascade index to short labels for "en" and "fr".
// Any tag with a base other than "en" or "fr" falls through to the en
// labels (x/text/feature/unit is not GA, so we hand-roll this table).
var byteUnits = []struct {
	en string
	fr string
}{
	{"B", "o"}, {"kB", "ko"}, {"MB", "Mo"}, {"GB", "Go"},
	{"TB", "To"}, {"PB", "Po"},
}

// formatBytes mirrors the JS cascadeBytes helper: pick the largest unit
// that keeps the displayed number >= 1 (and < 1024), then format the
// number locale-correctly and append the unit short label. Input is in
// the base unit (default "gigabyte" to match the JS side).
// Limitation: only en/fr unit suffixes; other locales fall back to en.
func formatBytes(tag language.Tag, v float64, baseUnit string) string {
	order := []string{"byte", "kilobyte", "megabyte", "gigabyte", "terabyte", "petabyte"}
	idx := -1
	for i, u := range order {
		if u == baseUnit {
			idx = i
			break
		}
	}
	if idx < 0 {
		return fmt.Sprintf("%v %s", v, baseUnit)
	}
	for v >= 1024 && idx < len(order)-1 {
		v /= 1024
		idx++
	}
	p := message.NewPrinter(tag)
	num := p.Sprintf("%v", number.Decimal(v, number.MaxFractionDigits(2)))
	unit := byteUnits[idx].en
	base, _ := tag.Base()
	if base.String() == "fr" {
		unit = byteUnits[idx].fr
	}
	// CLDR puts a NBSP between number and unit in most locales; match.
	return num + " " + unit
}

// formatDate formats a time.Time using a coarse style keyword.
// Supports "year" (locale-neutral numeric year) and "short" (ISO 8601);
// any other style returns RFC3339. medium/long/full are deferred pending
// a CLDR-driven date pattern dep (x/text does not yet provide one).
// tag is accepted for signature symmetry with the other helpers and for
// future use when CLDR date patterns land; intentionally unused today.
func formatDate(_ language.Tag, t time.Time, style string) string {
	switch style {
	case "year":
		return fmt.Sprintf("%d", t.Year())
	case "short":
		return t.Format("2006-01-02")
	default:
		return t.Format(time.RFC3339)
	}
}
