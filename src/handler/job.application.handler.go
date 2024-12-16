package handler

import (
	"encoding/json"
	"net/http"

	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/validator"
)

type JobApplication struct{}

func (JA *JobApplication) RegisterJobApplication(w http.ResponseWriter, r *http.Request) {

	var joabApplication models.JobApplication

	if err := json.NewDecoder(r.Body).Decode(&joabApplication); err != nil {
		logger.Error("Error parsing request body, valid JSON required for Job application", "error", err)
		response.SendBadRequestError(w, "Invalid request payload, valid JSON required for Job application")
		return
	}

	logger.Info("Successfully parsed request body of job application, Proceeding for Validation")

	validationErrors := validator.ValidateJobApplication(&joabApplication)

	if len(validationErrors) > 0 {
		logger.Error("Validation failed for job application attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
	}

	logger.Info("User attributes are valid Create Plan now")

	// // Call service layer to create the plan
	// err := service.CreatePlan(r.Context(), plan)

	// // If there was an error creating the plan, handle it appropriately
	// if err != nil {
	// 	logger.Error("Error creating plan", "error", err)

	// 	// Create an array of ErrorDetail with relevant fields
	// 	errorDetails := []response.ErrorDetail{
	// 		{
	// 			Field:   "Plan",                                           // The field that caused the error
	// 			Message: "Failed to create plan. Please try again later.", // The error message
	// 		},
	// 	}

	// 	// Pass the error details and the HTTP status code to SendError
	// 	response.SendError(w, errorDetails, http.StatusInternalServerError)
	// 	return
	// }

	// // Send a success response (you might want to return the created plan or a success message)
	// response.SendSuccess(w, "Plan created successfully", http.StatusCreated)

}

// func (JA *JobApplication) GetAlljobApplication(w http.ResponseWriter, r *http.Request) {
// 	// Create a context for database interaction
// 	ctx := r.Context()

// 	// Call the service to get all plans
// 	jobsApplication, err := service.GetAllJobApplication(ctx)
// 	if err != nil {
// 		logger.Error("Failed to retrieve jobs: " + err.Error())
// 		http.Error(w, "Failed to retrieve jobs", http.StatusInternalServerError)
// 		return
// 	}

// 	response.SendSuccess(w, jobsApplication, http.StatusOK)
// }

func (JA *JobApplication) UploadResumeAndCoverLetter(w http.ResponseWriter, r *http.Request) {

}
