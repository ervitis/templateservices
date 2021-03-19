package infra

import (
	"github.com/ervitis/logme"
	"github.com/ervitis/logme/config_loaders"
	"os"
)

var (
	Logger logme.Loggerme
)

type ConfigLogger struct {
	*config_loaders.EnvLoad
}

func NewConfig() (*ConfigLogger, error) {
	_ = os.Setenv("LOG_LEVEL", "DEBUG")
	_ = os.Setenv("LOG_FIELDS", "[component=grpc-basket-conn,service=grpc-basket]")

	cfg, err := config_loaders.NewEnvLogme()
	if err != nil {
		return nil, err
	}
	return &ConfigLogger{cfg}, nil
}

func NewLogger(config *ConfigLogger) logme.Loggerme {
	return logme.NewLogme(config)
}
