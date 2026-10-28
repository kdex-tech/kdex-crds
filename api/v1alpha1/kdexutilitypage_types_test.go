package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKDexUtilityPageGeneratedSchema_Consent(t *testing.T) {
	for _, file := range []string{
		"kdex.dev_kdexutilitypages.yaml",
		"kdex.dev_kdexclusterutilitypages.yaml",
		"kdex.dev_kdexinternalutilitypages.yaml",
	} {
		t.Run(file, func(t *testing.T) {
			spec := crdSpecSchema(t, file)
			typ := spec["properties"].(map[string]any)["type"].(map[string]any)
			// The type's enum is emitted once per marker, under allOf.
			allOf := typ["allOf"].([]any)
			require.NotEmpty(t, allOf)
			for _, e := range allOf {
				assert.Contains(t, e.(map[string]any)["enum"], "Consent")
			}

			rules, ok := spec["x-kubernetes-validations"].([]any)
			require.True(t, ok, "spec must carry validations")
			found := false
			for _, r := range rules {
				rm := r.(map[string]any)
				if rm["message"] == "a Consent page cannot reference scripts: no scriptLibraryRef and no appRef content" {
					found = true
					assert.Equal(t, "self.type != 'Consent' || (!has(self.scriptLibraryRef) && self.contentEntries.all(e, !has(e.appRef)))", rm["rule"])
				}
			}
			assert.True(t, found, "Consent scripts rule missing")
		})
	}
}

func TestKDexHostGeneratedSchema_ConsentRef(t *testing.T) {
	spec := crdSpecSchema(t, "kdex.dev_kdexhosts.yaml")
	up := spec["properties"].(map[string]any)["utilityPages"].(map[string]any)
	ref, ok := up["properties"].(map[string]any)["consentRef"].(map[string]any)
	require.True(t, ok, "utilityPages must expose consentRef")
	rules := ref["x-kubernetes-validations"].([]any)
	require.Len(t, rules, 1)
	rule := rules[0].(map[string]any)["rule"].(string)
	assert.Contains(t, rule, "KDexUtilityPage")
	assert.Contains(t, rule, "KDexClusterUtilityPage")
}

func TestKDexInternalHostGeneratedSchema_ConsentRef(t *testing.T) {
	spec := crdSpecSchema(t, "kdex.dev_kdexinternalhosts.yaml")
	_, ok := spec["properties"].(map[string]any)["consentRef"]
	assert.True(t, ok, "internal host spec must expose consentRef")
}
