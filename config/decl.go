package config

import (
	"github.com/caarlos0/env/v9"
	"strings"
)

type Base struct {
	Env      string `env:"ENV,required"`
	CliMode  bool   `env:"CLI_MODE,required"`
	LogLevel int8   `env:"LOG_LEVEL,required"`
}

func (b *Base) IsDev() bool {
	return strings.EqualFold(b.Env, "DEV")
}

type Rest struct {
	Port              uint16 `env:"REST_PORT,required"`
	CorsAllowedHosts  string `env:"REST_CORS_ALLOWED_HOSTS,required"`
	CorsMaxAgeSeconds uint32 `env:"REST_CORS_MAX_AGE_SECONDS,required"`
}

func ParseAPI() (*Main, error) {
	var cfg Main

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
