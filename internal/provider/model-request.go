package provider

type filter struct {
	OfferId    []string `json:"offer_id,omitempty"`
	ProductId  []string `json:"product_id,omitempty"`
	Visibility string   `json:"visibility,omitempty"`
}

type request struct {
	Filter filter `json:"filter"`
	LastId string `json:"last_id"`
	Limit  int    `json:"limit"`
}
