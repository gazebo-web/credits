package server

import (
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
)

type setupTestSuite struct {
	suite.Suite
	Logger *log.Logger
}

func TestSetupSuite(t *testing.T) {
	suite.Run(t, new(setupTestSuite))
}

func (s *setupTestSuite) SetupSuite() {
	s.Logger = log.New(os.Stdout, "[TestSetup] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
}

func (s *setupTestSuite) TearDownTest() {
	s.Require().NoError(os.Unsetenv("CREDITS_HTTP_SERVER_PORT"))

	s.Require().NoError(os.Unsetenv("CREDITS_CONVERSION_RATE"))

	s.Require().NoError(os.Unsetenv("CREDITS_DATABASE_USERNAME"))
	s.Require().NoError(os.Unsetenv("CREDITS_DATABASE_PASSWORD"))
	s.Require().NoError(os.Unsetenv("CREDITS_DATABASE_HOST"))
	s.Require().NoError(os.Unsetenv("CREDITS_DATABASE_PORT"))
	s.Require().NoError(os.Unsetenv("CREDITS_DATABASE_NAME"))
	s.Require().NoError(os.Unsetenv("CREDITS_DATABASE_CHARSET"))

}

func (s *setupTestSuite) TestSucceed() {
	s.Require().NoError(os.Setenv("CREDITS_HTTP_SERVER_PORT", "8001"))
	s.Require().NoError(os.Setenv("CREDITS_DATABASE_NAME", "db"))
	s.Require().NoError(os.Setenv("CREDITS_DATABASE_USERNAME", "root"))
	s.Require().NoError(os.Setenv("CREDITS_DATABASE_PASSWORD", "1234"))
	s.Require().NoError(os.Setenv("CREDITS_DATABASE_HOST", "localhost"))
	s.Require().NoError(os.Setenv("CREDITS_DATABASE_PORT", "3306"))
	s.Require().NoError(os.Setenv("CREDITS_DATABASE_CHARSET", "utf16"))
	s.Require().NoError(os.Setenv("CREDITS_CONVERSION_RATE", "2"))

	cfg, err := Setup(s.Logger)
	s.Require().NoError(err)

	// HTTP
	s.Assert().Equal(uint(8001), cfg.Port)

	// Conversion rate
	s.Assert().Equal(uint(2), cfg.ConversionRate)

	// DB
	s.Assert().Equal("root", cfg.Database.Username)
	s.Assert().Equal("1234", cfg.Database.Password)
	s.Assert().Equal("localhost", cfg.Database.Host)
	s.Assert().Equal(uint(3306), cfg.Database.Port)
	s.Assert().Equal("utf16", cfg.Database.Charset)
	s.Assert().Equal("db", cfg.Database.Name)

}

func (s *setupTestSuite) TestDefaultValues() {
	// Define required env vars
	s.Require().NoError(os.Setenv("CREDITS_CONVERSION_RATE", "2"))
	s.Require().NoError(os.Setenv("CREDITS_DATABASE_NAME", "db"))
	s.Require().NoError(os.Setenv("CREDITS_DATABASE_USERNAME", "root"))
	s.Require().NoError(os.Setenv("CREDITS_DATABASE_PASSWORD", "1234"))
	s.Require().NoError(os.Setenv("CREDITS_DATABASE_HOST", "localhost"))
	s.Require().NoError(os.Setenv("CREDITS_DATABASE_PORT", "3306"))

	cfg, err := Setup(s.Logger)
	s.Require().NoError(err)

	s.Assert().Equal(uint(80), cfg.Port)
	s.Assert().Equal("utf8", cfg.Database.Charset)
}

func (s *setupTestSuite) TestMissingEnvVars() {
	_, err := Setup(s.Logger)
	s.Assert().Error(err)
}

func (s *setupTestSuite) TestSetupWithErrors() {
	s.Require().NoError(os.Setenv("PAYMENTS_HTTP_SERVER_PORT", "ABCD"))

	_, err := Setup(s.Logger)
	s.Assert().Error(err)
}
