package parser

import (
	"fmt"
)

type Tokens struct {
	ch  chan Token
	buf []Token
	pos int
}

func New(input string) *Tokens {
	return &Tokens{
		ch:  Start(input),
		buf: make([]Token, 0),
		pos: 0,
	}
}

func (tokens *Tokens) Peek() Token {
	if tokens.pos >= len(tokens.buf) {
		tokens.buf = append(tokens.buf, <-tokens.ch)
	}
	return tokens.buf[tokens.pos]
}

func (tokens *Tokens) AcceptAny() Token {
	if tokens.pos >= len(tokens.buf) {
		tokens.buf = append(tokens.buf, <-tokens.ch)
	}
	tokens.pos++
	return tokens.buf[tokens.pos-1]
}

func (tokens *Tokens) Accept(kind Kind) (Token, error) {
	t := tokens.AcceptAny()
	switch t.Kind {
	case kind:
		return t, nil
	case Error:
		return Token{Kind: Error}, t.Error
	default:
		return Token{Kind: Error}, &LexError{
			Start: t.Start,
			Pos:   t.End,
			msg:   fmt.Sprintf("unexpected token (%s)", t.ErrString()),
		}
	}
}
