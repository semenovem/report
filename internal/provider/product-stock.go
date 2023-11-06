package provider

import (
	"context"
	"encoding/json"
	"github.com/semenovem/report/internal/model"
	"net/http"
	"time"
)

func (p *Provider) ProductInfoStock(
	ctx context.Context,
	marketID MarketID,
) ([]*model.ProductStock, error) {
	var (
		ll = p.logger.Named("ProductInfoStock")
		co = &conn{
			marketID: marketID,
			APIPath:  "/v3/product/info/stocks",
			Method:   http.MethodPost,
		}
		arg = request{
			Filter: filter{
				Visibility: "ALL",
			},
			LastId: "",
			Limit:  100,
		}
		items = make([]*model.ProductStock, 0)
	)

	for {
		b, err := p.request(ctx, co, arg)
		if err != nil {
			ll.Named("request").Debug(err.Error())
			return nil, err
		}

		var resp productInfoStocksResponse

		if err = json.Unmarshal(b, &resp); err != nil {
			ll.Named("Unmarshal").With("dat", string(b)).Error(err.Error())
			return nil, err
		}

		items = append(items, resp.Result.Items...)

		if len(resp.Result.Items) == 0 || resp.Result.LastID == "" || resp.Result.Total < arg.Limit {
			break
		}

		arg.LastId = resp.Result.LastID

		time.Sleep(time.Millisecond * 500)
	}

	return items, nil
}
