package controller

import (
	"github.com/labstack/echo/v4"
	"os"
)

func (ct *Controller) Index(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		ll  = ct.logger.Func(ctx, "Index")
	)

	dat, err := os.ReadFile(ct.config.IndexPath)
	if err != nil {
		ll.Named("os.ReadFile").Error(err.Error())
		return ct.errResponse(c, err)
	}

	ll.Debug("index data send")

	return c.HTML(200, string(dat))
}
