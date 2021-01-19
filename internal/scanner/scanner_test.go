package scanner_test

import (
	"strings"
	"testing"

	"github.com/osl4b/vally/internal/scanner"
	"github.com/stretchr/testify/require"
)

func TestScanner_Scan(t *testing.T) {
	tests := []struct {
		name      string
		inputExpr string
		want      []scanner.TokenType
	}{
		{
			name:      "empty input",
			inputExpr: "",
			want: []scanner.TokenType{
				scanner.EOF,
			},
		},
		{
			name:      "spaces and tabs only",
			inputExpr: "   \t",
			want: []scanner.TokenType{
				scanner.EOF,
			},
		},
		{
			name:      "illegal new lines",
			inputExpr: "   \n\r",
			want: []scanner.TokenType{
				scanner.ILLEGAL,
				scanner.EOF,
			},
		},
		{
			name:      "illegal input",
			inputExpr: "$",
			want: []scanner.TokenType{
				scanner.ILLEGAL,
				scanner.EOF,
			},
		},
		{
			name:      "terminators",
			inputExpr: "(),!",
			want: []scanner.TokenType{
				scanner.LPAREN,
				scanner.RPAREN,
				scanner.COMMA,
				scanner.LNOT,
				scanner.EOF,
			},
		},
		{
			name:      "terminators with spaces",
			inputExpr: "   ( )   ,       !          ",
			want: []scanner.TokenType{
				scanner.LPAREN,
				scanner.RPAREN,
				scanner.COMMA,
				scanner.LNOT,
				scanner.EOF,
			},
		},
		{
			name:      "logical AND operator",
			inputExpr: "&&",
			want: []scanner.TokenType{
				scanner.LAND,
				scanner.EOF,
			},
		},
		{
			name:      "invalid logical AND operator",
			inputExpr: "&",
			want: []scanner.TokenType{
				scanner.ILLEGAL,
				scanner.EOF,
			},
		},
		{
			name:      "logical OR operator",
			inputExpr: "||",
			want: []scanner.TokenType{
				scanner.LOR,
				scanner.EOF,
			},
		},
		{
			name:      "invalid logical OR operator",
			inputExpr: "|",
			want: []scanner.TokenType{
				scanner.ILLEGAL,
				scanner.EOF,
			},
		},
		{
			name:      "valid string",
			inputExpr: "'01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ _-+!@£$%^&*()[]{}<>:;\\|?/,.`~'",
			want: []scanner.TokenType{
				scanner.STRING,
				scanner.EOF,
			},
		},
		{
			name:      "unterminated string",
			inputExpr: "'unterminated",
			want: []scanner.TokenType{
				scanner.ILLEGAL,
				scanner.EOF,
			},
		},
		{
			name:      "valid identifier",
			inputExpr: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890",
			want: []scanner.TokenType{
				scanner.IDENT,
				scanner.EOF,
			},
		},
		{
			name:      "invalid identifier",
			inputExpr: "validUpToHere%",
			want: []scanner.TokenType{
				scanner.ILLEGAL,
				scanner.EOF,
			},
		},
		{
			name:      "valid positive integer",
			inputExpr: "123456",
			want: []scanner.TokenType{
				scanner.INT,
				scanner.EOF,
			},
		},
		{
			name:      "valid positive integer with sign",
			inputExpr: "+12345",
			want: []scanner.TokenType{
				scanner.INT,
				scanner.EOF,
			},
		},
		{
			name:      "valid negative integer",
			inputExpr: "-56789",
			want: []scanner.TokenType{
				scanner.INT,
				scanner.EOF,
			},
		},
		{
			name:      "valid 0 integer",
			inputExpr: "0",
			want: []scanner.TokenType{
				scanner.INT,
				scanner.EOF,
			},
		},
		{
			name:      "invalid integer",
			inputExpr: "012345",
			want: []scanner.TokenType{
				scanner.ILLEGAL,
				scanner.EOF,
			},
		},
		{
			name:      "valid positive float",
			inputExpr: "1.23456",
			want: []scanner.TokenType{
				scanner.FLOAT,
				scanner.EOF,
			},
		},
		{
			name:      "valid positive float with sign",
			inputExpr: "+0.12345",
			want: []scanner.TokenType{
				scanner.FLOAT,
				scanner.EOF,
			},
		},
		{
			name:      "valid negative float",
			inputExpr: "-6435.56789",
			want: []scanner.TokenType{
				scanner.FLOAT,
				scanner.EOF,
			},
		},
		{
			name:      "valid 0.0 float",
			inputExpr: "0.0",
			want: []scanner.TokenType{
				scanner.FLOAT,
				scanner.EOF,
			},
		},
		{
			name:      "simple function with single string argumant",
			inputExpr: "fn('string')",
			want: []scanner.TokenType{
				scanner.IDENT,
				scanner.LPAREN,
				scanner.STRING,
				scanner.RPAREN,
				scanner.EOF,
			},
		},
		{
			name:      "simple function with multiple arguments",
			inputExpr: "fn(.FieldName,'string',-12345,0.1234)",
			want: []scanner.TokenType{
				scanner.IDENT,
				scanner.LPAREN,
				scanner.FIELDREF,
				scanner.COMMA,
				scanner.STRING,
				scanner.COMMA,
				scanner.INT,
				scanner.COMMA,
				scanner.FLOAT,
				scanner.RPAREN,
				scanner.EOF,
			},
		},
		{
			name:      "complex expression",
			inputExpr: "fn1() && (!fn2('string') || fn3(12345))",
			want: []scanner.TokenType{
				scanner.IDENT,
				scanner.LPAREN,
				scanner.RPAREN,
				scanner.LAND,
				scanner.LPAREN,
				scanner.LNOT,
				scanner.IDENT,
				scanner.LPAREN,
				scanner.STRING,
				scanner.RPAREN,
				scanner.LOR,
				scanner.IDENT,
				scanner.LPAREN,
				scanner.INT,
				scanner.RPAREN,
				scanner.RPAREN,
				scanner.EOF,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := scanner.New(strings.NewReader(tt.inputExpr))
			got := make([]scanner.TokenType, 0, len(tt.want))
			for {
				tok := s.Scan()
				got = append(got, tok.Type)
				if tok.Type == scanner.EOF {
					break
				}
				t.Log(tok)
			}
			require.EqualValues(t, tt.want, got, "Scan() = %v, want %v", got, tt.want)
		})
	}
}
