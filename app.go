package amplifyx

import (
	"io"

	awsamplify "github.com/aws/aws-sdk-go-v2/service/amplify"
)

// App represents the amplifyx CLI application.
type App struct {
	writer        io.Writer
	amplifyClient *awsamplify.Client
}

// NewApp returns a new [App] that writes its output to writer and calls AWS Amplify through amplifyClient.
func NewApp(
	writer io.Writer,
	amplifyClient *awsamplify.Client,
) *App {
	return &App{
		writer:        writer,
		amplifyClient: amplifyClient,
	}
}
