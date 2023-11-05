package action

import (
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/lg"
	"github.com/semenovem/report/internal/provider"
)

type BL struct {
	config   *config.Main
	logger   *lg.Lg
	provider *provider.Provider
}

func New(c *config.Main, l *lg.Lg, p *provider.Provider) *BL {
	return &BL{
		config:   c,
		logger:   l,
		provider: p,
	}
}
