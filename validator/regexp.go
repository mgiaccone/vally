package validator

import (
	"context"
	"regexp"
)

// NewRegexpFunction builds a validator function that uses the given
// regular expression to validate the target value.
func NewRegexpFunction(re *regexp.Regexp, errCode ErrCode) Function {
	return &regexpMatchFunc{
		errCode: errCode,
		re:      re,
	}
}

// regexpMatchFunc evaluates the given regular expression against the target.
type regexpMatchFunc struct {
	errCode ErrCode
	re      *regexp.Regexp
}

// ErrCodes returns a list of all the error codes the function might return
func (f *regexpMatchFunc) ErrCodes() []ErrCode {
	return []ErrCode{f.errCode}
}

// Evaluate implements the Function interface
func (f *regexpMatchFunc) Evaluate(_ context.Context, eval EvalContext, target Target) (bool, error) {
	v, err := target.FieldRefValue(eval.TargetRef())
	if err != nil {
		return false, err
	}

	if f.re.MatchString(v.(string)) {
		return false, eval.NewFieldError(f.errCode)
	}
	return false, nil
}
