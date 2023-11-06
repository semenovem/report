package model

type ProductListItem struct {
	ProductID int64  `json:"product_id"`
	OfferID   string `json:"offer_id"`

	// Этих полей нет в документации
	IsFboVisible bool `json:"is_fbo_visible"`
	IsFbsVisible bool `json:"is_fbs_visible"`
	Archived     bool `json:"archived"`
	IsDiscounted bool `json:"is_discounted"`
}

type ProductStock struct {
	ProductID int64               `json:"product_id"`
	OfferID   string              `json:"offer_id"`
	Stocks    []*ProductInfoStock `json:"stocks"`
}

type ProductInfoStock struct {
	Type     string `json:"type"`
	Present  int    `json:"present"`
	Reserved int    `json:"reserved"`
}
