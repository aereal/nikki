//go:build wireinject

package entrypoint

import (
	"context"

	"github.com/aereal/nikki/backend/env"
	"github.com/aereal/nikki/backend/log"
	"github.com/aereal/nikki/backend/o11y"
	"github.com/aereal/nikki/backend/o11y/service"
	"github.com/google/wire"
)

func NewDevEntrypoint(_ context.Context) (*Entrypoint, error) {
	wire.Build(
		commonProvider,
		log.ProvideLogger,
		o11y.ProvideResource,
		o11y.ProvideSidecarExporter,
		provideDynamicServiceVersion,
		wire.Value(service.Environment("local")),
	)
	return nil, nil
}

func NewProductionEntrypoint(_ context.Context) (*Entrypoint, error) {
	wire.Build(
		commonProvider,
		env.ProvideGoogleCloudProject,
		env.ProvideServiceVersion,
		log.ProvideCloudTraceLinkedLogger,
		o11y.ProvideGoogleCloudRunResource,
		o11y.ProvideGoogleTelemetryTraceExporter,
		wire.Value(service.Environment("production")),
	)
	return nil, nil
}
