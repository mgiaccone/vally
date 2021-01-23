package builtin

import (
	"regexp"

	"github.com/osl4b/vally/sdk"
)

const (
	//	EmailInvalidDomain = sdk.ErrCode("email/invalid_domain")
	EmailInvalidFormat = sdk.ErrCode("email/invalid_format")
)

var (
	// email
	reEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	Email   = NewRegexpFunction(reEmail, EmailInvalidFormat)
)

// TODO: improve the email validation with optional MX checks.
//  It requires allowing flags in the expression scanner
