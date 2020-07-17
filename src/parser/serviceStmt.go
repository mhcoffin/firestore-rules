package parser

import (
	"fmt"
	"strings"
)

type ServiceStmt struct {
	Name Expr
	Statements []Stmt
}

func (ss *ServiceStmt) String() string {
	stmts := make([]string, len(ss.Statements))
	for k, v := range ss.Statements {
		stmts[k] = v.String()
	}
	return fmt.Sprintf("service %s { %s }", ss.Name, strings.Join(stmts, ""))
}

func ParseServiceStmt(tokens *Tokens) (*ServiceStmt, error) {
	_, err := tokens.Accept(Service)
	if err != nil {
		return nil, err
	}
	name, err := ParseServiceName(tokens)
	if err != nil {
		return nil, err
	}
	_, err = tokens.Accept(LeftBrace)
	if err != nil {
		return nil, err
	}

	stmts := make([]Stmt, 0)
	for {
		switch tokens.Peek().Kind {
		case RightBrace:
			tokens.AcceptAny()
			return &ServiceStmt{name, stmts}, nil
		case Function:
			f, err := ParseFunctionDef(tokens)
			if err != nil {
				return nil, err
			}
			stmts = append(stmts, f)
		case Match:
			m, err := ParseMatchStmt(tokens)
			if err != nil {
				return nil, err
			}
			stmts = append(stmts, m)
			fmt.Printf("finished match\n")
		default:
			return nil, LexError{
				Start: tokens.Peek().Start,
				Pos:   tokens.Peek().End,
				msg:   fmt.Sprintf("unexpected token (%s)", tokens.Peek()),
			}
		}
	}
}

type ServiceName struct {
	name1 Token
	name2 Token
}

func (sn *ServiceName) String() string {
	return fmt.Sprintf("%s.%s", sn.name1, sn.name2)
}


func ParseServiceName(tokens *Tokens) (*ServiceName, error) {
	name1, err := tokens.Accept(Identifier)
	if err != nil {
		return nil, err
	}
	_, err = tokens.Accept(Dot)
	if err != nil {
		return nil, err
	}
	name2, err := tokens.Accept(Identifier)
	if err != nil {
		return nil, err
	}
	return &ServiceName{name1, name2}, nil
}

