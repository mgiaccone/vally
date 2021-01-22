package validator

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/osl4b/vally/internal/ast"
	"github.com/osl4b/vally/internal/errutil"
	"github.com/osl4b/vally/internal/hashutil"
	"github.com/osl4b/vally/internal/parser"
	"github.com/osl4b/vally/internal/reflectutil"
	"github.com/osl4b/vally/internal/scanner"
)

const (
	_defaultStructTag = "vally"
)

type structEntry struct {
	fields    []fieldEntry
	validExpr ast.Node
}

type fieldEntry struct {
	Alias    string
	Expr     string
	Name     string
	FieldRef string
}

// Function represents a validator function signature
type Function func(ctx context.Context, args []Arg, t Target) (bool, error)

type Arg struct {
	raw ast.FunctionArg
}

// Target
type Target interface {
}

// Validator
type Validator struct {
	funcs           map[string]Function
	structCache     map[string]structEntry
	structCacheLock sync.Mutex
	structTag       string
}

// NewValidator returns a new instance of a Validator.
func NewValidator(opts ...Option) *Validator {
	v := &Validator{
		structCache: make(map[string]structEntry),
		structTag:   _defaultStructTag,
	}
	for _, apply := range opts {
		apply(v)
	}
	if v.funcs == nil {
		v.funcs = defaultFunctions()
	}
	return v
}

// MustRegister adds a validator function with the given id.
// It will panic if a function with the same id is already registered.
func (v *Validator) MustRegister(id string, fn Function) {
	errutil.Must(v.Register(id, fn))
}

// Register adds a validator function with the given id.
// It will return an error if a function with the same id is already registered, nil otherwise.
func (v *Validator) Register(id string, fn Function) error {
	if _, exists := v.funcs[id]; exists {
		return fmt.Errorf("function %q is already registered", id)
	}
	v.funcs[id] = fn
	return nil
}

// Replace adds or replace the validator function with the given id.
func (v *Validator) Replace(id string, fn Function) {
	v.funcs[id] = fn
}

// MustRegisterStruct
func (v *Validator) MustRegisterStruct(s interface{}) {
	errutil.Must(v.RegisterStruct(s))
}

// RegisterStruct
func (v *Validator) RegisterStruct(s interface{}) error {
	if _, err := v.retrieveOrBuildStructEntry(s); err != nil {
		return err
	}
	return nil
}

// ValidateStruct
func (v *Validator) ValidateStruct(ctx context.Context, val interface{}) error {
	if reflectutil.IsNil(val) {
		return fmt.Errorf("value must not be nil")
	}

	entry, err := v.retrieveOrBuildStructEntry(val)
	if err != nil {
		return fmt.Errorf("struct: %w", err)
	}

	ev := newEvalVisitor(ctx, v, val)
	if err = entry.validExpr.Visit(ev); err != nil {
		return fmt.Errorf("evaluate expression: %w", err)
	}

	result := ev.Result()
	fmt.Println(result)

	return ev.Err()
}

// ValidateValue
func (v *Validator) ValidateValue(ctx context.Context, expr string, value interface{}) error {
	// TODO: expr should not have .FieldRef or they should be stripped off

	r := strings.NewReader(expr)
	parsedExpr, err := parser.Parse(scanner.New(r))
	if err != nil {
		return err
	}
	ev := newEvalVisitor(ctx, v, newValueTarget(value))
	if err := parsedExpr.Visit(ev); err != nil {
		return err
	}
	if !ev.Result() {
		return fmt.Errorf("not valid")
	}
	return nil
}

// retrieveOrBuildStructEntry
// nolint:godox
// TODO: simplify this monster function
func (v *Validator) retrieveOrBuildStructEntry(s interface{}) (*structEntry, error) {
	if reflectutil.IsNil(s) {
		return nil, fmt.Errorf("value must not be nil")
	}

	var (
		err        error
		fieldAlias string
		fieldExpr  string
		fl         reflect.StructField
		buf        strings.Builder
		tag        string
		tagParts   []string
	)
	sv := reflect.ValueOf(s)
	st := sv.Type()
	if reflectutil.IsPointer(s) {
		st = st.Elem()
	}

	// nolint:godox
	// TODO: handle substructures
	fieldEntries := make([]fieldEntry, 0, st.NumField())
	for i := 0; i < st.NumField(); i++ {
		fl = st.Field(i)
		fieldAlias = fl.Name
		tag = fl.Tag.Get(v.structTag)
		if tag == "" || tag == "-" {
			continue
		}

		tagParts = strings.Split(tag, ";")
		if len(tagParts) == 1 {
			fieldExpr = tagParts[0]
		} else if len(tagParts) == 2 {
			fieldAlias = tagParts[0]
			fieldExpr = tagParts[1]
		} else {
			return nil, fmt.Errorf("invalid tag on field %q", fl.Name)
		}

		fieldEntries = append(fieldEntries, fieldEntry{
			Alias:    fieldAlias,
			Expr:     fieldExpr,
			Name:     fl.Name,
			FieldRef: "." + fl.Name,
		})
		buf.WriteString(fieldExpr)
	}

	v.structCacheLock.Lock()
	defer v.structCacheLock.Unlock()

	cacheKey := hashutil.Sha1Hex(buf.String())
	cacheEntry, exists := v.structCache[cacheKey]
	if exists {
		return &cacheEntry, nil
	}

	// reset string buffer for re-use
	buf.Reset()

	var (
		fp        *fieldEntry
		sc        *scanner.Scanner
		validExpr ast.Node
	)
	for i := 0; i < len(fieldEntries); i++ {
		fp = &fieldEntries[i]

		// validate field expression, this is only to verify the expression is formally correct;
		// each expression will be concatenated to create a global struct expression afterwards
		sc = scanner.New(strings.NewReader(fp.Expr))
		if _, err = parser.Parse(sc); err != nil {
			return nil, fmt.Errorf("field %q: %w", fp.Name, err)
		}

		if buf.Len() != 0 {
			buf.WriteString("&&")
		}
		// fp.Expr = patchExprRegex(fp.Expr, fp.FieldRef)
		fp.Expr, err = patchExprScanner(fp.Expr, fp.FieldRef)
		if err != nil {
			return nil, fmt.Errorf("field %q: %w", fp.Name, err)
		}

		buf.WriteString("(")
		buf.WriteString(fp.Expr)
		buf.WriteString(")")
	}

	sc = scanner.New(strings.NewReader(buf.String()))
	validExpr, err = parser.Parse(sc)
	if err != nil {
		return nil, fmt.Errorf("struct %q: %w", st.Name(), err)
	}

	e := structEntry{
		fields:    fieldEntries,
		validExpr: validExpr,
	}
	v.structCache[cacheKey] = e

	return &e, nil
}

// lookupFunction returns the validator function with the given id or an error if the function is not registered.
func (v *Validator) lookupFunction(id string) (Function, error) {
	if fn, exists := v.funcs[id]; exists {
		return fn, nil
	}
	return nil, ErrNotFound
}

func mapStructValues(s interface{}, prefix string) map[string]interface{} {
	return nil
}
