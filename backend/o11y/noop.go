package o11y

import (
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func ProvideNoopTracerProvider() trace.TracerProvider { return noop.NewTracerProvider() }
