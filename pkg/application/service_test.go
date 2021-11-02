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

type testIncreaseCreditsSuite struct {
	suite.Suite
	DB      *gorm.DB
	Logger  *log.Logger
	Service Service
	UserA   models.Customer
	UserB   models.Customer
	UserC   models.Customer
}

func TestIncreaseCredits(t *testing.T) {
	suite.Run(t, new(testIncreaseCreditsSuite))
}

func (s *testIncreaseCreditsSuite) SetupSuite() {
	s.Logger = log.New(os.Stdout, "[TestIncreaseCredits] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
}

func (s *testIncreaseCreditsSuite) SetupTest() {
	var err error

	var c conf.Config
	s.Require().NoError(c.Parse())

	s.DB, err = persistence.OpenConn(c.Database)
	s.Require().NoError(err)

	s.Require().NoError(persistence.MigrateTables(s.DB))

	s.Service = NewService(s.DB, s.Logger, 2)

	s.UserA = models.Customer{
		Model:       gorm.Model{},
		Handle:      "test1",
		Application: "fuel",
		Credits:     100,
	}
	s.UserA, err = persistence.CreateCustomer(s.DB, s.UserA)

	s.UserB = models.Customer{
		Model:       gorm.Model{},
		Handle:      "test2",
		Application: "cloudsim",
		Credits:     -100,
	}
	s.UserB, err = persistence.CreateCustomer(s.DB, s.UserB)

	s.UserC = models.Customer{
		Model:       gorm.Model{},
		Handle:      "test3",
		Application: "cloudsim",
		Credits:     0,
	}
	s.UserB, err = persistence.CreateCustomer(s.DB, s.UserB)
}

func (s *testIncreaseCreditsSuite) TestIncreaseCreditsFailsHandleNotSpecified() {
	_, err := s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Handle:      "",
		Amount:      10,
		Currency:    "usd",
		Application: "cloudsim",
	})
	s.Assert().Error(err)
}

func (s *testIncreaseCreditsSuite) TestIncreaseCreditsFailsAmountIsZero() {
	_, err := s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Handle:      "test1",
		Amount:      0,
		Currency:    "usd",
		Application: "cloudsim",
	})
	s.Assert().Error(err)
}

func (s *testIncreaseCreditsSuite) TestIncreaseCreditsFailsCurrencyIsInvalid() {
	_, err := s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Handle:      "test1",
		Amount:      10,
		Currency:    "",
		Application: "cloudsim",
	})
	s.Assert().Error(err)

	_, err = s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Handle:      "test1",
		Amount:      10,
		Currency:    "novalid",
		Application: "cloudsim",
	})
	s.Assert().Error(err)
}

func (s *testIncreaseCreditsSuite) TestIncreaseCreditsMissingApplication() {
	_, err := s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Handle:      "test1",
		Amount:      10,
		Currency:    "usd",
		Application: "",
	})
	s.Assert().Error(err)
}

func (s *testIncreaseCreditsSuite) TestIncreaseCreditsNonexistentHandleAndApplicationPair() {
	_, err := s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Handle:      "test1",
		Amount:      10,
		Currency:    "usd",
		Application: "cloudsim",
	})
	s.Assert().Error(err)
}

func (s *testIncreaseCreditsSuite) TestIncreaseCreditsConversionApplied() {
	before, err := persistence.GetCustomer(s.DB, "test1", "fuel")
	s.Require().NoError(err)
	s.Assert().Equal(100, before.Credits)

	_, err = s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Handle:      "test1",
		Amount:      10, // Conversion rate: 2 -> 10 usd = 5 credits
		Currency:    "usd",
		Application: "fuel",
	})
	s.Assert().NoError(err)

	after, err := persistence.GetCustomer(s.DB, "test1", "fuel")
	s.Require().NoError(err)
	s.Assert().Equal(before.Credits+5, after.Credits)
}

func (s *testIncreaseCreditsSuite) TestIncreaseCreditsToZero() {
	before, err := persistence.GetCustomer(s.DB, "test2", "cloudsim")
	s.Require().NoError(err)
	s.Assert().Equal(-100, before.Credits)

	_, err = s.Service.IncreaseCredits(context.Background(), api.IncreaseCreditsRequest{
		Handle:      "test2",
		Amount:      200, // Conversion rate: 2 -> 200 usd = 100 credits
		Currency:    "usd",
		Application: "cloudsim",
	})
	s.Assert().NoError(err)

	after, err := persistence.GetCustomer(s.DB, "test2", "cloudsim")
	s.Require().NoError(err)
	s.Assert().Equal(0, after.Credits)
}

func (s *testIncreaseCreditsSuite) TearDownTest() {
	s.Require().NoError(persistence.DropTables(s.DB))
}
