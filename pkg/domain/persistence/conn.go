package persistence

import (
	"gitlab.com/ignitionrobotics/billing/credits/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// OpenConn opens a database connection using the config provided from conf.Database.
func OpenConn(config conf.Database) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: config.ToDSN(),
	}))
	if err != nil {
		return nil, err
	}
	return db, nil
}
