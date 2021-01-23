package builtin

import (
	"context"

	"github.com/osl4b/vally/sdk"
)

var (
	// False is a validator function that always evaluates to false.
	Equal = &eqFunc{}
)

// eqFunc validator
type eqFunc struct{}

func (f *eqFunc) ArgTypes() []sdk.ArgType {
	return []sdk.ArgType{
		sdk.Any,
	}
}

func (f *eqFunc) Evaluate(_ context.Context, _ sdk.EvalContext, _ sdk.Target) (bool, error) {
	return false, nil
}
