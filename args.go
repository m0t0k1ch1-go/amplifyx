package amplifyx

import (
	"time"
)

type Args struct {
	Timeout time.Duration `default:"10m"`

	Deploy  DeployArgs  `cmd:""`
	Version VersionArgs `cmd:""`
}
