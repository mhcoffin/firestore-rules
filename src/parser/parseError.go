package parser

import "fmt"

type ParseError struct {
	StartPos InputPosition
	EndPos   InputPosition
	msg      string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("%s: %s", e.StartPos, e.msg)
}



