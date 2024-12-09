package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateJob validates the job request and returns error details
func ValidateJob(job *models.Job) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

	// Perform the validation
	err := v.Struct(job)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "max":
				message = e.Field() + " must not exceed " + e.Param() + " characters"
			case "min":
				message = e.Field() + " must have at least " + e.Param() + " items"
			case "oneof":
				message = e.Field() + " must be one of: " + e.Param()
			case "gt":
				message = e.Field() + " must be greater than " + e.Param()
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
