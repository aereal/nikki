package log

import (
	"io"
	"log/slog"
	"os"

	"github.com/aereal/nikki/backend/adapters/gcp/metadata"
	"github.com/aereal/nikki/backend/o11y/service"
)

type Output io.Writer

func ProvideStdout() Output { return os.Stdout }

func ProvideLogger(output Output, level slog.Level, project metadata.Project, version service.Version) *slog.Logger {
	h := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(newHandler(h, project, version))
}

type GlobalInstrumentationToken struct{}

func ProvideGlobalInstrumentation(l *slog.Logger) (_ GlobalInstrumentationToken) {
	slog.SetDefault(l)
	return
}
