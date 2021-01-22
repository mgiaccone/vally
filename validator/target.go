package validator

import (
	"fmt"

	"github.com/osl4b/vally/sdk"
)

var (
	_ sdk.Target = (*structTarget)(nil)
	_ sdk.Target = (*valueTarget)(nil)
)

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
