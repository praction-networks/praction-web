package models

type LoginUser struct {
	Username string `json:"username" validate:"required,email|username"` // Required, must be a valid email or username
	Password string `json:"password" validate:"required,min=8"`          // Required, minimum 8 characters
}
