package pointer

import "go/token"

type Variable struct {
	Name    string
	NamePos token.Pos
}

type Pointer struct {
	Name    string
	NamePos token.Pos
	IsNil   bool
}
