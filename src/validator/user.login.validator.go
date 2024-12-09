package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// Custom validation function to check if a value is a valid username
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// Username: 3-20 characters, letters, numbers, underscores allowed
	regex := `^[a-zA-Z0-9_]{3,20}$`
	return regexp.MustCompile(regex).MatchString(username)
}

func RegisterLoginValidators(v *validator.Validate) {
	// Register custom validation for username
	v.RegisterValidation("username", validateUsername)
}

func ValidateLoginUser(loginUser *models.LoginUser) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	v := validator.New()
	RegisterLoginValidators(v)

	err := v.Struct(loginUser)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "email":
				message = e.Field() + " should be a valid email address"
			case "username":
				message = e.Field() + " should be 3-20 characters long and only contain letters, numbers, or underscores"
			case "min":
				message = e.Field() + " should be at least " + e.Param() + " characters long"
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
