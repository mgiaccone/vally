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

	Any = FieldRef | String | Int | Float
)

// ArgTyper
type ArgTyper interface {
	ArgTypes() []ArgType
}

// ErrCode
type ErrCode string

// ErrCoder
type ErrCoder interface {
	ErrCodes() []ErrCode
}

// Function represents the minimal methods a validation function must implement.
type Function interface {
	Evaluate(ctx context.Context, ec EvalContext, t Target) (bool, error)
}

// EvalContext represents the evaluation context of the function being processed.
type EvalContext interface {
	FieldRef() string
	FunctionArgs() []ArgValue
	FunctionName() string
	TargetRef() string
}

// Target wraps the targer being validated
type Target interface {
	FieldRefValue(fieldRef string) (interface{}, error)
}

// TODO: This might need to be revisited, it probably needs more than
//  just the value or it has to be linked to the ArgType for checks?
//  Also it could be useful for returning the field error with names
//  for the argument values
//
// ArgValue represents the actual value of a function argument
type ArgValue interface{}
