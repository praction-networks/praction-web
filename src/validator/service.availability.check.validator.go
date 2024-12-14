package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateServiceAreaCheck validates the service area check request
func ValidateServiceAreaCheck(serviceCheck *models.ServiceCheck) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail
	var dynamicCheck map[string]interface{}

	// Check if Coordinates is present in the payload
	_, hasCoordinates := dynamicCheck["coordinates"]

	// Initialize validator
	v := validator.New()

	// If Coordinates exists, validate it
	if hasCoordinates {
		if err := v.Struct(serviceCheck.Coordinates); err != nil {
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
	}

	// Validate Pincode
	if err := v.Var(serviceCheck.Pincode, "omitempty,len=6"); err != nil {
		validationErrors = append(validationErrors, response.ErrorDetail{
			Field:   "Pincode",
			Message: "Pincode must be 6 digits long",
		})
	}

	return validationErrors
}
