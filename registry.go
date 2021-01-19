package vally

import (
	"fmt"

	"github.com/osl4b/vally/builtin"
	"github.com/osl4b/vally/internal/errutil"
)

var (
	_defaultRegistry = NewRegistry()
)

func init() {
	_defaultRegistry.MustRegister("false", builtin.False)
	_defaultRegistry.MustRegister("true", builtin.True)
	_defaultRegistry.MustRegister("required", builtin.Required)
	_defaultRegistry.MustRegister("email", builtin.Email)
}

// NewRegistry returns a new instance of a Registry.
func NewRegistry() *Registry {
	return &Registry{
		funcs: make(map[string]Func),
	}
}

// Registry manages a collection of validator functions.
//
// The registry is not safe for concurrent read/write operations. All custom functions
// must be registered before being assigned to a Validator.
type Registry struct {
	funcs map[string]Func
}

// Copy returns a current copy of the registry.
func (r *Registry) Copy() *Registry {
	funcs := make(map[string]Func, len(r.funcs))
	for k, v := range r.funcs {
		funcs[k] = v
	}
	return &Registry{funcs: funcs}
}

// Get returns the validator function with the given id or an error if the function is not registered.
func (r *Registry) Get(id string) (Func, error) {
	fn, exists := r.funcs[id]
	if !exists {
		return nil, fmt.Errorf("validator function %q is already registered", id)
	}
	return fn, nil
}

// Register adds a validator function with the given id.
// It will return an error if a function with the same id is already registered, nil otherwise.
func (r *Registry) Register(id string, fn Func) error {
	_, exists := r.funcs[id]
	if exists {
		return fmt.Errorf("function %q is already registered", id)
	}
	r.funcs[id] = fn
	return nil
}

// MustRegister adds a validator function with the given id.
// It will panic if a function with the same id is already registered.
func (r *Registry) MustRegister(id string, fn Func) {
	errutil.Must(r.Register(id, fn))
}

// Replace adds or replace the validator function with the given id.
func (r *Registry) Replace(id string, fn Func) {
	r.funcs[id] = fn
}
