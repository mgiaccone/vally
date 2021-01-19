package ast

var (
	_ Node = (*Function)(nil)
	_ Node = (*Group)(nil)
	_ Node = (*LogicAnd)(nil)
	_ Node = (*LogicNot)(nil)
	_ Node = (*LogicOr)(nil)
)

// Visitor represents the AST tree traversal visitor.
type Visitor interface {
	VisitFunction(*Function)
	VisitGroup(*Group)
	VisitLogicAnd(*LogicAnd)
	VisitLogicOr(*LogicOr)
	VisitLogicNot(*LogicNot)
}

// Node represents a generic AST node.
type Node interface {
	Visit(v Visitor)
}

// FunctionArg represents an argument to a function.
type FunctionArg struct {
	FieldRef    string  `json:"fieldRef,omitempty"`
	FloatValue  float64 `json:"floatValue,omitempty"`
	IntValue    int64   `json:"intValue,omitempty"`
	StringValue string  `json:"stringValue,omitempty"`
}

// Function represents an AST function.
type Function struct {
	Name string        `json:"name"`
	Args []FunctionArg `json:"args"`
}

// Visit implements the Node interface.
func (n *Function) Visit(v Visitor) {
	v.VisitFunction(n)
}

// Group represents an AST group.
type Group struct {
	Expression Node `json:"expression"`
}

// Visit implements the Node interface.
func (n *Group) Visit(v Visitor) {
	n.Expression.Visit(v)
	v.VisitGroup(n)
}

// LogicAnd represents an AST logical AND operator.
type LogicAnd struct {
	Left  Node `json:"left,omitempty"`
	Right Node `json:"right,omitempty"`
}

// Visit implements the Node interface.
func (n *LogicAnd) Visit(v Visitor) {
	n.Left.Visit(v)
	n.Right.Visit(v)
	v.VisitLogicAnd(n)
}

// LogicNot represents an AST logical NOT operator.
type LogicNot struct {
	Expression Node `json:"expression"`
}

// Visit implements the Node interface.
func (n *LogicNot) Visit(v Visitor) {
	n.Expression.Visit(v)
	v.VisitLogicNot(n)
}

// LogicOr represents an AST logical OR operator.
type LogicOr struct {
	Left  Node `json:"left,omitempty"`
	Right Node `json:"right,omitempty"`
}

// Visit implements the Node interface.
func (n *LogicOr) Visit(v Visitor) {
	n.Left.Visit(v)
	n.Right.Visit(v)
	v.VisitLogicOr(n)
}
