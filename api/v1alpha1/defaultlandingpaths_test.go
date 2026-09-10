package v1alpha1

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAuth_DefaultLandingPaths_OmittedWhenUnset(t *testing.T) {
	b, err := json.Marshal(Auth{})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(b), "defaultLandingPaths") {
		t.Fatalf("unset defaultLandingPaths must be omitted, got %s", b)
	}
}

func TestAuth_DefaultLandingPaths_RoundTrips(t *testing.T) {
	in := Auth{DefaultLandingPaths: []string{"/home", "/dashboard"}}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out Auth
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.DefaultLandingPaths) != 2 ||
		out.DefaultLandingPaths[0] != "/home" ||
		out.DefaultLandingPaths[1] != "/dashboard" {
		t.Fatalf("round-trip mismatch: %+v", out.DefaultLandingPaths)
	}
}
