package amplifyx

import (
	"io"
)

type App struct {
	writer        io.Writer
	amplifyClient AmplifyClient
}

func NewApp(
	writer io.Writer,
	amplifyClient AmplifyClient,
) *App {
	return &App{
		writer:        writer,
		amplifyClient: amplifyClient,
	}
}
