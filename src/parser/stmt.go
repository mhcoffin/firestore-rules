package parser

import "fmt"

// This is a marker for statement-level nodes: 'allow' and 'function'.
type Stmt interface{
	fmt.Stringer
}
