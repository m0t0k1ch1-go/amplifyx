package testutil

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

func Equal(want, got any, opts ...cmp.Option) error {
	if diff := cmp.Diff(want, got, opts...); len(diff) > 0 {
		return fmt.Errorf("must be equal to the expected value (-want +got):\n%s", diff)
	}

	return nil
}
