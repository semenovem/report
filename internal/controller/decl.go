package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/lg"
	"github.com/semenovem/report/internal/provider"
	"net/http"
)

const (
	stockFBS         = "fbs"
	stockFBO         = "fbo"
	stockCrossBorder = "crossborder"
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

func (ct *Controller) errResponse(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, err.Error())
}
