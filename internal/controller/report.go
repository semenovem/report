package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/semenovem/report/internal/spreadsheet"
)

func (cnt *Controller) Report1(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		ll  = cnt.logger.Func(ctx, "Report1")
	)

	ll.Infof("reported")

	t, err := cnt.bl.Report1()
	if err != nil {
		ll.Named("bl.Report1").Debug(err.Error())
		return cnt.errResponse(c, err)
	}

	b, err := spreadsheet.CreateCSV(t)
	if err != nil {
		ll.Named("spreadsheet.CreateCSV").Error(err.Error())
		return cnt.errResponse(c, err)
	}

	ll.Infof(">>> %+v", t)
	ll.Infof(">>> %+v", b.String())

	return c.JSON(200, "")
}
