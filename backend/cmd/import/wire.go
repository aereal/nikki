//go:build wireinject

package main

import (
	"context"

	"github.com/aereal/nikki/backend/domain"
	"github.com/aereal/nikki/backend/env"
	"github.com/aereal/nikki/backend/infra/db"
	"github.com/aereal/nikki/backend/infra/db/exec"
	"github.com/aereal/nikki/backend/log"
	"github.com/aereal/nikki/backend/o11y"
	"github.com/aereal/nikki/backend/o11y/service"
	"github.com/aereal/nikki/backend/usecases"
	"github.com/aereal/nikki/backend/usecases/interactions"
	"github.com/aereal/nikki/backend/usecases/unitofwork"
	"github.com/google/wire"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func build(_ context.Context) (*app, error) {
	wire.Build(
		db.ArticleIDGeneratorProvider,
		db.ArticleRevisionIDGeneratorProvider,
		db.CategoryIDGeneratorProvider,
		db.ProvideArticleRepository,
		db.ProvideCategoryRepository,
		db.ProvideDB,
		env.ProvideDBEndpoint,
		env.ProvideLogLevel,
		env.ProvideMTExportFileName,
		env.ProvideVariables,
		exec.ProvideRunner,
		interactions.ProvideImportMTExport,
		log.ProvideGlobalInstrumentation,
		log.ProvideLogger,
		log.ProvideStdout,
		o11y.ProvideResource,
		o11y.ProvideSidecarExporter,
		o11y.ProvideTracerProvider,
		provideApp,
		wire.Bind(new(domain.ArticleRepository), new(*db.ArticleRepository)),
		wire.Bind(new(domain.CategoryRepository), new(*db.CategoryRepository)),
		wire.Bind(new(exec.Context), new(*exec.Runner)),
		wire.Bind(new(trace.TracerProvider), new(*sdktrace.TracerProvider)),
		wire.Bind(new(unitofwork.Runner), new(*exec.Runner)),
		wire.Bind(new(usecases.ImportMTExport), new(*interactions.ImportMTExport)),
		wire.Value(service.Environment("local")),
		wire.Value(service.Version("latest")),
	)
	return nil, nil
}
