package predicate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfE(t *testing.T) {
	tests := []struct {
		name      string
		predicate bool
		trueVal   any
		elseVal   any
		want      any
	}{
		{
			name:      "string true",
			predicate: true,
			trueVal:   "yes",
			elseVal:   "no",
			want:      "yes",
		},
		{
			name:      "string false",
			predicate: false,
			trueVal:   "yes",
			elseVal:   "no",
			want:      "no",
		},
		{
			name:      "number true",
			predicate: true,
			trueVal:   10,
			elseVal:   20,
			want:      10,
		},
		{
			name:      "number false",
			predicate: false,
			trueVal:   10,
			elseVal:   20,
			want:      20,
		},
		{
			name:      "struct true",
			predicate: true,
			trueVal: struct {
				Foo string
			}{
				Foo: "foo",
			},
			elseVal: struct {
				Foo string
			}{
				Foo: "bar",
			},
			want: struct {
				Foo string
			}{
				Foo: "foo",
			},
		},
		{
			name:      "struct false",
			predicate: false,
			trueVal: struct {
				Foo string
			}{
				Foo: "foo",
			},
			elseVal: struct {
				Foo string
			}{
				Foo: "bar",
			},
			want: struct {
				Foo string
			}{
				Foo: "bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IfE(tt.predicate, tt.trueVal, tt.elseVal)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIfEF(t *testing.T) {
	tests := []struct {
		name      string
		predicate bool
		trueFunc  func() any
		elseFunc  func() any
		want      any
	}{
		{
			name:      "string true",
			predicate: true,
			trueFunc:  func() any { return "yes" },
			elseFunc:  func() any { return "no" },
			want:      "yes",
		},
		{
			name:      "string false",
			predicate: false,
			trueFunc:  func() any { return "yes" },
			elseFunc:  func() any { return "no" },
			want:      "no",
		},
		{
			name:      "number true",
			predicate: true,
			trueFunc:  func() any { return 10 },
			elseFunc:  func() any { return 20 },
			want:      10,
		},
		{
			name:      "number false",
			predicate: false,
			trueFunc:  func() any { return 10 },
			elseFunc:  func() any { return 20 },
			want:      20,
		},
		{
			name:      "struct true",
			predicate: true,
			trueFunc: func() any {
				return struct {
					Foo string
				}{
					Foo: "foo",
				}
			},
			elseFunc: func() any {
				return struct {
					Foo string
				}{
					Foo: "bar",
				}
			},
			want: struct {
				Foo string
			}{
				Foo: "foo",
			},
		},
		{
			name:      "struct false",
			predicate: false,
			trueFunc: func() any {
				return struct {
					Foo string
				}{
					Foo: "foo",
				}
			},
			elseFunc: func() any {
				return struct {
					Foo string
				}{
					Foo: "bar",
				}
			},
			want: struct {
				Foo string
			}{
				Foo: "bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IfEF(tt.predicate, tt.trueFunc, tt.elseFunc)
			assert.Equal(t, tt.want, got)
		})
	}
}
