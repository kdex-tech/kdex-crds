package v1alpha1

import (
	"testing"

	"github.com/kdex-tech/dmapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestKDexHostExtension_Decodes(t *testing.T) {
	var ext KDexHostExtension
	require.NoError(t, yaml.Unmarshal([]byte(`
apiVersion: kdex.dev/v1alpha1
kind: KDexHostExtension
metadata: { name: eum, namespace: site, labels: { kdex.dev/extension: eum } }
spec:
  hostRef: { name: site-host }
  weight: 100
  claimMappings:
    - sourceExpression: "has(self.eum_entitlements) ? self.eum_entitlements : []"
      targetPropPath: entitlements
  anonymousEntitlements: [ "functions:/eum/public:read" ]
`), &ext))
	assert.Equal(t, "site-host", ext.Spec.HostRef.Name)
	assert.Equal(t, int32(100), ext.Spec.Weight)
	require.Len(t, ext.Spec.ClaimMappings, 1)
	assert.Equal(t, "entitlements", ext.Spec.ClaimMappings[0].TargetPropPath)
	assert.Equal(t, []string{"functions:/eum/public:read"}, ext.Spec.AnonymousEntitlements)
}

func TestKDexHostExtensionGeneratedSchema(t *testing.T) {
	spec := translationCRDSpecSchema(t, "kdex.dev_kdexhostextensions.yaml")
	props := spec["properties"].(map[string]any)

	hostRef := props["hostRef"].(map[string]any)
	rules := hostRef["x-kubernetes-validations"].([]any)
	require.Len(t, rules, 1)
	assert.Equal(t, "hostRef.name must not be empty", rules[0].(map[string]any)["message"])
	assert.Contains(t, spec["required"].([]any), "hostRef")

	weight := props["weight"].(map[string]any)
	assert.EqualValues(t, -1000, weight["minimum"])
	assert.EqualValues(t, 1000, weight["maximum"])
	assert.EqualValues(t, 0, weight["default"])

	cm := props["claimMappings"].(map[string]any)
	assert.EqualValues(t, 16, cm["maxItems"])
	anon := props["anonymousEntitlements"].(map[string]any)
	assert.EqualValues(t, 64, anon["maxItems"])
	assert.EqualValues(t, 256, anon["items"].(map[string]any)["maxLength"])
}

func TestKDexHostGeneratedSchema_ExtensionSelector(t *testing.T) {
	spec := translationCRDSpecSchema(t, "kdex.dev_kdexhosts.yaml")
	_, ok := spec["properties"].(map[string]any)["extensionSelector"]
	assert.True(t, ok, "KDexHost spec must expose extensionSelector")
	assert.NotContains(t, spec["required"], "extensionSelector")
}

func TestKDexInternalHostGeneratedSchema_Extensions(t *testing.T) {
	spec := translationCRDSpecSchema(t, "kdex.dev_kdexinternalhosts.yaml")
	ext, ok := spec["properties"].(map[string]any)["extensions"].(map[string]any)
	require.True(t, ok, "KDexInternalHost spec must expose extensions")
	assert.EqualValues(t, 32, ext["maxItems"])
}

func rule(target, expr string) dmapper.MappingRule {
	return dmapper.MappingRule{SourceExpression: expr, TargetPropPath: target}
}

func TestEffectiveAuth_NoExtensionsReturnsHostAuth(t *testing.T) {
	a := &Auth{ClaimMappings: []dmapper.MappingRule{rule("entitlements", "self.x")}}
	s := &KDexInternalHostSpec{KDexHostSpec: KDexHostSpec{Auth: a}}
	assert.Same(t, a, s.EffectiveAuth())
}

func TestEffectiveAuth_NilAuthIgnoresExtensions(t *testing.T) {
	s := &KDexInternalHostSpec{Extensions: []InternalHostExtension{{Name: "x", AnonymousEntitlements: []string{"pages:/a:read"}}}}
	assert.Nil(t, s.EffectiveAuth())
}

func TestEffectiveAuth_ComposesInOrder(t *testing.T) {
	hostRule := rule("entitlements", "self.host_grants")
	extA := rule("entitlements", "self.a_grants")
	extB := rule("roles", "self.b_roles")
	s := &KDexInternalHostSpec{
		KDexHostSpec: KDexHostSpec{Auth: &Auth{
			ClaimMappings:         []dmapper.MappingRule{hostRule},
			AnonymousEntitlements: []string{"pages:/home:read"},
		}},
		Extensions: []InternalHostExtension{
			{Name: "a", ClaimMappings: []dmapper.MappingRule{extA}, AnonymousEntitlements: []string{"pages:/home:read", "functions:/a:read"}},
			{Name: "b", ClaimMappings: []dmapper.MappingRule{extB}, AnonymousEntitlements: []string{"functions:/b:read"}},
		},
	}
	got := s.EffectiveAuth()
	assert.Equal(t, []dmapper.MappingRule{hostRule, extA, extB}, got.ClaimMappings, "host rules first, then extensions in list order")
	assert.Equal(t, []string{"pages:/home:read", "functions:/a:read", "functions:/b:read"}, got.AnonymousEntitlements, "order-preserving, de-duplicated union")

	// The receiver is not mutated.
	assert.Equal(t, []dmapper.MappingRule{hostRule}, s.Auth.ClaimMappings)
	assert.Equal(t, []string{"pages:/home:read"}, s.Auth.AnonymousEntitlements)
}

func TestKDexHostExtensionGeneratedSchema_ValidationRules(t *testing.T) {
	props := translationCRDSpecSchema(t, "kdex.dev_kdexhostextensions.yaml")["properties"].(map[string]any)

	messages := func(items map[string]any) []string {
		rules := items["x-kubernetes-validations"].([]any)
		out := make([]string, 0, len(rules))
		for _, r := range rules {
			out = append(out, r.(map[string]any)["message"].(string))
		}
		return out
	}

	cmItems := props["claimMappings"].(map[string]any)["items"].(map[string]any)
	assert.ElementsMatch(t, []string{
		"an extension claimMapping cannot use merge: Replace",
		"an extension claimMapping must not target a reserved token claim",
	}, messages(cmItems))

	anonItems := props["anonymousEntitlements"].(map[string]any)["items"].(map[string]any)
	assert.Equal(t, []string{"anonymousEntitlements must be resource:name:verb with a name other than empty or *"}, messages(anonItems))
}
