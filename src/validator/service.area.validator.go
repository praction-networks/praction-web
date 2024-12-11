package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateServiceArea validates the ServiceArea model and returns error details
func ValidateServiceAreaCollection(serviceArea *models.FeatureCollection) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

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
	err := v.Struct(serviceArea)
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
			case "coordinatesRange":
				message = e.Field() + " must have valid coordinates"
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
