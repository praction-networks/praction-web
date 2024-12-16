package validator

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

// ValidateJobApplication validates the job application request and returns error details
func ValidateJobApplication(jobApplication *models.JobApplication) []response.ErrorDetail {
	var validationErrors []response.ErrorDetail

	// Initialize validator instance
	v := validator.New()

	// Register custom validation for IsCurrentEmployer and EndDate
	v.RegisterStructValidation(ExperienceDetailsValidation, models.ExperienceDetails{})

	// Perform the validation
	err := v.Struct(jobApplication)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var message string

			// Handle validation errors for specific tags
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "max":
				message = e.Field() + " must not exceed " + e.Param() + " characters"
			case "min":
				message = e.Field() + " must have at least " + e.Param() + " items"
			case "oneof":
				message = e.Field() + " must be one of: " + e.Param()
			case "url":
				message = e.Field() + " must be a valid URL"
			case "email":
				message = e.Field() + " must be a valid email address"
			case "e164":
				message = e.Field() + " must be a valid phone number in E.164 format"
			case "gtefield":
				message = e.Field() + " must be greater than or equal to " + e.Param()
			default:
				message = e.Field() + " validation failed on the '" + e.Tag() + "' tag"
			}

			// Add the error detail to the list
			validationErrors = append(validationErrors, response.ErrorDetail{
				Field:   e.Field(),
				Message: message,
			})
		}
	}

	// Additional custom validations (if needed)
	validationErrors = append(validationErrors, customValidations(jobApplication)...)

	return validationErrors
}

// customValidations performs additional checks that cannot be handled by tags
func customValidations(jobApplication *models.JobApplication) []response.ErrorDetail {
	var customErrors []response.ErrorDetail

	// Validate that at least one education detail exists
	if len(jobApplication.Education) == 0 {
		customErrors = append(customErrors, response.ErrorDetail{
			Field:   "Education",
			Message: "At least one education detail must be provided",
		})
	}

	// Validate that if IsCurrentEmployer is true, EndDate should not be set
	for i, experience := range jobApplication.Experience {
		if experience.IsCurrentEmployer && !experience.EndDate.Time.IsZero() {
			customErrors = append(customErrors, response.ErrorDetail{
				Field:   "Experience[" + strconv.Itoa(i) + "].EndDate",
				Message: "EndDate should not be set for the current employer",
			})
		}
	}

	return customErrors
}

// ExperienceDetailsValidation performs custom validation for ExperienceDetails
func ExperienceDetailsValidation(sl validator.StructLevel) {
	experience := sl.Current().Interface().(models.ExperienceDetails)

	// If IsCurrentEmployer is true, EndDate should not be set
	if experience.IsCurrentEmployer && !experience.EndDate.Time.IsZero() {
		sl.ReportError(experience.EndDate, "EndDate", "endDate", "currentEmployerEndDate", "")
	}

	// If IsCurrentEmployer is false, EndDate is required
	if !experience.IsCurrentEmployer && experience.EndDate.Time.IsZero() {
		sl.ReportError(experience.EndDate, "EndDate", "endDate", "requiredEndDate", "")
	}
}
