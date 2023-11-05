package controller

import (
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/lg"
	"github.com/semenovem/report/internal/provider"
)

type Controller struct {
	config   *config.Main
	logger   *lg.Lg
	provider *provider.Provider
}

func New(c *config.Main, l *lg.Lg, p *provider.Provider) *Controller {
	return &Controller{
		config:   c,
		logger:   l,
		provider: p,
	}
}
