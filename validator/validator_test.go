// +build unit

package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidator_RegisterStruct(t *testing.T) {
	type testStruct struct {
		Email           string `json:"email" vally:"email;email(.OtherField)"`
		Country         string `json:"country" vally:"country_code;required() && one_of('GB', 'IT', 'US')"`
		Other           int    `json:"other" vally:"required()"`
		DependOnCountry string `vally:"depend_on_country;(eq(.OtherField, 'GB') && required()) || true()"`
		NoTag           string `json:"no_tag"`
		Ignored         string `json:"ignored" vally:"-"`
		SelfReplaced    string `vally:"(eq(.OtherField, 'GB') && required(    )) || true() || eq('GB') && eq(1234)"`
	}
	err := NewValidator().RegisterStruct(testStruct{})
	require.NoError(t, err, "RegisterStruct() error = %v", err)
}
