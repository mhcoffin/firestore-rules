package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicValues(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result string
	}{
		{"identifier", "foobar", "foobar"},
		{"int", "34", "34"},
		{"float", "3.14159", "3.14159"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := New(test.input)
			expr, err := ParseExpr(tokens)
			assert.Nil(t, err)
			assert.Equal(t, expr.String(), test.result)
		})
	}
}

func TestBinaryOperations(t *testing.T) {
	tests := []struct {
		name string
		input string
		result string
	}{
		{"+", "a + b", "(a + b)"},
		{"-", "a - 3", "(a - 3)"},
		{"*", "a * b", "(a * b)"},
		{"/", "a/b", "(a / b)"},
		{"%", "a%2", "(a % 2)"},
		{"plus-left-associative", "a + b + c", "((a + b) + c)"},
		{"times-left-associative", "a * b * c", "((a * b) * c)"},
		{"<", "a < b", "(a < b)"},
		{"==", "a == b", "(a == b)"},
		{"||", "a || b || c", "((a || b) || c)"},
		{"&&", "a && b && c", "((a && b) && c)"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := New(test.input)
			expr, err := ParseExpr(tokens)
			assert.Nil(t, err)
			assert.Equal(t, test.result, expr.String())
		})
	}
}

func TestBinaryPrecedence(t *testing.T) {
	tests := []struct {
		name string
		expected string
	}{
		{"a + b * c", "(a + (b * c))"},
		{"a * b - c * d", "((a * b) - (c * d))"},
		{"a == b || c == d", "((a == b) || (c == d))"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := ParseExpr(New(test.name))
			assert.Nil(t, err)
			assert.Equal(t, test.expected, expr.String())
		})
	}
}

func TestUnary(t *testing.T) {
	tests := []struct {
		name string
		expected string
	}{
		{"!a", "(!a)"},
		{"-4", "(-4)"},
		{"--4", "(-(-4))"},
		{"!!true", "(!(!true))"},
		{"!!a==b", "((!(!a)) == b)"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := ParseExpr(New(test.name))
			assert.Nil(t, err)
			assert.Equal(t, test.expected, expr.String())
		})
	}
}

func TestTernaryExpr(t *testing.T) {
	tests := []struct {
		name string
		expected string
	}{
		{"a ? b : c", "(a ? b : c)"},
		{"a == 3 ? a+b : c+d", "((a == 3) ? (a + b) : (c + d))"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := ParseExpr(New(test.name))
			assert.Nil(t, err)
			assert.Equal(t, test.expected, expr.String())
		})
	}
}





