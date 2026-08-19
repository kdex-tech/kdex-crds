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
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
)

// TestOptionalFieldsAreOmittedWhenUnset asserts that an optional field really is
// absent from the marshalled object when the author never set it. encoding/json
// only honours omitempty for basic types, maps, slices, pointers and interfaces,
// so a value struct declared with omitempty is marshalled unconditionally — the
// stored CR then gains a key it never declared and CEL's has() is always true.
func TestOptionalFieldsAreOmittedWhenUnset(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		absentKeys []string
	}{
		{"Auth.jwt", Auth{}, []string{"jwt"}},
		{"KDexPageSpec.contact", KDexPageSpec{}, []string{"contact"}},
		{"NavigationHints.weight", NavigationHints{}, []string{"weight"}},
		{"Runtime.resources", Runtime{}, []string{"resources"}},
		{"KDexFunctionSpec.metadata+origin", KDexFunctionSpec{}, []string{"metadata", "origin"}},
		{"KDexHostSpec.registries", KDexHostSpec{}, []string{"registries"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := json.Marshal(tt.value)
			require.NoError(t, err)

			var got map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(raw, &got))

			for _, key := range tt.absentKeys {
				assert.NotContains(t, got, key,
					"%q must be absent when unset, got %s", key, raw)
			}
		})
	}
}

// TestNoOptionalValueStructFields is the regression guard for the class the test
// above samples: it walks every type registered in the scheme and fails on any
// field of ours that declares omitempty over a type encoding/json can never
// treat as empty. Add a pointer, not an allowlist entry, when this fires.
func TestNoOptionalValueStructFields(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, AddToScheme(scheme))

	pkgPath := reflect.TypeOf(KDexHostSpec{}).PkgPath()

	// Kubernetes API convention: every object carries ObjectMeta, every list
	// carries ListMeta, and every object's status is a value struct served
	// through the status subresource. These are upstream's shape, not ours.
	allowed := map[string]bool{
		"v1.ObjectMeta":               true,
		"v1.ListMeta":                 true,
		"v1alpha1.KDexObjectStatus":   true,
		"v1alpha1.KDexFunctionStatus": true,
	}

	seen := map[reflect.Type]bool{}
	var findings []string

	var walk func(reflect.Type)
	walk = func(t reflect.Type) {
		for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice ||
			t.Kind() == reflect.Array || t.Kind() == reflect.Map {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct || seen[t] || t.PkgPath() != pkgPath {
			return
		}
		seen[t] = true

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			if hasOmitempty(f.Tag.Get("json")) && neverEmpty(f.Type) &&
				!allowed[shortType(f.Type)] {
				findings = append(findings, fmt.Sprintf(
					"%s.%s is %s with omitempty — omitempty never fires for it",
					t.Name(), f.Name, shortType(f.Type)))
			}
			walk(f.Type)
		}
	}

	for _, t := range scheme.AllKnownTypes() {
		walk(t)
	}

	sort.Strings(findings)
	assert.Empty(t, findings, "value-struct fields carrying omitempty:\n%s",
		strings.Join(findings, "\n"))
}

func hasOmitempty(tag string) bool {
	parts := strings.Split(tag, ",")
	for _, opt := range parts[1:] {
		if opt == "omitempty" {
			return true
		}
	}
	return false
}

// neverEmpty reports whether encoding/json's isEmptyValue can never return true
// for this type. It covers bool, the numeric kinds, string, and anything with a
// length, plus pointers and interfaces — so struct and array kinds are stuck.
func neverEmpty(t reflect.Type) bool {
	return t.Kind() == reflect.Struct || t.Kind() == reflect.Array
}

func shortType(t reflect.Type) string {
	return strings.TrimPrefix(t.String(), "*")
}
