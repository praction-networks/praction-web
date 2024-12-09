package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateCreatePlan validates the create plan request and returns error details
func ValidateShowUserIntrest(userIntrest *models.UserInterest) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

	RegisterCustomValidatorsPinCode(v)

	// Perform the validation
	err := v.Struct(userIntrest)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "len":
				message = e.Field() + " must be maximum 10 digit"
			case "max":
				message = e.Field() + " must not exceed " + e.Param() + " characters"
			case "email":
				message = e.Field() + " should be a valid email address"
			default:
				message = e.Field() + " validation failed on the '" + e.Tag() + "' tag"
			}
			validationErrors = append(validationErrors, response.ErrorDetail{
				Field:   e.Field(),
				Message: message,
			})
		}
	}

	return validationErrors
}

// Custom OTP Length Validator
func pincodeCheck(fl validator.FieldLevel) bool {
	pincode := fl.Field().Int()
	// Check if OTP is a 6-digit number
	return pincode >= 100000 && pincode <= 999999
}

// Register the custom validation function
func RegisterCustomValidatorsPinCode(v *validator.Validate) {
	_ = v.RegisterValidation("pincode", pincodeCheck)
}
