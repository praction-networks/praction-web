package validator

import (
	"net/url"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// Custom validation function for UUID v4 or URL
func validateUUIDorURL(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Allow empty value (omitempty)
	}

	// Check if it's a valid UUID v4
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if uuidRegex.MatchString(value) {
		return true
	}

	// Check if it's a valid URL
	_, err := url.ParseRequestURI(value)
	return err == nil
}

// ValidateBlog validates the Blog model and returns detailed errors
func ValidateUpdateBlog(blog *models.BlogUpdate) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

	blog.BlogDescription = SanitizeHTML(blog.BlogDescription)

	// Register custom validation for BlogImage (UUID v4 or URL)
	v.RegisterValidation("uuid_or_url", validateUUIDorURL)

	// Register custom validation for "slug"
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
		return re.MatchString(fl.Field().String())
	})

	// Perform validation
	err := v.Struct(blog)
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
			case "url":
				message = e.Field() + " must be a valid URL"
			case "uuid4":
				message = e.Field() + " is not a valid uuid v4"
			case "uuid_or_url":
				message = e.Field() + " must be a valid UUID v4 or a valid URL"
			case "gte":
				message = e.Field() + " must be greater than or equal to " + e.Param()
			case "oneWord":
				message = e.Field() + " must not contain spaces"
			case "dive":
				message = e.Field() + " must have valid nested elements"
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
