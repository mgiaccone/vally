package main

import (
	"context"
	"fmt"

	"github.com/osl4b/vally"
)

type SampleStruct struct {
	Email string `vally:"email;email(target=.OtherField,strict=true)"`
	Other string `vally:"required()"`
}

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

	val := &SampleStruct{
		Email: "test@example.com",
		Other: "",
	}

	var err error
	for _, rawExpr := range exprs {
		err = vally.ValidateStruct(context.Background(), &val)
		fmt.Printf("%q => %v (%v)\n", rawExpr, err == nil, err)
	}
}
