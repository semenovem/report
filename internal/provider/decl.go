package provider

import (
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/lg"
)

type Provider struct {
	logger *lg.Lg
	config *config.Main
}

func New(config *config.Main, logger *lg.Lg) *Provider {
	return &Provider{
		config: config,
		logger: logger,
	}
}
