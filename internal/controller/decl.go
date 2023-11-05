package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/action"
	"github.com/semenovem/report/internal/lg"
	"net/http"
)

type Controller struct {
	config *config.Main
	logger *lg.Lg
	bl     *action.BL
}

func New(c *config.Main, l *lg.Lg, bl *action.BL) *Controller {
	return &Controller{
		config: c,
		logger: l,
		bl:     bl,
	}
}

func (cnt *Controller) errResponse(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, err.Error())
}
