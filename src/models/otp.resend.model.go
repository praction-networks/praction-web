package models

type UserOTPResend struct {
	Email  string `json:"email" validate:"email"`
	Mobile string `json:"mobile" validate:"len=10"`
	Resend string `json:"resend" validate:"required,oneof=text voice"` // Corrected closing double-quote
}
