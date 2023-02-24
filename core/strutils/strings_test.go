package strutils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spy16/forge/core/strutils"
)

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		str  string
		want string
	}{
		{"ABCD", "abcd"},
		{"AbcD", "abc_d"},
		{"helloWorld", "hello_world"},
		{"UseHTTPS", "use_https"},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			assert.Equalf(t, tt.want, strutils.SnakeCase(tt.str), "SnakeCase(\"%s\")", tt.str)
		})
	}
}

func TestOneOf(t *testing.T) {
	tests := []struct {
		arr  []string
		find string
		want bool
	}{
		{
			arr:  []string{},
			find: "hello",
			want: false,
		},
		{
			arr:  []string{"foo", "bar", "baz"},
			find: "Bar",
			want: false,
		},
		{
			arr:  []string{"foo", "bar", "baz"},
			find: "bar",
			want: true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Case#%d", i), func(t *testing.T) {
			assert.Equalf(t, tt.want, strutils.OneOf(tt.find, tt.arr), "OneOf(%v, %v)", tt.find, tt.arr)
		})
	}
}
