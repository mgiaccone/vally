package vally

import (
	"sync"

	"github.com/osl4b/vally/internal/ast"
)

type cache struct {
	cache     map[string]ast.Node
	cacheLock sync.RWMutex
}

func (c *cache) Put() error {
	return nil
}
