package amplifyx

import (
	"time"

	"github.com/samber/oops"
)

type Args struct {
	Timeout time.Duration `default:"10m"`

	Deploy  DeployArgs  `cmd:""`
	Version VersionArgs `cmd:""`
}

func (args Args) Validate() error {
	if args.Timeout < 0 {
		return oops.New("invalid timeout: must not be negative")
	}

	return nil
}
