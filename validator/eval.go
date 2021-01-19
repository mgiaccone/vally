package validator

import (
	"context"
	"fmt"

	"github.com/osl4b/vally/internal/ast"
)

var (
	_ ast.Visitor = (*evalVisitor)(nil)
)

func newEvalVisitor(ctx context.Context, v *Validator, t Target) *evalVisitor {
	return &evalVisitor{
		ctx: ctx,
		v:   v,
		t:   t,
	}
}

type evalVisitor struct {
	ctx    context.Context
	t      Target
	v      *Validator
	result []bool
}

func (e *evalVisitor) push(v bool) {
	e.result = append(e.result, v)
}

func (e *evalVisitor) pop() bool {
	if len(e.result) == 0 {
		return false
	}
	var v bool
	i := len(e.result) - 1
	v = e.result[i]
	e.result = e.result[:i]
	return v
}

func (e *evalVisitor) Result() bool {
	return e.pop()
}

func (e *evalVisitor) VisitFunction(fn *ast.Function) error {
	f, err := e.v.lookupFunction(fn.Name)
	if err != nil {
		return fmt.Errorf("evaluate %q: %w", fn.Name, err)
	}

	// FIXME: add evaluation target here
	res, err := f(e.ctx, nil)
	if err != nil {
		return fmt.Errorf("evaluate %q: %w", fn.Name, err)
	}

	e.push(res)
	return nil
}

func (e *evalVisitor) VisitGroup(_ *ast.Group) error {
	return nil
}

func (e *evalVisitor) VisitLogicAnd(_ *ast.LogicAnd) error {
	v1, v2 := e.pop(), e.pop()
	e.push(v1 && v2)
	return nil
}

func (e *evalVisitor) VisitLogicOr(_ *ast.LogicOr) error {
	v1, v2 := e.pop(), e.pop()
	e.push(v1 || v2)
	return nil
}

func (e *evalVisitor) VisitLogicNot(_ *ast.LogicNot) error {
	v := e.pop()
	e.push(!v)
	return nil
}
