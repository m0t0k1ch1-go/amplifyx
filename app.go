package amplifyx

import (
	"context"
	"io"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awsamplify "github.com/aws/aws-sdk-go-v2/service/amplify"
	"github.com/samber/oops"
)

type App struct {
	writer        io.Writer
	amplifyClient *awsamplify.Client
}

func NewApp(ctx context.Context) (*App, error) {
	awsConf, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, oops.Wrapf(err, "failed to load aws config")
	}

	return &App{
		writer:        os.Stdout,
		amplifyClient: awsamplify.NewFromConfig(awsConf),
	}, nil
}
