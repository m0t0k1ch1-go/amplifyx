//go:build wireinject

package amplifyx

import (
	"context"

	"github.com/goforj/wire"
)

// InitializeApp returns a new [App] whose AWS configuration is loaded from the environment.
func InitializeApp(ctx context.Context) (*App, error) {
	wire.Build(ProviderSet)

	return nil, nil
}
