package validator

import (
	"testing"
)

var (
	_registerStructBenchErr error
)

func BenchmarkValidator_RegisterStruct(b *testing.B) {
	b.ReportAllocs()

	type benchStruct struct {
		Email           string `json:"email" vally:"email;email(.OtherField)"`
		Country         string `json:"country" vally:"country_code;required() && one_of('GB', 'IT', 'US')"`
		Other           int    `json:"other" vally:"required()"`
		DependOnCountry string `vally:"depend_on_country;(eq(.OtherField, 'GB') && required()) || true()"`
		NoTag           string `json:"no_tag"`
		Ignored         string `json:"ignored" vally:"-"`
		SelfReplaced    string `vally:"(eq(.OtherField, 'GB') && required(    )) || true() || eq('GB') && eq(1234)"`
	}

	v := NewValidator()

	var err error
	for n := 0; n < b.N; n++ {
		err = v.RegisterStruct(benchStruct{})
	}
	_registerStructBenchErr = err
}
