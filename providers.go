package amplifyx

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awsamplify "github.com/aws/aws-sdk-go-v2/service/amplify"
	"github.com/goforj/wire"
	"github.com/samber/oops"
)

var (
	ProviderSet = wire.NewSet(
		provideAWSConfig,
		provideAWSAmplifyClient,

		provideApp,
	)
)

func provideAWSConfig(ctx context.Context) (aws.Config, error) {
	awsConf, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Config{}, oops.Wrapf(err, "failed to load aws config")
	}

	return awsConf, nil
}

func provideAWSAmplifyClient(awsConf aws.Config) *awsamplify.Client {
	return awsamplify.NewFromConfig(awsConf)
}

func provideApp(amplifyClient *awsamplify.Client) *App {
	return NewApp(
		os.Stdout,
		amplifyClient,
	)
}
