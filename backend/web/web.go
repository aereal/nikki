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

	"github.com/aereal/nikki/backend/log/attr"
)

type Port string

func ProvideServer(port Port) *Server {
	return &Server{
		port: port,
	}
}

type Server struct {
	port Port
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
	return mux
}
