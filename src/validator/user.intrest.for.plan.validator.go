package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

func ValidateAvailableUser(userRequest *models.AvailableUserRequest) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()
	RegisterCustomValidatorsPlan(v)

	// Perform the validation
	err := v.Struct(userRequest)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "email":
				message = e.Field() + " must be a valid email address"
			case "mobile":
				message = e.Field() + " must be a valid 10-digit mobile number"
			case "oneword":
				message = e.Field() + " must be a single word"
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

// Custom validation tags
func RegisterCustomValidatorsPlan(v *validator.Validate) {
	// Validate that a field contains only one word
	v.RegisterValidation("oneword", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString(`^\w+$`, fl.Field().String())
		return match
	})

	// Validate mobile numbers (assumes a 10-digit number for this example, adjust regex as per your requirements)
	v.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString(`^\d{10}$`, fl.Field().String())
		return match
	})
}
