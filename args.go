package amplifyx

import (
	"time"
)

type Args struct {
	Timeout time.Duration `default:"5m"`

	Deploy  DeployArgs  `cmd:""`
	Version VersionArgs `cmd:""`
}
