package amplifyx

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/samber/oops"
)

type VersionArgs struct{}

func (app App) Version(ctx context.Context, args VersionArgs) error {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return oops.New("failed to read build info")
	}

	fmt.Fprintln(app.writer, info.Main.Version)

	return nil
}
