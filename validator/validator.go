package validator

import (
	"context"
	"fmt"
	"io"

	"github.com/osl4b/vally/internal/errutil"
	"github.com/osl4b/vally/internal/parser"
	"github.com/osl4b/vally/internal/scanner"
)

type Target interface {
}

type Function func(ctx context.Context, t Target) (bool, error)

type Validator struct {
	funcs map[string]Function
}

// NewValidator returns a new instance of a Validator.
func NewValidator(opts ...Option) *Validator {
	v := Validator{}
	for _, apply := range opts {
		apply(&v)
	}

	// if no validator function was registered use the default functions
	if v.funcs == nil {
		(&v).funcs = defaultFuncs()
	}
	return &v
}

// MustRegister adds a validator function with the given id.
// It will panic if a function with the same id is already registered.
func (v *Validator) MustRegister(id string, fn Function) {
	errutil.Must(v.Register(id, fn))
}

// Register adds a validator function with the given id.
// It will return an error if a function with the same id is already registered, nil otherwise.
func (v *Validator) Register(id string, fn Function) error {
	if _, exists := v.funcs[id]; exists {
		return fmt.Errorf("function %q is already registered", id)
	}
	v.funcs[id] = fn
	return nil
}

// Replace adds or replace the validator function with the given id.
func (v *Validator) Replace(id string, fn Function) {
	v.funcs[id] = fn
}

func (v *Validator) MustRegisterStruct(s interface{}) {
	errutil.Must(v.RegisterStruct(s))
}

func (v *Validator) RegisterStruct(s interface{}) error {
	// nolint:godox
	// FIXME: missing implementation
	return nil
}

func (v *Validator) ValidateStruct(ctx context.Context, s interface{}) error {
	// nolint:godox
	// FIXME: missing implementation
	return nil
}

func (v *Validator) ValidateValue(ctx context.Context, expr io.Reader, value interface{}) error {
	parsedExpr, err := parser.Parse(scanner.New(expr))
	if err != nil {
		return err
	}

	ev := newEvalVisitor(ctx, v, newValueTarget(value))
	if err := parsedExpr.Visit(ev); err != nil {
		return err
	}
	if !ev.Result() {
		return fmt.Errorf("not valid")
	}
	return nil
}

// lookupFunction returns the validator function with the given id or an error if the function is not registered.
func (v *Validator) lookupFunction(id string) (Function, error) {
	fn, exists := v.funcs[id]
	if !exists {
		return nil, ErrNotFound
	}
	return fn, nil
}
