package validator

import (
	"testing"
)

var (
	_patchExprBenchString string
)

func Benchmark_patchExprRegex(b *testing.B) {
	b.ReportAllocs()

	testExpr := "(eq(.OtherField, 'GB') && required(    )) || true() || eq('GB') && eq(1234)"
	testFieldRef := ".ReplacedRef"
	var res string
	for n := 0; n < b.N; n++ {
		res = patchExprRegex(testExpr, testFieldRef)
	}
	_patchExprBenchString = res
}

func Benchmark_patchExprScanner(b *testing.B) {
	b.ReportAllocs()

	testExpr := "(eq(.OtherField, 'GB') && required(    )) || true() || eq('GB') && eq(1234)"
	testFieldRef := ".ReplacedRef"
	var res string
	for n := 0; n < b.N; n++ {
		res, _ = patchExprScanner(testExpr, testFieldRef)
	}
	_patchExprBenchString = res
}
