package fake

import (
	"context"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/api"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/client"
)

// mock is a fake client.Client implementation.
type mock struct{}

// IncreaseCredits mocks a valid call to the Credits API.
func (m *mock) IncreaseCredits(ctx context.Context, req api.IncreaseCreditsRequest) (api.IncreaseCreditsResponse, error) {
	return api.IncreaseCreditsResponse{}, nil
}

// DecreaseCredits mocks a valid call to the Credits API.
func (m *mock) DecreaseCredits(ctx context.Context, req api.DecreaseCreditsRequest) (api.DecreaseCreditsResponse, error) {
	return api.DecreaseCreditsResponse{}, nil
}

// GetBalance mocks a valid call to the Credits API.
func (m *mock) GetBalance(ctx context.Context, req api.GetBalanceRequest) (api.GetBalanceResponse, error) {
	return api.GetBalanceResponse{}, nil
}

// ConvertCurrency mocks a valid call to the Credits API.
func (m *mock) ConvertCurrency(ctx context.Context, req api.ConvertCurrencyRequest) (api.ConvertCurrencyResponse, error) {
	return api.ConvertCurrencyResponse{}, nil
}

// NewClient initializes a fake client.Client implementation.
func NewClient() client.Client {
	return &mock{}
}
