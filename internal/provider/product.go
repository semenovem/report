package provider

import (
	"context"
	"encoding/json"
	"github.com/semenovem/report/internal/model"
	"net/http"
	"time"
)

func (p *Provider) ProductList(
	ctx context.Context,
	marketID MarketID,
) ([]*model.ProductListItem, error) {
	var (
		ll = p.logger.Named("ReportProducts")
		co = &conn{
			marketID: marketID,
			APIPath:  "/v2/product/list",
			Method:   http.MethodPost,
		}
		arg = request{
			Filter: filter{
				Visibility: "ALL",
			},
			LastId: "",
			Limit:  100,
		}
		items = make([]*model.ProductListItem, 0)
	)

	for {
		b, err := p.request(ctx, co, arg)
		if err != nil {
			ll.Named("request").Debug(err.Error())
			return nil, err
		}

		var resp productListResponse

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
