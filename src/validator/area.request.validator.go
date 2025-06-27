package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateAreaCategory validates the AreaCategory model and returns error details
func ValidateAreaPage(areaCategory *models.ServiceAreaPage) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	v := validator.New()

	// Perform the validation
	err := v.Struct(areaCategory)
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
			case "url":
				message = e.Field() + " must be a valid URL"
			case "uuid4":
				message = e.Field() + " must be a valid UUIDv4"
			case "dive":
				message = e.Field() + " contains invalid entries"
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
