package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestPolicyRule_OpaqueScopesDecode(t *testing.T) {
	specYaml := `
hostRef: { name: test-host }
rules:
  - scopes: [vector_stores_create]
  - resources: [vector_stores]
    resourceNames: ["vs_alice"]
    verbs: [read]
`
	var spec KDexRoleSpec
	err := yaml.Unmarshal([]byte(specYaml), &spec)
	assert.NoError(t, err)

	assert.Len(t, spec.Rules, 2)
	assert.Equal(t, []string{"vector_stores_create"}, spec.Rules[0].Scopes)
	assert.Empty(t, spec.Rules[0].Resources)
	assert.Empty(t, spec.Rules[0].Verbs)
	assert.Equal(t, []string{"vector_stores"}, spec.Rules[1].Resources)
	assert.Equal(t, []string{"read"}, spec.Rules[1].Verbs)
}
