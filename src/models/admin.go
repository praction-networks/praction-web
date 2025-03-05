package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Admin struct with mobile validation
type Admin struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Username        string             `json:"username" bson:"username" validate:"required,min=3,max=20"`                     // Required, minimum 3, and maximum 20 characters
	Password        string             `json:"password" bson:"password" validate:"required,min=8"`                            // Required, minimum 8 characters
	ConfirmPassword string             `json:"confirmPassword" bson:"_,omitempty" validate:"required,min=8,eqfield=Password"` // Required same as password
	Salt            string             `json:"saltstr" bson:"saltstr"`                                                        // No validation, managed internally
	Mobile          string             `json:"mobile" bson:"mobile" validate:"required,len=10"`                               // Validation for 10 digits starting with 6, 7, 8, or 9
	Email           string             `json:"email" bson:"email" validate:"required,email"`                                  // Required and must be a valid email
	FirstName       string             `json:"first_name" bson:"first_name" validate:"required"`                              // Required
	LastName        string             `json:"last_name" bson:"last_name" validate:"required"`                                // Required
	Role            string             `json:"role" bson:"role" validate:"required,oneof=admin user manager"`                 // Required, must be either "admin", "user" or "manager"
	PasswordToken   string             `bson:"passwordToken,omitempty"`
}

type ResponseAdmin struct {
	ID        primitive.ObjectID `json:"id"`
	Username  string             `json:"username"`
	Mobile    string             `json:"mobile"`     // Required, minimum 3, and maximum 20 characters
	Email     string             `json:"email"`      // Required and must be a valid email
	FirstName string             `json:"first_name"` // Required
	LastName  string             `json:"last_name"`  // Required
	Role      string             `json:"role"`       // Required, must be either "admin", "user" or "manager"
}

type UpdateAdmin struct {
	Username  string `json:"username" bson:"username" validate:"required,min=3,max=20"`     // Required, minimum 3, and maximum 20 characters
	Mobile    string `json:"mobile" bson:"mobile" validate:"required,len=10"`               // Validation for 10 digits starting with 6, 7, 8, or 9
	Email     string `json:"email" bson:"email" validate:"required,email"`                  // Required and must be a valid email
	FirstName string `json:"first_name" bson:"first_name" validate:"required"`              // Required
	LastName  string `json:"last_name" bson:"last_name" validate:"required"`                // Required
	Role      string `json:"role" bson:"role" validate:"required,oneof=admin user manager"` // Required, must be either "admin", "user" or "manager"
}
