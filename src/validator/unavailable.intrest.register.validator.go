package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateUnAvailableArea validates the UnAvailableArea structure
func ValidateUnAvailableArea(unAvailableArea *models.UnAvailableArea) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()
	RegisterCustomValidatorsForUnAvailableArea(v)

	// Perform the validation
	err := v.Struct(unAvailableArea)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "len":
				message = e.Field() + " must have a length of " + e.Param()
			case "latitude":
				message = e.Field() + " must be a valid latitude (-90 to 90)"
			case "longitude":
				message = e.Field() + " must be a valid longitude (-180 to 180)"
			case "oneword":
				message = e.Field() + " must be a single word"
			case "email":
				message = e.Field() + " must be a valid email address"
			case "mobile":
				message = e.Field() + " must be a valid 10-digit mobile number"
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

// RegisterCustomValidatorsForUnAvailableArea registers custom validation tags
func RegisterCustomValidatorsForUnAvailableArea(v *validator.Validate) {
	// Validate that a field contains only one word
	v.RegisterValidation("oneword", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString(`^\w+$`, fl.Field().String())
		return match
	})

	// Validate mobile numbers (assumes a 10-digit number for this example, adjust regex as per your requirements)
	v.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString(`^\d{10}$`, fl.Field().String())
		return match
	})

	// Validate latitude (between -90 and 90)
	v.RegisterValidation("latitude", func(fl validator.FieldLevel) bool {
		value := fl.Field().Float()
		return value >= -90 && value <= 90
	})

	// Validate longitude (between -180 and 180)
	v.RegisterValidation("longitude", func(fl validator.FieldLevel) bool {
		value := fl.Field().Float()
		return value >= -180 && value <= 180
	})
}
