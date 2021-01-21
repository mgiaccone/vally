package parser

import (
	"fmt"
	"strconv"

	"github.com/osl4b/vally/internal/ast"
	"github.com/osl4b/vally/internal/scanner"
)

// Scanner represents the lexical scanner
type Scanner interface {
	Scan() scanner.Token
}

// Parse parses a lexical token stream and builds the AST for the given expression
func Parse(s Scanner) (ast.Node, error) {
	root, err := parseExpression(s, false)
	if err != nil {
		return nil, err
	}
	if root == nil {
		return nil, fmt.Errorf("empty expression")
	}
	return root, nil
}

// nolint:gocognit,gocyclo
func parseExpression(s Scanner, inGroup bool) (ast.Node, error) {
	var (
		node     ast.Node
		err      error
		breakout bool
	)
	for {
		tok := s.Scan()

		switch tok.Type {
		case scanner.EOF:
			breakout = true
		case scanner.ILLEGAL:
			return nil, fmt.Errorf(tok.Literal)
		case scanner.LNOT:
			node, err = parseLogicNot(s)
			if err != nil {
				return nil, fmt.Errorf("logic not: %w", err)
			}
		case scanner.LPAREN:
			node, err = parseGroup(s)
			if err != nil {
				return nil, fmt.Errorf("group: %w", err)
			}
			if inGroup {
				breakout = true
			}
		case scanner.RPAREN:
			if inGroup {
				breakout = true
			}
		case scanner.IDENT:
			node, err = parseFunction(s, tok.Literal)
			if err != nil {
				return nil, fmt.Errorf("expression: %w", err)
			}
		case scanner.LAND:
			if node == nil {
				return nil, fmt.Errorf("expression: unexpected token %s", tok.Type)
			}
			node, err = parseLogicAnd(s, node)
			if err != nil {
				return nil, fmt.Errorf("logic and: %w", err)
			}
		case scanner.LOR:
			if node == nil {
				return nil, fmt.Errorf("expression: unexpected token %s", tok.Type)
			}
			node, err = parseLogicOr(s, node)
			if err != nil {
				return nil, fmt.Errorf("logic or: %w", err)
			}
		default:
			return nil, fmt.Errorf("expression: unexpected token %s", tok.Type)
		}

		if breakout {
			break
		}
	}

	return node, nil
}

func parseFunction(s Scanner, name string) (ast.Node, error) {
	var (
		args     []ast.FunctionArg
		err      error
		breakout bool
	)
	for {
		tok := s.Scan()

		switch tok.Type {
		case scanner.ILLEGAL:
			return nil, fmt.Errorf(tok.Literal)
		case scanner.LPAREN:
			args, err = parseFunctionArgs(s)
			if err != nil {
				return nil, fmt.Errorf("function %q: %w", name, err)
			}
			breakout = true
		case scanner.RPAREN:
			breakout = true
		default:
			return nil, fmt.Errorf("function %q: unexpected token %s", name, tok.Type)
		}

		if breakout {
			break
		}
	}

	return &ast.Function{
		Name: name,
		Args: args,
	}, nil
}

// nolint:gocognit,gocyclo
func parseFunctionArgs(s Scanner) ([]ast.FunctionArg, error) {
	var (
		args       []ast.FunctionArg
		requireArg bool
		breakout   bool
		argCount   int
	)
	for {
		tok := s.Scan()

		switch tok.Type {
		case scanner.ILLEGAL:
			return nil, fmt.Errorf(tok.Literal)
		case scanner.COMMA:
			if argCount == 0 {
				return nil, fmt.Errorf("argument #%d: missing value", argCount)
			}
			if requireArg {
				return nil, fmt.Errorf("argument #%d: missing argument", argCount)
			}
			requireArg = true
		case scanner.FIELDREF:
			args = append(args, ast.FunctionArg{
				FieldRef: tok.Literal,
			})
			argCount++
			requireArg = false
		case scanner.STRING:
			args = append(args, ast.FunctionArg{
				StringValue: tok.Literal[1 : len(tok.Literal)-1],
			})
			argCount++
			requireArg = false
		case scanner.INT:
			v, err := strconv.ParseInt(tok.Literal, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("argument #%d: %w", argCount, err)
			}
			args = append(args, ast.FunctionArg{
				IntValue: v,
			})
			argCount++
			requireArg = false
		case scanner.FLOAT:
			v, err := strconv.ParseFloat(tok.Literal, 64)
			if err != nil {
				return nil, fmt.Errorf("argument #%d: %w", argCount, err)
			}
			args = append(args, ast.FunctionArg{
				FloatValue: v,
			})
			argCount++
			requireArg = false
		case scanner.RPAREN:
			if requireArg {
				return nil, fmt.Errorf("argument #%d: missing argument", argCount)
			}
			breakout = true
		default:
			return nil, fmt.Errorf("argument #%d: unexpected token %s", argCount, tok.Type)
		}

		if breakout {
			break
		}
	}

	return args, nil
}

func parseGroup(s Scanner) (ast.Node, error) {
	expr, err := parseExpression(s, true)
	if err != nil {
		return nil, err
	}

	return &ast.Group{
		Expression: expr,
	}, nil
}

func parseLogicNot(s Scanner) (ast.Node, error) {
	expr, err := parseExpression(s, false)
	if err != nil {
		return nil, err
	}

	return &ast.LogicNot{
		Expression: expr,
	}, nil
}

func parseLogicAnd(s Scanner, leftNode ast.Node) (ast.Node, error) {
	expr, err := parseExpression(s, false)
	if err != nil {
		return nil, err
	}

	// restore operator precedence when expr is an OR operatore
	if orExpr, ok := expr.(*ast.LogicOr); ok {
		return &ast.LogicOr{
			Left: &ast.LogicAnd{
				Left:  leftNode,
				Right: orExpr.Left,
			},
			Right: orExpr.Right,
		}, nil
	}

	return &ast.LogicAnd{
		Left:  leftNode,
		Right: expr,
	}, nil
}

func parseLogicOr(s Scanner, leftNode ast.Node) (ast.Node, error) {
	expr, err := parseExpression(s, false)
	if err != nil {
		return nil, err
	}

	return &ast.LogicOr{
		Left:  leftNode,
		Right: expr,
	}, nil
}
