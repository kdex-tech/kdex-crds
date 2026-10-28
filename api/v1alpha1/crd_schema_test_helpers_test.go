/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

// crdVersion loads a generated CRD from config/crd/bases and returns its first
// version entry.
func crdVersion(t *testing.T, file string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "config", "crd", "bases", file))
	require.NoError(t, err)
	var crd map[string]any
	require.NoError(t, yaml.Unmarshal(raw, &crd))
	versions := crd["spec"].(map[string]any)["versions"].([]any)
	require.NotEmpty(t, versions)
	return versions[0].(map[string]any)
}

// crdSpecSchema returns the openAPIV3Schema of .spec from a generated CRD, so
// tests assert on the generated schema rather than on the kubebuilder markers.
func crdSpecSchema(t *testing.T, file string) map[string]any {
	t.Helper()
	schema := crdVersion(t, file)["schema"].(map[string]any)["openAPIV3Schema"].(map[string]any)
	return schema["properties"].(map[string]any)["spec"].(map[string]any)
}

// decodeSchema re-decodes a generic schema map into a typed struct subset.
func decodeSchema(t *testing.T, in map[string]any, out any) {
	t.Helper()
	b, err := json.Marshal(in)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(b, out))
}
