package v1alpha1

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMintToken_CapabilityTTLCapSeconds_OmittedWhenUnset(t *testing.T) {
	b, err := json.Marshal(MintToken{Enabled: true})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(b), "capabilityTtlCapSeconds") {
		t.Fatalf("unset capabilityTtlCapSeconds must be omitted, got %s", b)
	}
}

func TestMintToken_CapabilityTTLCapSeconds_RoundTrips(t *testing.T) {
	in := MintToken{Enabled: true, TTLCapSeconds: 60, CapabilityTTLCapSeconds: 600}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out MintToken
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.CapabilityTTLCapSeconds != 600 {
		t.Fatalf("round-trip mismatch: got %d, want 600", out.CapabilityTTLCapSeconds)
	}
	// The REST cap is independent of the MCP ttlCapSeconds — both survive.
	if out.TTLCapSeconds != 60 {
		t.Fatalf("ttlCapSeconds clobbered: got %d, want 60", out.TTLCapSeconds)
	}
}
