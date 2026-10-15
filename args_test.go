package amplifyx_test

import (
	"testing"
	"time"

	"github.com/alecthomas/kong"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/amplifyx/v2"
	"github.com/m0t0k1ch1-go/amplifyx/v2/internal/testutil"
)

func TestArgs_Validate(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []string
			want string
		}{
			{
				"with a negative timeout",
				[]string{"version", "--timeout=-1s"},
				"invalid timeout: must not be negative",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var args amplifyx.Args

				k, err := kong.New(&args)
				require.NoError(t, err)

				_, err = k.Parse(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []string
			want amplifyx.Args
		}{
			{
				"with the default timeout",
				[]string{"version"},
				amplifyx.Args{
					Timeout: 10 * time.Minute,
				},
			},
			{
				"with a zero timeout",
				[]string{"version", "--timeout=0"},
				amplifyx.Args{
					Timeout: 0,
				},
			},
			{
				"with a positive timeout",
				[]string{"version", "--timeout=1m"},
				amplifyx.Args{
					Timeout: time.Minute,
				},
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var args amplifyx.Args

				k, err := kong.New(&args)
				require.NoError(t, err)

				_, err = k.Parse(tc.in)
				require.NoError(t, err)
				testutil.RequireEqual(t, tc.want, args, cmpopts.IgnoreFields(amplifyx.Args{}, "Deploy", "Version"))
			})
		}
	})
}
