package parser

import (
	"fmt"
	"strings"
)

type Expr interface {
	String() string
}

type TernaryExpr struct {
	cond Expr
	t    Expr
	f    Expr
}

func (ter *TernaryExpr) String() string {
	return fmt.Sprintf("(%s ? %s : %s)", ter.cond, ter.t, ter.f)
}

type BinaryExpr struct {
	op  Token
	lhs Expr
	rhs Expr
}

func (be *BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", be.lhs.String(), be.op.String(), be.rhs.String())
}

type UnaryExpr struct {
	op      Token
	operand Expr
}

func (ue *UnaryExpr) String() string {
	return fmt.Sprintf("%s%s", ue.op, ue.operand)
}

type Id struct {
	id Token
}

func (id *Id) String() string {
	return fmt.Sprintf("%s", id.id.Value)
}

type Literal struct {
	value Token
}

func (lit *Literal) String() string {
	return lit.value.String()
}

type FunctionCall struct {
	fn   Expr
	args []Expr
}

func (fc *FunctionCall) String() string {
	argList := make([]string, len(fc.args))
	for k, arg := range fc.args {
		argList[k] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", fc.fn, strings.Join(argList, ", "))
}

func ParseExpr(tokens *Tokens) (Expr, error) {
	return ParseTernary(tokens)
}

func ParseTernary(tokens *Tokens) (Expr, error) {
	x, err := ParseOrExprList(tokens)
	if err != nil {
		return nil, err
	}
	if tokens.Peek().Kind == QuestionMark {
		tokens.AcceptAny()
		y, err := ParseExpr(tokens)
		if err != nil {
			return nil, err
		}
		_, err = tokens.Accept(Colon)
		if err != nil {
			return nil, err
		}
		z, err := ParseExpr(tokens)
		if err != nil {
			return nil, err
		}
		return &TernaryExpr{x, y, z}, nil
	} else {
		return x, nil
	}
}

func ParseOrExprList(tokens *Tokens) (Expr, error) {
	result, err := ParseAndExprList(tokens)
	if err != nil {
		return nil, err
	}
	for {
		switch tokens.Peek().Kind {
		case OrOr:
			op := tokens.AcceptAny()
			rhs, err := ParseAndExprList(tokens)
			if err != nil {
				return nil, err
			}
			result = &BinaryExpr{op, result, rhs}
		default:
			return result, nil
		}
	}
}

func ParseAndExprList(tokens *Tokens) (Expr, error) {
	result, err := ParseEqExprList(tokens)
	if err != nil {
		return nil, err
	}
	for {
		switch tokens.Peek().Kind {
		case AndAnd:
			op := tokens.AcceptAny()
			rhs, err := ParseEqExprList(tokens)
			if err != nil {
				return nil, err
			}
			result = &BinaryExpr{op, result, rhs}
		default:
			return result, nil
		}
	}
}

func ParseEqExprList(tokens *Tokens) (Expr, error) {
	result, err := ParseInExprList(tokens)
	if err != nil {
		return nil, err
	}
	for {
		switch tokens.Peek().Kind {
		case EqEq, NotEq:
			op := tokens.AcceptAny()
			rhs, err := ParseInExprList(tokens)
			if err != nil {
				return nil, err
			}
			result = &BinaryExpr{op, result, rhs}
		default:
			return result, nil
		}
	}
}

func ParseInExprList(tokens *Tokens) (Expr, error) {
	result, err := ParseRelationalExprList(tokens)
	if err != nil {
		return nil, err
	}
	for {
		switch tokens.Peek().Kind {
		case In, Is:
			op := tokens.AcceptAny()
			rhs, err := ParseRelationalExprList(tokens)
			if err != nil {
				return nil, err
			}
			result = &BinaryExpr{op, result, rhs}
		default:
			return result, nil
		}
	}
}

func ParseRelationalExprList(tokens *Tokens) (Expr, error) {
	result, err := ParseAdditiveExprList(tokens)
	if err != nil {
		return nil, err
	}
	for {
		switch tokens.Peek().Kind {
		case Less, LessEq, Greater, GreaterEq:
			op := tokens.AcceptAny()
			rhs, err := ParseAdditiveExprList(tokens)
			if err != nil {
				return nil, err
			}
			result = &BinaryExpr{op, result, rhs}
		default:
			return result, nil
		}
	}
}

func ParseAdditiveExprList(tokens *Tokens) (Expr, error) {
	result, err:= ParseMultiplicativeExprList(tokens)
	if err != nil {
		return nil, err
	}
	for {
		switch tokens.Peek().Kind {
		case Plus, Minus:
			op := tokens.AcceptAny()
			rhs, err := ParseMultiplicativeExprList(tokens)
			if err != nil {
				return nil, err
			}
			result = &BinaryExpr{op, result, rhs}
		default:
			return result, nil
		}
	}
}

func ParseMultiplicativeExprList(tokens *Tokens) (Expr, error) {
	result, err := ParseUnaryExpr(tokens)
	if err != nil {
		return nil, err
	}
	for {
		switch tokens.Peek().Kind {
		case Star, Slash, Percent:
			op := tokens.AcceptAny()
			rhs, err := ParseUnaryExpr(tokens)
			if err != nil {
				return nil, err
			}
			result = &BinaryExpr{op, result, rhs}
		default:
			return result, nil
		}
	}
}

func ParseUnaryExpr(tokens *Tokens) (Expr, error) {
	for {
		switch tokens.Peek().Kind {
		case Minus, Bang:
			op := tokens.AcceptAny()
			operand, err := ParseUnaryExpr(tokens)
			if err != nil {
				return nil, err
			}
			return &UnaryExpr{op, operand}, nil
		default:
			return ParseTerm(tokens)
		}
	}
}

func ParseTerm(tokens *Tokens) (Expr, error) {
	result, err := ParseBasicTerm(tokens)
	if err != nil {
		return nil, err
	}

	for {
		switch tokens.Peek().Kind {
		case Dot:
			dot := tokens.AcceptAny()
			rhs, err := tokens.Accept(Identifier)
			if err != nil {
				return nil, err
			}
			result = &BinaryExpr{dot, result, rhs}
		case LeftSquareBracket:
			brace := tokens.AcceptAny()
			index, err := ParseExpr(tokens)
			if err != nil {
				return nil, err
			}
			_, _ = tokens.Accept(RightSquareBracket)
			result = &BinaryExpr{brace, result, index}
		case LeftParen:
			tokens.AcceptAny()
			argList, err := ParseExprList(tokens)
			if err != nil {
				return nil, err
			}
			_, err = tokens.Accept(RightParen)
			if err != nil {
				return nil, err
			}
			result = &FunctionCall{result, argList}
		default:
			return result, nil
		}
	}
}

func ParseBasicTerm(tokens *Tokens) (Expr, error) {
	switch tokens.Peek().Kind {
	case Identifier:
		return &Id{tokens.AcceptAny()}, nil
	case StringLiteral, IntLiteral, FloatLiteral, Bytes, True, False:
		return &Literal{tokens.AcceptAny()}, nil
	case LeftParen:
		tokens.AcceptAny()
		result, err := ParseExpr(tokens)
		if err != nil {
			return nil, err
		}
		_, err = tokens.Accept(RightParen)
		if err != nil {
			return nil, err
		}
		return result, nil
	case LeftSquareBracket:
		result, err := ParseArrayLiteral(tokens)
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		nextToken := tokens.Peek()
		return nil, ParseError{
			StartPos: nextToken.Start,
			EndPos: nextToken.End,
			msg: fmt.Sprintf("unexpected token in input: %s", nextToken),
		}
	}
}

func ParseExprList(tokens *Tokens) ([]Expr, error) {
	result := make([]Expr, 0)
	for {
		switch tokens.Peek().Kind {
		case RightParen:
			return result, nil
		default:
			x, err := ParseExpr(tokens)
			if err != nil {
				return nil, err
			}
			result = append(result, x)
			if tokens.Peek().Kind == Comma {
				tokens.AcceptAny()
			}
		}
	}
}
