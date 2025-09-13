package module

import "time"

// Module represents a module entity in the system.
//
// This model is used across all layers of the application.
// It contains the core properties of a module that are persisted to the database.
//
// Example:
//
//	{
//	  "id": 123,
//	  "name": "Inventory",
//	  "description": "Handles product stock management",
//	  "isActive": true,
//	  "createdAt": "2023-08-15T14:30:00Z"
//	}
type Module struct {
	// Unique identifier for the module
	ID int `json:"id" gorm:"primaryKey"`

	// Name of the module (3-50 characters, required)
	// Business Rule: Must be unique across active modules
	Name string `json:"name" gorm:"size:50;not null;uniqueIndex:idx_name_active"`

	// Description of what the module does (max 200 characters)
	Description string `json:"description" gorm:"size:200"`

	// Indicates if the module is currently active
	IsActive bool `json:"isActive" gorm:"default:true"`

	// Timestamp when the module was created
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// ModuleRequest represents the payload for creating a new module.
//
// This DTO is used by the presentation layer to validate incoming requests.
//
// Example:
//
//	{
//	  "name": "Inventory",
//	  "description": "Handles product stock management",
//	  "isActive": true
//	}
type ModuleRequest struct {
	// Name of the module (3-50 characters, required)
	// Validation: Must be 3-50 characters, alphanumeric with spaces
	Name string `json:"name" binding:"required,min=3,max=50"`

	// Description of what the module does (max 200 characters)
	Description string `json:"description" binding:"max=200"`

	// Indicates if the module should be active upon creation
	IsActive bool `json:"isActive"`
}

// ModuleResponse represents the response structure for module operations.
//
// This DTO is used to format responses from the API.
//
// Example:
//
//	{
//	  "id": 123,
//	  "name": "Inventory",
//	  "description": "Handles product stock management",
//	  "isActive": true,
//	  "createdAt": "2023-08-15T14:30:00Z"
//	}
type ModuleResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
}
