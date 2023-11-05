package provider

import "github.com/semenovem/report/config"

type Provider struct {
	config *config.Main
}

func New(config *config.Main) *Provider {
	return &Provider{
		config: config,
	}
}
