package main

import (
	"fmt"
	"reflect"

	"github.com/osl4b/vally/internal/reflectutil"
)

const (
	structTag = "vally"
)

type SampleStruct struct {
	Email           string `json:"email" vally:"email;email(.OtherField,strict=true)"`
	Country         string `json:"country" vally:"country_code;required() && one_of('GB', 'IT', 'US')"`
	Other           int    `json:"other" vally:"required()"`                                           // always required
	DependOnCountry string `vally:"depend_on_country;(eq(.OtherField, 'GB') && required()) || true()"` // required if country=GB
}

func main() {
	val := SampleStruct{
		Email: "test@example.com",
		Other: 0,
	}

	fmt.Println("OUT:", parse(&val))
}

func parse(s interface{}) error {
	if reflectutil.IsNil(s) || !reflectutil.IsPointer(s) {
		return fmt.Errorf("must be a valid struct pointer")
	}

	ns := reflectutil.StructKey(s)
	fmt.Println("NS:", ns)

	sv := reflect.ValueOf(s)
	st := sv.Type()
	fmt.Println("ST", st)

	el := st.Elem()

	for i := 0; i < el.NumField(); i++ {
		fl := el.Field(i)
		tag := fl.Tag.Get(structTag)
		fmt.Println("FIELD:", fl.Name, "TYPE:", fl.Type.String(), "TAG:", tag)
	}

	return nil
}

func parseTag(name string, s interface{}) {

}
