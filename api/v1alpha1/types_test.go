package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssets_String(t *testing.T) {
	tests := []struct {
		name   string
		assets Assets
		want   string
	}{
		{
			name:   "empty",
			assets: Assets{},
			want:   ``,
		},
		{
			name: "style",
			assets: Assets{
				{
					StyleDef: StyleDef{
						Style: Ptr(`color: #fff;`),
					},
				},
			},
			want: `<style>
color: #fff;
</style>`,
		},
		{
			name: "style with attributes",
			assets: Assets{
				{
					StyleDef: StyleDef{
						Attributes: map[string]string{
							"data-foo": "some data",
						},
						Style: Ptr(`color: #fff;`),
					},
				},
			},
			want: `<style data-foo="some data">
color: #fff;
</style>`,
		},
		{
			name: "link href",
			assets: Assets{
				{
					LinkDef: LinkDef{
						LinkHref: Ptr("/some/path"),
					},
				},
			},
			want: `<link href="/some/path"/>`,
		},
		{
			name: "link href with attributes",
			assets: Assets{
				{
					LinkDef: LinkDef{
						Attributes: map[string]string{
							"data-foo": "some data",
							"href":     "/some/path",
						},
						LinkHref: Ptr("/some/path"),
					},
				},
			},
			want: `<link data-foo="some data" href="/some/path"/>`,
		},
		{
			name: "both link href and style",
			assets: Assets{
				{
					LinkDef: LinkDef{
						Attributes: map[string]string{
							"data-foo": "some data",
						},
						LinkHref: Ptr("/some/path"),
					},
				},
				{
					StyleDef: StyleDef{
						Attributes: map[string]string{
							"data-foo": "some data",
						},
						Style: Ptr(`color: #fff;`),
					},
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
			got := tt.assets.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPackageReference_ToScriptTag(t *testing.T) {
	tests := []struct {
		name       string
		packageRef PackageReference
		want       string
	}{
		{
			name: "basic",
			packageRef: PackageReference{
				Name:    "test",
				Version: "1.0.0",
			},
			want: `<script type="module">
  import "test";
</script>`,
		},
		{
			name: "export mapping",
			packageRef: PackageReference{
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

func TestScript_ToHeadTag(t *testing.T) {
	tests := []struct {
		name   string
		script ScriptDef
		want   string
	}{
		{
			name: "basic",
			script: ScriptDef{
				Script: Ptr("test"),
			},
			want: `<script>
test
</script>`,
		},
		{
			name: "no foot script",
			script: ScriptDef{
				Script:     Ptr("test"),
				FootScript: true,
			},
			want: "",
		},
		{
			name: "foot script",
			script: ScriptDef{
				Script: Ptr("test"),
			},
			want: `<script>
test
</script>`,
		},
		{
			name: "script attributes",
			script: ScriptDef{
				Attributes: map[string]string{
					"src":      "/some/path",
					"data-foo": "some data",
				},
				Script: Ptr("test"),
			},
			want: `<script data-foo="some data">
test
</script>`,
		},
		{
			name: "script src",
			script: ScriptDef{
				ScriptSrc: Ptr("/some/path"),
			},
			want: `<script src="/some/path"></script>`,
		},
		{
			name: "script src attributes",
			script: ScriptDef{
				Attributes: map[string]string{
					"src":      "/bad/path",
					"data-foo": "some data",
				},
				ScriptSrc: Ptr("/some/path"),
			},
			want: `<script data-foo="some data" src="/some/path"></script>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.script.ToHeadTag()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestScript_ToFootTag(t *testing.T) {
	tests := []struct {
		name   string
		script ScriptDef
		want   string
	}{
		{
			name: "basic",
			script: ScriptDef{
				Script:     Ptr("test"),
				FootScript: true,
			},
			want: `<script>
test
</script>`,
		},
		{
			name: "no foot script",
			script: ScriptDef{
				Script: Ptr("test"),
			},
			want: "",
		},
		{
			name: "foot script",
			script: ScriptDef{
				Script:     Ptr("test"),
				FootScript: true,
			},
			want: `<script>
test
</script>`,
		},
		{
			name: "script attributes",
			script: ScriptDef{
				Attributes: map[string]string{
					"src":      "/some/path",
					"data-foo": "some data",
				},
				FootScript: true,
				Script:     Ptr("test"),
			},
			want: `<script data-foo="some data">
test
</script>`,
		},
		{
			name: "script src",
			script: ScriptDef{
				ScriptSrc:  Ptr("/some/path"),
				FootScript: true,
			},
			want: `<script src="/some/path"></script>`,
		},
		{
			name: "script src attributes",
			script: ScriptDef{
				Attributes: map[string]string{
					"src":      "/bad/path",
					"data-foo": "some data",
				},
				FootScript: true,
				ScriptSrc:  Ptr("/some/path"),
			},
			want: `<script data-foo="some data" src="/some/path"></script>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.script.ToFootTag()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Ptr[T any](v T) *T {
	return &v
}
