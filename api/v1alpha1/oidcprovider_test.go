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
)

// TestOIDCProviderScopesSerializeAsScopes pins the rename in kdex-crds#21.
//
// The Go field is Scopes, the consumer appends the values to the OAuth 2.0
// scope list, and the JSON key was `roles`. An author asking for a scope had
// to write it under a key named for something else, and an author who wrote
// an actual role name there had it sent to the IdP as a scope -- which most
// providers reject the whole authorization request over. host-manager#190
// makes that worse by giving the field a real, documented use
// (`offline_access`), so the misleading name now has a reason to be typed.
func TestOIDCProviderScopesSerializeAsScopes(t *testing.T) {
	raw, err := json.Marshal(OIDCProvider{
		OIDCProviderURL: "https://accounts.google.com",
		Scopes:          []string{"offline_access"},
	})
	require.NoError(t, err)

	var got map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw, &got))

	assert.Contains(t, got, "scopes",
		"the OAuth scope list must serialize under `scopes` (#21), got %s", raw)
	assert.NotContains(t, got, "roles",
		"`roles` was never what this field meant (#21), got %s", raw)

	var decoded OIDCProvider
	require.NoError(t, json.Unmarshal([]byte(`{"scopes":["offline_access"]}`), &decoded))
	assert.Equal(t, []string{"offline_access"}, decoded.Scopes,
		"a CR written against the documented key must decode")
}
