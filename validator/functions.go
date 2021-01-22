package validator

import (
	"context"
	"regexp"
)

var (
	// email
	reEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	Email   = RegexpMatch(reEmail)
)

func defaultFunctions() map[string]Function {
	return map[string]Function{
		"false":    False,
		"true":     True,
		"required": Required,
		"email":    Email,
	}
}

// False always evaluates to false
func False(_ context.Context, _ []Arg, _ Target) (bool, error) {
	return false, nil
}

// True always evaluates to true
func True(_ context.Context, _ []Arg, _ Target) (bool, error) {
	return true, nil
}

// Required evaluates to true if the value does not equal to the type's default value, false otherwise.
func Required(ctx context.Context, args []Arg, t Target) (bool, error) {
	// nolint:godox
	// FIXME: missing implementation
	return false, nil
}

// RegexpMatch returns a validator functions that uses the given regexp to validate the target.
func RegexpMatch(re *regexp.Regexp) Function {
	return func(ctx context.Context, args []Arg, t Target) (bool, error) {
		// nolint:godox
		// TODO: get value as string

		// re.M

		// nolint:godox
		// FIXME: missing implementation
		return false, nil
	}
}
