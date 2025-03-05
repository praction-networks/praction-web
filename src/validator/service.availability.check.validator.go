package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateServiceAreaCheck validates the service area check request
func ValidateServiceAreaCheck(serviceCheck *models.ServiceCheck) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail
	// Initialize validator
	v := validator.New()

	// Validate Pincode
	if err := v.Var(serviceCheck.Pincode, "omitempty,len=6"); err != nil {
		validationErrors = append(validationErrors, response.ErrorDetail{
			Field:   "Pincode",
			Message: "Pincode must be 6 digits long",
		})
	}

	return validationErrors
}
