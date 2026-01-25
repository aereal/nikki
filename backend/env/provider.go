package env

import (
	"log/slog"

	"github.com/aereal/nikki/backend/web"
)

func ProvidePort(vars Variables) (web.Port, error) {
	scan := scanOrElse(scannerWithParse(stringAs[web.Port]), "8080")
	var port web.Port
	if err := scan(vars, "PORT", &port); err != nil {
		return "", err
	}
	return port, nil
}

func ProvideLogLevel(vars Variables) (slog.Level, error) {
	scan := scanOrElse(scannerWithParse(parseLogLevel), slog.LevelInfo)
	var level slog.Level
	if err := scan(vars, "LOG_LEVEL", &level); err != nil {
		return 0, err
	}
	return level, nil
}
