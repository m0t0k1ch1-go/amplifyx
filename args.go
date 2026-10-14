package amplifyx

import (
	"time"

	"github.com/samber/oops"
)

// Args represents the command line arguments of the amplifyx CLI.
type Args struct {
	Timeout time.Duration `default:"10m"`

	Deploy  DeployArgs  `cmd:""`
	Version VersionArgs `cmd:""`
}

// Validate returns an error if args is invalid.
func (args Args) Validate() error {
	if args.Timeout < 0 {
		return oops.New("invalid timeout: must not be negative")
	}

	return nil
}
