package web

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aereal/nikki/backend/log/attr"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Port string

func ProvideServer(tp trace.TracerProvider, port Port, db *sql.DB) *Server {
	return &Server{
		port:   port,
		tp:     tp,
		tracer: tp.Tracer("github.com/aereal/nikki/backend/web.Server"),
		db:     db,
	}
}

type Server struct {
	port   Port
	tp     trace.TracerProvider
	tracer trace.Tracer
	db     *sql.DB
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
	mux.Handle("POST /init", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := doInit(r.Context(), s.db); err != nil {
			handleError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	mux.Handle("GET /counter/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, err := doGetCounter(r.Context(), s.db, r.PathValue("id"))
		if err != nil {
			handleError(w, err)
			return
		}
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]int64{"value": value})
	}))
	mux.Handle("PATCH /counter/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := doIncrement(r.Context(), s.db, r.PathValue("id")); err != nil {
			handleError(w, err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
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

func doInit(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `create table counter (id text primary key, value integer not null)`)
	return err
}

func doGetCounter(ctx context.Context, db *sql.DB, id string) (int64, error) {
	row := db.QueryRowContext(ctx, `select value from counter where id = ?`, id)
	if err := row.Err(); err != nil {
		return 0, err
	}
	var value int64
	if err := row.Scan(&value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("sql.Row.Scan: %w", err)
	}
	return value, nil
}

func doIncrement(ctx context.Context, db *sql.DB, id string) error {
	_, err := db.ExecContext(ctx, `insert into counter (id, value) values (?, ?) on conflict (id) do update set value = excluded.value + 1`, id, 1)
	return err
}

func handleError(w http.ResponseWriter, err error) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}) //nolint:errchkjson // ignore JSON write failure
}
