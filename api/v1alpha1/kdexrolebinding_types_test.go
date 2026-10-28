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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestKDexRoleBinding_SubjectKindsDecodes(t *testing.T) {
	var rb KDexRoleBinding
	require.NoError(t, yaml.Unmarshal([]byte(`
apiVersion: kdex.dev/v1alpha1
kind: KDexRoleBinding
metadata: { name: rb, namespace: site }
spec:
  hostRef: { name: site-host }
  roles: [reader]
  subject: "*"
  subjectKinds: [User, Client]
`), &rb))

	assert.Equal(t, []SubjectKind{SubjectKindUser, SubjectKindClient}, rb.Spec.SubjectKinds)
}

func TestKDexRoleBinding_WithoutSubjectKindsIsUnchanged(t *testing.T) {
	var rb KDexRoleBinding
	require.NoError(t, yaml.Unmarshal([]byte(`
apiVersion: kdex.dev/v1alpha1
kind: KDexRoleBinding
metadata: { name: rb, namespace: site }
spec:
  hostRef: { name: site-host }
  roles: [reader]
  subject: alice
`), &rb))

	assert.Nil(t, rb.Spec.SubjectKinds)
	out, err := json.Marshal(rb.Spec)
	require.NoError(t, err)
	assert.NotContains(t, string(out), "subjectKinds")
}

func TestSubjectKind_Known(t *testing.T) {
	assert.True(t, SubjectKindUser.Known())
	assert.True(t, SubjectKindClient.Known())
	for _, k := range []SubjectKind{"", "ServiceAccount", "user"} {
		assert.Falsef(t, k.Known(), "%q must not be a known kind", k)
	}
}

// Asserts on the generated CRD so a marker that silently fails to apply is
// caught.
func TestKDexRoleBindingGeneratedSchema_SubjectKinds(t *testing.T) {
	spec := crdSpecSchema(t, "kdex.dev_kdexrolebindings.yaml")
	sk, ok := spec["properties"].(map[string]any)["subjectKinds"].(map[string]any)
	require.True(t, ok, "KDexRoleBinding spec must expose subjectKinds")

	rawEnum := sk["items"].(map[string]any)["enum"].([]any)
	enum := make([]string, 0, len(rawEnum))
	for _, e := range rawEnum {
		enum = append(enum, e.(string))
	}
	want := make([]string, 0, len(SubjectKinds))
	for _, k := range SubjectKinds {
		want = append(want, string(k))
	}
	assert.Equal(t, want, enum)

	assert.EqualValues(t, 1, sk["minItems"])
	assert.EqualValues(t, 8, sk["maxItems"])
	assert.Equal(t, "set", sk["x-kubernetes-list-type"])

	required, _ := spec["required"].([]any)
	assert.NotContains(t, required, "subjectKinds")
}
