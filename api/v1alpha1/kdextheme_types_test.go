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

func TestKDexThemeSpec_String(t *testing.T) {
	tests := []struct {
		name   string
		assets []v1alpha1.Asset
		want   string
	}{
		{
			name:   "empty",
			assets: []v1alpha1.Asset{},
			want:   ``,
		},
		{
			name: "style",
			assets: []v1alpha1.Asset{
				{
					Style: `color: #fff;`,
				},
			},
			want: `<style>
color: #fff;
</style>`,
		},
		{
			name: "style with attributes",
			assets: []v1alpha1.Asset{
				{
					Attributes: map[string]string{
						"data-foo": "some data",
					},
					Style: `color: #fff;`,
				},
			},
			want: `<style data-foo="some data">
color: #fff;
</style>`,
		},
		{
			name: "link href",
			assets: []v1alpha1.Asset{
				{
					LinkHref: "/some/path",
				},
			},
			want: `<link href="/some/path"/>`,
		},
		{
			name: "link href with attributes",
			assets: []v1alpha1.Asset{
				{
					Attributes: map[string]string{
						"data-foo": "some data",
						"href":     "/some/path",
					},
					LinkHref: "/some/path",
				},
			},
			want: `<link data-foo="some data" href="/some/path"/>`,
		},
		{
			name: "both link href and style",
			assets: []v1alpha1.Asset{
				{
					Attributes: map[string]string{
						"data-foo": "some data",
					},
					LinkHref: "/some/path",
				},
				{
					Attributes: map[string]string{
						"data-foo": "some data",
					},
					Style: `color: #fff;`,
				},
			},
			want: `<link data-foo="some data" href="/some/path"/>
<style data-foo="some data">
color: #fff;
</style>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := v1alpha1.KDexThemeSpec{
				Assets: tt.assets,
			}
			got := theme.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
