package provider

import (
	"errors"
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/zoo/lg"
)

type MarketID string

const (
	Ozon1 MarketID = "ozon1"
	Ozon2 MarketID = "ozon2"
)

func (m MarketID) Name() string {
	switch m {
	case Ozon1:
		return "ozon1"
	case Ozon2:
		return "ozon2"
	default:
		return "unknown"
	}
}

var (
	ErrUnknownMarket = errors.New("unknown market")
	ErrNot200        = errors.New("failure")
)

type Provider struct {
	logger *lg.Lg
	config *config.Main
}

type conn struct {
	marketID MarketID
	APIPath  string
	Method   string
}

func New(config *config.Main, logger *lg.Lg) *Provider {
	return &Provider{
		config: config,
		logger: logger,
	}
}
