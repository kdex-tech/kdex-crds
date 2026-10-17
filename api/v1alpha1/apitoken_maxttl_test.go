package v1alpha1

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An unset maxTTL keeps today's unbounded mint, so it must not appear on the
// wire (an empty string would fail the duration rule). See host-manager#238.
func TestAPIToken_MaxTTL_OmittedWhenUnset(t *testing.T) {
	b, err := json.Marshal(APIToken{TokenPrefix: "acme_pat_"})
	require.NoError(t, err)
	assert.False(t, strings.Contains(string(b), "maxTTL"), "unset maxTTL must be omitted, got %s", b)
}

func TestAPIToken_MaxTTL_RoundTrips(t *testing.T) {
	b, err := json.Marshal(APIToken{MaxTTL: "2160h"})
	require.NoError(t, err)
	var out APIToken
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, "2160h", out.MaxTTL)
}

// Asserts on the generated CRDs so a marker that silently fails to apply is
// caught. CEL rejection itself is exercised by kdex-nexus-manager's envtest.
func TestAPITokenGeneratedSchema_MaxTTL(t *testing.T) {
	for _, file := range []string{"kdex.dev_kdexhosts.yaml", "kdex.dev_kdexinternalhosts.yaml"} {
		t.Run(file, func(t *testing.T) {
			spec := translationCRDSpecSchema(t, file)
			auth := spec["properties"].(map[string]any)["auth"].(map[string]any)
			apiToken := auth["properties"].(map[string]any)["apiToken"].(map[string]any)
			maxTTL, ok := apiToken["properties"].(map[string]any)["maxTTL"].(map[string]any)
			require.True(t, ok, "apiToken must expose maxTTL")

			assert.Equal(t, "string", maxTTL["type"])
			assert.EqualValues(t, 32, maxTTL["maxLength"])
			rules := maxTTL["x-kubernetes-validations"].([]any)
			require.Len(t, rules, 1)
			assert.Equal(t, "duration(self) > duration('0s')", rules[0].(map[string]any)["rule"])
		})
	}
}
