//go:build wireinject

package entrypoint

import (
	"context"

	"github.com/aereal/nikki/backend/env"
	"github.com/aereal/nikki/backend/infra/db"
	"github.com/aereal/nikki/backend/log"
	"github.com/aereal/nikki/backend/o11y"
	"github.com/aereal/nikki/backend/o11y/service"
	"github.com/aereal/nikki/backend/web"
	"github.com/google/wire"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func NewDevEntrypoint(_ context.Context) (*Entrypoint, error) {
	wire.Build(
		db.ProvideDB,
		env.ProvideDBEndpoint,
		env.ProvideLogLevel,
		env.ProvidePort,
		env.ProvideVariables,
		log.ProvideGlobalInstrumentation,
		log.ProvideLogger,
		log.ProvideStdout,
		o11y.ProvideResource,
		o11y.ProvideSidecarExporter,
		o11y.ProvideTracerProvider,
		provideDynamicServiceVersion,
		provideEntrypoint,
		web.ProvideServer,
		wire.Bind(new(trace.TracerProvider), new(*sdktrace.TracerProvider)),
		wire.Value(log.GoogleCloudProject("dummy")),
		wire.Value(service.Environment("local")),
	)
	return nil, nil
}
