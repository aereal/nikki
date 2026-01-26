package env_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/aereal/nikki/backend/env"
	"github.com/aereal/nikki/backend/infra/db"
	"github.com/aereal/nikki/backend/web"
	"github.com/google/go-cmp/cmp"
)

type testCaseRoot struct {
	port       map[string]result[web.Port]
	logLevel   map[string]result[slog.Level]
	dbEndpoint map[string]result[db.Endpoint]
}

func TestProviders(t *testing.T) {
	t.Parallel()

	root := testCaseRoot{
		port: map[string]result[web.Port]{
			"ok": {
				provideFunc: env.ProvidePort,
				want:        "8888",
				wantErr:     nil,
				variables:   env.Variables{"PORT": "8888"},
			},
			"default": {
				provideFunc: env.ProvidePort,
				want:        "8080",
				wantErr:     nil,
				variables:   env.Variables{},
			},
		},
		logLevel: map[string]result[slog.Level]{
			"ok": {
				provideFunc: env.ProvideLogLevel,
				want:        slog.LevelWarn,
				wantErr:     nil,
				variables:   env.Variables{"LOG_LEVEL": "WARN"},
			},
			"default": {
				provideFunc: env.ProvideLogLevel,
				want:        slog.LevelInfo,
				wantErr:     nil,
				variables:   env.Variables{},
			},
			"invalid": {
				provideFunc: env.ProvideLogLevel,
				want:        0,
				wantErr:     literalError(`slog: level string "abc": unknown name`),
				variables:   env.Variables{"LOG_LEVEL": "abc"},
			},
		},
		dbEndpoint: map[string]result[db.Endpoint]{
			"ok": {
				provideFunc: env.ProvideDBEndpoint,
				want:        &db.FileEndpoint{Path: "a.db", Params: &db.ParameterSet{Cache: db.CacheModeShared}},
				wantErr:     nil,
				variables:   env.Variables{"DB_FILE": "a.db"},
			},
		},
	}
	for name, tc := range root.port {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assertProvider(t, &tc)
		})
	}
	for name, tc := range root.logLevel {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assertProvider(t, &tc)
		})
	}
	for name, tc := range root.dbEndpoint {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assertProvider(t, &tc)
		})
	}
}

type result[T any] struct {
	provideFunc func(env.Variables) (T, error)
	want        T
	wantErr     error
	variables   env.Variables
}

func assertProvider[T any](t *testing.T, tc *result[T]) {
	t.Helper()

	got, gotErr := tc.provideFunc(tc.variables)
	if !errors.Is(tc.wantErr, gotErr) {
		t.Errorf("error: want=%s got=%s", tc.wantErr, gotErr)
	}
	if gotErr != nil {
		return
	}
	if diff := cmp.Diff(tc.want, got); diff != "" {
		t.Errorf("value (-want, +got):\n%s", diff)
	}
}

type literalError string

func (l literalError) Error() string { return string(l) }

func (l literalError) Is(other error) bool {
	if other == nil {
		return false
	}
	return string(l) == other.Error()
}
