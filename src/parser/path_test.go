package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPaths(t *testing.T) {
	tests := []struct {
		name string
		input string
	}{
		{"simple", "/foo/bar/baz"},
		{"wildcard", "/foo/{bar}/{baz}/zoo"},
		{"recursive", "/foo/{bar}/{baz=**}"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := New(test.input)
			p, err := ParsePath(tokens)
			assert.Nil(t, err)
			assert.Equal(t, test.input, p.String())
		})
	}
}

