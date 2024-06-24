package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsZero(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    bool
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsZero(tt.value)
			require.Equal(t, tt.wantErr, err != nil, "IsZero() error = %v, wantErr %v", err, tt.wantErr)
			require.Equal(t, tt.want, got, "IsZero() got = %v, want %v", got, tt.want)
		})
	}
}
