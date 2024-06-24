package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/osl4b/vally/internal/ast"
	"github.com/osl4b/vally/internal/scanner"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		inputExpr string
		want      ast.Node
		wantErr   bool
	}{
		{
			name:      "empty expression",
			inputExpr: "",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "function without arguments",
			inputExpr: "fn()",
			want: &ast.Function{
				Name: "fn",
				Args: nil,
			},
			wantErr: false,
		},
		{
			name:      "function with missing first argument",
			inputExpr: "fn(,1)",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "function with unspecified argument after COMMA (end)",
			inputExpr: "fn(1,)",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "function with unspecified argument after COMMA (middle)",
			inputExpr: "fn(1,,2)",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "function with single argument (fieldref)",
			inputExpr: "fn(.Field.Name)",
			want: &ast.Function{
				Name: "fn",
				Args: []ast.FunctionArg{
					{
						FieldRef: ".Field.Name",
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "function with single argument (string)",
			inputExpr: "fn('string')",
			want: &ast.Function{
				Name: "fn",
				Args: []ast.FunctionArg{
					{
						StringValue: "string",
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "function with single argument (int=0)",
			inputExpr: "fn(0)",
			want: &ast.Function{
				Name: "fn",
				Args: []ast.FunctionArg{
					{
						IntValue: 0,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "function with single argument (positive int)",
			inputExpr: "fn(12345)",
			want: &ast.Function{
				Name: "fn",
				Args: []ast.FunctionArg{
					{
						IntValue: 12345,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "function with single argument (negative int)",
			inputExpr: "fn(-12345)",
			want: &ast.Function{
				Name: "fn",
				Args: []ast.FunctionArg{
					{
						IntValue: -12345,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "function with single argument (float=0.0)",
			inputExpr: "fn(0.0)",
			want: &ast.Function{
				Name: "fn",
				Args: []ast.FunctionArg{
					{
						FloatValue: 0,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "function with single argument (positive float)",
			inputExpr: "fn(0.98765)",
			want: &ast.Function{
				Name: "fn",
				Args: []ast.FunctionArg{
					{
						FloatValue: 0.98765,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "function with single argument (negative float)",
			inputExpr: "fn(-0.98765)",
			want: &ast.Function{
				Name: "fn",
				Args: []ast.FunctionArg{
					{
						FloatValue: -0.98765,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "function with multiple arguments",
			inputExpr: "fn(.FieldRef, 'string_value', 12345, -12345,0.98765,-0.98765)",
			want: &ast.Function{
				Name: "fn",
				Args: []ast.FunctionArg{
					{
						FieldRef: ".FieldRef",
					},
					{
						StringValue: "string_value",
					},
					{
						IntValue: 12345,
					},
					{
						IntValue: -12345,
					},
					{
						FloatValue: 0.98765,
					},
					{
						FloatValue: -0.98765,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "negate direct function",
			inputExpr: "!fn()",
			want: &ast.LogicNot{
				Expression: &ast.Function{
					Name: "fn",
					Args: nil,
				},
			},
			wantErr: false,
		},
		{
			name:      "negate grouped function",
			inputExpr: "!(fn())",
			want: &ast.LogicNot{
				Expression: &ast.Group{
					Expression: &ast.Function{
						Name: "fn",
						Args: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "root group of negated function",
			inputExpr: "(!fn())",
			want: &ast.Group{
				Expression: &ast.LogicNot{
					Expression: &ast.Function{
						Name: "fn",
						Args: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "operator logic and with missing left node",
			inputExpr: "&& fn()",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple logical and expression",
			inputExpr: "fn1() && fn2()",
			want: &ast.LogicAnd{
				Left: &ast.Function{
					Name: "fn1",
					Args: nil,
				},
				Right: &ast.Function{
					Name: "fn2",
					Args: nil,
				},
			},
			wantErr: false,
		},
		{
			name:      "grouped logical and expression",
			inputExpr: "(fn1() && fn2())",
			want: &ast.Group{
				Expression: &ast.LogicAnd{
					Left: &ast.Function{
						Name: "fn1",
						Args: nil,
					},
					Right: &ast.Function{
						Name: "fn2",
						Args: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "negated grouped logical and expression",
			inputExpr: "!(fn1() && fn2())",
			want: &ast.LogicNot{
				Expression: &ast.Group{
					Expression: &ast.LogicAnd{
						Left: &ast.Function{
							Name: "fn1",
							Args: nil,
						},
						Right: &ast.Function{
							Name: "fn2",
							Args: nil,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "operator logic or with missing left node",
			inputExpr: "|| fn()",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple logical or expression",
			inputExpr: "fn1() || fn2()",
			want: &ast.LogicOr{
				Left: &ast.Function{
					Name: "fn1",
					Args: nil,
				},
				Right: &ast.Function{
					Name: "fn2",
					Args: nil,
				},
			},
			wantErr: false,
		},
		{
			name:      "grouped logical or expression",
			inputExpr: "(fn1() || fn2())",
			want: &ast.Group{
				Expression: &ast.LogicOr{
					Left: &ast.Function{
						Name: "fn1",
						Args: nil,
					},
					Right: &ast.Function{
						Name: "fn2",
						Args: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "simple logical or expression",
			inputExpr: "!(fn1() || fn2())",
			want: &ast.LogicNot{
				Expression: &ast.Group{
					Expression: &ast.LogicOr{
						Left: &ast.Function{
							Name: "fn1",
							Args: nil,
						},
						Right: &ast.Function{
							Name: "fn2",
							Args: nil,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "default logical operator precedence expression #1",
			inputExpr: "fn1() || fn2() && fn3()",
			want: &ast.LogicOr{
				Left: &ast.Function{
					Name: "fn1",
					Args: nil,
				},
				Right: &ast.LogicAnd{
					Left: &ast.Function{
						Name: "fn2",
						Args: nil,
					},
					Right: &ast.Function{
						Name: "fn3",
						Args: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "default logical operator precedence expression #2",
			inputExpr: "fn1() && fn2() || fn3()",
			want: &ast.LogicOr{
				Left: &ast.LogicAnd{
					Left: &ast.Function{
						Name: "fn1",
						Args: nil,
					},
					Right: &ast.Function{
						Name: "fn2",
						Args: nil,
					},
				},
				Right: &ast.Function{
					Name: "fn3",
					Args: nil,
				},
			},
			wantErr: false,
		},
		{
			name:      "group overrides default logical operator precedence",
			inputExpr: "fn1() && (fn2() || fn3())",
			want: &ast.LogicAnd{
				Left: &ast.Function{
					Name: "fn1",
					Args: nil,
				},
				Right: &ast.Group{
					Expression: &ast.LogicOr{
						Left: &ast.Function{
							Name: "fn2",
							Args: nil,
						},
						Right: &ast.Function{
							Name: "fn3",
							Args: nil,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := scanner.New(strings.NewReader(tt.inputExpr))
			got, err := Parse(s)
			require.Equal(t, tt.wantErr, err != nil, "Parse() error = %v, wantErr %v", err, tt.wantErr)
			require.True(t, reflect.DeepEqual(got, tt.want), "Parse() got = %v, want %v", spew.Sdump(got), spew.Sdump(tt.want))
		})
	}
}
