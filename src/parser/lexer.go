package parser

//go:generate stringer -type=Kind

import (
	"fmt"
	"os"
	"unicode"
)

type Kind int

const (
	Error Kind = iota
	Eof
	Word
	Dot
	IntLiteral
	FloatLiteral
	LeftBrace
	RightBrace
	LeftParen
	RightParen
	StringLiteral
	Slash
	Minus
	Plus
	Comma
	Eq
	EqEq
	SemiColon
	LeftSquareBracket
	RightSquareBracket
	Less
	LessEq
	Greater
	GreaterEq
	Colon
	QuestionMark
	And
	AndAnd
	Or
	OrOr
	NotEq
	Bang
	Star
	StarStar
	Bytes
	Percent
	Identifier

	// reserved words
	Service
	Match
	Allow
	Create
	Update
	Delete
	Write
	Get
	List
	Read
	If
	Function
	True
	False
	In
	Is
	Return
	Let
	RulesVersion
)

var reservedWords = map[string]Kind{
	"service":  Service,
	"match":    Match,
	"allow":    Allow,
	"create":   Create,
	"update":   Update,
	"delete":   Delete,
	"write":    Write,
	"get":      Get,
	"list":     List,
	"read":     Read,
	"if":       If,
	"function": Function,
	"true":     True,
	"false":    False,
	"in":       In,
	"is":       Is,
	"return":   Return,
	"let":      Let,
	"rules_version": RulesVersion,
}

func IsAction(kind Kind) bool {
	return kind == Read || kind == Write || kind == Get || kind == List || kind == Create || kind == Update || kind == Delete
}

const EOF = -1

type InputPosition struct {
	Pos  int
	Line int
	Col  int
}

func (ip InputPosition) String() string {
	return fmt.Sprintf("line %d col %d", ip.Line+1, ip.Col+1)
}

type LexError struct {
	Start InputPosition
	Pos InputPosition
	msg string
}

func (le LexError) Error() string {
	return fmt.Sprintf("%s: %s", le.Start, le.msg)
}

type Token struct {
	Kind  Kind
	Value string
	Start InputPosition
	End   InputPosition
	Error error
}

func (token Token) String() string {
	return token.Value
}

func (token Token) ErrString() string {
	switch token.Kind {
	case Eof:
		return "EOF"
	default:
		return token.Value
	}
}

type Lexer struct {
	Input []rune
	start InputPosition
	pos   InputPosition
	Chan  chan Token
}

func Start(s string) chan Token {
	ch := make(chan Token)
	lexer := &Lexer{
		Input: []rune(s),
		Chan:  ch,
		start: InputPosition{},
		pos:   InputPosition{},
	}
	go Run(lexer)
	return ch
}

func (lexer *Lexer) peek() rune {
	if len(lexer.Input) == lexer.pos.Pos {
		return EOF
	}
	return lexer.Input[lexer.pos.Pos]
}

func (lexer *Lexer) expect(r rune) error {
	if lexer.peek() != r {
		return LexError{lexer.start, lexer.pos, fmt.Sprintf("expected %s", string(r))}
	}
	lexer.acceptChar()
	return nil
}

func (lexer *Lexer) acceptCharAndGenerate(kind Kind) {
	lexer.acceptChar()
	lexer.generate(kind)
}

func (lexer *Lexer) acceptChar() {
	c := lexer.Input[lexer.pos.Pos]
	lexer.pos.Pos++
	if c == '\n' {
		lexer.pos.Line++
		lexer.pos.Col = 0
	} else {
		lexer.pos.Col++
	}
}

func (lexer *Lexer) current() string {
	return string(lexer.Input[lexer.start.Pos:lexer.pos.Pos])
}

func (lexer *Lexer) generate(kind Kind) {
	if kind == Word {
		id, ok := reservedWords[lexer.current()]
		if ok {
			lexer.Chan <- Token{Kind: id, Value: lexer.current(), Start: lexer.start, End: lexer.pos}
		} else {
			lexer.Chan <- Token{Kind: Identifier, Value: lexer.current(), Start: lexer.start, End: lexer.pos}
		}
	} else {
		lexer.Chan <- Token{Kind: kind, Value: lexer.current(), Start: lexer.start, End: lexer.pos}
	}
}

type stFunction func(*Lexer) (stFunction, error)

