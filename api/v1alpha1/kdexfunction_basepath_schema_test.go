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
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

// TestKDexFunctionBasePathSchema asserts on the *generated* CRD rather than on
// the kubebuilder marker, so that a marker which silently fails to take is
// caught here. spec.api.basePath is interpolated by kdex-host-manager into the
// RFC 9728 resource_metadata parameter of a WWW-Authenticate header (an HTTP
// quoted-string), so the pattern must be anchored at both ends and must admit
// no quote, backslash, or whitespace.
func TestKDexFunctionBasePathSchema(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "config", "crd", "bases", "kdex.dev_kdexfunctions.yaml"))
	require.NoError(t, err)

	var crd struct {
		Spec struct {
			Versions []struct {
				Name   string `json:"name"`
				Schema struct {
					OpenAPIV3Schema struct {
						Properties struct {
							Spec struct {
								Properties struct {
									API struct {
										Properties struct {
											BasePath struct {
												Pattern   string `json:"pattern"`
												MaxLength int    `json:"maxLength"`
											} `json:"basePath"`
										} `json:"properties"`
									} `json:"api"`
								} `json:"properties"`
							} `json:"spec"`
						} `json:"properties"`
					} `json:"openAPIV3Schema"`
				} `json:"schema"`
			} `json:"versions"`
		} `json:"spec"`
	}
	require.NoError(t, yaml.Unmarshal(raw, &crd))
	require.NotEmpty(t, crd.Spec.Versions)

	basePath := crd.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties.Spec.Properties.API.Properties.BasePath
	assert.Equal(t, "v1alpha1", crd.Spec.Versions[0].Name)
	assert.Equal(t, `^(/\w[\w-]*){2,8}$`, basePath.Pattern, "generated CRD lost the basePath pattern")
	assert.Equal(t, 64, basePath.MaxLength, "generated CRD lost the basePath maxLength")

	// The apiserver applies `pattern` with Go's regexp (RE2, MatchString), where
	// ^ and $ bind to the ends of the text. Compile it exactly as it ships.
	re, err := regexp.Compile(basePath.Pattern)
	require.NoError(t, err)

	// Every basePath in real use across the kdex/RSI corpus must keep validating.
	// Narrowing this set is a breaking change to live KDexFunction CRs.
	inUse := []string{
		"/v1/users", "/v1/auth", "/v1/admin", "/v1/chat",
		"/tenant/v1", "/invite/v1", "/billing/v1", "/feedback/v1",
		"/contracts/v1", "/delta/v1", "/api/echo",
		"/v1/credential-check", // hyphen: \w alone does not cover it
		"/api/v1/files", "/api/v1/mcp", "/api/v1/search",
		"/api/v1/events", "/api/v1/ingest", "/api/v1/uploads",
		"/api/v1/vector_stores", // longest in use, 21 chars
	}
	for _, v := range inUse {
		assert.Truef(t, re.MatchString(v), "in-use basePath %q must remain valid", v)
		assert.LessOrEqualf(t, len(v), basePath.MaxLength, "in-use basePath %q exceeds maxLength", v)
	}

	rejected := map[string]string{
		`/a/b",resource_metadata="https://attacker.example/x`: "WWW-Authenticate quoted-string injection",
		`/v1/users"`:           "bare quote",
		`/v1/users\`:           "backslash",
		"/v1/users x":          "whitespace",
		"/v1/users\nX-Evil: 1": "newline / header split",
		"/foo":                 "single segment",
		"/":                    "root",
		"foo/bar":              "not absolute",
		"/-/openapi":           "host-reserved /-/ prefix",
		"/v1/../../etc":        "dot segments",
	}
	for v, why := range rejected {
		assert.Falsef(t, re.MatchString(v), "basePath %q must be rejected (%s)", v, why)
	}
}
