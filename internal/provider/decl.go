package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/semenovem/report/config"
	"github.com/semenovem/report/internal/lg"
	"io"
	"net/http"
)

type MarketID string

const (
	Ozon1 MarketID = "ozon1"
	Ozon2 MarketID = "ozon2"
)

func (m MarketID) Name() string {
	switch m {
	case Ozon1:
		return "ozon1"
	case Ozon2:
		return "ozon2"
	default:
		return "unknown"
	}
}

var (
	ErrUnknownMarket = errors.New("unknown market")
	ErrNot200        = errors.New("failure")
)

type Provider struct {
	logger *lg.Lg
	config *config.Main
}

type conn struct {
	marketID MarketID
	APIPath  string
	Method   string
}

func New(config *config.Main, logger *lg.Lg) *Provider {
	return &Provider{
		config: config,
		logger: logger,
	}
}

func (p *Provider) makeHTTPClient(ctx context.Context, conn *conn, b []byte) (*http.Request, error) {
	var path = conn.APIPath

	switch conn.marketID {
	case Ozon1, Ozon2:
		path = p.config.Ozon.Path + conn.APIPath
	default:
		p.logger.With("conn", conn).Error("unknown marketID")
		return nil, ErrUnknownMarket
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewBuffer(b))
	if err != nil {
		p.logger.Named("makeHTTPClient").With("conn", conn).Error(err.Error())
		return nil, err
	}

	switch conn.marketID {
	case Ozon1:
		req.Header.Set("Client-Id", p.config.Ozon.ClientID1)
		req.Header.Set("Api-Key", p.config.Ozon.APIKey1)
	case Ozon2:
		req.Header.Set("Client-Id", p.config.Ozon.ClientID2)
		req.Header.Set("Api-Key", p.config.Ozon.APIKey2)
	default:
		p.logger.With("conn", conn).Error(ErrUnknownMarket.Error())
		return nil, ErrUnknownMarket
	}

	return req, nil
}

func (p *Provider) request(ctx context.Context, conn *conn, body any) ([]byte, error) {
	ll := p.logger.Named("request").
		With("marketID", conn.marketID).
		With("apiPath", conn.APIPath)

	b, err := json.Marshal(body)
	if err != nil {
		ll.Named("json.Marshal").Error(err.Error())
		return nil, err
	}

	request, err := p.makeHTTPClient(ctx, conn, b)
	if err != nil {
		ll.Named("makeHTTPClient").Debug(err.Error())
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ll.Named("client.Do").Error(err.Error())
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		ll.Error("response failed")

		bodyResp, err := io.ReadAll(response.Body)
		if err != nil {
			ll.Named("io.ReadAll").Error(err.Error())
			return nil, err
		}

		ll.Named("body").Debug(string(bodyResp))

		return nil, ErrNot200
	}

	bodyResp, err := io.ReadAll(response.Body)
	if err != nil {
		ll.Named("io.ReadAll").Error(err.Error())
		return nil, err
	}

	return bodyResp, nil
}
