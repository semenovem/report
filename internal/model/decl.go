package model

type ProductListResponse struct {
	Result struct {
		Items  []ProductListItem `json:"items"`
		Total  int               `json:"total"`
		LastID string            `json:"last_id"`
	} `json:"result"`
}

type ProductListItem struct {
	ProductID    int64  `json:"product_id"`
	OfferID      string `json:"offer_id"`
	IsFboVisible bool   `json:"is_fbo_visible"`
	IsFbsVisible bool   `json:"is_fbs_visible"`
	Archived     bool   `json:"archived"`
	IsDiscounted bool   `json:"is_discounted"`
}
