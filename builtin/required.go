package builtin

import (
	"context"
	"fmt"

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
	fmt.Println("FIELD: ", eval.FieldRef(), " | FN: ", eval.FunctionName(), " | ARGS: ", eval.FunctionArgs())

	v, _ := target.FieldRefValue(eval.FieldRef())
	fmt.Println("VALUE: ", v)

	// FIXME: missing implementation
	return false, nil
}
