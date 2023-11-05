package main

import (
	"context"
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/lg"
	"github.com/semenovem/report/internal/provider"
	"github.com/semenovem/report/internal/router"
)

type app struct {
	ctx      context.Context
	logger   *lg.Lg
	router   *router.Router
	config   *config.Main
	provider *provider.Provider
}

func newApp(ctx context.Context, logger *lg.Lg, config *config.Main) error {
	var (
		ll = logger.Named("newApp")

		a = app{
			ctx:      ctx,
			logger:   logger,
			router:   nil,
			config:   config,
			provider: provider.New(config),
		}

		err error
	)

	a.router, err = router.New(ctx, logger, config, a.provider)
	if err != nil {
		ll.Named("router.New").Debug(err.Error())
		return err
	}

	go a.router.Start()

	return nil

}
