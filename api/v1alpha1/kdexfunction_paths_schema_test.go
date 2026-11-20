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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Asserts on the generated CRD so a marker that silently fails to apply is
// caught. The paths bound is the only gate on how many routes a function can
// declare, and it feeds the CEL cost estimate of the API-level startsWith rule,
// so a change here must be deliberate. See #35.
func TestKDexFunctionGeneratedSchema_PathsBounds(t *testing.T) {
	spec := crdSpecSchema(t, "kdex.dev_kdexfunctions.yaml")
	api := spec["properties"].(map[string]any)["api"].(map[string]any)
	paths, ok := api["properties"].(map[string]any)["paths"].(map[string]any)
	require.True(t, ok, "spec.api must expose paths")

	assert.EqualValues(t, 1, paths["minProperties"])
	assert.EqualValues(t, 32, paths["maxProperties"])
}
