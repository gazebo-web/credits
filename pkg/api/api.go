package api

import "context"

// CreditsV1 holds the methods that allow managing user credits.
type CreditsV1 interface {
	// IncreaseCredits increases the amount of credits for a given user.
	IncreaseCredits(ctx context.Context, req IncreaseCreditsRequest) (IncreaseCreditsResponse, error)

	// DecreaseCredits decreases the amount of credits for a given user.
	DecreaseCredits(ctx context.Context, req DecreaseCreditsRequest) (DecreaseCreditsResponse, error)

	// GetBalance returns the current amount of credits of a given user.
	GetBalance(ctx context.Context, req GetBalanceRequest) (GetBalanceResponse, error)

	// ConvertCurrency converts a certain amount of FIAT currency in USD to credits.
	ConvertCurrency(ctx context.Context, req ConvertCurrencyRequest) (ConvertCurrencyResponse, error)
}

type IncreaseCreditsRequest struct{}

type IncreaseCreditsResponse struct{}

type DecreaseCreditsRequest struct{}

type DecreaseCreditsResponse struct{}

type GetBalanceRequest struct{}

type GetBalanceResponse struct{}

type ConvertCurrencyRequest struct{}

type ConvertCurrencyResponse struct{}
