package provider

import "github.com/semenovem/report/internal/model"

type productListResponse struct {
	Result struct {
		Items  []*model.ProductListItem `json:"items"`
		Total  int                      `json:"total"`
		LastID string                   `json:"last_id"`
	} `json:"result"`
}

type productInfoStocksResponse struct {
	Result struct {
		Items  []*model.ProductStock `json:"items"`
		Total  int                   `json:"total"`
		LastID string                `json:"last_id"`
	} `json:"result"`
}
