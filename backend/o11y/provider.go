package o11y

import (
	"context"
	"fmt"
	"iter"
	"slices"

	"github.com/aereal/nikki/backend/adapters/gcp/metadata"
	"github.com/aereal/nikki/backend/o11y/service"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
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

func ProvideGoogleTelemetryTraceExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	creds, err := oauth.NewApplicationDefault(ctx)
	if err != nil {
		return nil, fmt.Errorf("oauth.NewApplicationDefault: %w", err)
	}
	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("telemetry.googleapis.com:443"),
		otlptracegrpc.WithDialOption(grpc.WithPerRPCCredentials(creds)),
	)
}

func ProvideResource(ctx context.Context, version service.Version, env service.Environment) (*resource.Resource, error) {
	return resource.New(ctx, slices.Collect(commonResourceOpts(version, env))...)
}

func ProvideGoogleCloudRunResource(ctx context.Context, version service.Version, env service.Environment, project metadata.Project) (*resource.Resource, error) {
	opts := slices.Collect(commonResourceOpts(version, env))
	opts = append(opts,
		resource.WithAttributes(attribute.String("gcp.project_id", string(project))), // this is mandatory for Google Cloud Telemetry API
		resource.WithDetectors(gcp.NewDetector()),
	)
	return resource.New(ctx, opts...)
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
