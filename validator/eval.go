package validator

import (
	"context"
	"fmt"

	"github.com/osl4b/vally/internal/ast"
)

var (
	_ ast.Visitor = (*evalVisitor)(nil)
	_ EvalContext = (*evalContext)(nil)
)

type evalContext struct {
	rawFunc *ast.Function
}

func newEvalContext(fn *ast.Function) (*evalContext, error) {
	ec := evalContext{
		rawFunc: fn,
	}
	if len(fn.Args) == 0 || fn.Args[0].FieldRef == "" {
		return nil, fmt.Errorf("eval context requires a .FieldRef")
	}

	return &ec, nil
}

func (ec *evalContext) FieldRef() string {
	return ec.rawFunc.Args[0].FieldRef
}

func (ec *evalContext) FunctionName() string {
	return ec.rawFunc.Name
}

type evalVisitor struct {
	ctx  context.Context
	t    Target
	v    *Validator
	res  []bool
	errs []error
}

func newEvalVisitor(ctx context.Context, v *Validator, t Target) *evalVisitor {
	return &evalVisitor{
		ctx: ctx,
		v:   v,
		t:   t,
	}
}

func (ev *evalVisitor) Err() error {
	if len(ev.errs) > 0 {
		return &Error{FieldErrs: ev.errs}
	}
	return nil
}

func (ev *evalVisitor) Result() bool {
	return ev.res[0]
}

func (ev *evalVisitor) VisitFunction(fn *ast.Function) error {
	f, err := ev.v.lookupFunction(fn.Name)
	if err != nil {
		return fmt.Errorf("eval lookup %q: %w", fn.Name, err)
	}

	ec, err := newEvalContext(fn)
	if err != nil {
		return fmt.Errorf("eval context %q: %w", fn.Name, err)
	}

	// TODO: check argument types against required types

	res, err := f.Evaluate(ev.ctx, ec, ev.t)
	if err != nil {
		fe, ok := err.(*FieldError)
		if !ok {
			return fmt.Errorf("eval execute %q: %w", fn.Name, err)
		}
		ev.errs = append(ev.errs, fe)
	}

	ev.push(res)
	return nil
}

func (ev *evalVisitor) VisitGroup(_ *ast.Group) error {
	return nil
}

func (ev *evalVisitor) VisitLogicAnd(_ *ast.LogicAnd) error {
	v1, v2 := ev.pop(), ev.pop()
	ev.push(v1 && v2)
	return nil
}

func (ev *evalVisitor) VisitLogicOr(_ *ast.LogicOr) error {
	v1, v2 := ev.pop(), ev.pop()
	ev.push(v1 || v2)
	return nil
}

func (ev *evalVisitor) VisitLogicNot(_ *ast.LogicNot) error {
	v := ev.pop()
	ev.push(!v)
	return nil
}

func (ev *evalVisitor) push(v bool) {
	ev.res = append(ev.res, v)
}

func (ev *evalVisitor) pop() bool {
	if len(ev.res) == 0 {
		return false
	}
	var v bool
	i := len(ev.res) - 1
	v, ev.res = ev.res[i], ev.res[:i]
	return v
}
