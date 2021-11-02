package persistence

import (
	"gitlab.com/ignitionrobotics/billing/credits/pkg/domain/models"
	"gorm.io/gorm"
)

// MigrateTables migrates all the model tables.
func MigrateTables(db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&models.Customer{},
	)
}

// DropTables drops all the model tables.
func DropTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.Customer{},
	)
}
