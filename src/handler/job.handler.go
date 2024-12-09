package handler

import (
	"encoding/json"
	"net/http"

	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
	"github.com/praction-networks/quantum-ISP365/webapp/src/validator"
)

type JobHandler struct{}

// CreateJobHandler handles the POST request for creating a job posting.
func (h *JobHandler) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	var job models.Job

	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	// Validate the job object
	validationErrors := validator.ValidateJob(&job)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed for admin user attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	logger.Info("User attributes are valid Create Plan now")

	// Call service layer to create the plan
	err := service.CreateJob(r.Context(), job)

	// If there was an error creating the plan, handle it appropriately
	if err != nil {
		logger.Error("Error creating job", "error", err)

		// Create an array of ErrorDetail with relevant fields
		errorDetails := []response.ErrorDetail{
			{
				Field:   "Job",                                           // The field that caused the error
				Message: "Failed to create job. Please try again later.", // The error message
			},
		}

		// Pass the error details and the HTTP status code to SendError
		response.SendError(w, errorDetails, http.StatusInternalServerError)
		return
	}

	// Send a success response (you might want to return the created plan or a success message)
	response.SendSuccess(w, "Job created successfully", http.StatusCreated)
}

func (p *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	// Create a context for database interaction
	ctx := r.Context()

	// Call the service to get all plans
	jobs, err := service.GetAllJobs(ctx)
	if err != nil {
		logger.Error("Failed to retrieve All Jobs: " + err.Error())
		http.Error(w, "Failed to retrieve All Jobs", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, jobs, http.StatusOK)
}
