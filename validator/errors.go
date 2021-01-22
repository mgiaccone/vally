package validator

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
)

type Error struct {
	FieldErrs []error
}

func (e *Error) Error() string {
	return fmt.Sprintf("validation failed with %d error(s)", len(e.FieldErrs))
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
	// FIXME: missing implementation
	return ""
}
