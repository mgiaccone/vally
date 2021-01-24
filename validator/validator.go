package validator

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"unicode"

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

var (
	ErrNotFound = errors.New("not found")
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

// Validator contains the implementation of the validation logic.
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
		v.funcs = builtInFunctions()
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

// ValidateStruct processes a struct value and applies the validation
// expressions as given in each property tag.
func (v *Validator) ValidateStruct(ctx context.Context, target interface{}) error {
	if reflectutil.IsNil(target) {
		return fmt.Errorf("target must not be nil")
	}

	entry, err := v.retrieveOrBuildStructEntry(target)
	if err != nil {
		return err
	}

	t, err := newStructTarget(target)
	if err != nil {
		return err
	}

	ev := newEvalVisitor(ctx, v, t)
	if err = entry.validExpr.Visit(ev); err != nil {
		return err
	}

	return ev.Err()
}

// ValidateValue applies the validation expression to the given value.
//
// It returns an error of type validator.ValidationError validation fails or a generic error
// if other issues are detected.
func (v *Validator) ValidateValue(ctx context.Context, expr string, value interface{}) error {
	expr, err := patchExprScanner(expr, ".Value")
	if err != nil {
		return err
	}
	r := strings.NewReader(expr)

	parsedExpr, err := parser.Parse(scanner.New(r))
	if err != nil {
		return err
	}

	t, err := newValueTarget(value)
	if err != nil {
		return fmt.Errorf("target: %w", err)
	}

	ev := newEvalVisitor(ctx, v, t)
	if err := parsedExpr.Visit(ev); err != nil {
		return err
	}

	return ev.Err()
}

// retrieveOrBuildStructEntry
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

// TODO: review docs
// patchExprScanner rewrites the expression adding the given fieldRef and targetRef to any function
// that doesn't explicitly declare a targetRef. For instance the given the fieldRef ".SomeField"
// and the expression "require()", the output expression will become "require(.SomeField,.SomeField)".
//
// It uses a scanner to perform the job, this implementation is about 5x faster
// than the same function implemented with regular expressions.
//
func patchExprScanner(expr, fieldRef string) (string, error) {
	var (
		buf strings.Builder
		ch  rune
		err error
	)
	br := bufio.NewReaderSize(strings.NewReader(expr), 256)
	for {
		ch, _, err = br.ReadRune()
		if errors.Is(err, io.EOF) {
			break
		}

		// start ident
		if unicode.IsLower(ch) {
			buf.WriteRune(ch)

			if err = consumeTo(br, &buf, '('); err == io.EOF {
				return "", fmt.Errorf("consume ident: unexpected EOF")
			}

			// in argument list from here
			ch, err = consumePrefixArg(br)
			if err == io.EOF {
				return "", fmt.Errorf("consume args: unexpected EOF")
			}

			// end of arg body
			if ch == ')' {
				buf.WriteString(fieldRef)
				buf.WriteRune(',')
				buf.WriteString(fieldRef)
				buf.WriteRune(ch)
				continue
			}

			// existing field ref, consume and move on
			if ch == '.' {
				buf.WriteString(fieldRef)
				buf.WriteRune(',')
				buf.WriteRune(ch)
			}

			// variable declaration
			if ch == '\'' || unicode.IsDigit(ch) {
				buf.WriteString(fieldRef)
				buf.WriteRune(',')
				buf.WriteString(fieldRef)
				buf.WriteRune(',')
				buf.WriteRune(ch)
			}

			if err = consumeTo(br, &buf, ')'); err == io.EOF {
				return "", fmt.Errorf("consume ref args: unexpected EOF")
			}
		}

		if ch == '(' || ch == ')' || ch == ' ' || ch == '\t' || ch == '&' || ch == '|' {
			buf.WriteRune(ch)
			continue
		}
	}

	return buf.String(), nil
}

// consumeTo consumes the reader up to and including the first instance of the given stopCh rune is found.
func consumeTo(in *bufio.Reader, out *strings.Builder, stopCh rune) error {
	var (
		ch  rune
		err error
	)
	for {
		// consume ident up to open left parenthesis
		ch, _, err = in.ReadRune()
		if err == io.EOF {
			return err
		}
		out.WriteRune(ch)
		if ch == stopCh {
			break
		}
	}
	return nil
}

// consumePrefixArg consumes the reader from the beginning of the argument body to either the end of the arg
// body or to a valid argument. It is used when rewriting the expression to remove trailing whitespaces.
func consumePrefixArg(in *bufio.Reader) (rune, error) {
	var (
		ch  rune
		err error
	)
	for {
		ch, _, err = in.ReadRune()
		if err == io.EOF {
			return 0, nil
		}
		if ch == ')' || ch == '.' || ch == ',' || ch == '\'' || unicode.IsDigit(ch) {
			return ch, nil
		}
	}
}
