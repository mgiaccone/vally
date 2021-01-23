package ast

var (
	_ Node = (*Function)(nil)
	_ Node = (*Group)(nil)
	_ Node = (*LogicAnd)(nil)
	_ Node = (*LogicNot)(nil)
	_ Node = (*LogicOr)(nil)
)

type ArgType int

const (
	FieldRef ArgType = iota + 1
	FloatValue
	IntValue
	StringValue
)

// Visitor represents the AST tree traversal visitor.
type Visitor interface {
	VisitFunction(*Function) error
	VisitGroup(*Group) error
	VisitLogicAnd(*LogicAnd) error
	VisitLogicOr(*LogicOr) error
	VisitLogicNot(*LogicNot) error
}

// Node represents a generic AST node.
type Node interface {
	Visit(v Visitor) error
}

// FunctionArg represents an argument to a function.
type FunctionArg struct {
	Type        ArgType `json:"type"`
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
func (n *Function) Visit(v Visitor) error {
	return v.VisitFunction(n)
}

// Group represents an AST group.
type Group struct {
	Expression Node `json:"expression"`
}

// Visit implements the Node interface.
func (n *Group) Visit(v Visitor) (err error) {
	if err = n.Expression.Visit(v); err != nil {
		return err
	}
	return v.VisitGroup(n)
}

// LogicAnd represents an AST logical AND operator.
type LogicAnd struct {
	Left  Node `json:"left,omitempty"`
	Right Node `json:"right,omitempty"`
}

// Visit implements the Node interface.
func (n *LogicAnd) Visit(v Visitor) (err error) {
	if err = n.Left.Visit(v); err != nil {
		return err
	}
	if err = n.Right.Visit(v); err != nil {
		return err
	}
	return v.VisitLogicAnd(n)
}

// LogicNot represents an AST logical NOT operator.
type LogicNot struct {
	Expression Node `json:"expression"`
}

// Visit implements the Node interface.
func (n *LogicNot) Visit(v Visitor) (err error) {
	if err = n.Expression.Visit(v); err != nil {
		return err
	}
	return v.VisitLogicNot(n)
}

// LogicOr represents an AST logical OR operator.
type LogicOr struct {
	Left  Node `json:"left,omitempty"`
	Right Node `json:"right,omitempty"`
}

// Visit implements the Node interface.
func (n *LogicOr) Visit(v Visitor) (err error) {
	if err = n.Left.Visit(v); err != nil {
		return err
	}
	if err = n.Right.Visit(v); err != nil {
		return err
	}
	return v.VisitLogicOr(n)
}
