package builtin

import (
	"regexp"

	"github.com/osl4b/vally/sdk"
)

const (
	ErrCodeISBN10Invalid = sdk.ErrCode("isbn10/invalid")
	ErrCodeISBN13Invalid = sdk.ErrCode("isbn13/invalid")
)

var (
	// FIXME: missing regular expressions
	reISBN10 = regexp.MustCompile(``)
	reISBN13 = regexp.MustCompile(``)

	// ISBN10 validates an ISBN10 only
	ISBN10 = NewRegexpFunction(reISBN10, ErrCodeISBN10Invalid)

	// ISBN13 validates an ISBN13 only
	ISBN13 = NewRegexpFunction(reISBN13, ErrCodeISBN13Invalid)
)
