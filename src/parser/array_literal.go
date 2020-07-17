package parser

import (
	"fmt"
	"strings"
)

type ArrayLiteral struct {
	elements []Expr
}

func (al *ArrayLiteral) String() string {
	result := make([]string, len(al.elements))
	for k, v := range al.elements {
		result[k] = v.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(result, ", "))
}

func ParseArrayLiteral(tokens *Tokens) (*ArrayLiteral, error) {
	elts := make([]Expr, 0)
	_, err := tokens.Accept(LeftSquareBracket)
	if err != nil {
		return nil, err
	}
	_, err = ParseExpr(tokens)
	if err != nil {
		return nil, err
	}
	for {
		switch tokens.Peek().Kind {
		case Comma:
			_, _ = tokens.Accept(Comma)
			exp, err := ParseExpr(tokens)
			if err != nil {
				return nil, err
			}
			elts = append(elts, exp)
		case RightSquareBracket:
			_, _ = tokens.Accept(RightSquareBracket)
			return &ArrayLiteral{elts}, nil
		default:
			t := tokens.Peek()
			return nil, LexError{
				Start: t.Start,
				Pos:   t.End,
				msg:   fmt.Sprintf("unexpected token in array literal (%s)", t),
			}
		}

	}
}
