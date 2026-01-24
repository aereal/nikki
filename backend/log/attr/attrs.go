package attr

import (
	"fmt"
	"log/slog"
	"reflect"

	"go.opentelemetry.io/otel/trace"
)

const (
	KeyTrace          = "logging.googleapis.com/trace"
	KeySpanID         = "logging.googleapis.com/spanId"
	KeyTraceSampled   = "logging.googleapis.com/trace_sampled"
	KeySourceLocation = "logging.googleapis.com/sourceLocation"
	KeyError          = "error"
	KeyServiceVersion = "service.version"
)

func Trace(project string, id trace.TraceID) slog.Attr {
	return slog.String(KeyTrace, fmt.Sprintf("projects/%s/traces/%s", project, id))
}

func SpanID(id trace.SpanID) slog.Attr {
	return slog.String(KeySpanID, id.String())
}

func TraceSampled(v bool) slog.Attr { return slog.Bool(KeyTraceSampled, v) }

func SourceLocation(src *slog.Source) slog.Attr {
	if src == nil {
		return slog.Attr{}
	}
	return slog.GroupAttrs(
		KeySourceLocation,
		slog.String("file", src.File),
		slog.Int("line", src.Line),
		slog.String("function", src.Function),
	)
}

func Error(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	typErr := reflect.TypeOf(err)
	attrs := make([]slog.Attr, 0, 3)
	attrs = append(attrs,
		slog.String("type", typErr.String()),
		slog.String("msg", err.Error()),
	)
	if pkg := typErr.PkgPath(); pkg != "" {
		attrs = append(attrs, slog.String("pkg", pkg))
	}
	return slog.GroupAttrs(KeyError, attrs...)
}

func ServiceVersion(v string) slog.Attr {
	return slog.String(KeyServiceVersion, v)
}
