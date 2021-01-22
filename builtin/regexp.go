package builtin

import (
	"context"
	"fmt"
	"regexp"

	"github.com/osl4b/vally/sdk"
)

// regexpMatch evaluates the given regular expression against the target.
type regexpMatch struct {
	re *regexp.Regexp
}

func (f *regexpMatch) ArgTypes() []sdk.ArgType { return nil }

func (f *regexpMatch) Evaluate(ctx context.Context, ec sdk.EvalContext, t sdk.Target) (bool, error) {
	fmt.Println("FN: ", ec.FunctionName(), " | FIELD: ", ec.FieldRef())

	// TODO: get value as string

	// re.M

	// FIXME: missing implementation
	return false, nil
}

// NewRegexpFunction builds a validator function that uses a regular expression to validate the target.
func NewRegexpFunction(re *regexp.Regexp) sdk.Function {
	return &regexpMatch{re: re}
}
