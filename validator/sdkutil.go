package validator

import (
	"fmt"
	"math"
	"reflect"
)

func IsZero(v interface{}) (bool, error) {
	vv := reflect.ValueOf(v)
	switch vv.Kind() {
	case reflect.Bool:
		return !vv.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return vv.Int() == 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return vv.Uint() == 0, nil
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(vv.Float()) == 0, nil
	case reflect.String:
		return vv.Len() == 0, nil
	case reflect.Ptr:
		if vv.IsNil() {
			return true, nil
		}
		return IsZero(vv.Elem().Interface())
	// case Array:
	// 	for i := 0; i < v.Len(); i++ {
	// 		if !v.Index(i).IsZero() {
	// 			return false
	// 		}
	// 	}
	// 	return true
	// case Interface, Map, Ptr, Slice:
	// 	return v.IsNil()

	// case Struct:
	// 	for i := 0; i < v.NumField(); i++ {
	// 		if !v.Field(i).IsZero() {
	// 			return false
	// 		}
	// 	}
	// 	return true
	default:
		return false, fmt.Errorf("unsupported value type %q", vv.Kind().String())
	}
}
