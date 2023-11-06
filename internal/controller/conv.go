package controller

import "github.com/semenovem/report/internal/model"

func mapProductStock(a []*model.ProductStock) map[int64]map[string]model.ProductInfoStock {
	b := make(map[int64]map[string]model.ProductInfoStock)

	for _, v := range a {
		b[v.ProductID] = make(map[string]model.ProductInfoStock)

		for _, v2 := range v.Stocks {
			b[v.ProductID][v2.Type] = *v2
		}
	}

	return b
}
