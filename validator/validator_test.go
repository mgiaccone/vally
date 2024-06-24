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
		Email      string   `vally:"email;email()"`
		EmailSlice []string `vally:"emails;email()"`
		Required   string   `vally:"required;required()"`
		Ignored    string   `vally:"-"`
		NoTag      string
	}
	v := testStruct{
		Email:    "test@example.com",
		Required: "GB",
		NoTag:    "NoTag",
		Ignored:  "Ignore",
	}
	err := NewValidator().ValidateStruct(context.Background(), v)
	require.NoError(t, err, "ValidateStruct() error = %v", err)
}

func TestValidator_ValidateValue(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		value   interface{}
		wantErr bool
	}{
		{
			name:    "empty expression",
			expr:    "",
			wantErr: true,
		},
		{
			name:    "parse error",
			expr:    "bad_expression",
			wantErr: true,
		},
		{
			name:    "simple expression evaluation",
			expr:    "true()",
			wantErr: false,
		},
		{
			name:    "value required (fail)",
			expr:    "required()",
			value:   "",
			wantErr: true,
		},
		{
			name:    "value required (success)",
			expr:    "required()",
			value:   "aaa",
			wantErr: false,
		},
		// {
		// 	name:    "value equal (success)",
		// 	expr:    "eq('fail')",
		// 	value:   "do_fail",
		// 	wantErr: true,
		// },
		// {
		// 	name:    "value equal (success)",
		// 	expr:    "eq('success')",
		// 	value:   "success",
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := NewValidator().ValidateValue(context.Background(), tt.expr, tt.value); (err != nil) != tt.wantErr {
				t.Errorf("ValidateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_patchExprScanner(t *testing.T) {
	testExpr := "(eq(.OtherField, 'GB') && required(    )) || true(    ) || eq('GB') && eq(1234) && eq(.OtherField,1234)"
	testFieldRef := ".TestFieldRef"
	expectedExpr := "(eq(.TestFieldRef,.OtherField, 'GB') && required(.TestFieldRef,.TestFieldRef)) || true(.TestFieldRef,.TestFieldRef) || eq(.TestFieldRef,.TestFieldRef,'GB') && eq(.TestFieldRef,.TestFieldRef,1234) && eq(.TestFieldRef,.OtherField,1234)"
	got, err := patchExprScanner(testExpr, testFieldRef)
	require.NoError(t, err)
	require.Equal(t, expectedExpr, got)
}
