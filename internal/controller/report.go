package controller

import (
	"github.com/labstack/echo/v4"
)

func (cnt *Controller) Report1(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		ll  = cnt.logger.Func(ctx, "Report1")
	)

	ll.Infof("reported")

	return c.JSON(200, "")
}
