package provider

import (
	"context"
	"encoding/json"
	"github.com/semenovem/report/internal/model"
	"net/http"
)

func (p *Provider) ProductList(
	ctx context.Context,
	marketID MarketID,
) (*model.ProductListResponse, error) {
	var (
		ll = p.logger.Named("ProductList")
		co = &conn{
			marketID: marketID,
			APIPath:  "/v2/product/list",
			Method:   http.MethodPost,
		}
	)

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

	var (
		arg = request{
			Filter: filter{
				Visibility: "ALL",
			},
			LastId: "",
			Limit:  100,
		}
		m = model.ProductListResponse{}
	)

	b, err := p.request(ctx, co, arg)
	if err != nil {
		ll.Named("request").Debug(err.Error())
		return nil, err
	}

	if err = json.Unmarshal(b, &m); err != nil {
		ll.Named("Unmarshal").With("dat", string(b)).Error(err.Error())
		return nil, err
	}

	return &m, nil
}
