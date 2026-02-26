package amplifyx

import (
	"io"

	awsamplify "github.com/aws/aws-sdk-go-v2/service/amplify"
)

type App struct {
	writer        io.Writer
	amplifyClient *awsamplify.Client
}

func NewApp(
	writer io.Writer,
	amplifyClient *awsamplify.Client,
) *App {
	return &App{
		writer:        writer,
		amplifyClient: amplifyClient,
	}
}
