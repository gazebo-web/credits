package api

import (
	"context"
	"errors"
)

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

var (
	// ErrHandleNotProvided is returned when the handler is not provided in the request.
	ErrHandleNotProvided = errors.New("handler not provided")
	// ErrInvalidAmount is returned when an invalid amount is passed in the request.
	ErrInvalidAmount = errors.New("invalid amount")
	// ErrInvalidCurrencyFormat is returned when the currency format is invalid.
	ErrInvalidCurrencyFormat = errors.New("invalid currency format")
	// ErrMissingApplication is returned when there's no application defined in a request.
	ErrMissingApplication = errors.New("missing application")
)

// Transaction is an operation made with credits. It's usually used to increase and decrease the amount of credits of certain models.Customer.
type Transaction struct {
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

// Validate validates the current transaction is valid.
func (t Transaction) Validate() error {
	if len(t.Handle) == 0 {
		return ErrHandleNotProvided
	}
	if t.Amount == 0 {
		return ErrInvalidAmount
	}
	// TODO: Add check for valid currency values as well as lowercase.
	if len(t.Currency) == 0 || len(t.Currency) > 3 {
		return ErrInvalidCurrencyFormat
	}
	if len(t.Application) == 0 {
		return ErrMissingApplication
	}
	return nil
}

// IncreaseCreditsRequest is the input for the CreditsV1.IncreaseCredits method.
type IncreaseCreditsRequest struct {
	Transaction
}

// IncreaseCreditsResponse is the output of the CreditsV1.IncreaseCredits method.
type IncreaseCreditsResponse struct{}

// DecreaseCreditsRequest is the input for the CreditsV1.DecreaseCredits method.
type DecreaseCreditsRequest struct {
	Transaction
}

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
