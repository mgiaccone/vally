package validator

import (
	"context"
	"testing"
)

func TestEmail_Evaluate(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		value   interface{}
		wantErr bool
	}{
		{
			name:    "success",
			expr:    "email()",
			value:   "domain.user@example.com",
			wantErr: false,
		},
		// {
		// 	name:    "parse error",
		// 	expr:    "bad_expression",
		// 	wantErr: true,
		// },
		// {
		// 	name:    "simple expression evaluation",
		// 	expr:    "true()",
		// 	wantErr: false,
		// },
		// {
		// 	name:    "value required (fail)",
		// 	expr:    "required()",
		// 	value:   "",
		// 	wantErr: true,
		// },
		// {
		// 	name:    "value required (success)",
		// 	expr:    "required()",
		// 	value:   "aaa",
		// 	wantErr: false,
		// },
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
			v := NewValidator()
			if err := v.ValidateValue(context.Background(), tt.expr, tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Email() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
