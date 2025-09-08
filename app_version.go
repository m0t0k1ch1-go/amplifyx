package amplifyx

import (
	"context"
	"fmt"
	"runtime/debug"
)

type VersionArgs struct{}

func (app *App) Version(ctx context.Context, args VersionArgs) error {
	var version string
	{
		info, ok := debug.ReadBuildInfo()
		if !ok {
			version = "unknown"
		} else {
			version = info.Main.Version
		}
	}

	fmt.Fprintln(app.stdout, version)

	return nil
}
