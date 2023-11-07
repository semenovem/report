package controller

import (
	"context"
	"github.com/semenovem/report/internal/provider"
	"strconv"
)

func (ct *Controller) mining(
	ctx context.Context,
	marketID provider.MarketID,
	header bool,
) ([][]string, error) {
	ll := ct.logger.Named("mining")

	// Получение данных
	// ----------------------------------------------
	products, err := ct.provider.ProductList(ctx, marketID)
	if err != nil {
		ll.Named("provider.ReportProducts").Debug(err.Error())
		return nil, err
	}

	stocks, err := ct.provider.ProductInfoStock(ctx, provider.Ozon1)
	if err != nil {
		ll.Named("provider.ReportProducts").Debug(err.Error())
		return nil, err
	}

	stockMap := mapProductStock(stocks)

	// Таблица
	// ---------------------------------------------
	table := make([][]string, 0, len(products)+1)

	if header {
		table = append(table, []string{
			"market",
			"num",
			"product_id",
			"offer_id",
			"is_fbo_visible",
			"is_fbs_visible",
			"archived",
			"is_discounted",
			// stockMap-fbs
			"stock_fbs_present",
			"stock_fbs_reserved",
			// stockMap-fbo
			"stock_fbo_present",
			"stock_fbo_reserved",
			// stockMap-cross-border
			"stock_crossborder_present",
			"stock_crossborder_reserved",
		})
	}

	for i, m := range products {
		row := []string{
			marketID.Name(),
			strconv.Itoa(i + 1),
			strconv.FormatInt(m.ProductID, 10),
			m.OfferID,
			strconv.FormatBool(m.IsFboVisible),
			strconv.FormatBool(m.IsFbsVisible),
			strconv.FormatBool(m.Archived),
			strconv.FormatBool(m.IsDiscounted),
			// stockMap-fbs
			strconv.Itoa(stockMap[m.ProductID][stockFBS].Present),
			strconv.Itoa(stockMap[m.ProductID][stockFBS].Reserved),
			// stockMap-fbo
			strconv.Itoa(stockMap[m.ProductID][stockFBO].Present),
			strconv.Itoa(stockMap[m.ProductID][stockFBO].Reserved),
			// stockMap-cross-border
			strconv.Itoa(stockMap[m.ProductID][stockCrossBorder].Present),
			strconv.Itoa(stockMap[m.ProductID][stockCrossBorder].Reserved),
		}
		table = append(table, row)
	}

	return table, nil
}
