package validator

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type ValidationError struct {
	ErrCode       string
	FieldFullPath string
	FieldName     string
	FieldPath     string
	FieldValue    interface{}
	FunctionArgs  map[string]interface{}
	FunctionName  string
}

func (e *ValidationError) Error() string {
	return ""
}

func (e *ValidationError) Unwrap() error {
	return nil
}
