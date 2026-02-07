package entrypoint

import (
	"context"
	"errors"
	"log/slog"

	"github.com/aereal/nikki/backend/log"
	"github.com/aereal/nikki/backend/log/attr"
	"github.com/aereal/nikki/backend/web"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func provideEntrypoint(ctx context.Context, _ log.GlobalInstrumentationToken, s *web.Server, tp *sdktrace.TracerProvider) *Entrypoint {
	return &Entrypoint{ctx: ctx, server: s, tp: tp}
}

type Entrypoint struct {
	ctx    context.Context //nolint:containedctx // allow contained ctx
	server *web.Server
	tp     *sdktrace.TracerProvider
}

func ExitCodeOf(err error) (code int) {
	defer func() {
		level := slog.LevelDebug
		if code > 0 {
			level = slog.LevelError
		}
		attrs := make([]slog.Attr, 1, 2)
		attrs[0] = slog.Int("exit_code", code)
		if err != nil {
			attrs = append(attrs, attr.Error(err))
		}
		slog.LogAttrs(context.Background(), level, "application exited", attrs...)
	}()

	if err == nil {
		return 0
	}
	var hasExitCode interface{ ExitCode() int }
	if errors.As(err, &hasExitCode) {
		return hasExitCode.ExitCode()
	}
	return 1
}
