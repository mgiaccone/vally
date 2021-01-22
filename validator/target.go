package validator

import (
	"fmt"
)

var (
	_ Target = (*structTarget)(nil)
	_ Target = (*valueTarget)(nil)
)

// Target wraps the targer being validated
type Target interface {
	ValueOf(fieldRef string) (interface{}, error)
}

type structTarget struct {
	values map[string]interface{}
}

func newStructTarget(s interface{}) *structTarget {
	values := make(map[string]interface{})

	// TODO: implement struct value extraction

	return &structTarget{values: values}
}

func (t *structTarget) ValueOf(fieldRef string) (interface{}, error) {
	v, ok := t.values[fieldRef]
	if !ok {
		return nil, fmt.Errorf("field %q not found", fieldRef)
	}
	return v, nil
}

type valueTarget struct {
	v interface{}
}

func newValueTarget(v interface{}) *valueTarget {
	return &valueTarget{v: v}
}

func (t *valueTarget) ValueOf(_ string) (interface{}, error) {
	return t.v, nil
}
