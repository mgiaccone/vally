package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_patchExprRegex(t *testing.T) {
	testExpr := "(eq(.OtherField, 'GB') && required(    )) || true() || eq('GB') && eq(1234)"
	testFieldRef := ".ReplacedRef"
	expectedExpr := "(eq(.OtherField, 'GB') && required(.ReplacedRef)) || true(.ReplacedRef) || eq(.ReplacedRef,'GB') && eq(.ReplacedRef,1234)"
	got := patchExprRegex(testExpr, testFieldRef)
	require.Equal(t, expectedExpr, got)
}

func Test_patchExprScanner(t *testing.T) {
	testExpr := "(eq(.OtherField, 'GB') && required(    )) || true(    ) || eq('GB') && eq(1234)"
	testFieldRef := ".ReplacedRef"
	expectedExpr := "(eq(.OtherField, 'GB') && required(.ReplacedRef)) || true(.ReplacedRef) || eq(.ReplacedRef,'GB') && eq(.ReplacedRef,1234)"
	got, err := patchExprScanner(testExpr, testFieldRef)
	require.NoError(t, err)
	require.Equal(t, expectedExpr, got)
}
