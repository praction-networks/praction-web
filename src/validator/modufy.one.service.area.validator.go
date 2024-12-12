package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateModifyArea validates the Feature model during area modification and returns error details
func ValidateModifyArea(modifyArea *models.UpdateOneArea) []response.ErrorDetail {
	logger.Info("Starting validation for modifyArea") // Debug log
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()
	logger.Info("Validator instance initialized") // Debug log
	// Perform the validation
	err := v.Struct(modifyArea)
	if err != nil {
		logger.Error("Validation errors detected", "errors", err) // Debug log
		for _, e := range err.(validator.ValidationErrors) {
			logger.Info("Processing validation error", "field", e.Field()) // Debug log
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "oneof":
				message = e.Field() + " must be one of: " + e.Param()
			case "coordinatesRange":
				message = e.Field() + " contains invalid GeoJSON coordinates"
			default:
				message = e.Field() + " validation failed on the '" + e.Tag() + "' tag"
			}
			logger.Info("Validation error message created", "message", message) // Debug log
			validationErrors = append(validationErrors, response.ErrorDetail{
				Field:   e.Field(),
				Message: message,
			})
		}
	} else {
		logger.Info("Validation passed with no errors") // Debug log
	}

	logger.Info("Validation errors collected", "errors", validationErrors) // Debug log
	return validationErrors
}
