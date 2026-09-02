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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

// TestKDexPageSpec_TextPageFields decode-tests the localized/mimeType/body
// fields added so a KDexPage can serve a raw text document (robots.txt,
// llms.txt, sitemap.xml, ...) instead of composing HTML from an archetype.
// See docs-app/.../2026-09-02-text-pages-and-enumerated-l10n-routing.
func TestKDexPageSpec_TextPageFields(t *testing.T) {
	t.Run("text page: mimeType+body with no contentEntries decodes", func(t *testing.T) {
		specYaml := `
hostRef: { name: test-host }
label: robots.txt
basePath: /robots.txt
mimeType: txt
body: "User-agent: *\nDisallow:\n"
`
		var spec KDexPageSpec
		require.NoError(t, yaml.Unmarshal([]byte(specYaml), &spec))
		assert.Equal(t, "txt", spec.MimeType)
		assert.Equal(t, "User-agent: *\nDisallow:\n", spec.Body)
		assert.Empty(t, spec.ContentEntries)
	})

	t.Run("localized omitted decodes as nil (the apiserver applies the default, not encoding/json)", func(t *testing.T) {
		var spec KDexPageSpec
		require.NoError(t, yaml.Unmarshal([]byte(`hostRef: { name: test-host }`), &spec))
		assert.Nil(t, spec.Localized)
	})

	t.Run("localized explicit false round-trips through JSON", func(t *testing.T) {
		f := false
		spec := KDexPageSpec{Localized: &f}
		raw, err := json.Marshal(spec)
		require.NoError(t, err)

		var got map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(raw, &got))
		require.Contains(t, got, "localized")
		assert.JSONEq(t, "false", string(got["localized"]))
	})

	t.Run("localized omitted is absent from JSON (omitempty)", func(t *testing.T) {
		raw, err := json.Marshal(KDexPageSpec{})
		require.NoError(t, err)

		var got map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(raw, &got))
		assert.NotContains(t, got, "localized")
		assert.NotContains(t, got, "mimeType")
		assert.NotContains(t, got, "body")
	})
}

// TestKDexPageGeneratedSchema asserts on the *generated* CRD rather than on the
// kubebuilder marker, so a marker that silently fails to take (or a protobuf
// field-number collision) is caught here -- mirroring
// TestKDexFunctionBasePathSchema's pattern in kdexfunction_basepath_schema_test.go.
//
// kdex-crds has no envtest (see docs/superpowers/plans/2026-05-24-kdexfunction-service-backend.md:
// "the kdex-crds repo has no envtest, so CEL behavior is covered in the
// host-manager tests"), so this cannot exercise CEL *rejection* -- only that the
// rule text made it into the generated OpenAPI schema unchanged. Actual
// admission (including proving the two new spec-level CEL rules stay within the
// apiserver's per-rule cost budget) is exercised by kdex-host-manager's envtest
// suite in a later phase of this plan.
func TestKDexPageGeneratedSchema(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "config", "crd", "bases", "kdex.dev_kdexpages.yaml"))
	require.NoError(t, err)

	var crd struct {
		Spec struct {
			Versions []struct {
				Name   string `json:"name"`
				Schema struct {
					OpenAPIV3Schema struct {
						Properties struct {
							Spec struct {
								Required   []string `json:"required"`
								Properties struct {
									Localized struct {
										Type    string `json:"type"`
										Default bool   `json:"default"`
									} `json:"localized"`
									MimeType struct {
										Type string   `json:"type"`
										Enum []string `json:"enum"`
									} `json:"mimeType"`
									Body struct {
										Type      string `json:"type"`
										MaxLength int    `json:"maxLength"`
									} `json:"body"`
									ContentEntries struct {
										MinItems *int `json:"minItems"`
										MaxItems int  `json:"maxItems"`
									} `json:"contentEntries"`
								} `json:"properties"`
								XKubernetesValidations []struct {
									Rule    string `json:"rule"`
									Message string `json:"message"`
								} `json:"x-kubernetes-validations"`
							} `json:"spec"`
						} `json:"properties"`
					} `json:"openAPIV3Schema"`
				} `json:"schema"`
			} `json:"versions"`
		} `json:"spec"`
	}
	require.NoError(t, yaml.Unmarshal(raw, &crd))
	require.NotEmpty(t, crd.Spec.Versions)

	spec := crd.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties.Spec
	assert.Equal(t, "v1alpha1", crd.Spec.Versions[0].Name)

	assert.Equal(t, "boolean", spec.Properties.Localized.Type, "localized must be a boolean")
	assert.True(t, spec.Properties.Localized.Default, "localized must default to true")

	assert.Equal(t, "string", spec.Properties.MimeType.Type)
	assert.ElementsMatch(t, []string{"txt", "json", "yaml", "markdown", "xml"}, spec.Properties.MimeType.Enum,
		"mimeType enum must match the brief exactly")

	assert.Equal(t, "string", spec.Properties.Body.Type)
	assert.Equal(t, 65536, spec.Properties.Body.MaxLength)

	assert.Nil(t, spec.Properties.ContentEntries.MinItems, "contentEntries must no longer require minItems")
	assert.Equal(t, 32, spec.Properties.ContentEntries.MaxItems, "contentEntries must keep maxItems=32")
	assert.NotContains(t, spec.Required, "contentEntries", "contentEntries must no longer be spec-required")

	rules := make([]string, 0, len(spec.XKubernetesValidations))
	for _, r := range spec.XKubernetesValidations {
		rules = append(rules, r.Rule)
	}
	assert.Contains(t, rules, `has(self.mimeType) == has(self.body)`,
		"mimeType/body must-be-set-together CEL must be present")
	assert.Contains(t, rules, `has(self.mimeType) || (has(self.contentEntries) && self.contentEntries.exists(x, x.slot == 'main'))`,
		"main-slot-unless-text-page CEL must be present")
}
