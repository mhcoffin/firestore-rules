package parser

import (
	"fmt"
	"strings"
)

type Path []Component

func (p Path) String() string {
	s := make([]string, len(p))
	for k, v := range p {
		s[k] = v.String()
	}
	return "/" + strings.Join(s, "/")
}

type Component struct {
	wildcard bool
	recursive bool
	literal Token
}

func (c Component) String() string {
	if c.recursive {
		return fmt.Sprintf("{%s=**}", c.literal)
	} else if c.wildcard {
		return fmt.Sprintf("{%s}", c.literal)
	} else {
		return fmt.Sprintf("%s", c.literal)
	}
}

func ParsePath(tokens *Tokens) (Path, error) {
	components := make([]Component, 0)
	for {
		if tokens.Peek().Kind == Slash {
			tokens.AcceptAny()
		} else if len(components) == 0 {
			return nil, fmt.Errorf("%s: Path expected", tokens.Peek().Start)
		} else {
			return components, nil
		}

		switch tokens.Peek().Kind {
		case Identifier:
			components = append(components, Component{literal: tokens.AcceptAny()})
		case LeftBrace:
			tokens.AcceptAny()
			name, err := tokens.Accept(Identifier)
			if err != nil {
				return nil, err
			}
			recursive := false
			if tokens.Peek().Kind == Eq {
				tokens.AcceptAny()
				tokens.Accept(StarStar)
				recursive = true
			}
			components = append(components, Component{literal: name, wildcard: true, recursive: recursive})
			_, err = tokens.Accept(RightBrace)
			if err != nil {
				return nil, err
			}
		case Slash:
			tokens.AcceptAny()
		default:
			return nil, fmt.Errorf("%s: unexpected token: %s", tokens.Peek().Start, tokens.Peek())
		}
	}
}