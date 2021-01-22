package validator

import (
	"context"
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

func TestValidator_ValidateStruct(t *testing.T) {
	type testStruct struct {
		Email           string `vally:"email;required() && email()"`
		Country         string `vally:"country_code;required() && one_of('GB', 'IT', 'US')"`
		Other           int    `vally:"required()"`
		DependOnCountry string `vally:"depend_on_country;(eq(.OtherField, 'GB') && required()) || true()"`
		NoTag           string
		Ignored         string `vally:"-"`
		SelfReplaced    string `vally:"(eq(.OtherField, 'GB') && required(    )) || true() || eq('GB') && eq(1234)"`
	}
	v := testStruct{
		Email:           "test@example.com",
		Country:         "GB",
		Other:           100,
		DependOnCountry: "Test",
		NoTag:           "Ignore",
		Ignored:         "Ignore",
		SelfReplaced:    "X",
	}
	err := NewValidator().ValidateStruct(context.Background(), v)
	require.NoError(t, err, "ValidateStruct() error = %v", err)
}
