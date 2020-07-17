package parser

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func nows(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func TestParseMatchStmt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trivial",
			input:    `match /foo/bar { }`,
			expected: "match /foo/bar {}",
		},
		{
			name:     "func",
			input:    "match /foo/bar {function foo() {return 2+3;}}",
			expected: "match /foo/bar {function foo () {return (2+3);}}",
		},
		{
			name:     "two functions",
			input:    "match /foo/{bar} {function foo() {return 2+3;} function bar(a) {return a+1;}}",
			expected: "match /foo/{bar} {function foo() {return (2+3); } function bar(a) {return (a+1);}}",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := New(test.input)
			ms, err := ParseMatchStmt(tokens)
			assert.Nil(t, err)
			assert.Equal(t, nows(test.expected), nows(ms.String()))
		})
	}
}
