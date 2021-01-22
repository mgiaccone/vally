package vally

import (
	"context"

	"github.com/osl4b/vally/sdk"
	"github.com/osl4b/vally/validator"
)

var (
	_defaultValidator = validator.NewValidator()
)

// MustRegister
func MustRegister(name string, fn sdk.Function) {
	_defaultValidator.MustRegister(name, fn)
}

func Register(name string, fn sdk.Function) error {
	return _defaultValidator.Register(name, fn)
}

func MustRegisterStruct(s interface{}) {
	_defaultValidator.MustRegisterStruct(s)
}

func RegisterStruct(s interface{}) error {
	return _defaultValidator.RegisterStruct(s)
}

func ValidateStruct(ctx context.Context, s interface{}) error {
	return _defaultValidator.ValidateStruct(ctx, s)
}

func ValidateValue(ctx context.Context, expr string, value interface{}) error {
	return _defaultValidator.ValidateValue(ctx, expr, value)
}
