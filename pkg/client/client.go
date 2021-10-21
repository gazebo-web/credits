package client

import (
	"context"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/api"
)

// client contains the HTTP client to connect to the credits API.
type client struct{}

// IncreaseCredits performs an HTTP request to increase the given user credits.
func (c *client) IncreaseCredits(ctx context.Context, req api.IncreaseCreditsRequest) (api.IncreaseCreditsResponse, error) {
	panic("implement me")
}

// DecreaseCredits performs an HTTP request to decrease the given user credits.
func (c *client) DecreaseCredits(ctx context.Context, req api.DecreaseCreditsRequest) (api.DecreaseCreditsResponse, error) {
	panic("implement me")
}

// GetBalance performs an HTTP request to get the user's balance.
func (c *client) GetBalance(ctx context.Context, req api.GetBalanceRequest) (api.GetBalanceResponse, error) {
	panic("implement me")
}

// ConvertCurrency performs an HTTP request to convert a certain FIAT currency into credits units.
func (c *client) ConvertCurrency(ctx context.Context, req api.ConvertCurrencyRequest) (api.ConvertCurrencyResponse, error) {
	panic("implement me")
}

// Client holds methods to interact with the api.CreditsV1.
type Client interface {
	api.CreditsV1
}

// NewClient initializes a new api.CreditsV1 client implementation using an HTTP client.
func NewClient() Client {
	return &client{}
}
