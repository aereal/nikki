//go:build wireinject

package entrypoint

import (
	"context"

	"github.com/aereal/nikki/backend/env"
	"github.com/aereal/nikki/backend/log"
	"github.com/aereal/nikki/backend/o11y"
	"github.com/aereal/nikki/backend/web"
	"github.com/google/wire"
)

func NewDevEntrypoint(_ context.Context) (*Entrypoint, error) {
	wire.Build(
		env.ProvideLogLevel,
		env.ProvidePort,
		env.ProvideVariables,
		log.ProvideGlobalInstrumentation,
		log.ProvideLogger,
		log.ProvideStdout,
		o11y.ProvideTracerProvider,
		provideDynamicServiceVersion,
		provideEntrypoint,
		web.ProvideServer,
		wire.Value(log.GoogleCloudProject("dummy")),
	)
	return nil, nil
}
