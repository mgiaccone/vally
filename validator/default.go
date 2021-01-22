package validator

import (
	"github.com/osl4b/vally/builtin"
	"github.com/osl4b/vally/sdk"
)

func defaultFunctions() map[string]sdk.Function {
	return map[string]sdk.Function{
		"false":    builtin.False,
		"true":     builtin.True,
		"required": builtin.Required,
		"email":    builtin.Email,
	}
}
