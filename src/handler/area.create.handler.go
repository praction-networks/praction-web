package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"
	"github.com/praction-networks/quantum-ISP365/webapp/src/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AreaPageHandler struct{}

// CreateJobHandler handles the POST request for creating a job posting.
func (h *AreaPageHandler) CreateAreaPageHandler(w http.ResponseWriter, r *http.Request) {

	var area models.ServiceAreaPage

	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&area); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	// Validate the job object
	validationErrors := validator.ValidateAreaPage(&area)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed for Area attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	logger.Info("area attributes are valid Create a new area now")

	err := service.CreateAreaPage(r.Context(), area)

	// If there was an error creating the plan, handle it appropriately
	if err != nil {
		logger.Error("Error creating blog category", "error", err)

		// Create an array of ErrorDetail with relevant fields
		errorDetails := []response.ErrorDetail{
			{
				Field:   "category",                                                              // The field that caused the error
				Message: "Failed to create blog category. Please try again later." + err.Error(), // The error message
			},
		}

		// Pass the error details and the HTTP status code to SendError
		response.SendError(w, errorDetails, http.StatusInternalServerError)
		return
	}

	// Send a success response (you might want to return the created plan or a success message)
	response.SendSuccess(w, "Page Meta created successfully", http.StatusCreated)

}

func (h *AreaPageHandler) GetAllPageMetaHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters using the common function
	params, err := utils.ParseQueryParams(r.URL.Query())
	if err != nil {
		logger.Error("Error parsing query parameters", "Error", err)
		response.SendBadRequestError(w, "Invalid query parameters")
		return
	}

	// Fetch blog categories using the parsed parameters
	areaPages, err := service.GetAllPageAreaService(r.Context(), params)
	if err != nil {
		logger.Error("Error fetching blog categories", "Error", err)
		response.SendNotFoundError(w, "Failed to fetch blog categories")
		return
	}

	// Return the fetched categories
	response.SendSuccess(w, areaPages, http.StatusOK)
}

func (h *AreaPageHandler) DeletePageMetaByID(w http.ResponseWriter, r *http.Request) {
	logger.Info("Initiating Page meta deletion process")

	// Extract tag ID from the request
	pagemetaID := chi.URLParam(r, "id")
	if pagemetaID == "" {
		logger.Error("Missing Category ID in request")
		response.SendBadRequestError(w, "Missing Category ID in request")
		return
	}

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(pagemetaID)
	if err != nil {
		logger.Error("Invalid tag ID format", "pagemetaID ID", pagemetaID, "error", err)
		response.SendBadRequestError(w, "Invalid Category ID format, should be a valid MongoDB ObjectID")
		return
	}

	ctx := r.Context()

	// Call service to delete blog tag
	err = service.DeleteAreaPage(ctx, objID)
	if err != nil {
		if strings.Contains(err.Error(), "no page meta category found") {
			logger.Info("page meta not found", "meta", pagemetaID)
			response.SendNotFoundError(w, fmt.Sprintf("No blog category found with ID: %s", pagemetaID))
			return
		}

		logger.Error("Failed to delete page", "pagemetaID", pagemetaID, "error", err)
		response.SendInternalServerError(w, "Internal Server Error: Failed to delete page meta")
		return
	}

	logger.Info("Page Meta successfully deleted", "pagemetaID", pagemetaID)
	response.SendSuccess(w, "Page Meta deleted successfully", http.StatusOK)
}
