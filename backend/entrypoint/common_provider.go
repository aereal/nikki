package entrypoint

import (
	"github.com/aereal/nikki/backend/env"
	"github.com/aereal/nikki/backend/graph"
	"github.com/aereal/nikki/backend/graph/resolvers"
	"github.com/aereal/nikki/backend/infra/db"
	"github.com/aereal/nikki/backend/log"
	"github.com/aereal/nikki/backend/o11y"
	"github.com/aereal/nikki/backend/web"
	"github.com/google/wire"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var commonProvider = wire.NewSet(
	db.ProvideDB,
	env.ProvideDBEndpoint,
	env.ProvideLogLevel,
	env.ProvidePort,
	env.ProvideVariables,
	graph.ProviveHandler,
	log.ProvideGlobalInstrumentation,
	log.ProvideStdout,
	o11y.ProvideTracerProvider,
	provideEntrypoint,
	resolvers.ProvideResolver,
	web.ProvideServer,
	wire.Bind(new(trace.TracerProvider), new(*sdktrace.TracerProvider)),
)
