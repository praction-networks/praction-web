package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateCreatePlan validates the create plan request and returns error details
func ValidateCreatePlan(createPlan *models.Plan) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

	// Perform the validation
	err := v.Struct(createPlan)
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
