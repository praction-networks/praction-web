package models

type User struct {
	Username  string `json:"username" bson:"username" validate:"required,min=3,max=20"` // Required, minimum 3, and maximum 20 characters
	Password  string `json:"password" bson:"password" validate:"required,min=8"`        // Required, minimum 8 characters
	Salt      string `json:"saltstr" bson:"saltstr"`                                    // No validation, managed internally
	Email     string `json:"email" bson:"email" validate:"required,email"`              // Required and must be a valid email
	FirstName string `json:"first_name" bson:"first_name" validate:"required"`          // Required
	LastName  string `json:"last_name" bson:"last_name" validate:"required"`            // Required
	Role      string `json:"role" bson:"role" validate:"required,oneof=admin user"`     // Required, must be either "admin" or "user"
}
