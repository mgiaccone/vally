package builtin

import (
	"context"

	"github.com/osl4b/vally/sdk"
)

const (
	ErrCodeRequired = sdk.ErrCode("required")
)

var (
	// Required evaluates to true if the value does not equal to the type's default value, false otherwise.
	Required = &requiredFunc{}
)

type requiredFunc struct{}

func (f *requiredFunc) Evaluate(_ context.Context, eval sdk.EvalContext, target sdk.Target) (bool, error) {
	v, err := target.FieldRefValue(eval.TargetRef())
	if err != nil {
		return false, err
	}

	isZero, err := sdk.IsZero(v)
	if err != nil {
		return false, err
	}
	if isZero {
		return false, &sdk.FieldError{ErrCode: ErrCodeRequired}
	}
	return false, nil
}
