package validator

var (
	_ Target = (*structTarget)(nil)
	_ Target = (*valueTarget)(nil)
)

type structTarget struct {
	s interface{}
}

func newStructTarget(s interface{}) *structTarget {
	return &structTarget{s: s}
}

type valueTarget struct {
	v interface{}
}

func newValueTarget(v interface{}) *valueTarget {
	return &valueTarget{v: v}
}
