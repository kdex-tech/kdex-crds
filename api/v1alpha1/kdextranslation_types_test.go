package v1alpha1

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestKDexTranslation_HostRefDecodes(t *testing.T) {
	var tr KDexTranslation
	require.NoError(t, yaml.Unmarshal([]byte(`
apiVersion: kdex.dev/v1alpha1
kind: KDexTranslation
metadata: { name: shop-strings, namespace: site }
spec:
  hostRef: { name: site-host }
  translations:
    - lang: en
      keysAndValues: { shop.title: Shop }
`), &tr))

	require.NotNil(t, tr.Spec.HostRef)
	assert.Equal(t, "site-host", tr.Spec.HostRef.Name)
	require.Len(t, tr.Spec.Translations, 1)
	assert.Equal(t, "Shop", tr.Spec.Translations[0].KeysAndValues["shop.title"])
}

// A manifest written before hostRef existed must decode unchanged and must not
// grow a hostRef key when re-serialized.
func TestKDexTranslation_WithoutHostRefIsUnchanged(t *testing.T) {
	var tr KDexTranslation
	require.NoError(t, yaml.Unmarshal([]byte(`
apiVersion: kdex.dev/v1alpha1
kind: KDexTranslation
metadata: { name: legacy, namespace: site }
spec:
  translations:
    - lang: en
      keysAndValues: { a: b }
`), &tr))

	assert.Nil(t, tr.Spec.HostRef)
	out, err := yaml.Marshal(tr.Spec)
	require.NoError(t, err)
	assert.NotContains(t, string(out), "hostRef")
}

func translationCRDSpecSchema(t *testing.T, file string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "config", "crd", "bases", file))
	require.NoError(t, err)
	var crd map[string]any
	require.NoError(t, yaml.Unmarshal(raw, &crd))
	versions := crd["spec"].(map[string]any)["versions"].([]any)
	schema := versions[0].(map[string]any)["schema"].(map[string]any)["openAPIV3Schema"].(map[string]any)
	return schema["properties"].(map[string]any)["spec"].(map[string]any)
}

// Asserts on the generated CRD so a marker that silently fails to apply is
// caught. CEL rejection itself is exercised by kdex-nexus-manager's envtest
// (kdex-crds has no envtest).
func TestKDexTranslationGeneratedSchema_HostRef(t *testing.T) {
	spec := translationCRDSpecSchema(t, "kdex.dev_kdextranslations.yaml")
	props := spec["properties"].(map[string]any)

	hostRef, ok := props["hostRef"].(map[string]any)
	require.True(t, ok, "KDexTranslation spec must expose hostRef")
	rules := hostRef["x-kubernetes-validations"].([]any)
	require.Len(t, rules, 1)
	assert.Equal(t, "self.name.size() > 0", rules[0].(map[string]any)["rule"])
	assert.Equal(t, "hostRef.name must not be empty", rules[0].(map[string]any)["message"])

	required, _ := spec["required"].([]any)
	assert.NotContains(t, required, "hostRef", "hostRef must stay optional")
	assert.Contains(t, required, "translations")
}

func TestKDexClusterTranslationGeneratedSchema_HasNoHostRef(t *testing.T) {
	spec := translationCRDSpecSchema(t, "kdex.dev_kdexclustertranslations.yaml")
	_, has := spec["properties"].(map[string]any)["hostRef"]
	assert.False(t, has, "a cluster-scoped translation cannot attach to a namespaced host")
}
