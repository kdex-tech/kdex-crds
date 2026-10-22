package v1alpha1

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An unset hsts keeps the default header, so it must not appear on the wire.
// See host-manager#248.
func TestRouting_HSTS_OmittedWhenUnset(t *testing.T) {
	b, err := json.Marshal(Routing{Domains: []string{"a.example"}})
	require.NoError(t, err)
	assert.False(t, strings.Contains(string(b), `"hsts"`), "unset hsts must be omitted, got %s", b)
}

// An explicit maxAge of 0 is the back-out and must survive a round trip.
func TestRouting_HSTS_RoundTrips(t *testing.T) {
	in := Routing{
		Domains: []string{"a.example"},
		HSTS:    &HSTS{MaxAge: new(int32(0)), IncludeSubDomains: true},
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out Routing
	require.NoError(t, json.Unmarshal(b, &out))
	require.NotNil(t, out.HSTS)
	require.NotNil(t, out.HSTS.MaxAge)
	assert.EqualValues(t, 0, *out.HSTS.MaxAge)
	assert.True(t, out.HSTS.IncludeSubDomains)
	assert.False(t, out.HSTS.Preload)
}

// Asserts on the generated CRDs so a marker that silently fails to apply is
// caught. CEL rejection itself is exercised by kdex-nexus-manager's envtest.
func TestRoutingGeneratedSchema_HSTS(t *testing.T) {
	for _, file := range []string{"kdex.dev_kdexhosts.yaml", "kdex.dev_kdexinternalhosts.yaml"} {
		t.Run(file, func(t *testing.T) {
			spec := translationCRDSpecSchema(t, file)
			routing := spec["properties"].(map[string]any)["routing"].(map[string]any)
			hsts, ok := routing["properties"].(map[string]any)["hsts"].(map[string]any)
			require.True(t, ok, "routing must expose hsts")

			maxAge := hsts["properties"].(map[string]any)["maxAge"].(map[string]any)
			assert.EqualValues(t, 0, maxAge["minimum"])
			assert.EqualValues(t, 63072000, maxAge["maximum"])

			rules := hsts["x-kubernetes-validations"].([]any)
			require.Len(t, rules, 1)
			assert.Equal(t,
				"!has(self.preload) || !self.preload || (has(self.includeSubDomains) && self.includeSubDomains && (!has(self.maxAge) || self.maxAge >= 31536000))",
				rules[0].(map[string]any)["rule"])
		})
	}
}
