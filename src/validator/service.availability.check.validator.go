package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidatePointCheck validates the point check request (latitude and longitude) and returns error details.
func ValidatePointCheck(pointRequest *models.PointRequest) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

	// Perform the validation
	err := v.Struct(pointRequest)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "latitude":
				message = e.Field() + " must be a valid latitude between -90 and 90"
			case "longitude":
				message = e.Field() + " must be a valid longitude between -180 and 180"
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
