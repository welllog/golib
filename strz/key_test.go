package strz

import (
	"testing"

	"github.com/welllog/golib/testz"
)

func TestNewKeyGenerator(t *testing.T) {
	tests := []struct {
		delimiter string
		prefix    []string
		parts     []string
		expect    string
	}{
		{
			":",
			[]string{"a", "b"},
			[]string{"1", "2"},
			"a:b:1:2",
		},
		{
			":",
			[]string{},
			[]string{"1", "2"},
			"1:2",
		},
		{
			"@",
			[]string{"a"},
			[]string{"1"},
			"a@1",
		},
		{
			"@",
			[]string{"a"},
			[]string{},
			"a@",
		},
		{
			"@",
			[]string{},
			[]string{},
			"",
		},
	}
	for _, tt := range tests {
		testz.Equal(t, tt.expect, NewKeyGenerator(tt.delimiter, tt.prefix...).Generate(tt.parts...))
	}
}

func TestKeyGeneratorWith(t *testing.T) {
	tests := []struct {
		delimiter string
		prefix    []string
		with      []string
		parts     []string
		expect    string
	}{
		{
			":",
			[]string{"a", "b"},
			[]string{},
			[]string{"1", "2"},
			"a:b:1:2",
		},
		{
			":",
			[]string{"a"},
			[]string{},
			[]string{"1", "2"},
			"a:1:2",
		},
		{
			":",
			[]string{},
			[]string{"c"},
			[]string{"1", "2"},
			"c:1:2",
		},
		{
			":",
			[]string{},
			[]string{"c"},
			[]string{},
			"c:",
		},
		{
			"@",
			[]string{"a", "b"},
			[]string{},
			[]string{"1"},
			"a@b@1",
		},
		{
			"@",
			[]string{},
			[]string{},
			[]string{},
			"",
		},
	}
	for _, tt := range tests {
		testz.Equal(t, tt.expect, NewKeyGenerator(tt.delimiter, tt.prefix...).With(tt.with...).Generate(tt.parts...))
	}
}
