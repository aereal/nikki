package o11y

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func ProvideTracerProvider() *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider() // TODO
}
