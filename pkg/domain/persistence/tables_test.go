package persistence

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/billing/credits/internal/conf"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/domain/models"
	"gorm.io/gorm"
	"testing"
)

func TestTables(t *testing.T) {
	suite.Run(t, new(testTablesSuite))
}

type testTablesSuite struct {
	suite.Suite
	DB *gorm.DB
}

func (s *testTablesSuite) SetupTest() {
	var cfg conf.Database
	s.Require().NoError(cfg.Parse())

	var err error
	s.DB, err = OpenConn(cfg)
	s.Require().NoError(err)
}

func (s *testTablesSuite) TearDownTest() {
	_ = DropTables(s.DB)
}

func (s *testTablesSuite) TestMigrateTables() {
	s.Require().False(s.DB.Migrator().HasTable(&models.Customer{}))
	s.Assert().NoError(MigrateTables(s.DB))
	s.Assert().True(s.DB.Migrator().HasTable(&models.Customer{}))
}

func (s *testTablesSuite) TestDropTables() {
	s.Require().NoError(MigrateTables(s.DB))
	s.Assert().NoError(DropTables(s.DB))
}
