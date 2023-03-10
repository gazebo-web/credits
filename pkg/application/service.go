package application

import (
	"context"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/api"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/domain/persistence"
	"gorm.io/gorm"
	"io"
	"log"
	"math"
)

// service contains the business logic to manage credits.
type service struct {
	logger         *log.Logger
	db             *gorm.DB
	conversionRate uint
}

// GetUnitPrice returns the value of how much a credit costs.
func (s *service) GetUnitPrice(ctx context.Context, req api.GetUnitPriceRequest) (api.GetUnitPriceResponse, error) {
	// TODO: Add check for valid currency values as well as lowercase.
	if len(req.Currency) == 0 || len(req.Currency) > 3 {
		return api.GetUnitPriceResponse{}, api.ErrInvalidCurrencyFormat
	}
	return api.GetUnitPriceResponse{
		Amount:   s.conversionRate,
		Currency: "usd",
	}, nil
}

// IncreaseCredits increases the amount of service for a given user.
func (s *service) IncreaseCredits(ctx context.Context, req api.IncreaseCreditsRequest) (api.IncreaseCreditsResponse, error) {
	if err := req.Validate(); err != nil {
		return api.IncreaseCreditsResponse{}, err
	}

	value := s.calculateCredits(req.Amount, req.Currency)

	if err := persistence.UpdateCredits(s.db, req.Handle, req.Application, int(value)); err != nil {
		return api.IncreaseCreditsResponse{}, err
	}

	return api.IncreaseCreditsResponse{}, nil
}

// DecreaseCredits decreases the amount of service for a given user.
func (s *service) DecreaseCredits(ctx context.Context, req api.DecreaseCreditsRequest) (api.DecreaseCreditsResponse, error) {
	if err := req.Validate(); err != nil {
		return api.DecreaseCreditsResponse{}, err
	}

	value := s.calculateCredits(req.Amount, req.Currency)

	if err := persistence.UpdateCredits(s.db, req.Handle, req.Application, -1*int(value)); err != nil {
		return api.DecreaseCreditsResponse{}, err
	}

	return api.DecreaseCreditsResponse{}, nil
}

// GetBalance returns the current amount of service of a given user.
func (s *service) GetBalance(ctx context.Context, req api.GetBalanceRequest) (api.GetBalanceResponse, error) {
	if len(req.Handle) == 0 {
		s.logger.Println("No handle provided")
		return api.GetBalanceResponse{}, api.ErrHandleNotProvided
	}
	if len(req.Application) == 0 {
		s.logger.Println("Missing application")
		return api.GetBalanceResponse{}, api.ErrMissingApplication
	}
	c, err := persistence.GetCustomer(s.db, req.Handle, req.Application)
	if err != nil {
		return api.GetBalanceResponse{}, err
	}

	return api.GetBalanceResponse{
		Handle:      c.Handle,
		Application: c.Application,
		Credits:     c.Credits,
	}, nil
}

// ConvertCurrency converts a certain amount of FIAT currency in USD to service.
func (s *service) ConvertCurrency(ctx context.Context, req api.ConvertCurrencyRequest) (api.ConvertCurrencyResponse, error) {
	if len(req.Currency) == 0 || len(req.Currency) > 3 {
		s.logger.Println("Invalid currency format")
		return api.ConvertCurrencyResponse{}, api.ErrInvalidCurrencyFormat
	}
	return api.ConvertCurrencyResponse{
		Credits: s.calculateCredits(req.Amount, req.Currency),
	}, nil
}

// calculateCredits applies the conversion rate to amount in a certain currency and returns a credits value.
// It rounds up the output value to the closest integer.
func (s *service) calculateCredits(amount uint, currency string) uint {
	return uint(math.Ceil(float64(amount) / float64(s.conversionRate)))
}

// Service holds the methods of the service in charge of managing user credits.
type Service interface {
	api.CreditsV1
}

// NewCreditsService initializes a new api.CreditsV1 service implementation.
func NewCreditsService(db *gorm.DB, logger *log.Logger, rate uint) Service {
	if logger == nil {
		logger = log.New(io.Discard, "", log.LstdFlags)
	}
	return &service{
		db:             db,
		logger:         logger,
		conversionRate: rate,
	}
}
