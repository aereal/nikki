//go:build wireinject

package entrypoint

import (
	"context"

	"github.com/aereal/nikki/backend/adapters/gcp/metadata"
	"github.com/aereal/nikki/backend/env"
	"github.com/aereal/nikki/backend/graph"
	"github.com/aereal/nikki/backend/graph/resolvers"
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
		graph.ProviveHandler,
		log.ProvideGlobalInstrumentation,
		log.ProvideLogger,
		log.ProvideStdout,
		o11y.ProvideResource,
		o11y.ProvideSidecarExporter,
		o11y.ProvideTracerProvider,
		provideDynamicServiceVersion,
		provideEntrypoint,
		resolvers.ProvideResolver,
		web.ProvideServer,
		wire.Bind(new(trace.TracerProvider), new(*sdktrace.TracerProvider)),
		wire.Value(metadata.Project("dummy")),
		wire.Value(service.Environment("local")),
	)
	return nil, nil
}

func NewProductionEntrypoint(_ context.Context) (*Entrypoint, error) {
	wire.Build(
		db.ProvideDB,
		env.ProvideDBEndpoint,
		env.ProvideGoogleCloudProject,
		env.ProvideLogLevel,
		env.ProvidePort,
		env.ProvideServiceVersion,
		env.ProvideVariables,
		graph.ProviveHandler,
		log.ProvideGlobalInstrumentation,
		log.ProvideLogger,
		log.ProvideStdout,
		o11y.ProvideGoogleCloudRunResource,
		o11y.ProvideGoogleTelemetryTraceExporter,
		o11y.ProvideTracerProvider,
		provideEntrypoint,
		resolvers.ProvideResolver,
		web.ProvideServer,
		wire.Bind(new(trace.TracerProvider), new(*sdktrace.TracerProvider)),
		wire.Value(service.Environment("production")),
	)
	return nil, nil
}
