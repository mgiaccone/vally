package validator

import (
	"fmt"
)

type ValidationError struct {
	FieldErrs []error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed with %d error(s)", len(e.FieldErrs))
}

type FieldError struct {
	ErrCode       ErrCode
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
