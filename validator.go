package vally

import (
	"github.com/osl4b/vally/internal/errutil"
)

var (
	_default = NewValidator()
)

type Func func() error

// NewValidator returns a new instance of a Validator.
func NewValidator() *Validator {
	return &Validator{
		r: _defaultRegistry.Copy(),
	}
}

type Validator struct {
	c *cache
	r *Registry
}

func (v *Validator) MustRegisterFunc(id string, fn Func) {
	errutil.Must(v.RegisterFunc(id, fn))
}

func (v *Validator) MustRegisterStruct(s interface{}) {
	errutil.Must(v.RegisterStruct(s))
}

func (v *Validator) RegisterFunc(id string, fn Func) error {
	return v.r.Register(id, fn)
}

func (v *Validator) RegisterStruct(s interface{}) error {
	// nolint:godox
	// FIXME: missing implementation
	return nil
}

func (v *Validator) ValidateStruct(s interface{}) error {
	// nolint:godox
	// FIXME: missing implementation
	return nil
}

func MustRegisterFunc(name string, fn Func) {
	_default.MustRegisterFunc(name, fn)
}

func MustRegisterStruct(s interface{}) {
	_default.MustRegisterStruct(s)
}

func RegisterFunc(name string, fn Func) error {
	return _default.RegisterFunc(name, fn)
}

func RegisterStruct(s interface{}) error {
	return _default.RegisterStruct(s)
}

func ValidateStruct(s interface{}) error {
	return _default.ValidateStruct(s)
}
