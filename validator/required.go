package validator

import (
	"context"
)

const (
	ErrCodeRequired = ErrCode("required")
)

var (
	// Required evaluates to true if the value does not equal to the type's default value, false otherwise.
	Required = &requiredFunc{}
)

type requiredFunc struct{}

func (f *requiredFunc) Evaluate(_ context.Context, eval EvalContext, target Target) (bool, error) {
	v, err := target.FieldRefValue(eval.TargetRef())
	if err != nil {
		return false, err
	}

	isZero, err := IsZero(v)
	if err != nil {
		return false, err
	}
	if isZero {
		return false, eval.NewFieldError(ErrCodeRequired)
	}
	return false, nil
}
