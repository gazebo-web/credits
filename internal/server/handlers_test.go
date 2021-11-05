package server

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/billing/credits/internal/conf"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/api"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/application"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/domain/models"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/domain/persistence"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type handlersTestSuite struct {
	suite.Suite
	Logger           *log.Logger
	DB               *gorm.DB
	Service          application.Service
	Server           *Server
	Handler          http.HandlerFunc
	CustomerA        models.Customer
	CustomerB        models.Customer
	CustomerC        models.Customer
	ResponseRecorder *httptest.ResponseRecorder
}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(handlersTestSuite))
}

func (s *handlersTestSuite) SetupSuite() {
	s.Logger = log.New(os.Stdout, "[TestHandlers] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)

	var c conf.Config
	s.Require().NoError(c.Parse())

	var err error
	s.DB, err = persistence.OpenConn(c.Database)
	s.Require().NoError(err)

	s.Service = application.NewCreditsService(s.DB, s.Logger, 2)

	s.Server = NewServer(Options{
		config:  c,
		credits: s.Service,
		logger:  s.Logger,
	})
}

func (s *handlersTestSuite) SetupTest() {
	s.Require().NoError(persistence.MigrateTables(s.DB))

	var err error
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

	s.ResponseRecorder = httptest.NewRecorder()
}

func (s *handlersTestSuite) TearDownTest() {
	s.Require().NoError(persistence.DropTables(s.DB))
}

func (s *handlersTestSuite) TearDownSuite() {
	db, err := s.DB.DB()
	s.Require().NoError(err)
	s.Require().NoError(db.Close())
}

func (s *handlersTestSuite) TestGetBalanceOK() {
	s.Handler = s.Server.GetBalance

	in := api.GetBalanceRequest{
		Handle:      "test1",
		Application: "fuel",
	}
	request := s.setupRequest(in, http.MethodGet)

	s.Handler.ServeHTTP(s.ResponseRecorder, request)

	s.Require().Equal(http.StatusOK, s.ResponseRecorder.Code)

	var out api.GetBalanceResponse
	s.parseResponseJSON(&out)

	s.Assert().Equal(100, out.Credits)
}

func (s *handlersTestSuite) TestIncreaseCreditsOK() {
	s.Handler = s.Server.IncreaseCredits

	before, err := persistence.GetCustomer(s.DB, "test2", "cloudsim")
	s.Require().NoError(err)

	in := api.IncreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test2",
			Application: "cloudsim",
			Amount:      200,
			Currency:    "usd",
		},
	}
	request := s.setupRequest(in, http.MethodGet)

	s.Handler.ServeHTTP(s.ResponseRecorder, request)

	s.Require().Equal(http.StatusOK, s.ResponseRecorder.Code)

	var out api.IncreaseCreditsResponse
	s.parseResponseJSON(&out)

	after, err := persistence.GetCustomer(s.DB, "test2", "cloudsim")
	s.Require().NoError(err)

	s.Assert().Equal(0, after.Credits)
	s.Assert().Equal(before.Credits+100, after.Credits)
}

func (s *handlersTestSuite) TestDecreaseCreditsOK() {
	s.Handler = s.Server.DecreaseCredits

	before, err := persistence.GetCustomer(s.DB, "test1", "fuel")
	s.Require().NoError(err)

	in := api.DecreaseCreditsRequest{
		Transaction: api.Transaction{
			Handle:      "test1",
			Application: "fuel",
			Amount:      200,
			Currency:    "usd",
		},
	}
	request := s.setupRequest(in, http.MethodGet)

	s.Handler.ServeHTTP(s.ResponseRecorder, request)

	s.Require().Equal(http.StatusOK, s.ResponseRecorder.Code)

	var out api.DecreaseCreditsResponse
	s.parseResponseJSON(&out)

	after, err := persistence.GetCustomer(s.DB, "test1", "fuel")
	s.Require().NoError(err)

	s.Assert().Equal(0, after.Credits)
	s.Assert().Equal(before.Credits-100, after.Credits)
}

func (s *handlersTestSuite) TestConvertCurrencyOK() {
	s.Handler = s.Server.ConvertCurrency

	in := api.ConvertCurrencyRequest{
		Amount:   100,
		Currency: "usd",
	}

	request := s.setupRequest(in, http.MethodGet)

	s.Handler.ServeHTTP(s.ResponseRecorder, request)

	s.Require().Equal(http.StatusOK, s.ResponseRecorder.Code)

	var out api.ConvertCurrencyResponse
	s.parseResponseJSON(&out)

	s.Assert().Equal(uint(50), out.Credits)
}

func (s *handlersTestSuite) setupRequest(in interface{}, method string) *http.Request {
	body, err := json.Marshal(in)
	s.Require().NoError(err)

	buff := bytes.NewBuffer(body)

	request, err := http.NewRequest(method, "/", buff)
	s.Require().NoError(err)

	return request
}

func (s *handlersTestSuite) parseResponseJSON(out interface{}) {
	body, err := io.ReadAll(s.ResponseRecorder.Body)
	s.Require().NoError(err)
	s.Require().NoError(json.Unmarshal(body, out))
}
