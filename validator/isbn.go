package validator

import (
	"regexp"
)

const (
	ErrCodeISBN10Invalid = ErrCode("isbn10/invalid")
	ErrCodeISBN13Invalid = ErrCode("isbn13/invalid")
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
