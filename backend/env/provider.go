package env

import (
	"log/slog"

	"github.com/aereal/nikki/backend/adapters/gcp/metadata"
	"github.com/aereal/nikki/backend/infra/db"
	"github.com/aereal/nikki/backend/o11y/service"
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

func ProvideDBEndpoint(vars Variables) (db.Endpoint, error) {
	scan := scanOrElse(scanString, "local.db")
	endpoint := &db.FileEndpoint{Params: &db.ParameterSet{Cache: db.CacheModeShared}}
	if err := scan(vars, "DB_FILE", &endpoint.Path); err != nil {
		return nil, err
	}
	return endpoint, nil
}

func ProvideGoogleCloudProject(vars Variables) (metadata.Project, error) {
	var project metadata.Project
	if err := scannerWithParse(stringAs[metadata.Project])(vars, "GOOGLE_CLOUD_PROJECT", &project); err != nil {
		return "", err
	}
	return project, nil
}

func ProvideServiceVersion(vars Variables) (service.Version, error) {
	var version service.Version
	if err := scannerWithParse(stringAs[service.Version])(vars, "SERVICE_VERSION", &version); err != nil {
		return "", err
	}
	return version, nil
}
