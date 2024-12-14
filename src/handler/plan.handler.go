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

type Plan struct{}

// CreatePlan godoc
// @Summary Create a new plan
// @Description This endpoint allows the creation of a new plan in the system
// @Tags Plan
// @Accept json
// @Produce json
// @Param plan body models.Plan true "Create Plan"
// @Success 201 {string} string "Plan created successfully"
// @Failure 400 {array} response.ErrorDetail "Invalid request payload or validation errors"
// @Failure 500 {array} response.ErrorDetail "Error creating plan"
// @Router /web/v1/plan [post]
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
		logger.Error("Error creating plan", "error", err)

		// Create an array of ErrorDetail with relevant fields
		errorDetails := []response.ErrorDetail{
			{
				Field:   "Plan",                                           // The field that caused the error
				Message: "Failed to create plan. Please try again later.", // The error message
			},
		}

		// Pass the error details and the HTTP status code to SendError
		response.SendError(w, errorDetails, http.StatusInternalServerError)
		return
	}

	// Send a success response (you might want to return the created plan or a success message)
	response.SendSuccess(w, "Plan created successfully", http.StatusCreated)

}

// GetAllPlan godoc
// @Summary Get all plans
// @Description Retrieve all the plans available in the system
// @Tags Plan
// @Produce json
// @Success 200 {array} models.Plan "List of plans"
// @Failure 500 {array} response.ErrorDetail "Error retrieving plans"
// @Router /web/v1/plans [get]
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
