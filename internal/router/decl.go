package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/controller"
	"github.com/semenovem/report/internal/zoo/lg"
	"net/http"
)

type Router struct {
	ctx         context.Context
	logger      *lg.Lg
	server      *echo.Echo
	config      *config.Main
	auth, admin *echo.Group
	cnt         *controller.Controller
}

func (r *Router) Start() {
	go func() {
		<-r.ctx.Done()
		if err := r.server.Close(); err != nil {
			r.logger.Named("Close").Error(err.Error())
		}
	}()

	r.logger.Infof("router start on %d", r.config.Rest.Port)

	r.server.HidePort = true
	r.server.HideBanner = true

	addr := fmt.Sprintf(":%d", r.config.Rest.Port)

	if err := r.server.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		r.logger.Named("Start").Error(err.Error())
	}
}
