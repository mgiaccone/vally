package builtin

import (
	"context"
	"fmt"
	"regexp"

	"github.com/osl4b/vally/sdk"
)

// regexpMatch evaluates the given regular expression against the target.
type regexpMatch struct {
	errCode sdk.ErrCode
	re      *regexp.Regexp
}

func (f *regexpMatch) ErrCodes() []sdk.ErrCode {
	return []sdk.ErrCode{f.errCode}
}

func (f *regexpMatch) ArgTypes() []sdk.ArgType { return nil }

func (f *regexpMatch) Evaluate(_ context.Context, ec sdk.EvalContext, t sdk.Target) (bool, error) {
	fmt.Println("FN: ", ec.FunctionName(), " | FIELD: ", ec.FieldRef())

	v, _ := t.FieldRefValue(ec.FieldRef())
	fmt.Println("V: ", v)

	// FIXME: missing implementation
	return false, nil
}

// NewRegexpFunction builds a validator function that uses a regular expression to validate the target.
func NewRegexpFunction(re *regexp.Regexp, errCode sdk.ErrCode) sdk.Function {
	return &regexpMatch{
		errCode: errCode,
		re:      re,
	}
}
