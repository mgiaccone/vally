package reflectutil

import (
	"reflect"
)

func StructKey(s interface{}) string {
	rt := reflect.TypeOf(s)
	name := rt.String()
	star := ""
	if rt.Name() == "" {
		if pt := rt; pt.Kind() == reflect.Ptr {
			star = "*"
			rt = pt.Elem()
		}
	}
	if rt.Name() != "" {
		if rt.PkgPath() == "" {
			name = star + rt.Name()
		} else {
			name = star + rt.PkgPath() + "#" + rt.Name()
		}
	}
	return name
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}

func IsPointer(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Ptr
}

// // IsZero return true if underlying type is equal to its zero value
// func (v Value) IsZero() bool {
// 	if zeroer, k := v.Zeroer(); k {
// 		return zeroer.IsZero()
// 	}
// 	if v.IsNil() {
// 		return true
// 	}
// 	t := v.Type()
// 	switch t.Kind() {
// 	case reflect.Map:
// 		return v.Len() == 0
// 	case reflect.Chan:
// 		return v.Len() == 0
// 	case reflect.Slice:
// 		s := t.NewSlice()
// 		return reflect.DeepEqual(v.Interface(), s.Interface())
// 	case reflect.Ptr:
// 		return v.Indirect().IsZero()
// 	default:
// 		return reflect.DeepEqual(v.Interface(), t.Zero().InterfaceOrNil())
// 	}
// }
//
// func StructType(s interface{}) reflect.Type {
// 	sv := reflect.ValueOf(s)
// 	st := sv.Type()
//
// 	for st.Kind() == reflect.Ptr {
// 		st = st.Elem()
// 	}
//
// 	if reflectutil.IsPointer(s) {
// 		st = st.Elem()
// 	}
//
// 	v := reflect.ValueOf(s)
//
//
// 	if v.Kind() != reflect.Struct {
// 		panic("not struct")
// 	}
//
// 	return v
// }
//
// func StructType(s interface{}) reflect.Value {
// 	v := reflect.ValueOf(s)
//
// 	// if pointer get the underlying element≤
// 	for v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}
//
// 	if v.Kind() != reflect.Struct {
// 		panic("not struct")
// 	}
//
// 	return v
// }
