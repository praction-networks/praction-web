package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateCreateAdmin validates the create admin request and returns error details
func ValidateCreateAdmin(createAdmin *models.Admin) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

	// Perform the validation
	err := v.Struct(createAdmin)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "gt":
				message = e.Field() + " must be greater than 0"
			case "min":
				message = e.Field() + " must be at least " + e.Param()
			case "oneof":
				message = e.Field() + " must be one of: " + e.Param()
			case "len":
				message = e.Field() + " must have exactly " + e.Param() + " characters"
			case "regexp":
				message = e.Field() + " must match the required pattern"
			default:
				message = e.Field() + " validation failed on the '" + e.Tag() + "' tag"
			}
			validationErrors = append(validationErrors, response.ErrorDetail{
				Field:   e.Field(),
				Message: message,
			})
		}
	}
	// Additional custom validation for Mobile (10 digits starting with 6,7,8,9)
	mobileRegexp := `^[6-9]\d{9}$`
	matched, err := regexp.MatchString(mobileRegexp, createAdmin.Mobile)
	if err != nil || !matched {
		validationErrors = append(validationErrors, response.ErrorDetail{
			Field:   "mobile",
			Message: "Mobile must be a 10-digit number starting with 6, 7, 8, or 9",
		})
	}

	return validationErrors
}
