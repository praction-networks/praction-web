package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateBlog validates the Blog model and returns detailed errors
func ValidateBlogComments(blog *models.Comments) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

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
			case "numeric":
				message = e.Field() + " must be a digit only"
			case "gte":
				message = e.Field() + " must be greater than or equal to " + e.Param()
			case "mongodb":
				message = e.Field() + " must be a valid MongoDB ObjectID"
			case "len":
				message = e.Field() + " must be exectly " + e.Param() + "digit long"
			case "uuid4":
				message = e.Field() + " is not a valid uuid v4"
			case "email":
				message = e.Field() + " must ve a valid email"
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
