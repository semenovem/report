package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/semenovem/report/internal/provider"
	"github.com/semenovem/report/internal/zoo/spreadsheet"
	"net/http"
	"time"
)

func (ct *Controller) ReportProducts(c echo.Context) error {
	var (
		ctx       = c.Request().Context()
		sessionID = c.FormValue(SessionIDName)
		ll        = ct.logger.Func(ctx, "ReportProducts").With("sessionID", sessionID)
		table     = make([][]string, 0, 100)
	)

	session, ok := ct.sessions[sessionID]
	if !ok {
		ll.Named("session not found")
		return c.HTML(403, "forbidden")
	}

	if session.fileDownload {
		return c.HTML(http.StatusTooManyRequests, "the report has already been downloaded")
	}

	session.fileDownload = true

	for i, marketID := range []provider.MarketID{provider.Ozon1, provider.Ozon2} {
		tab, err := ct.dataMining.ProductListReport(ctx, marketID, i == 0)
		if err != nil {
			ll.Named("mining").Debug(err.Error())
			return ct.errResponse(c, err)
		}
		table = append(table, tab...)
	}

	b, err := spreadsheet.CreateCSV(table)
	if err != nil {
		ll.Named("spreadsheet.CreateCSV").Error(err.Error())
		return ct.errResponse(c, err)
	}

	var (
		contentType = spreadsheet.ContentTypeTextCSV
		fileExt     = spreadsheet.FileExtensionCSV
		filename    = fmt.Sprintf(
			"product_list.%s.%s",
			time.Now().Format("2006-01-02"),
			fileExt,
		)
	)

	c.Response().Header().Set("Content-Description", "File Transfer")
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
	c.Response().Header().Set("Cache-Control", "No-Cache")

	ll.Debug("file product list downloaded")

	return c.Blob(http.StatusOK, contentType, b.Bytes())
}
