package builtin

import (
	"regexp"
)

var (
	// email
	reEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	Email   = NewRegexpFunction(reEmail)
)

// TODO: improve the email validation with optional MX checks.
//  It requires allowing flags in the expression scanner