func startState(lexer *Lexer) (stFunction, error) {
	lexer.start = lexer.pos
	c := lexer.peek()
	switch {
	case c == '\n':
		lexer.acceptChar()
	case unicode.IsSpace(c):
		lexer.acceptChar()
	case c == EOF:
		lexer.generate(Eof)
		close(lexer.Chan)
		return nil, nil
	case c == 'b':
		lexer.acceptChar()
		if lexer.peek() == '\'' {
			err := acceptStringLiteral(lexer, lexer.peek())
			if err != nil {
				return nil, fmt.Errorf("error in bytes literal: %w", err)
			}
			lexer.generate(Bytes)
		} else {
			acceptIdentifier(lexer)
			lexer.generate(Word)
		}
	case c == '.':
		lexer.acceptCharAndGenerate(Dot)
	case unicode.IsLetter(c):
		acceptIdentifier(lexer)
		lexer.generate(Word)
	case c == '"':
		err := acceptStringLiteral(lexer, c)
		if err != nil {
			return nil, fmt.Errorf("error in string literal: %w", err)
		}
		lexer.generate(StringLiteral)
	case c == '\'':
		err := acceptStringLiteral(lexer, c)
		if err != nil {
			return nil, fmt.Errorf("error in string literal: %w", err)
		}
		lexer.generate(StringLiteral)
	case c == ',':
		lexer.acceptCharAndGenerate(Comma)
	case c == '=':
		lexer.acceptChar()
		if lexer.peek() == '=' {
			lexer.acceptCharAndGenerate(EqEq)
		} else {
			lexer.generate(Eq)
		}
	case c == '+':
		lexer.acceptCharAndGenerate(Plus)
	case c == '-':
		lexer.acceptCharAndGenerate(Minus)
	case unicode.IsDigit(c):
		acceptDigits(lexer)
		if lexer.peek() == '.' {
			lexer.acceptChar()
			acceptDigits(lexer)
			lexer.generate(FloatLiteral)
		} else {
			lexer.generate(IntLiteral)
		}
	case c == '(':
		lexer.acceptCharAndGenerate(LeftParen)
	case c == ')':
		lexer.acceptCharAndGenerate(RightParen)
	case c == '{':
		lexer.acceptCharAndGenerate(LeftBrace)
	case c == '}':
		lexer.acceptCharAndGenerate(RightBrace)
	case c == '/':
		lexer.acceptCharAndGenerate(Slash)
	case c == ';':
		lexer.acceptCharAndGenerate(SemiColon)
	case c == '[':
		lexer.acceptCharAndGenerate(LeftSquareBracket)
	case c == ']':
		lexer.acceptCharAndGenerate(RightSquareBracket)
	case c == '<':
		lexer.acceptChar()
		if lexer.peek() == '=' {
			lexer.acceptCharAndGenerate(LessEq)
		} else {
			lexer.generate(Less)
		}
	case c == '>':
		lexer.acceptChar()
		if lexer.peek() == '=' {
			lexer.acceptCharAndGenerate(GreaterEq)
		} else {
			lexer.generate(Greater)
		}
	case c == ':':
		lexer.acceptCharAndGenerate(Colon)
	case c == '?':
		lexer.acceptCharAndGenerate(QuestionMark)
	case c == '&':
		lexer.acceptChar()
		if lexer.peek() == '&' {
			lexer.acceptCharAndGenerate(AndAnd)
		} else {
			lexer.generate(And)
		}
	case c == '|':
		lexer.acceptChar()
		if lexer.peek() == '|' {
			lexer.acceptCharAndGenerate(OrOr)
		} else {
			lexer.generate(Or)
		}
	case c == '!':
		lexer.acceptChar()
		if lexer.peek() == '=' {
			lexer.acceptCharAndGenerate(NotEq)
		} else {
			lexer.generate(Bang)
		}
	case c == '*':
		lexer.acceptChar()
		if lexer.peek() == '*' {
			lexer.acceptCharAndGenerate(StarStar)
		} else {
			lexer.generate(Star)
		}
	case c == '%':
		lexer.acceptCharAndGenerate(Percent)
	default:
		return nil, LexError{lexer.start, lexer.pos, fmt.Sprintf("unexpected character: %c", c)}
	}
	return startState, nil
}

func acceptDigits(lexer *Lexer) {
	for {
		switch {
		case unicode.IsDigit(lexer.peek()):
			lexer.acceptChar()
		default:
			return
		}
	}
}

func acceptIdentifier(lexer *Lexer) {
	for c := lexer.peek(); unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_'; c = lexer.peek() {
		lexer.acceptChar()
	}
}

func acceptStringLiteral(lexer *Lexer, delimiter rune) error {
	err := lexer.expect(delimiter)
	if err != nil {
		return err
	}
	for {
		switch lexer.peek() {
		case delimiter:
			lexer.acceptChar()
			return nil
		case -1, '\n':
			return LexError{lexer.start,lexer.pos, "unclosed string literal"}
		case '\\':
			lexer.acceptChar()
			c := lexer.peek()
			if c == '\n' || c == -1 {
				return LexError{lexer.start, lexer.pos, "cannot escape newline or EOF"}
			}
			lexer.acceptChar()
		default:
			lexer.acceptChar()
		}
	}
}

func Run(lexer *Lexer) {
	for state := startState; state != nil; {
		nextState, err := state(lexer)
		if err != nil {
			lexError := err.(*LexError)
			_, _ = fmt.Fprintf(os.Stderr, "lexical error: %s\n", lexError.Error())
			errToken := Token{
				Kind: Error,
				Value: string(lexer.Input[lexer.start.Pos: lexer.pos.Pos]),
				Start: lexError.Start,
				End: lexError.Pos,
				Error: lexError,
			}
			lexer.Chan <- errToken
			close(lexer.Chan)
		}
		state = nextState
	}
}
