package parser

import (
	"strings"
	"testing"
)

func TestEmpty(t *testing.T) {
	ch := setup("")

	token := <-ch
	if token.Kind != Eof {
		t.Fail()
	}
}

func TestDot(t *testing.T) {
	ch := setup(".")
	token := <-ch
	if token.Kind != Dot {
		t.Fail()
	}
	token = <-ch
	if token.Kind != Eof {
		t.Fail()
	}
}

func TestOperators(t *testing.T) {
	tests := []struct {
		name  string
		input string
		kind  Kind
	}{
		{"dot", ".", Dot},
		{"slash", "/", Slash},
		{"minus", "-", Minus},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ch := setup(test.input)
			token := <-ch
			if token.Kind != test.kind {
				t.Errorf("expected %v but got %v", test.kind, token.Kind)
			}
			if token.Value != test.input {
				t.Errorf("did not consume Input: %v", token.Value)
			}

		})
	}
}

func TestWord(t *testing.T) {
	ch := setup("hello")
	token := <-ch
	if token.Kind != Identifier || token.Value != "hello" {
		t.Fail()
	}
}

func TestWordDotWord(t *testing.T) {
	ch := setup("foo42.ba3r")
	token := <-ch
	if token.Kind != Identifier || token.Value != "foo42" {
		t.Error("expected word1")
	}
	token = <-ch
	if token.Kind != Dot || token.Value != "." {
		t.Error("expected dot")
	}
	token = <-ch
	if token.Kind != Identifier || token.Value != "ba3r" {
		t.Error("expected ba3r but got", token.Value)
	}
	token = <-ch
	if token.Kind != Eof || token.Value != "" {
		t.Error("expected eof")
	}
}

func TestStringLiterals(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"basic", `"foo"`},
		{"escape", `"foo\"bar"`},
		{"escape2", `" \\ foo \\ "`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ch := setup(test.input)
			token := <-ch
			if token.Kind != StringLiteral || token.Value != test.input {
				t.Errorf("expected %v but got %v", test.input, token.Value)
			}
		})
	}
}

func TestSingleTokens(t *testing.T) {
	tests := []struct {
		name  string
		kind  Kind
		input string
	}{
		{"word", Identifier, "foobar"},
		{"dot", Dot, "."},
		{"int", IntLiteral, "123"},
		{"float", FloatLiteral, "123.456"},
		{"{", LeftBrace, "{"},
		{"}", RightBrace, "}"},
		{"(", LeftParen, "("},
		{")", RightParen, ")"},
		{"[", LeftSquareBracket, "["},
		{"]", RightSquareBracket, "]"},
		{"String literal", StringLiteral, `'here is a string'`},
		{"Slash", Slash, "/"},
		{"Minus", Minus, "-"},
		{"Plus", Plus, "+"},
		{"Comma", Comma, ","},
		{"=", Eq, `=`},
		{"==", EqEq, "=="},
		{";", SemiColon, ";"},
		{"<", Less, "<"},
		{"<=", LessEq, "<="},
		{">", Greater, ">"},
		{">=", GreaterEq, ">="},
		{"!", Bang, "!"},
		{"!=", NotEq, "!="},
		{":", Colon, ":"},
		{"&", And, "&"},
		{"&&", AndAnd, "&&"},
		{"|", Or, "|"},
		{"||", OrOr, "||"},
		{"*", Star, "*"},
		{"**", StarStar, "**"},
		{"bytes", Bytes, "b'abc'"},
		{"percent", Percent, "%"},
		{"if", If, "if"},
		{"service", Service, "service"},
		{"match", Match, "match"},
		{"allow", Allow, "allow"},
		{"create", Create, "create"},
		{"delete", Delete, "delete"},
		{"write", Write, "write"},
		{"get", Get, "get"},
		{"list", List, "list"},
		{"read", Read, "read"},
		{"if", If, "if"},
		{"function", Function, "function"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ch := setup(test.input)
			token := <-ch
			if token.Kind != test.kind || token.Value != test.input {
				t.Fatalf("expected %s(%s) but got %s", test.kind, test.input, token)
			}
			token = <-ch
			if token.Kind != Eof {
				t.Fatalf("expected EOF")
			}
		})
	}
}

const multilineString = `
Here is a multi line string
with at least
  three characters
  and some lines indented
`

func TestPositions(t *testing.T) {
	tests := []struct {
		name string
		input string
	}{
		{"one line", "apple.pear(banana, peach, plum) == tomato"},
		{"multiline", multilineString},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := New(test.input)
			for {
				token := tokens.AcceptAny()
				if token.Kind == Eof {
					break
				}
				if token.Kind != Identifier {
					continue
				}
				start := strings.Index(test.input, token.Value)
				if start != token.Start.Pos {
					t.Errorf("token %s", token )
				}
				newLines := strings.Count(test.input[0:start], "\n")
				if newLines != token.Start.Line {
					t.Errorf("token %s: expected line %d but got %d", token, newLines, token.Start.Line)
				}
				col := start - strings.LastIndex(test.input[0:start], "\n") - 1
				if col != token.Start.Col {
					t.Errorf("token %s: expected col %d but got %d", token, col, token.Start.Col)
				}
			}

		})
	}
}

func TestLexicalErrors(t *testing.T) {
	tests := []struct {
		name string
		input string
	}{
		{
			name: "unclosed string literal",
			input: `foo := "bar || fred"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := New(test.input)
			for {
				tok := tokens.AcceptAny()
				if tok.Kind == Eof {
					break
				}
			}
		})
	}
}

func setup(s string) chan Token {
	ch := make(chan Token)
	lexer := Lexer{
		Input: []rune(s),
		Chan:  ch,
	}
	go Run(&lexer)
	return ch
}
