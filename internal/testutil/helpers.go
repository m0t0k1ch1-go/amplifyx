package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func RequireEqual(t *testing.T, want, got any, opts ...cmp.Option) {
	t.Helper()

	err := Equal(want, got, opts...)
	require.NoError(t, err)
}
