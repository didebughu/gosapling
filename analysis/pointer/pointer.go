package pointer

import "go/token"

type Pointer struct {
	name    string
	namePos token.Pos
	isNil   bool
}
