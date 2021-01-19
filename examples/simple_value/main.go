package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/osl4b/vally"
)

func main() {
	exprs := []string{
		"false() || false()",
		"false() || true()",
		"true() || false()",
		"true() || true()",
		"false() && false()",
		"false() && true()",
		"true() && false()",
		"true() && true()",
	}

	var err error
	for _, rawExpr := range exprs {
		err = vally.ValidateValue(context.Background(), strings.NewReader(rawExpr), nil)
		fmt.Printf("%q => %v (%v)\n", rawExpr, err == nil, err)
	}
}
