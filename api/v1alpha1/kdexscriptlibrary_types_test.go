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

package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"kdex.dev/crds/api/v1alpha1"
)

func TestPackageReference_ToScriptTag(t *testing.T) {
	tests := []struct {
		name       string
		packageRef v1alpha1.PackageReference
		want       string
	}{
		{
			name: "basic",
			packageRef: v1alpha1.PackageReference{
				Name:    "test",
				Version: "1.0.0",
			},
			want: `<script type="module">
  import "test";
</script>`,
		},
		{
			name: "export mapping",
			packageRef: v1alpha1.PackageReference{
				ExportMapping: "{ test }",
				Name:          "test",
				Version:       "1.0.0",
			},
			want: `<script type="module">
  import { test } from "test";
</script>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.packageRef.ToScriptTag()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestScript_ToScriptTag(t *testing.T) {
	tests := []struct {
		name       string
		script     v1alpha1.Script
		footScript bool
		want       string
	}{
		{
			name: "basic",
			script: v1alpha1.Script{
				Script: "test",
			},
			footScript: false,
			want: `<script>
test
</script>`,
		},
		{
			name: "no foot script",
			script: v1alpha1.Script{
				Script: "test",
			},
			footScript: true,
			want:       ``,
		},
		{
			name: "foot script",
			script: v1alpha1.Script{
				Script: "test",
			},
			footScript: false,
			want: `<script>
test
</script>`,
		},
		{
			name: "script attributes",
			script: v1alpha1.Script{
				Attributes: map[string]string{
					"src":      "/some/path",
					"data-foo": "some data",
				},
				Script: "test",
			},
			footScript: false,
			want: `<script data-foo="some data">
test
</script>`,
		},
		{
			name: "script src",
			script: v1alpha1.Script{
				ScriptSrc: "/some/path",
			},
			footScript: false,
			want:       `<script src="/some/path"></script>`,
		},
		{
			name: "script src attributes",
			script: v1alpha1.Script{
				Attributes: map[string]string{
					"src":      "/bad/path",
					"data-foo": "some data",
				},
				ScriptSrc: "/some/path",
			},
			footScript: false,
			want:       `<script data-foo="some data" src="/some/path"></script>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.script.ToScriptTag(tt.footScript)
			assert.Equal(t, tt.want, got)
		})
	}
}
