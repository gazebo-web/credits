package application

import (
	"context"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/api"
	"io"
	"log"
)

// credits contains the business logic to manage credits.
type credits struct {
	logger *log.Logger
}

// IncreaseCredits increases the amount of credits for a given user.
func (s *credits) IncreaseCredits(ctx context.Context, req api.IncreaseCreditsRequest) (api.IncreaseCreditsResponse, error) {
	panic("implement me")
}

// DecreaseCredits decreases the amount of credits for a given user.
func (s *credits) DecreaseCredits(ctx context.Context, req api.DecreaseCreditsRequest) (api.DecreaseCreditsResponse, error) {
	panic("implement me")
}

// GetBalance returns the current amount of credits of a given user.
func (s *credits) GetBalance(ctx context.Context, req api.GetBalanceRequest) (api.GetBalanceResponse, error) {
	panic("implement me")
}

// ConvertCurrency converts a certain amount of FIAT currency in USD to credits.
func (s *credits) ConvertCurrency(ctx context.Context, req api.ConvertCurrencyRequest) (api.ConvertCurrencyResponse, error) {
	panic("implement me")
}

// Service holds the methods of the service in charge of managing user credits.
type Service interface {
	api.CreditsV1
}

// NewService initializes a new api.CreditsV1 service implementation.
func NewService(logger *log.Logger) Service {
	if logger == nil {
		logger = log.New(io.Discard, "", log.LstdFlags)
	}
	return &credits{
		logger: logger,
	}
}
