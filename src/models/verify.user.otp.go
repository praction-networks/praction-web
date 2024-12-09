package models

type UserOTPVerify struct {
	Email  string `json:"email" validate:"email"`
	Mobile string `json:"mobile" validate:"len=10"`
	OTP    int64  `json:"otp" validate:"required,otp_len"`
}
