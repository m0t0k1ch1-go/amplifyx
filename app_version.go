package amplifyx

import (
	"context"
	"fmt"
	"runtime/debug"
)

// VersionArgs represents the arguments of the version command.
type VersionArgs struct{}

// Version writes the version of the running binary to the writer of app.
func (app App) Version(ctx context.Context, args VersionArgs) error {
	var version string
	{
		info, ok := debug.ReadBuildInfo()
		if !ok {
			version = "unknown"
		} else {
			version = info.Main.Version
		}
	}

	fmt.Fprintln(app.writer, version)

	return nil
}
