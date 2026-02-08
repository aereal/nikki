package web

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aereal/nikki/backend/graph"
	"github.com/aereal/nikki/backend/log/attr"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Port string

func ProvideServer(tp trace.TracerProvider, port Port, gh graph.Handler) *Server {
	return &Server{
		port:           port,
		tp:             tp,
		tracer:         tp.Tracer("github.com/aereal/nikki/backend/web.Server"),
		graphqlHandler: gh,
	}
}

type Server struct {
	port           Port
	tp             trace.TracerProvider
	tracer         trace.Tracer
	graphqlHandler graph.Handler
}

func (s *Server) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:              net.JoinHostPort("", string(s.port)),
		Handler:           s.handler(),
		ReadHeaderTimeout: time.Second * 5,
	}

	sigCtx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()
	go func() {
		<-sigCtx.Done()
		shutdownGrace := time.Second * 5
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(sigCtx), shutdownGrace)
		defer cancel()
		slog.DebugContext(shutdownCtx, "shutting down server", slog.Duration("grace", shutdownGrace))
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.WarnContext(shutdownCtx, "server has stopped inelegantly", attr.Error(err))
		}
	}()

	slog.InfoContext(ctx, "starting a server", slog.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /graphql", s.graphqlHandler)
	return otelhttp.NewMiddleware("default",
		otelhttp.WithTracerProvider(s.tp),
		otelhttp.WithPropagators(propagation.TraceContext{}),
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			if r.Pattern != "" {
				return r.Pattern
			}
			return r.Method + " " + r.URL.Path
		}),
	)(mux)
}
