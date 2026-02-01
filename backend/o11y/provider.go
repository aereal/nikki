package o11y

import (
	"context"
	"iter"
	"slices"

	"github.com/aereal/nikki/backend/o11y/service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func ProvideTracerProvider(exporter *otlptrace.Exporter, res *resource.Resource) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
}

func ProvideSidecarExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
}

func ProvideResource(ctx context.Context, version service.Version, env service.Environment) (*resource.Resource, error) {
	return resource.New(ctx, slices.Collect(commonResourceOpts(version, env))...)
}

func commonResourceOpts(version service.Version, env service.Environment) iter.Seq[resource.Option] {
	return func(yield func(resource.Option) bool) {
		attrs := make([]attribute.KeyValue, 0, 3)
		attrs = append(attrs,
			semconv.ServiceName(service.ServiceName),
			semconv.ServiceVersion(string(version)),
			semconv.DeploymentEnvironmentName(string(env)),
		)
		if !yield(resource.WithAttributes(attrs...)) {
			return
		}
	}
}
