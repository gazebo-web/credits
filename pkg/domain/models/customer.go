package models

import "gorm.io/gorm"

// Customer is a single entity in a certain application that wants to use credits to consume different services.
// The definition of an entity as a customer is up to the applications.
// This model only serves as a simple credit store, any relationship between customers and different entities
// (e.g. users and organizations) should be handled by applications.
type Customer struct {
	gorm.Model

	// Handle contains the customer handle. This handle is specific to the Application.
	Handle string

	// Application is the application that the credits are being tracked for.
	Application string

	// Credits is the amount of credits this Customer can use in services provided by Application.
	Credits int
}
