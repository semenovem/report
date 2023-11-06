package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/semenovem/report/internal/provider"
	"github.com/semenovem/report/internal/spreadsheet"
	"net/http"
	"strconv"
	"time"
)

func (ct *Controller) ProductList(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		ll  = ct.logger.Func(ctx, "ProductList")
	)

	ll.Infof("reported")

	response, err := ct.provider.ProductList(ctx, provider.Ozon1)
	if err != nil {
		ll.Named("provider.ProductList").Debug(err.Error())
		return ct.errResponse(c, err)
	}

	table := make([][]string, 0, len(response.Result.Items)+1)

	table = append(table, []string{
		"num", "product_id", "offer_id", "is_fbo_visible", "is_fbs_visible", "archived", "is_discounted",
	})

	for i, m := range response.Result.Items {
		row := []string{
			strconv.Itoa(i + 1),
			strconv.FormatInt(m.ProductID, 10),
			m.OfferID,
			strconv.FormatBool(m.IsFboVisible),
			strconv.FormatBool(m.IsFbsVisible),
			strconv.FormatBool(m.Archived),
			strconv.FormatBool(m.IsDiscounted),
		}
		table = append(table, row)
	}

	b, err := spreadsheet.CreateCSV(table)
	if err != nil {
		ll.Named("spreadsheet.CreateCSV").Error(err.Error())
		return ct.errResponse(c, err)
	}

	var (
		contentType = spreadsheet.ContentTypeTextCSV
		fileExt     = spreadsheet.FileExtensionCSV
		filename    = fmt.Sprintf("product_list.%s.%s", time.Now().Format("2006-01-02"), fileExt)
	)

	c.Response().Header().Set("Content-Description", "File Transfer")
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
	c.Response().Header().Set("Cache-Control", "No-Cache")

	ll.Debug("file product list downloaded")

	return c.Blob(http.StatusOK, contentType, b.Bytes())
}
