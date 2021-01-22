package builtin

import (
	"context"

	"github.com/osl4b/vally/sdk"
)

var (
	// False is a validator function that always evaluates to false.
	False = &falseFunc{}

	// True is a validator function that always evaluates to true.
	True = &trueFunc{}
)

// falseFunc validator
type falseFunc struct{}

func (f *falseFunc) ArgTypes() []sdk.ArgType { return nil }

func (f *falseFunc) Evaluate(_ context.Context, _ sdk.EvalContext, _ sdk.Target) (bool, error) {
	return false, nil
}

// trueFunc validator
type trueFunc struct{}

func (f *trueFunc) ArgTypes() []sdk.ArgType { return nil }

func (f *trueFunc) Evaluate(_ context.Context, _ sdk.EvalContext, _ sdk.Target) (bool, error) {
	return true, nil
}
