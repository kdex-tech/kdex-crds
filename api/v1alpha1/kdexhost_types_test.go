package v1alpha1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestKDexHostSpec_Helm(t *testing.T) {
	specYaml := `
brandName: "Test Brand"
organization: "KDex Tech"
routing:
  domains: ["test.kdex.dev"]
helm:
  hostManager:
    values: |
      replicaCount: 2
      valkey:
        enabled: false
  companionCharts:
    - name: my-companion
      chart: some-chart
      repository: https://charts.example.com
      version: 1.2.3
      values: |
        foo: bar
`
	var spec KDexHostSpec
	err := yaml.Unmarshal([]byte(specYaml), &spec)
	assert.NoError(t, err)

	assert.NotNil(t, spec.Helm)
	assert.NotNil(t, spec.Helm.HostManager)
	assert.Equal(t, "replicaCount: 2\nvalkey:\n  enabled: false\n", spec.Helm.HostManager.Values)

	assert.Len(t, spec.Helm.CompanionCharts, 1)
	assert.Equal(t, "my-companion", spec.Helm.CompanionCharts[0].Name)
	assert.Equal(t, "some-chart", spec.Helm.CompanionCharts[0].Chart)
	assert.Equal(t, "https://charts.example.com", spec.Helm.CompanionCharts[0].Repository)
	assert.Equal(t, "1.2.3", spec.Helm.CompanionCharts[0].Version)
	assert.Equal(t, "foo: bar\n", spec.Helm.CompanionCharts[0].Values)

	// Verify JSON marshaling/unmarshaling
	jsonData, err := json.Marshal(spec)
	assert.NoError(t, err)

	var spec2 KDexHostSpec
	err = json.Unmarshal(jsonData, &spec2)
	assert.NoError(t, err)
	assert.Equal(t, spec.Helm, spec2.Helm)
}
