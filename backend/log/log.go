package log

import (
	"io"
	"log/slog"
	"os"
)

type Output io.Writer

func ProvideStdout() Output { return os.Stdout }

type GoogleCloudProject string

type ServiceVersion string

func ProvideLogger(output Output, level slog.Level, project GoogleCloudProject, version ServiceVersion) *slog.Logger {
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
