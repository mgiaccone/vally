package validator

import (
	"context"
	"fmt"
	"regexp"
)

var (
	// email
	reEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	Email   = RegexpMatch(reEmail)
)

func defaultFunctions() map[string]Function {
	return map[string]Function{
		"false":    &False{},
		"true":     &True{},
		"required": &Required{},
		"email":    Email,
	}
}

// False always evaluates to false
type False struct{}

func (f *False) ArgTypes() []ArgType { return nil }

func (f *False) Evaluate(_ context.Context, _ EvalContext, _ Target) (bool, error) {
	return false, nil
}

// True always evaluates to true
type True struct{}

func (f *True) ArgTypes() []ArgType { return nil }

func (f *True) Evaluate(ctx context.Context, ec EvalContext, t Target) (bool, error) {
	return true, nil
}

// Required evaluates to true if the value does not equal to the type's default value, false otherwise.
type Required struct{}

func (f *Required) ArgTypes() []ArgType {
	return []ArgType{Any}
}

func (f *Required) Evaluate(ctx context.Context, ec EvalContext, t Target) (bool, error) {
	fmt.Println("FN: ", ec.FunctionName(), " | FIELD: ", ec.FieldRef())

	// FIXME: missing implementation
	return false, nil
}

// regex evaluates the given regular expression against the target.
type regex struct {
	re *regexp.Regexp
}

func (f *regex) ArgTypes() []ArgType { return nil }

func (f *regex) Evaluate(ctx context.Context, ec EvalContext, t Target) (bool, error) {
	fmt.Println("FN: ", ec.FunctionName(), " | FIELD: ", ec.FieldRef())

	// TODO: get value as string

	// re.M

	// FIXME: missing implementation
	return false, nil
}

// RegexpMatch builds a validator function that uses a regular expression to validate the target.
func RegexpMatch(re *regexp.Regexp) Function {
	return &regex{re: re}
}
