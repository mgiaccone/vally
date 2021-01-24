package validator

import (
	"context"
)

var (
	// False is a validator function that always evaluates to false.
	False = &falseFunc{}

	// True is a validator function that always evaluates to true.
	True = &trueFunc{}
)

// falseFunc validator
type falseFunc struct{}

func (f *falseFunc) Evaluate(_ context.Context, _ EvalContext, _ Target) (bool, error) {
	return false, nil
}

// trueFunc validator
type trueFunc struct{}

func (f *trueFunc) Evaluate(_ context.Context, _ EvalContext, _ Target) (bool, error) {
	return true, nil
}
