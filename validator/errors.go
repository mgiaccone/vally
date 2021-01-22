package validator

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type Error struct {
	FieldErrs []FieldError
}

func (e *Error) Error() string {
	return ""
}

func (e *Error) Unwrap() error {
	return nil
}

type FieldError struct {
	ErrCode       string
	FieldFullPath string
	FieldName     string
	FieldPath     string
	FieldValue    interface{}
	FunctionArgs  map[string]interface{}
	FunctionName  string
}

func (fe *FieldError) Error() string {
	return ""
}

func (fe *FieldError) Unwrap() error {
	return nil
}
