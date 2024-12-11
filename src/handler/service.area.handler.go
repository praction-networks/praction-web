package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"
	"github.com/praction-networks/quantum-ISP365/webapp/src/validator"
)

type ServiceAreaHandler struct{}

// CreateServiceArea handles the POST request for creating a service area collection.
func (sa *ServiceAreaHandler) CreateServiceArea(w http.ResponseWriter, r *http.Request) {
	var serviceAreaCollection *models.FeatureCollection

	// Parse the request body into ServiceAreaCollection
	if err := json.NewDecoder(r.Body).Decode(&serviceAreaCollection); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	// Validate the FeatureCollection object (this will include validation for each Feature)
	validationErrors := validator.ValidateServiceAreaCollection(serviceAreaCollection)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	err := service.CreateServiceArea(r.Context(), serviceAreaCollection)

	if err != nil {
		logger.Error("Error saving service area collection", "error", err)
		response.SendInternalServerError(w, "Failed to save service area collection")
		return
	}

	// Save the service area collection (if needed) or just respond with success
	logger.Info("Service area collection created successfully", "serviceAreaCollection", serviceAreaCollection)
	response.SendSuccess(w, "Service area collection created successfully", http.StatusCreated)
}

func (sa *ServiceAreaHandler) GetAllServiceArea(w http.ResponseWriter, r *http.Request) {

	params, err := utils.ParseQueryParams(r.URL.Query())
	if err != nil {
		logger.Error("Error parsing query parameters", "Error", err)
		response.SendBadRequestError(w, "Invalid query parameters")
		return
	}

	// Call the service to get all plans
	blogs, err := service.GetAllServiceAreaService(r.Context(), params)
	if err != nil {
		logger.Error("Failed to retrieve Service Area: " + err.Error())
		response.SendInternalServerError(w, "Failed to retrieve Service Area")
		return
	}

	response.SendSuccess(w, blogs, http.StatusOK)

}

func (sa *ServiceAreaHandler) CheckServiceArea(w http.ResponseWriter, r *http.Request) {
	var point *models.PointRequest

	// Parse the request body into ServiceAreaCollection
	if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	// Validate the FeatureCollection object (this will include validation for each Feature)
	validationErrors := validator.ValidatePointCheck(point)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	err := service.CheckServiceAvailability(r.Context(), point)

	if err != nil {
		// Handle different types of errors based on their content
		switch {
		case err.Error() == "no service areas found":
			// Return 404 if no service areas are found
			logger.Error("No service area found for the provided point")
			response.SendNotFoundError(w, "No service area found for the provided point")
		case err.Error() == fmt.Sprintf("the point with latitude %f and longitude %f is not within any service area", point.Latitude, point.Longitude):
			// Return 403 if the point is not within any service area
			logger.Error("Point is not within any service area", "lat", point.Latitude, "lon", point.Longitude)
			response.SendUnauthorizedError(w, "Point is not within any service area")
		default:
			// Return 500 for other unexpected errors
			logger.Error("Unexpected error while checking service availability", "error", err)
			response.SendInternalServerError(w, "Internal server error")
		}
		return
	}

	payload := map[string]interface{}{
		"message": "OK",
	}

	response.SendSuccess(w, payload, http.StatusOK)

}
