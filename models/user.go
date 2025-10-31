package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the user model in the database
// It contains all user-related information including authentication data
type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`                                // MongoDB unique identifier
	FirstName     *string            `json:"first_name" validate:"required,min=2,max=100"` // User's first name
	LastName      *string            `json:"last_name" validate:"required,min=2,max=100"`  // User's last name
	Email         *string            `json:"email" validate:"required,email"`              // User's email (unique)
	Password      *string            `json:"password" validate:"required,min=6"`           // Hashed password
	Phone         *string            `json:"phone" validate:"required"`                    // Phone number (unique)
	Token         *string            `json:"token,omitempty"`                              // JWT access token
	Role          *string            `json:"role" validate:"required,eq=ADMIN|eq=USER"`    // User role (ADMIN or USER)
	Refresh_token *string            `json:"refresh_token,omitempty"`                      // JWT refresh token
	Created_at    time.Time          `json:"created_at"`                                   // Account creation timestamp
	Updated_at    time.Time          `json:"updated_at"`                                   // Last update timestamp
	User_id       string             `json:"user_id"`                                      // String representation of user ID
}
