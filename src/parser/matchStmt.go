package parser

import (
	"fmt"
	"strings"
)

type MatchStmt struct {
	Stmt
	Path       Path
	Components []Stmt
}

func (ms *MatchStmt) String() string {
	comp := make([]string, len(ms.Components))
	for k, v := range ms.Components {
		comp[k] = v.String()
	}
	return fmt.Sprintf("match %s {%s}", ms.Path, strings.Join(comp, "  "))
}

func ParseMatchStmt(tokens *Tokens) (*MatchStmt, error) {
	_, err := tokens.Accept(Match)
	if err != nil {
		return nil, err
	}
	c := make([]Stmt, 0)
	path, err := ParsePath(tokens)
	if err != nil {
		return nil, err
	}
	_, err = tokens.Accept(LeftBrace)
	if err != nil {
		return nil, err
	}

	for {
		switch tokens.Peek().Kind {
		case Function:
			fn, err := ParseFunctionDef(tokens)
			if err != nil {
				return nil, err
			}
			c = append(c, fn)
		case Match:
			m, err := ParseMatchStmt(tokens)
			if err != nil {
				return nil, err
			}
			c = append(c, m)
		case Allow:
			a, err := ParseAllowStmt(tokens)
			if err != nil {
				return nil, err
			}
			c = append(c, a)
		case RightBrace:
			tokens.AcceptAny()
			return &MatchStmt{Path: path, Components: c}, nil
		default:
			return nil, fmt.Errorf("unexpected token: %s", tokens.Peek())
		}
	}
}
