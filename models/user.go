package models

import "gorm.io/gorm"


// Define role constants
const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

// @title User Struct
// @notice Represents a user account in the application.
// This structure maps directly to the 'users' table in the database and is used by GORM.
type User struct {
	// @dev gorm.Model is embedded to provide standard ID, CreatedAt, UpdatedAt, and DeletedAt fields.
	gorm.Model

	// @notice The full name of the user.
	Name string

	// @notice The primary identifier for the user. Must be unique across all records.
	// @custom:gorm:unique ensures database-level uniqueness.
	Email string `gorm:"unique"`

	// @notice The hashed password for the user's account.
	// @dev Note: Lowercase 'password' makes this field unexported, limiting access outside the package.
	Password string

	// @notice The access level or role of the user (e.g., USER, ADMIN).
	// @custom:gorm:default sets the initial value; NOT NULL ensures it is always present.
	Role string `gorm:"default:USER; NOT NULL"`

	Links []Link
}