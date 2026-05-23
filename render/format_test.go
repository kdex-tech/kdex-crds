package render

import (
	"testing"
	"time"

	"golang.org/x/text/language"
)

// NBSP is U+00A0 NO-BREAK SPACE. x/text emits this between groups and
// before unit/percent suffixes for French locales (matching CLDR).
const NBSP = " "

func TestFormatNumber(t *testing.T) {
	cases := []struct {
		tag      string
		v        float64
		expected string
	}{
		{"en-CA", 1234567.89, "1,234,567.89"},
		// fr-CA: x/text uses U+00A0 NBSP as group separator.
		{"fr-CA", 1234567.89, "1" + NBSP + "234" + NBSP + "567,89"},
	}
	for _, c := range cases {
		got := formatNumber(language.MustParse(c.tag), c.v)
		if got != c.expected {
			t.Errorf("formatNumber(%s, %v) = %q, want %q", c.tag, c.v, got, c.expected)
		}
	}
}

func TestFormatCurrency(t *testing.T) {
	cases := []struct {
		tag      string
		v        float64
		code     string
		expected string
	}{
		// x/text emits the currency symbol followed by a space then the
		// amount; in fr-CA the symbol form is "$ US" with an internal
		// NBSP (U+00A0). Note this diverges from JS Intl, which puts the
		// amount before the symbol for fr-CA. Server-side templates that
		// surface currency should be aware of this drift; ideally callers
		// stick to numbers + a hand-written symbol until we get a CLDR-
		// driven currency formatter.
		{"en-CA", 15, "USD", "US$ 15.00"},
		{"fr-CA", 15, "USD", "$" + NBSP + "US 15,00"},
	}
	for _, c := range cases {
		got := formatCurrency(language.MustParse(c.tag), c.v, c.code)
		if got != c.expected {
			t.Errorf("formatCurrency(%s, %v, %s) = %q, want %q", c.tag, c.v, c.code, got, c.expected)
		}
	}
}

func TestFormatPercent(t *testing.T) {
	cases := []struct {
		tag      string
		v        float64
		expected string
	}{
		{"en-CA", 0.2, "20%"},
		// fr-CA: NBSP (U+00A0) between number and percent sign.
		{"fr-CA", 0.2, "20" + NBSP + "%"},
	}
	for _, c := range cases {
		got := formatPercent(language.MustParse(c.tag), c.v)
		if got != c.expected {
			t.Errorf("formatPercent(%s, %v) = %q, want %q", c.tag, c.v, got, c.expected)
		}
	}
}

func TestFormatBytes_CascadeMatchesJS(t *testing.T) {
	cases := []struct {
		tag      string
		v        float64
		base     string
		expected string
	}{
		{"en-CA", 1, "gigabyte", "1 GB"},
		{"fr-CA", 1, "gigabyte", "1 Go"},
		{"en-CA", 1500, "gigabyte", "1.46 TB"},
		{"en-CA", 1048576, "gigabyte", "1 PB"},
		{"fr-CA", 1048576, "gigabyte", "1 Po"},
	}
	for _, c := range cases {
		got := formatBytes(language.MustParse(c.tag), c.v, c.base)
		if got != c.expected {
			t.Errorf("formatBytes(%s, %v, %s) = %q, want %q", c.tag, c.v, c.base, got, c.expected)
		}
	}
}

func TestFormatDate_YearStyleIsLocaleNeutral(t *testing.T) {
	d := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	for _, tag := range []string{"en-CA", "fr-CA"} {
		got := formatDate(language.MustParse(tag), d, "year")
		if got != "2026" {
			t.Errorf("formatDate(%s, year) = %q, want %q", tag, got, "2026")
		}
	}
}

// TestFormatFuncMap_EndToEnd proves the FuncMap registration in
// RenderOne resolves the parsed language.Tag and wires through to the
// formatBytes helper. K-CNAS templates use [[ ]] delimiters.
func TestFormatFuncMap_EndToEnd(t *testing.T) {
	r := &Renderer{Language: "fr-CA"}
	got, err := r.RenderOne("smoke", `[[ bytes 1.0 "gigabyte" ]]`, TemplateData{})
	if err != nil {
		t.Fatalf("RenderOne returned error: %v", err)
	}
	want := "1 Go"
	if got != want {
		t.Errorf("RenderOne bytes = %q, want %q", got, want)
	}
}
