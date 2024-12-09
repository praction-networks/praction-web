package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateBlogCategory validates the BlogCategory model and returns error details
func ValidateBlogTag(blogtag *models.BlogTag) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

	// Register custom validation for "slug" if required
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
		return re.MatchString(fl.Field().String())
	})

	// Register custom validation for "oneWord"
	v.RegisterValidation("oneWord", func(fl validator.FieldLevel) bool {
		// Check if the field contains only one word (no spaces)
		return !strings.Contains(fl.Field().String(), " ")
	})

	// Perform the validation
	err := v.Struct(blogtag)
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
			case "oneWord":
				message = e.Field() + " must be a single word (no spaces)"
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
