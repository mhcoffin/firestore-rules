package parser

import (
	"fmt"
	"strings"
)

type AllowStmt struct {
	Stmt
	actions   []Token
	condition Expr
}

func (as *AllowStmt) String() string {
	a := make([]string, len(as.actions))
	for k, v := range as.actions {
		a[k] = v.String()
	}
	return fmt.Sprintf("allow %s: if %s;", strings.Join(a, ", "), as.condition)
}

func ParseAllowStmt(tokens *Tokens) (*AllowStmt, error) {
	_, err := tokens.Accept(Allow)
	if err != nil {
		return nil, err
	}
	actions, err := ParseActionNameList(tokens)
	if err != nil {
		return nil, err
	}
	_, err = tokens.Accept(Colon)
	if err != nil {
		return nil, err
	}
	_, err = tokens.Accept(If)
	if err != nil {
		return nil, err
	}
	expr, err := ParseExpr(tokens)
	if err != nil {
		return nil, err
	}
	_, err = tokens.Accept(SemiColon)
	if err != nil {
		return nil, err
	}
	return &AllowStmt{actions: actions, condition: expr}, nil
}

func ParseActionNameList(tokens *Tokens) ([]Token, error) {
	actions := make([]Token, 0)
	first, err := ParseActionName(tokens)
	if err != nil {
		return nil, err
	}
	actions = append(actions, first)
	for {
		switch tokens.Peek().Kind {
		case Comma:
			tokens.Accept(Comma)
			action, err := ParseActionName(tokens)
			if err != nil {
				return nil, err
			}
			actions = append(actions, action)
		case Colon:
			return actions, nil
		default:
			badToken := tokens.AcceptAny()
			return nil, LexError{
				Start: badToken.Start,
				Pos:   badToken.End,
				msg:   fmt.Sprintf("unexpected token in accept list: (%s)", badToken),
			}
		}
	}
}

func ParseActionName(tokens *Tokens) (Token, error) {
	result := tokens.AcceptAny()
	if IsAction(result.Kind) {
		return result, nil
	} else {
		return Token{Kind: Error}, &LexError{
			Start: result.Start,
			Pos:   result.End,
			msg:   fmt.Sprintf("unexpected token in accept list (%s)", result),
		}
	}
}
