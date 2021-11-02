package persistence

import (
	"gitlab.com/ignitionrobotics/billing/credits/pkg/domain/models"
	"gorm.io/gorm"
)

// CreateCustomer creates a new customer.
func CreateCustomer(db *gorm.DB, customer models.Customer) (models.Customer, error) {
	if err := db.Model(&models.Customer{}).Create(&customer).Error; err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

// UpdateCredits increases or decreases a certain amount of credits to a specific customer given by its handle
// for the given application.
// This operation will not create a new `Customer` if it is not found.
func UpdateCredits(db *gorm.DB, handle, application string, value int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var c models.Customer

		err := tx.
			Model(&models.Customer{}).Where("handle = ? AND application = ?", handle, application).
			First(&c).Error
		if err != nil {
			return err
		}

		result := tx.
			Model(&c).
			Updates(map[string]interface{}{"credits": c.Credits + value})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

// GetCustomer returns a customer based on the given handle and application.
func GetCustomer(db *gorm.DB, handle, application string) (models.Customer, error) {
	var result models.Customer
	err := db.Model(&models.Customer{}).
		Where("handle = ? AND application = ?", handle, application).
		First(&result).Error
	if err != nil {
		return models.Customer{}, err
	}
	return result, nil
}
