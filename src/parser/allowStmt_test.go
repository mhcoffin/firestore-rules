package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAllowStmt(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple", "allow read, write: if false;"},
		{"simple2", "allow get, list: if (a == b);"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := New(test.input)
			as, err := ParseAllowStmt(tokens)
			assert.Nil(t, err)
			assert.Equal(t, test.input, as.String())
		})
	}
}

func TestAllowStmtErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{
			name:  "bad action",
			input: "allow foo: if false;",
			err:   "line 1 col 7: unexpected token in accept list (foo)",
		},
		{
			name:  "missing comma",
			input: "allow read write: if true;",
			err:   "line 1 col 12: unexpected token in accept list: (write)",
		},
		{
			name: "bad expression",
			input: "allow read: if a b c;",
			err: "line 1 col 18: unexpected token (b)",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := New(test.input)
			_, err := ParseAllowStmt(tokens)
			assert.NotNil(t, err)
			assert.Equal(t, test.err, err.Error())
		})
	}
}
