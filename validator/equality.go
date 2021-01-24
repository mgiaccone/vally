package validator

import (
	"context"
)

var (
	// False is a validator function that always evaluates to false.
	Equal = &eqFunc{}
)

// eqFunc validator
type eqFunc struct{}

func (f *eqFunc) ArgTypes() []ArgType {
	return []ArgType{
		Any,
	}
}

func (f *eqFunc) Evaluate(_ context.Context, _ EvalContext, _ Target) (bool, error) {
	return false, nil
}
