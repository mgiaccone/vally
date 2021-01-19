package reflectutil_test

import (
	"testing"

	"github.com/osl4b/vally/internal/reflectutil"
)

type testStruct struct{}

func TestStructKey(t *testing.T) {
	type funcTestStruct struct{}

	t.Log(reflectutil.StructKey(testStruct{}))
	t.Log(reflectutil.StructKey(&testStruct{}))
	t.Log(reflectutil.StructKey(funcTestStruct{}))
	t.Log(reflectutil.StructKey(&funcTestStruct{}))

	// tests := []struct {
	// 	name string
	// 	args args
	// 	want string
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		if got := StructKey(tt.args.s); got != tt.want {
	// 			t.Errorf("StructKey() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}
