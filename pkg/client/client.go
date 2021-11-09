package client

import (
	"context"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/api"
	"gitlab.com/ignitionrobotics/web/ign-go/encoders"
	"gitlab.com/ignitionrobotics/web/ign-go/net"
	"net/http"
	"net/url"
	"time"
)

// client contains the HTTP client to connect to the credits API.
type client struct {
	client net.Client
}

// IncreaseCredits performs an HTTP request to increase the given user credits.
func (c *client) IncreaseCredits(ctx context.Context, in api.IncreaseCreditsRequest) (api.IncreaseCreditsResponse, error) {
	var out api.IncreaseCreditsResponse
	if err := c.client.Call(ctx, "IncreaseCredits", &in, &out); err != nil {
		return api.IncreaseCreditsResponse{}, err
	}
	return out, nil
}

// DecreaseCredits performs an HTTP request to decrease the given user credits.
func (c *client) DecreaseCredits(ctx context.Context, in api.DecreaseCreditsRequest) (api.DecreaseCreditsResponse, error) {
	var out api.DecreaseCreditsResponse
	if err := c.client.Call(ctx, "DecreaseCredits", &in, &out); err != nil {
		return api.DecreaseCreditsResponse{}, err
	}
	return out, nil
}

// GetBalance performs an HTTP request to get the customer's balance.
func (c *client) GetBalance(ctx context.Context, in api.GetBalanceRequest) (api.GetBalanceResponse, error) {
	var out api.GetBalanceResponse
	if err := c.client.Call(ctx, "GetBalance", &in, &out); err != nil {
		return api.GetBalanceResponse{}, err
	}
	return out, nil
}

// ConvertCurrency performs an HTTP request to convert a certain FIAT currency into credits units.
func (c *client) ConvertCurrency(ctx context.Context, in api.ConvertCurrencyRequest) (api.ConvertCurrencyResponse, error) {
	var out api.ConvertCurrencyResponse
	if err := c.client.Call(ctx, "ConvertCurrency", &in, &out); err != nil {
		return api.ConvertCurrencyResponse{}, err
	}
	return out, nil
}

// Client holds methods to interact with the api.CreditsV1.
type Client interface {
	api.CreditsV1
}

// NewCreditsClientV1 initializes a new api.CreditsV1 client implementation using an HTTP client.
func NewCreditsClientV1(baseURL *url.URL, timeout time.Duration) Client {
	endpoints := map[string]net.EndpointHTTP{
		"IncreaseCredits": {
			Method: http.MethodPost,
			Path:   "/credits/increase",
		},
		"DecreaseCredits": {
			Method: http.MethodPost,
			Path:   "/credits/decrease",
		},
		"GetBalance": {
			Method: http.MethodGet,
			Path:   "/credits",
		},
		"ConvertCurrency": {
			Method: http.MethodPost,
			Path:   "/credits/convert",
		},
	}
	return &client{
		client: net.NewClient(net.NewCallerHTTP(baseURL, endpoints, timeout), encoders.JSON),
	}
}
