package action

import (
	"github.com/semenovem/report/internal/provider"
	"github.com/semenovem/report/internal/zoo/lg"
)

const (
	stockFBS         = "fbs"
	stockFBO         = "fbo"
	stockCrossBorder = "crossborder"
)

type DataMining struct {
	logger   *lg.Lg
	provider *provider.Provider
}

func NewDataMining(l *lg.Lg, p *provider.Provider) *DataMining {
	return &DataMining{
		logger:   l,
		provider: p,
	}
}
