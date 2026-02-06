package main

import (
	"context"
	"os"
	"time"

	"github.com/aereal/nikki/backend/entrypoint"
	"github.com/aereal/nikki/backend/log"
	"github.com/aereal/nikki/backend/usecases"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()
	a, err := build(ctx)
	if err != nil {
		return entrypoint.ExitCodeOf(err)
	}
	if err := a.run(ctx); err != nil {
		return entrypoint.ExitCodeOf(err)
	}
	return 0
}

func provideApp(_ log.GlobalInstrumentationToken, u usecases.ImportMTExport, tp *sdktrace.TracerProvider) *app {
	return &app{ImportMTExport: u, tp: tp}
}

type app struct {
	usecases.ImportMTExport

	tp *sdktrace.TracerProvider
}

func (a *app) run(ctx context.Context) error {
	defer func() {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second*5)
		defer cancel()
		_ = a.tp.Shutdown(ctx)
	}()
	return a.ImportMTExport.ImportMTExport(ctx)
}
