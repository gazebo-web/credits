package fake

import (
	"context"
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/api"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/client"
)

var _ client.Client = (*Fake)(nil)

// Fake is a fake client.Client implementation.
type Fake struct {
	mock.Mock
}

// GetUnitPrice mocks a call to the Credits API.
func (c *Fake) GetUnitPrice(ctx context.Context, req api.GetUnitPriceRequest) (api.GetUnitPriceResponse, error) {
	args := c.Called(ctx, req)
	res := args.Get(0).(api.GetUnitPriceResponse)
	return res, args.Error(1)
}

// IncreaseCredits mocks a call to the Credits API.
func (c *Fake) IncreaseCredits(ctx context.Context, req api.IncreaseCreditsRequest) (api.IncreaseCreditsResponse, error) {
	args := c.Called(ctx, req)
	res := args.Get(0).(api.IncreaseCreditsResponse)
	return res, args.Error(1)
}

// DecreaseCredits mocks a call to the Credits API.
func (c *Fake) DecreaseCredits(ctx context.Context, req api.DecreaseCreditsRequest) (api.DecreaseCreditsResponse, error) {
	args := c.Called(ctx, req)
	res := args.Get(0).(api.DecreaseCreditsResponse)
	return res, args.Error(1)
}

// GetBalance mocks a valid to the Credits API.
func (c *Fake) GetBalance(ctx context.Context, req api.GetBalanceRequest) (api.GetBalanceResponse, error) {
	args := c.Called(ctx, req)
	res := args.Get(0).(api.GetBalanceResponse)
	return res, args.Error(1)
}

// ConvertCurrency mocks a call to the Credits API.
func (c *Fake) ConvertCurrency(ctx context.Context, req api.ConvertCurrencyRequest) (api.ConvertCurrencyResponse, error) {
	args := c.Called(ctx, req)
	res := args.Get(0).(api.ConvertCurrencyResponse)
	return res, args.Error(1)
}

// NewClient initializes a fake client.Client implementation.
func NewClient() *Fake {
	return &Fake{}
}
