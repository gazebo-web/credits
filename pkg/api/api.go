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

// IncreaseCreditsRequest is the input for the CreditsV1.IncreaseCredits method.
type IncreaseCreditsRequest struct {
	// User is the username of the User that will receive the credits.
	User string

	// Amount is the money that the user paid in cents.
	Amount uint

	// Currency is the currency the user paid in.
	Currency string
}

// IncreaseCreditsResponse is the output of the CreditsV1.IncreaseCredits method.
type IncreaseCreditsResponse struct{}

// DecreaseCreditsRequest is the input for the CreditsV1.DecreaseCredits method.
type DecreaseCreditsRequest struct{}

// DecreaseCreditsResponse is the output of the CreditsV1.DecreaseCredits method.
type DecreaseCreditsResponse struct{}

// GetBalanceRequest is the input for the CreditsV1.GetBalance method.
type GetBalanceRequest struct{}

// GetBalanceResponse is the output of the CreditsV1.GetBalance method.
type GetBalanceResponse struct{}

// ConvertCurrencyRequest is the input for the CreditsV1.ConvertCurrency method.
type ConvertCurrencyRequest struct{}

// ConvertCurrencyResponse is the output of the CreditsV1.ConvertCurrency method.
type ConvertCurrencyResponse struct{}
