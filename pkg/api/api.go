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
	// Handle is the username of the customer that will receive the credits.
	Handle string

	// Amount is the money that the user paid in the minimum currency value (e.g. cents for USD) that should be converted
	// to credits.
	Amount uint

	// Currency is the ISO 4217 currency code in lowercase format.
	Currency string

	// Application is the application that credits are tracked for.
	Application string
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
