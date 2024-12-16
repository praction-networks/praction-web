package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
	"github.com/praction-networks/quantum-ISP365/webapp/src/validator"
)

type Plan struct{}

func (p *Plan) CreatePlan(w http.ResponseWriter, r *http.Request) {

	var plan models.Plan

	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload, valid JSON required for Plan Create")
		return
	}

	logger.Info("Successfully parsed request body of plan, Proceeding for Validation")

	validationErrors := validator.ValidateCreatePlan(&plan)

	if len(validationErrors) > 0 {
		logger.Error("Validation failed for admin user attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
	}

	logger.Info("User attributes are valid Create Plan now")

	// Call service layer to create the plan
	err := service.CreatePlan(r.Context(), plan)

	// If there was an error creating the plan, handle it appropriately
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDuplicateKey):
			// Handle duplicate key error
			logger.Error("Duplicate plan detected")
			// Return HTTP 409 Conflict to the user
			response.SendConflictError(w, "Conflict: Duplicate plan")
		case errors.Is(err, service.ErrDatabaseInsert):
			// Handle database insertion error
			logger.Error("Database insertion failed")
			// Return HTTP 500 Internal Server Error
			response.SendInternalServerError(w, "Internal Server Error")
		case errors.Is(err, service.ErrDatabaseInternal):
			// Handle database connection error
			logger.Error("Database connection issue")
			// Return HTTP 503 Service Unavailable
			response.SendServiceUnavailableError(w, "Database connection issue")
		default:
			// Handle other errors
			logger.Error("Unknown error occurred")
			response.SendInternalServerError(w, "Internal Server Error")
		}
		return
	}

	// Send a success response (you might want to return the created plan or a success message)
	response.SendSuccess(w, "Plan created successfully", http.StatusCreated)

}

func (p *Plan) GetAllPlan(w http.ResponseWriter, r *http.Request) {
	// Create a context for database interaction
	ctx := r.Context()

	// Call the service to get all plans
	plans, err := service.GetAllPlans(ctx)
	if err != nil {
		logger.Error("Failed to retrieve plans: " + err.Error())
		http.Error(w, "Failed to retrieve plans", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, plans, http.StatusOK)
}
