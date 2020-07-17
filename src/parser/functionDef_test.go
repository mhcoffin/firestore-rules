package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunctionDef(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trivial",
			input:    "function f() { return 0; }",
			expected: "function f () {  return 0; }",
		},
		{
			name:     "args",
			input:    "function f(a, b) { return a + b; }",
			expected: "function f (a, b) {  return (a + b); }",
		},
		{
			name:     "let",
			input:    "function f() { let foo = 2 + 3; return foo - 7; }",
			expected: "function f () { let foo = (2 + 3); return (foo - 7); }",
		},
		{
			name:     "let2",
			input:    "function f() { let a = 3; let b = foobar(5); return a+b;}",
			expected: "function f () { let a = 3;let b = foobar(5); return (a + b); }",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := New(test.input)
			def, err := ParseFunctionDef(tokens)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, def.String())
		})
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		error string
	}{
		{
			name:  "unclosed parens",
			input: "function foo(a, b {return 0}",
			error: "line 1 col 19: unexpected token in parameter list ( \"{\" )",
		},
		{
			name:  "missing semicolon",
			input: "function foo() { let b = 1 let c=4; return 7;}",
			error: "line 1 col 28: unexpected token (let)",
		},
		{
			name: "missing closing brace",
			input: `function foo(a, b) { let x = a; let y = b; return LatLon(x, y); `,
			error: "line 1 col 65: unexpected token (EOF)",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := New(test.input)
			_, err := ParseFunctionDef(tokens)
			assert.NotNil(t, err)
			assert.Equal(t, test.error, err.Error())
		})
	}
}
