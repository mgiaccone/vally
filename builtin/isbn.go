package builtin

import (
	"regexp"

	"github.com/osl4b/vally/sdk"
)

const (
	ISBNInvalid   = sdk.ErrCode("isbn/invalid")
	ISBN10Invalid = sdk.ErrCode("isbn10/invalid")
	ISBN13Invalid = sdk.ErrCode("isbn13/invalid")
)

var (
	// FIXME: missing regular expressions
	reISBN   = regexp.MustCompile(``)
	reISBN10 = regexp.MustCompile(``)
	reISBN13 = regexp.MustCompile(``)

	// ISBN validates an ISBN10 or ISBN13
	ISBN = NewRegexpFunction(reISBN, ISBNInvalid)

	// ISBN10 validates an ISBN10 only
	ISBN10 = NewRegexpFunction(reISBN10, ISBN10Invalid)

	// ISBN13 validates an ISBN13 only
	ISBN13 = NewRegexpFunction(reISBN13, ISBN13Invalid)
)
