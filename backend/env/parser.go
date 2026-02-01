package env

import (
	"log/slog"
)

func stringAs[T ~string](s string) (T, error) { return T(s), nil }

func parseLogLevel(s string) (slog.Level, error) {
	var level slog.Level
	if err := (&level).UnmarshalText([]byte(s)); err != nil {
		return 0, err
	}
	return level, nil
}
