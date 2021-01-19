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
