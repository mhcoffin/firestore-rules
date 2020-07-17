package parser

import (
	"fmt"
	"strings"
)

type FunctionDef struct {
	Stmt
	functionName  Token
	paramNames    []Param
	letStatements []LetDef
	returnStmt    Expr
}

func (fd *FunctionDef) String() string {
	paramNames := make([]string, len(fd.paramNames))
	for k, v := range fd.paramNames {
		paramNames[k] = v.name.Value
	}

	letStmts := make([]string, len(fd.letStatements))
	for k, v := range fd.letStatements {
		letStmts[k] = v.String()
	}

	retStmt := fmt.Sprintf("return %s;", fd.returnStmt.String())

	return fmt.Sprintf("function %s (%s) { %s %s }",
		fd.functionName.Value,
		strings.Join(paramNames, ", "),
		strings.Join(letStmts, ""),
		retStmt,
	)
}

type LetDef struct {
	name  Token
	value Expr
}

func (ld *LetDef) String() string {
	return fmt.Sprintf("let %s = %s;", ld.name.Value, ld.value.String())
}

type Param struct {
	name Token
}

func ParseFunctionDef(tokens *Tokens) (*FunctionDef, error) {
	_, err := tokens.Accept(Function)
	if err != nil {
		return nil, err
	}
	name, err := tokens.Accept(Identifier)
	if err != nil {
		return nil, err
	}
	params, err := ParseParamList(tokens)
	if err != nil {
		return nil, err
	}
	_, err = tokens.Accept(LeftBrace)
	if err != nil {
		return nil, err
	}
	letStmts, err := ParseLetStmtList(tokens)
	if err != nil {
		return nil, err
	}
	retStmt, err := ParseReturnStmt(tokens)
	if err != nil {
		return nil, err
	}
	_, err = tokens.Accept(RightBrace)
	if err != nil {
		return nil, err
	}
	return &FunctionDef{
		functionName:  name,
		paramNames:    params,
		letStatements: letStmts,
		returnStmt:    retStmt,
	}, nil
}

func ParseReturnStmt(tokens *Tokens) (Expr, error) {
	_, err := tokens.Accept(Return)
	if err != nil {
		return nil, err
	}
	result, err := ParseExpr(tokens)
	if err != nil {
		return nil, err
	}
	_, err = tokens.Accept(SemiColon)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ParseLetStmtList(tokens *Tokens) ([]LetDef, error) {
	result := make([]LetDef, 0)
	for {
		switch tokens.Peek().Kind {
		case Let:
			letDef, err := ParseLetStmt(tokens)
			if err != nil {
				return nil, err
			}
			result = append(result, *letDef)
		default:
			return result, nil
		}
	}
}

func ParseLetStmt(tokens *Tokens) (*LetDef, error) {
	_, err := tokens.Accept(Let)
	if err != nil {
		return nil, err
	}
	id, err := tokens.Accept(Identifier)
	if err != nil {
		return nil, err
	}
	_, err = tokens.Accept(Eq)
	if err != nil {
		return nil, err
	}
	rhs, err := ParseExpr(tokens)
	if err != nil {
		return nil, err
	}
	_, err = tokens.Accept(SemiColon)
	if err != nil {
		return nil, err
	}
	return  &LetDef{id, rhs}, nil
}

func ParseParamList(tokens *Tokens) ([]Param, error) {
	result := make([]Param, 0)

	_, err := tokens.Accept(LeftParen)
	if err != nil {
		return nil, err
	}
	if tokens.Peek().Kind == RightParen {
		_, err := tokens.Accept(RightParen)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	result = append(result, Param{tokens.AcceptAny()})
	for {
		switch tokens.Peek().Kind {
		case Comma:
			_, err := tokens.Accept(Comma)
			if err != nil {
				return nil, err
			}
			id, err := tokens.Accept(Identifier)
			if err != nil {
				return nil, err
			}
			result = append(result, Param{id})
		case RightParen:
			tokens.AcceptAny()
			return result, nil
		default:
			return nil, fmt.Errorf(
				`%s: unexpected token in parameter list ( "%s" )`,
				tokens.Peek().Start,
				tokens.Peek())
		}
	}
}
