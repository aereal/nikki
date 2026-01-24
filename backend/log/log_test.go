package log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/aereal/nikki/backend/log"
	"github.com/aereal/nikki/backend/log/attr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.opentelemetry.io/otel/trace"
)

type logEntry = map[string]any

var (
	traceID    = trace.TraceID{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	strTraceID = "30313233343536373839616263646566"
	spanID     = trace.SpanID{'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n'}
	strSpanID  = "6768696a6b6c6d6e"
)

func TestLogger_handler(t *testing.T) {
	t.Parallel()

	errOops := errors.New("oops") //nolint:err113
	testCases := []struct {
		name  string
		do    func(l *slog.Logger)
		level slog.Level
		want  []logEntry
	}{
		{
			name: "only message",
			do: func(l *slog.Logger) {
				l.Info("msg")
			},
			want: []logEntry{
				{
					slog.MessageKey:   "msg",
					slog.LevelKey:     "INFO",
					"service.version": "latest",
				},
			},
		},
		{
			name: "with error",
			do: func(l *slog.Logger) {
				l.Error("msg", attr.Error(errOops))
			},
			want: []logEntry{
				{
					slog.MessageKey:   "msg",
					slog.LevelKey:     "ERROR",
					"service.version": "latest",
					"error": map[string]any{
						"msg":  "oops",
						"type": "*errors.errorString",
					},
				},
			},
		},
		{
			name: "log in the OpenTelemetry spans",
			do:   doOtelWithAttrs(),
			want: []logEntry{
				{
					slog.MessageKey:                        "msg",
					slog.LevelKey:                          "INFO",
					"service.version":                      "latest",
					"logging.googleapis.com/trace":         "projects/dummy/traces/" + strTraceID,
					"logging.googleapis.com/spanId":        strSpanID,
					"logging.googleapis.com/trace_sampled": true,
				},
			},
		},
		{
			name: "log with other attribute in the OpenTelemetry spans",
			do:   doOtelWithAttrs(slog.Bool("debug", false)),
			want: []logEntry{
				{
					slog.MessageKey:                        "msg",
					slog.LevelKey:                          "INFO",
					"service.version":                      "latest",
					"debug":                                false,
					"logging.googleapis.com/trace":         "projects/dummy/traces/" + strTraceID,
					"logging.googleapis.com/spanId":        strSpanID,
					"logging.googleapis.com/trace_sampled": true,
				},
			},
		},
		{
			name: "log with existence otel group in the OpenTelemetry spans",
			do:   doOtelWithAttrs(slog.GroupAttrs("otel", slog.Bool("traced", true))),
			want: []logEntry{
				{
					slog.MessageKey:   "msg",
					slog.LevelKey:     "INFO",
					"service.version": "latest",
					"otel": map[string]any{
						"traced": true,
					},
					"logging.googleapis.com/trace":         "projects/dummy/traces/" + strTraceID,
					"logging.googleapis.com/spanId":        strSpanID,
					"logging.googleapis.com/trace_sampled": true,
				},
			},
		},
		{
			name: "log with existence otel scalar attribute in the OpenTelemetry spans",
			do:   doOtelWithAttrs(slog.Bool("otel", true)),
			want: []logEntry{
				{
					slog.MessageKey:                        "msg",
					slog.LevelKey:                          "INFO",
					"service.version":                      "latest",
					"otel":                                 true,
					"logging.googleapis.com/trace":         "projects/dummy/traces/" + strTraceID,
					"logging.googleapis.com/spanId":        strSpanID,
					"logging.googleapis.com/trace_sampled": true,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			buf := new(bytes.Buffer)
			logger := log.ProvideLogger(buf, tc.level, "dummy", "latest")
			tc.do(logger)
			got := make([]logEntry, 0)
			for line := range bytes.Lines(buf.Bytes()) {
				e := make(logEntry)
				if err := json.Unmarshal(line, &e); err != nil {
					t.Fatal(err)
				}
				got = append(got, e)
			}
			if diff := diffLogEntries(tc.want, got); diff != "" {
				t.Errorf("(-want, +got):\n%s", diff)
			}
		})
	}
}

func diffLogEntries(want, got []logEntry) string {
	return cmp.Diff(
		want, got,
		cmpopts.IgnoreMapEntries(func(key string, _ any) bool {
			return key == slog.TimeKey || key == attr.KeySourceLocation
		}),
	)
}

func doOtelWithAttrs(attrs ...slog.Attr) func(*slog.Logger) {
	return func(l *slog.Logger) {
		sc := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		})
		ctx := trace.ContextWithSpanContext(context.Background(), sc)
		l.LogAttrs(ctx, slog.LevelInfo, "msg", attrs...)
	}
}
