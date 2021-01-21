package scanner

import (
	"fmt"
)

// TokenType represents a lexical Token.
type TokenType int

// Constants representing lexical tokens
const (
	// Special tokens
	EOF TokenType = iota
	ILLEGAL

	// Operators and delimiters
	LPAREN // (
	RPAREN // )
	COMMA  // ,
	LAND   // &&
	LNOT   // !
	LOR    // ||

	// Identifiers and type literals
	IDENT    // fn0
	STRING   // 'string'
	FLOAT    // 1.23434
	INT      // 123456
	FIELDREF // .Name
	FLAG     // flag1
)

var _tokenTypes = []string{
	"EOF",
	"ILLEGAL",
	"LPAREN",
	"RPAREN",
	"COMMA",
	"LAND",
	"LNOT",
	"LOR",
	"IDENT",
	"STRING",
	"FLOAT",
	"INT",
	"FIELDREF",
	"FLAG",
}

func (t TokenType) String() string {
	return _tokenTypes[t]
}

// Token represents a Token returned from the scanner.
type Token struct {
	Type     TokenType
	Literal  string
	StartPos int
	EndPos   int
}

// newToken creates a new token with the given values
func newToken(typ TokenType, literal string, start, end int) Token {
	return Token{
		Type:     typ,
		Literal:  literal,
		StartPos: start,
		EndPos:   end,
	}
}

func (t Token) String() string {
	var pos string
	if t.StartPos == t.EndPos {
		pos = fmt.Sprintf("1:%d", t.StartPos)
	} else {
		pos = fmt.Sprintf("1:%d-%d", t.StartPos, t.EndPos)
	}

	return fmt.Sprintf("%-12s%-10s %s", pos, t.Type, t.Literal)
}
