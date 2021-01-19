package vally

import (
	"github.com/osl4b/vally/internal/ast"
)

var (
	_ ast.Visitor = (*eval)(nil)
)

type evalContext struct {
}

type eval struct {
	r   *Registry
	ctx *evalContext
	res []bool
}

func (e *eval) push(v bool) {
	e.res = append(e.res, v)
}

func (e *eval) pop() bool {
	if len(e.res) == 0 {
		return false
	}
	var v bool
	i := len(e.res) - 1
	v, e.res = e.res[i-1], e.res[:i]
	return v
}

func (e *eval) Result() bool {
	return e.pop()
}

func (e *eval) VisitFunction(fn *ast.Function) {
	if fn.Name == "true" {
		e.push(true)
	}
	if fn.Name == "false" {
		e.push(false)
	}
}

func (e *eval) VisitGroup(_ *ast.Group) {
	// nop
}

func (e *eval) VisitLogicAnd(_ *ast.LogicAnd) {
	v1, v2 := e.pop(), e.pop()
	e.push(v1 && v2)
}

func (e *eval) VisitLogicOr(_ *ast.LogicOr) {
	v1, v2 := e.pop(), e.pop()
	e.push(v1 || v2)
}

func (e *eval) VisitLogicNot(_ *ast.LogicNot) {
	e.push(!e.pop())
}
