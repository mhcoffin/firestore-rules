package parser

import "fmt"

type Rules struct {
	version Token
	service *ServiceStmt
}

func (rules *Rules) String() string {
	return fmt.Sprintf(`
rules_version = %s

%s
`, rules.version, rules.service)
}

func ParseRules(tokens *Tokens) (*Rules, error) {
	version, err := ParseRulesVersion(tokens)
	if err != nil {
		return nil, err
	}
	service, err := ParseServiceStmt(tokens)
	if err != nil {
		return nil, err
	}
	return &Rules{version: version, service: service}, nil
}

func ParseRulesVersion(tokens *Tokens) (Token, error) {
	switch tokens.Peek().Kind {
	case RulesVersion:
		tokens.AcceptAny()
		_, err := tokens.Accept(Eq)
		if err != nil {
			return Token{}, err
		}
		v, err := tokens.Accept(StringLiteral)
		if err != nil {
			return Token{}, err
		}
		_, err = tokens.Accept(SemiColon)
		if err != nil {
			return Token{}, err
		}
		return v, nil
	default:
		return Token{}, fmt.Errorf("must start with rules_version")
	}
}
