package application

import (
	"context"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/billing/credits/internal/conf"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/api"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/domain/models"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/domain/persistence"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

type testManageCreditsSuite struct {
	suite.Suite
	DB        *gorm.DB
	Logger    *log.Logger
	Service   Service
	CustomerA models.Customer
	CustomerB models.Customer
	CustomerC models.Customer
}

func TestManageCredits(t *testing.T) {
	suite.Run(t, new(testManageCreditsSuite))
}

func (s *testManageCreditsSuite) SetupSuite() {
	s.Logger = log.New(os.Stdout, "[TestManageCredits] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
	var err error

	var c conf.Config
	s.Require().NoError(c.Parse())

	s.DB, err = persistence.OpenConn(c.Database)
	s.Require().NoError(err)

	s.Require().NoError(persistence.DropTables(s.DB))
}

func (s *testManageCreditsSuite) SetupTest() {
	s.Require().NoError(persistence.MigrateTables(s.DB))
	var err error

	s.Service = NewCreditsService(s.DB, s.Logger, 2)

	s.CustomerA = models.Customer{
		Handle:      "test1",
		Application: "fuel",
		Credits:     100,
	}
	s.CustomerA, err = persistence.CreateCustomer(s.DB, s.CustomerA)
	s.Require().NoError(err)

	s.CustomerB = models.Customer{
		Handle:      "test2",
		Application: "cloudsim",
		Credits:     -100,
	}
	s.CustomerB, err = persistence.CreateCustomer(s.DB, s.CustomerB)
	s.Require().NoError(err)

	s.CustomerC = models.Customer{
		Handle:      "test3",
		Application: "cloudsim",
		Credits:     0,
	}
	s.CustomerC, err = persistence.CreateCustomer(s.DB, s.CustomerC)
	s.Require().NoError(err)
}

func (s *testManageCreditsSuite) TearDownTest() {
	s.Require().NoError(persistence.DropTables(s.DB))
}

func (s *testManageCreditsSuite) TestManageCreditsFailsHandleNotSpecified() {
	_, err := s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "",
			Amount:      10,
			Currency:    "usd",
			Application: "cloudsim",
		},
	})
	s.Assert().Error(err)

	_, err = s.Service.DecreaseCredits(context.Background(), api.DecreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "",
			Amount:      10,
			Currency:    "usd",
			Application: "cloudsim",
		},
	})
	s.Assert().Error(err)
}

func (s *testManageCreditsSuite) TestManageCreditsFailsAmountIsZero() {
	_, err := s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      0,
			Currency:    "usd",
			Application: "cloudsim",
		},
	})
	s.Assert().Error(err)

	_, err = s.Service.DecreaseCredits(context.Background(), api.DecreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      0,
			Currency:    "usd",
			Application: "cloudsim",
		},
	})
	s.Assert().Error(err)
}

func (s *testManageCreditsSuite) TestManageCreditsFailsCurrencyIsInvalid() {
	_, err := s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      10,
			Currency:    "",
			Application: "cloudsim",
		},
	})
	s.Assert().Error(err)

	_, err = s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      10,
			Currency:    "novalid",
			Application: "cloudsim",
		},
	})
	s.Assert().Error(err)

	_, err = s.Service.DecreaseCredits(context.Background(), api.DecreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      10,
			Currency:    "",
			Application: "cloudsim",
		},
	})
	s.Assert().Error(err)

	_, err = s.Service.DecreaseCredits(context.Background(), api.DecreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      10,
			Currency:    "novalid",
			Application: "cloudsim",
		},
	})
	s.Assert().Error(err)
}

func (s *testManageCreditsSuite) TestManageCreditsMissingApplication() {
	_, err := s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      10,
			Currency:    "usd",
			Application: "",
		},
	})
	s.Assert().Error(err)

	_, err = s.Service.DecreaseCredits(context.Background(), api.DecreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      10,
			Currency:    "usd",
			Application: "",
		},
	})
	s.Assert().Error(err)
}

func (s *testManageCreditsSuite) TestManageCreditsNonexistentHandleAndApplicationPair() {
	_, err := s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      10,
			Currency:    "usd",
			Application: "cloudsim",
		},
	})
	s.Assert().Error(err)

	_, err = s.Service.DecreaseCredits(context.Background(), api.DecreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      10,
			Currency:    "usd",
			Application: "cloudsim",
		},
	})
	s.Assert().Error(err)
}

func (s *testManageCreditsSuite) TestIncreaseCreditsConversionApplied() {
	before, err := persistence.GetCustomer(s.DB, "test1", "fuel")
	s.Require().NoError(err)
	s.Assert().Equal(100, before.Credits)

	_, err = s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      10, // Conversion rate: 2 -> 10 usd = 5 credits
			Currency:    "usd",
			Application: "fuel",
		},
	})
	s.Assert().NoError(err)

	after, err := persistence.GetCustomer(s.DB, "test1", "fuel")
	s.Require().NoError(err)
	s.Assert().Equal(before.Credits+5, after.Credits)
}

