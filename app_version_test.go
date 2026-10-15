package amplifyx_test

import (
	"bytes"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/amplifyx/v2"
	"github.com/m0t0k1ch1-go/amplifyx/v2/internal/testutil"
)

func TestApp_Version(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := t.Context()

		var buf bytes.Buffer

		app := amplifyx.NewApp(&buf, nil)

		err := app.Version(ctx, amplifyx.VersionArgs{})
		require.NoError(t, err)

		info, ok := debug.ReadBuildInfo()
		require.True(t, ok)
		testutil.RequireEqual(t, info.Main.Version+"\n", buf.String())
	})
}
