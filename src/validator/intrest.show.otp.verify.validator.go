package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateShowUserIntrestOTP validates the OTP verification request
func ValidateShowUserIntrestOTP(userIntrest *models.UserOTPVerify) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

	// Register custom validators
	RegisterCustomValidators(v)

	// Perform basic struct validation
	err := v.Struct(userIntrest)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "len":
				message = e.Field() + " must be exactly " + e.Param() + " characters"
			case "email":
				message = e.Field() + " should be a valid email address"
			case "otp_len":
				message = e.Field() + " must be a 6-digit OTP"
			default:
				message = e.Field() + " validation failed on the '" + e.Tag() + "' tag"
			}
			validationErrors = append(validationErrors, response.ErrorDetail{
				Field:   e.Field(),
				Message: message,
			})
		}
	}

	// Custom check: Ensure either Email or Mobile is provided
	if userIntrest.Email == "" && userIntrest.Mobile == "" {
		validationErrors = append(validationErrors, response.ErrorDetail{
			Field:   "email or mobile",
			Message: "At least one of Email or Mobile is required",
		})
	}

	return validationErrors
}

// Custom OTP Length Validator
func otpLength(fl validator.FieldLevel) bool {
	otp := fl.Field().Int()
	// Check if OTP is a 6-digit number
	return otp >= 100000 && otp <= 999999
}

// Register the custom validation function
func RegisterCustomValidators(v *validator.Validate) {
	_ = v.RegisterValidation("otp_len", otpLength)
}
