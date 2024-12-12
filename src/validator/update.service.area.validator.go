package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateUpdateArea validates the UpdateArea request and returns error details
func ValidateUpdateArea(updateArea *models.UpdateFeture) []response.ErrorDetail {
	logger.Info("Starting validation for UpdateArea") // Debug log
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()
	logger.Info("Validator instance initialized") // Debug log

	// Register custom validation for coordinates range
	v.RegisterValidation("coordinatesRange", func(fl validator.FieldLevel) bool {
		// Get the coordinates as an array of floats
		switch field := fl.Field().Interface().(type) {
		case []float64:
			// Validate GeoJSONPoint coordinates
			if len(field) != 2 {
				return false
			}
			longitude, latitude := field[0], field[1]
			// Longitude should be between -180 and 180
			if longitude < -180 || longitude > 180 {
				return false
			}
			// Latitude should be between -90 and 90
			if latitude < -90 || latitude > 90 {
				return false
			}
		case [][][]float64:
			// Validate GeoJSONPolygon coordinates
			for _, ring := range field {
				for _, coordinate := range ring {
					if len(coordinate) != 2 {
						return false
					}
					longitude, latitude := coordinate[0], coordinate[1]
					// Longitude should be between -180 and 180
					if longitude < -180 || longitude > 180 {
						return false
					}
					// Latitude should be between -90 and 90
					if latitude < -90 || latitude > 90 {
						return false
					}
				}
			}
		}
		return true
	})

	// Perform the validation
	err := v.Struct(updateArea)
	if err != nil {
		logger.Error("Validation errors detected", "errors", err) // Debug log
		for _, e := range err.(validator.ValidationErrors) {
			logger.Info("Processing validation error", "field", e.Field()) // Debug log
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "dive":
				message = e.Field() + " contains invalid entries"
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
