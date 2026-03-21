package v1alpha1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestServiceAccountSecrets_Find(t *testing.T) {
	now := time.Now()
	s1 := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "s1",
			CreationTimestamp: metav1.NewTime(now.Add(-1 * time.Hour)),
		},
	}
	s2 := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "s2",
			CreationTimestamp: metav1.NewTime(now),
		},
	}

	secrets := ServiceAccountSecrets{s1, s2}

	// Should find s2 because it's newer and we sort descending
	found := secrets.Find(func(s corev1.Secret) bool {
		return true
	})

	assert.NotNil(t, found)
	assert.Equal(t, "s2", found.Name)

	// Verify it's not a pointer to a loop variable by checking the address
	// In the original buggy version, it might return the same address for different finds if they were in the same loop,
	// but here we just want to ensure it's pointing into the slice.
	found.Annotations = map[string]string{"modified": "true"}
	// Find again - if it was a copy, the original in the slice wouldn't be modified.
	// Note: Find sorts the slice in place! So s2 is at secrets[0] now.
	assert.Equal(t, "true", secrets[0].Annotations["modified"])
}

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
					Style: `color: #fff;`,
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
			assets: Assets{
				{
					LinkHref: "/some/path",
				},
			},
			want: `<link href="/some/path"/>`,
		},
		{
			name: "link href with attributes",
			assets: Assets{
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
			assets: Assets{
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
				Script: "test",
			},
			want: `<script>
test
</script>`,
		},
		{
			name: "no foot script",
			script: ScriptDef{
				Script:     "test",
				FootScript: true,
			},
			want: "",
		},
		{
			name: "foot script",
			script: ScriptDef{
				Script: "test",
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
				Script: "test",
			},
			want: `<script data-foo="some data">
test
</script>`,
		},
		{
			name: "script src",
			script: ScriptDef{
				ScriptSrc: "/some/path",
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
				ScriptSrc: "/some/path",
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
				Script:     "test",
				FootScript: true,
			},
			want: `<script>
test
</script>`,
		},
		{
			name: "no foot script",
			script: ScriptDef{
				Script: "test",
			},
			want: "",
		},
		{
			name: "foot script",
			script: ScriptDef{
				Script:     "test",
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
				Script:     "test",
			},
			want: `<script data-foo="some data">
test
</script>`,
		},
		{
			name: "script src",
			script: ScriptDef{
				ScriptSrc:  "/some/path",
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
				ScriptSrc:  "/some/path",
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
