package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateBlogCategory validates the BlogCategory model and returns error details
func ValidateBlogCategory(blogCategory *models.BlogCategory) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

	// Register custom validation for "slug" if required
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
		return re.MatchString(fl.Field().String())
	})

	// Perform the validation
	err := v.Struct(blogCategory)
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
			case "slug":
				message = e.Field() + " must be a valid slug (lowercase, hyphen-separated)"
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
