package v1alpha1

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

// Every generated CRD must be listed in config/crd/kustomization.yaml, or
// `make build-installer` ships a release install.yaml without it and any
// controller that watches the missing kind crash-loops.
func TestCRDKustomizationListsEveryBase(t *testing.T) {
	crdDir := filepath.Join("..", "..", "config", "crd")

	raw, err := os.ReadFile(filepath.Join(crdDir, "kustomization.yaml"))
	require.NoError(t, err)
	var k struct {
		Resources []string `json:"resources"`
	}
	require.NoError(t, yaml.Unmarshal(raw, &k))

	bases, err := filepath.Glob(filepath.Join(crdDir, "bases", "*.yaml"))
	require.NoError(t, err)
	require.NotEmpty(t, bases)

	for _, b := range bases {
		want := "bases/" + filepath.Base(b)
		assert.Contains(t, k.Resources, want, "config/crd/kustomization.yaml resources must list %s", want)
	}
}
