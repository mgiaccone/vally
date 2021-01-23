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

func newStructTarget(s interface{}) (*structTarget, error) {
	values := make(map[string]interface{})

	// TODO: implement struct to map
	// out := make(map[string]interface{})
	// s.FillMap(out)
	// return out

	return &structTarget{values: values}, nil
}

func (t *structTarget) FieldRefValue(fieldRef string) (interface{}, error) {
	v, ok := t.values[fieldRef]
	if !ok {
		return nil, fmt.Errorf("field %q not found", fieldRef)
	}
	return v, nil
}

// valueTarget is a target implementation that works on a primitive value
type valueTarget struct {
	v interface{}
}

func newValueTarget(v interface{}) (*valueTarget, error) {
	// TODO: check if not primitive (?)
	//  value type must be one of int|float|string|etc....
	return &valueTarget{v: v}, nil
}

func (t *valueTarget) FieldRefValue(_ string) (interface{}, error) {
	return t.v, nil
}
