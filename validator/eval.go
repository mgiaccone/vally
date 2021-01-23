package validator

import (
	"context"
	"fmt"

	"github.com/osl4b/vally/internal/ast"
	"github.com/osl4b/vally/sdk"
)

// FIXME: How to implement && and || short-circuit evaluation?
//  leaving it unimplemented will lead to inconsistencies with error messages
//  i.e. v1() && v2() <- if v1 adds an error message, v2 should not be run or should
//  be prevented from adding another error message (ignore it?).
//  || might even be harder as the whole chain should be checked and
//  messages should be removed instead
//  plus cutting evaluation when possible will make it faster as well

var (
	_ ast.Visitor     = (*evalVisitor)(nil)
	_ sdk.EvalContext = (*evalContext)(nil)
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

func (ec *evalContext) TargetRef() string {
	return ec.rawFunc.Args[1].FieldRef
}

func (ec *evalContext) FunctionArgs() []sdk.ArgValue {
	var args []sdk.ArgValue
	for i, rawArg := range ec.rawFunc.Args {
		// skip first fieldref
		if i == 0 {
			// FIXME: might need to skip to to account for fieldRef (the field the expression belongs to )
			//  and funcTargetRef (the field the function needs to evaluate against)
			continue
		}

		switch rawArg.Type {
		case ast.FieldRef:
			args = append(args, rawArg.FieldRef)
		case ast.FloatValue:
			args = append(args, rawArg.FloatValue)
		case ast.IntValue:
			args = append(args, rawArg.IntValue)
		case ast.StringValue:
			args = append(args, rawArg.StringValue)
		}
	}
	return args
}

func (ec *evalContext) FunctionName() string {
	return ec.rawFunc.Name
}

type evalVisitor struct {
	ctx  context.Context
	t    sdk.Target
	v    *Validator
	res  []bool
	errs []error // FIXME: This might need to be treated as a stack as well
}

func newEvalVisitor(ctx context.Context, v *Validator, t sdk.Target) *evalVisitor {
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

	// if the function implements sdk.ArgTyper we can check
	// its argument types against the expected types
	if ft, ok := f.(sdk.ArgTyper); ok {
		_ = ft
		// TODO: check argument types against required types
	}

	ec, err := newEvalContext(fn)
	if err != nil {
		return fmt.Errorf("eval context %q: %w", fn.Name, err)
	}

	res, err := f.Evaluate(ev.ctx, ec, ev.t)
	if err != nil {
		// FIXME: this is probably not correct, a way to fix the &&/|| issue might be to simply
		//  return the error from here and check what to do with it in the VisitAnd/VisitOr? Maybe...
		//  using a stack for errors the same way we pop results could work even better...
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
