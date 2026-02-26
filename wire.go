//go:build wireinject

package amplifyx

import (
	"context"

	"github.com/goforj/wire"
)

func InitializeApp(ctx context.Context) (*App, error) {
	wire.Build(ProviderSet)

	return nil, nil
}
