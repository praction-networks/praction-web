package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateShowUserIntrestOTP validates the OTP verification request
func ValidateShowUserIntrestOTPResnd(userIntrest *models.UserOTPResend) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

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
			case "oneof":
				message = e.Field() + " must be one of: " + e.Param()
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
