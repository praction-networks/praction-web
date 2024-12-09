package validator

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateUserRefrence validates the user reference request and returns error details
func ValidateUserRefrence(userRefrence *models.UserRefrence) []response.ErrorDetail {
	log.Printf("Validating UserRefrence: %+v", userRefrence)
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

	// Register custom validators
	RegisterCustomValidatorsPinCode(v)
	RegisterCustomValidatorsMobile(v)
	RegisterCustomValidatorsUnique(v)

	// Perform the validation
	err := v.Struct(userRefrence)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "len":
				message = e.Field() + " must be maximum 10 digits"
			case "max":
				message = e.Field() + " must not exceed " + e.Param() + " characters"
			case "email":
				message = e.Field() + " should be a valid email address"
			case "min":
				message = e.Field() + " must have a minimum of " + e.Param() + " items"
			case "dive":
				message = e.Field() + " must be a valid user"
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

// Custom Mobile Number Validator
func mobileCheck(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// Check if the mobile number is 10 digits
	return len(mobile) == 10
}

// Register the custom validation for Mobile Number
func RegisterCustomValidatorsMobile(v *validator.Validate) {
	_ = v.RegisterValidation("mobile", mobileCheck)
}

func validateUniqueEmailsAndMobiles(fl validator.FieldLevel) bool {
	referrels, ok := fl.Field().Interface().([]models.UserType)
	if !ok {
		log.Println("Field is not of type []UserType")
		return false
	}

	emailSet := make(map[string]bool)
	mobileSet := make(map[string]bool)

	for i, referrel := range referrels {
		if emailSet[referrel.Email] {
			log.Printf("Duplicate email found: %s at index %d", referrel.Email, i)
			return false
		}
		emailSet[referrel.Email] = true

		if mobileSet[referrel.Mobile] {
			log.Printf("Duplicate mobile found: %s at index %d", referrel.Mobile, i)
			return false
		}
		mobileSet[referrel.Mobile] = true
	}

	log.Println("Validation passed for uniqueEmailsAndMobiles")
	return true
}

// Register the custom validation for Unique Emails and Mobiles
func RegisterCustomValidatorsUnique(v *validator.Validate) {
	_ = v.RegisterValidation("uniqueEmailsAndMobiles", validateUniqueEmailsAndMobiles)
}
