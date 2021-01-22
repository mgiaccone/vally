package sdk

import (
	"context"
)

// ArgType
type ArgType int

// TODO: bitmask operations examples here:
//  https://yourbasic.org/golang/bitmask-flag-set-clear/
const (
	FieldRef ArgType = 1 << iota
	String
	Int
	Float
	StringArray
	IntArray
	FloatArray
	Any
)

// Function represents the minimal methods a validation function must implement.
type Function interface {
	ArgTypes() []ArgType
	Evaluate(ctx context.Context, ec EvalContext, t Target) (bool, error)
}

// FunctionMeta
type FunctionMeta interface {
	ErrCodes() []string
}

// EvalContext represents the evaluation context of the function being processed.
type EvalContext interface {
	FieldRef() string
	FunctionName() string
}

// Target wraps the targer being validated
type Target interface {
	ValueOf(fieldRef string) (interface{}, error)
}