func (s *testManageCreditsSuite) TestDecreaseCreditsConversionApplied() {
	before, err := persistence.GetCustomer(s.DB, "test1", "fuel")
	s.Require().NoError(err)
	s.Assert().Equal(100, before.Credits)

	_, err = s.Service.DecreaseCredits(context.Background(), api.DecreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      10, // Conversion rate: 2 -> 10 usd = 5 credits
			Currency:    "usd",
			Application: "fuel",
		},
	})
	s.Assert().NoError(err)

	after, err := persistence.GetCustomer(s.DB, "test1", "fuel")
	s.Require().NoError(err)
	s.Assert().Equal(before.Credits-5, after.Credits)
}

func (s *testManageCreditsSuite) TestGetBalanceMissingAttributes() {
	_, err := s.Service.GetBalance(context.Background(), api.GetBalanceRequest{
		Handle:      "",
		Application: "notfound",
	})
	s.Assert().Error(err)
	s.Assert().Equal(api.ErrHandleNotProvided, err)

	_, err = s.Service.GetBalance(context.Background(), api.GetBalanceRequest{
		Handle:      "test5",
		Application: "",
	})
	s.Assert().Error(err)
	s.Assert().Equal(api.ErrMissingApplication, err)
}

func (s *testManageCreditsSuite) TestGetBalanceNotFound() {
	_, err := s.Service.GetBalance(context.Background(), api.GetBalanceRequest{
		Handle:      "test1",
		Application: "notfound",
	})
	s.Assert().Error(err)

	_, err = s.Service.GetBalance(context.Background(), api.GetBalanceRequest{
		Handle:      "test5",
		Application: "fuel",
	})
	s.Assert().Error(err)
}

func (s *testManageCreditsSuite) TestGetBalance() {
	c, err := persistence.GetCustomer(s.DB, "test1", "fuel")
	s.Require().NoError(err)

	balance := c.Credits

	res, err := s.Service.GetBalance(context.Background(), api.GetBalanceRequest{
		Handle:      "test1",
		Application: "fuel",
	})
	s.Require().NoError(err)
	s.Assert().Equal(balance, res.Credits)
}

func (s *testManageCreditsSuite) TestConvertCreditsAmountIsZero() {
	res, err := s.Service.ConvertCurrency(context.Background(), api.ConvertCurrencyRequest{
		Amount:   0,
		Currency: "usd",
	})
	s.Require().NoError(err)
	s.Assert().Equal(uint(0), res.Credits)
}

func (s *testManageCreditsSuite) TestConvertCreditsInvalidCurrency() {
	_, err := s.Service.ConvertCurrency(context.Background(), api.ConvertCurrencyRequest{
		Amount:   10,
		Currency: "",
	})
	s.Assert().Error(err)
	s.Assert().Equal(api.ErrInvalidCurrencyFormat, err)

	_, err = s.Service.ConvertCurrency(context.Background(), api.ConvertCurrencyRequest{
		Amount:   10,
		Currency: "novalid",
	})
	s.Assert().Error(err)
	s.Assert().Equal(api.ErrInvalidCurrencyFormat, err)
}

func (s *testManageCreditsSuite) TestConvertCreditsConversionRateApplied() {
	res, err := s.Service.ConvertCurrency(context.Background(), api.ConvertCurrencyRequest{
		Amount:   10,
		Currency: "usd",
	})
	s.Require().NoError(err)

	s.Assert().Equal(uint(5), res.Credits)
}

func (s *testManageCreditsSuite) TestIncreaseCreditsToZero() {
	before, err := persistence.GetCustomer(s.DB, "test2", "cloudsim")
	s.Require().NoError(err)
	s.Assert().Equal(-100, before.Credits)

	_, err = s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test2",
			Amount:      200, // Conversion rate: 2 -> 200 usd = 100 credits
			Currency:    "usd",
			Application: "cloudsim",
		},
	})
	s.Assert().NoError(err)

	after, err := persistence.GetCustomer(s.DB, "test2", "cloudsim")
	s.Require().NoError(err)
	s.Assert().Equal(0, after.Credits)
}

func (s *testManageCreditsSuite) TestDecreaseCreditsToZero() {
	before, err := persistence.GetCustomer(s.DB, "test1", "fuel")
	s.Require().NoError(err)
	s.Assert().Equal(100, before.Credits)

	_, err = s.Service.DecreaseCredits(context.Background(), api.DecreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Amount:      200, // Conversion rate: 2 -> 200 usd = 100 credits
			Currency:    "usd",
			Application: "fuel",
		},
	})
	s.Assert().NoError(err)

	after, err := persistence.GetCustomer(s.DB, "test1", "fuel")
	s.Require().NoError(err)
	s.Assert().Equal(0, after.Credits)
}

func (s *testManageCreditsSuite) TestGetUnitPriceValidationFails() {
	_, err := s.Service.GetUnitPrice(context.Background(), api.GetUnitPriceRequest{Currency: ""})
	s.Assert().Error(err)
}

func (s *testManageCreditsSuite) TestGetUnitPrice() {
	const currency = "usd"
	res, err := s.Service.GetUnitPrice(context.Background(), api.GetUnitPriceRequest{Currency: currency})
	s.Require().NoError(err)

	// The amount should be the conversion rate
	s.Assert().Equal(uint(2), res.Amount)
	s.Assert().Equal(currency, res.Currency)
}
