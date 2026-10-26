package openapi

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"

	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	kdexv1alpha1 "kdex.dev/crds/api/v1alpha1"
	"sigs.k8s.io/yaml"
)

func TestBuildOpenAPI(t *testing.T) {
	serverUrl := "http://test"
	tests := []struct {
		name         string
		functionName string
		paths        map[string]PathInfo
		filter       Filter
		assertions   func(t *testing.T, o *openapi.T)
	}{
		{
			name:         "no paths",
			functionName: "no-paths",
			assertions: func(t *testing.T, o *openapi.T) {
				assert.NotNil(t, o)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Builder{
				TypesToInclude: []PathType{
					BackendPathType, FunctionPathType, PagePathType, SystemPathType,
				},
			}
			got := b.BuildOpenAPI(serverUrl, tt.functionName, tt.paths, tt.filter)
			tt.assertions(t, got)
		})
	}
}

// TestBuildOneOffByteIdentical pins the one-off document byte-for-byte against
// the openapi-spec.json the in-cluster codegen Job committed for real
// KDexFunctions (host-manager internal/generate marshals it with
// json.MarshalIndent(doc, "", "  "); entry-point.sh's echo adds the trailing
// newline). Client-side codegen depends on this: if the document drifts, a
// function regenerated locally differs from the one the cluster generates.
//
//   - tenancy-service: schemas, plus the used-scheme filter (only the
//     referenced bearer scheme and the Unauthorized response are emitted).
//   - user-service-auth: a public function that references no scheme.
func TestBuildOneOffByteIdentical(t *testing.T) {
	const serverURL = "https://dev.knowdrive.ai"

	for _, name := range []string{"tenancy-service", "user-service-auth"} {
		t.Run(name, func(t *testing.T) {
			crBytes, err := os.ReadFile(filepath.Join("testdata", "oneoff", name+".kdexfunction.yaml"))
			require.NoError(t, err)
			var fn kdexv1alpha1.KDexFunction
			require.NoError(t, yaml.Unmarshal(crBytes, &fn))

			got, err := json.MarshalIndent(NewOneOffBuilder().BuildOneOff(serverURL, &fn), "", "  ")
			require.NoError(t, err)
			got = append(got, '\n')

			want, err := os.ReadFile(filepath.Join("testdata", "oneoff", name+".openapi-spec.json"))
			require.NoError(t, err)
			assert.Equal(t, string(want), string(got))
		})
	}
}

// The one-off document must carry only the schemes its operations reference,
// even though the builder knows all of them.
func TestBuildOneOffEmitsOnlyReferencedSchemes(t *testing.T) {
	fn := &kdexv1alpha1.KDexFunction{}
	fn.Name = "fn"
	fn.Spec.API.BasePath = "/v1/fn"
	item := kdexv1alpha1.PathItem{}
	item.SetGet(&openapi.Operation{
		OperationID: "getThing",
		Security:    &openapi.SecurityRequirements{{"apiKeyHeader": {"things:*:read"}}},
		Responses:   openapi.NewResponses(),
	})
	fn.Spec.API.Paths = map[string]kdexv1alpha1.PathItem{"/v1/fn/thing": item}

	doc := NewOneOffBuilder().BuildOneOff("http://test", fn)

	assert.Equal(t, []string{"apiKeyHeader"}, slices.Collect(maps.Keys(doc.Components.SecuritySchemes)))
	assert.NotNil(t, doc.Components.Responses["Unauthorized"])
	assert.Empty(t, doc.Security, "a one-off document declares no document-level security")
}

func TestSecuritySchemes(t *testing.T) {
	schemes := SecuritySchemes()
	names := slices.Sorted(maps.Keys(*schemes))
	assert.Equal(t, []string{"apiKeyCookie", "apiKeyHeader", "apiKeyQuery", "bearer", "oauth2", "oidc"}, names)
	assert.Equal(t, OAuthAuthorizePath, (*schemes)["oauth2"].Value.Flows.AuthorizationCode.AuthorizationURL)
}
