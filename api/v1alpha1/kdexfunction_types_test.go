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
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_Regex(t *testing.T) {
	tests := []struct {
		name       string
		re         regexp.Regexp
		assertions func(t *testing.T, re regexp.Regexp)
	}{
		{
			name: "basepath / basic match",
			re:   (&API{}).BasePathRegex(),
			assertions: func(t *testing.T, re regexp.Regexp) {
				assert.True(t, re.MatchString("/foo/bar"))
			},
		},
		{
			name: "basepath / named capturing group",
			re:   (&API{}).BasePathRegex(),
			assertions: func(t *testing.T, re regexp.Regexp) {
				match := re.FindStringSubmatch("/foo/bar")
				result := make(map[string]string)

				if len(match) > 0 {
					for i, name := range re.SubexpNames() {
						if i != 0 && name != "" {
							result[name] = match[i]
						}
					}
				}

				assert.Equal(t, "/foo/bar", result["basePath"])
			},
		},
		{
			name: "basepath / not a match / a",
			re:   (&API{}).BasePathRegex(),
			assertions: func(t *testing.T, re regexp.Regexp) {
				assert.False(t, re.MatchString("foo/bar"))
			},
		},
		{
			name: "basepath / not a match / b",
			re:   (&API{}).BasePathRegex(),
			assertions: func(t *testing.T, re regexp.Regexp) {
				assert.False(t, re.MatchString("/foo"))
			},
		},
		{
			name: "itempath / basic match",
			re:   (&API{}).ItemPathRegex(),
			assertions: func(t *testing.T, re regexp.Regexp) {
				assert.True(t, re.MatchString("/foo/bar/fiz"))
			},
		},
		{
			name: "itempath / named capturing group",
			re:   (&API{}).ItemPathRegex(),
			assertions: func(t *testing.T, re regexp.Regexp) {
				match := re.FindStringSubmatch("/foo/bar/fiz")
				result := make(map[string]string)

				if len(match) > 0 {
					for i, name := range re.SubexpNames() {
						if i != 0 && name != "" {
							result[name] = match[i]
						}
					}
				}

				assert.Equal(t, "/foo/bar", result["basePath"])
			},
		},
		{
			name: "itempath / not a match / a",
			re:   (&API{}).ItemPathRegex(),
			assertions: func(t *testing.T, re regexp.Regexp) {
				assert.False(t, re.MatchString("foo/bar"))
			},
		},
		{
			name: "itempath / not a match / b",
			re:   (&API{}).ItemPathRegex(),
			assertions: func(t *testing.T, re regexp.Regexp) {
				assert.False(t, re.MatchString("/foo"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertions(t, tt.re)
		})
	}
}
