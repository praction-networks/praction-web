package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// Custom validation for date format (YYYY-MM-DD)
func validateDateFormat(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Allow empty values (omitempty)
	}
	layout := "2006-01-02"
	_, err := time.Parse(layout, value)
	return err == nil
}

// ValidateUserInterestUpdate validates the UserInterestUpdate model
func ValidateUserInterestUpdate(userInterestUpdate *models.UserInterestUpdate) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

	// Register custom validation for Date format
	v.RegisterValidation("date_format", validateDateFormat)

	// Perform validation
	err := v.Struct(userInterestUpdate)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "min":
				message = e.Field() + " must be at least " + e.Param() + " characters"
			case "max":
				message = e.Field() + " must not exceed " + e.Param() + " characters"
			case "oneof":
				message = e.Field() + " must be one of [" + e.Param() + "]"
			case "omitempty":
				message = e.Field() + " is optional but must be valid if provided"
			case "date_format":
				message = e.Field() + " must be a valid date in YYYY-MM-DD format"
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
