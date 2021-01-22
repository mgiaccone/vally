package builtin

import (
	"context"
	"fmt"

	"github.com/osl4b/vally/sdk"
)

var (
	// Required evaluates to true if the value does not equal to the type's default value, false otherwise.
	Required = &requiredFunc{}
)

type requiredFunc struct{}

func (f *requiredFunc) ArgTypes() []sdk.ArgType {
	return []sdk.ArgType{sdk.Any}
}

func (f *requiredFunc) Evaluate(ctx context.Context, ec sdk.EvalContext, t sdk.Target) (bool, error) {
	fmt.Println("FN: ", ec.FunctionName(), " | FIELD: ", ec.FieldRef())

	// FIXME: missing implementation
	return false, nil
}
