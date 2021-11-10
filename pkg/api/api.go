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

	// GetUnitPrice returns the amount of currency needed to buy 1 credit.
	GetUnitPrice(ctx context.Context, req GetUnitPriceRequest) (GetUnitPriceResponse, error)
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
	Handle string `json:"handle"`

	// Amount is the money that the user paid in the minimum currency value (e.g. cents for USD) that should be converted
	// to credits.
	Amount uint `json:"amount"`

	// Currency is the ISO 4217 currency code in lowercase format.
	Currency string `json:"currency"`

	// Application is the application that credits are tracked for.
	Application string `json:"application"`
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
type GetBalanceRequest struct {
	// Handle is the username of the customer that should receive the balance summary.
	Handle string `json:"handle"`

	// Application is the application that credits are tracked for.
	Application string `json:"application"`
}

// GetBalanceResponse is the output of the CreditsV1.GetBalance method.
type GetBalanceResponse struct {
	// Handle is the username of the customer that is receiving the balance summary.
	Handle string `json:"handle"`

	// Application is the application that credits are tracked for.
	Application string `json:"application"`

	// Credits is the amount of credits that the customer identified by Handle has.
	Credits int `json:"credits"`
}

// ConvertCurrencyRequest is the input for the CreditsV1.ConvertCurrency method.
type ConvertCurrencyRequest struct {
	// Amount is the money in the minimum currency value (e.g. cents for USD) that should be converted
	// to credits.
	Amount uint `json:"amount"`

	// Currency is the ISO 4217 currency code in lowercase format.
	Currency string `json:"currency"`
}

// ConvertCurrencyResponse is the output of the CreditsV1.ConvertCurrency method.
type ConvertCurrencyResponse struct {
	// Credits contains the result of converting a certain currency value into credits.
	Credits uint `json:"credits"`
}

// GetUnitPriceRequest is the input for the CreditsV1.GetUnitPrice method.
type GetUnitPriceRequest struct {
	// Currency is the ISO 4217 currency code in lowercase format.
	Currency string `json:"currency"`
}

// GetUnitPriceResponse is the output of the CreditsV1.GetUnitPrice method.
type GetUnitPriceResponse struct {
	// Amount is the money in the minimum currency value (e.g. cents for USD) of how much a credit cost.
	Amount uint `json:"amount"`

	// Currency is the ISO 4217 currency code in lowercase format.
	Currency string `json:"currency"`
}
