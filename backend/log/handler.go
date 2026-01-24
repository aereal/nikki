package log

import (
	"context"
	"log/slog"

	"github.com/aereal/nikki/backend/log/attr"
	"go.opentelemetry.io/otel/trace"
)

func newHandler(base slog.Handler, project GoogleCloudProject, version ServiceVersion) *handler {
	return &handler{base: base, project: project, version: version}
}

type handler struct {
	base    slog.Handler
	project GoogleCloudProject
	version ServiceVersion
}

var _ slog.Handler = (*handler)(nil)

func (h *handler) Handle(ctx context.Context, orig slog.Record) error {
	ret := slog.NewRecord(orig.Time, orig.Level, orig.Message, orig.PC)
	for a := range orig.Attrs {
		if a.Key == slog.SourceKey {
			continue
		}
		ret.AddAttrs(a)
	}
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		ret.AddAttrs(
			attr.Trace(string(h.project), sc.TraceID()),
			attr.SpanID(sc.SpanID()),
			attr.TraceSampled(sc.IsSampled()),
		)
	}
	ret.AddAttrs(
		attr.SourceLocation(ret.Source()),
		attr.ServiceVersion(string(h.version)),
	)
	return h.base.Handle(ctx, ret)
}

func (h *handler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.base.Enabled(ctx, l)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newHandler(h.base.WithAttrs(attrs), h.project, h.version)
}

func (h *handler) WithGroup(name string) slog.Handler {
	return newHandler(h.base.WithGroup(name), h.project, h.version)
}
