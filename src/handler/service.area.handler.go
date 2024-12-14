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
	"go.mongodb.org/mongo-driver/mongo"
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
	var serviceCheck *models.ServiceCheck

	// Parse the request body into ServiceAreaCollection
	if err := json.NewDecoder(r.Body).Decode(&serviceCheck); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	// Check if both Coordinates and Pincode are empty
	if (serviceCheck.Coordinates == (models.PointRequest{}) ||
		(serviceCheck.Coordinates.Latitude == 0 && serviceCheck.Coordinates.Longitude == 0)) &&
		serviceCheck.Pincode == "" {
		logger.Error("At least one of the fields 'Coordinates' or 'Pincode' is required to locate service")
		response.SendBadRequestError(w, "At least one of the fields 'Coordinates' or 'Pincode' is required")
		return
	}
	// Validate the FeatureCollection object (this will include validation for each Feature)
	validationErrors := validator.ValidateServiceAreaCheck(serviceCheck)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	if serviceCheck.Coordinates != (models.PointRequest{}) {
		err := service.CheckServiceAvailability(r.Context(), &serviceCheck.Coordinates)

		if err != nil {
			// Handle different types of errors based on their content
			switch {
			case err.Error() == "no service areas found":
				// Return 404 if no service areas are found
				logger.Error("No service area found for the provided point")
				response.SendNotFoundError(w, "No service area found for the provided point")
			case err.Error() == fmt.Sprintf("the point with latitude %f and longitude %f is not within any service area", serviceCheck.Coordinates.Latitude, serviceCheck.Coordinates.Longitude):
				// Return 403 if the point is not within any service area
				logger.Error("Point is not within any service area", "lat", serviceCheck.Coordinates.Latitude, "lon", serviceCheck.Coordinates.Longitude)
				payload := map[string]interface{}{
					"status":    "success",
					"available": false,
					"message":   "Sorry, we are not available in your area.",
					"data":      nil,
				}
				response.SendSuccess(w, payload, http.StatusAccepted)
			default:
				// Return 500 for other unexpected errors
				logger.Error("Unexpected error while checking service availability", "error", err)
				response.SendInternalServerError(w, "Internal server error")
			}
			return
		}

		payload := map[string]interface{}{
			"status":    "success",
			"available": true,
			"message":   "We are available in your area!",
			"data":      nil,
		}

		response.SendSuccess(w, payload, http.StatusOK)
		return
	}

	// Fetch area data based on pincode
	areaData, err := service.CheckServiceByPinCode(r.Context(), serviceCheck.Pincode)
	if err != nil {
		// Handle error scenarios (already checks for len(areaData) == 0)
		logger.Error("Error checking service area", "error", err)

		if err.Error() == fmt.Sprintf("no service areas found for pincode %s", serviceCheck.Pincode) {
			response.SendNotFoundError(w, fmt.Sprintf("No service available for pincode %s", serviceCheck.Pincode))
		} else {
			response.SendInternalServerError(w, "An error occurred while checking service availability")
		}
		return
	}

	// Respond with service area details (no need to check len(areaData) again)
	logger.Info("Service areas found for pincode", "pincode", serviceCheck.Pincode)

	payload := map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("Service is available for pincode %s", serviceCheck.Pincode),
		"data":    areaData,
	}

	response.SendSuccess(w, payload, http.StatusOK)

}

func (sa *ServiceAreaHandler) UpdateServiceArea(w http.ResponseWriter, r *http.Request) {

	var areaUpdate models.UpdateFeture

	// Parse the request body into ServiceAreaCollection
	if err := json.NewDecoder(r.Body).Decode(&areaUpdate); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	logger.Info("Parsed UpdateArea request", "updateArea", areaUpdate)

	// Ensure at least one of AddArea or RemoveArea is provided
	if len(areaUpdate.AddArea) == 0 && len(areaUpdate.RemoveArea) == 0 {
		logger.Error("Validation failed: both AddArea and RemoveArea are empty")
		response.SendError(w, []response.ErrorDetail{{
			Field:   "AddArea, RemoveArea",
			Message: "At least one of AddArea or RemoveArea must be provided",
		}}, http.StatusBadRequest)
		return
	}

	// Validate the FeatureCollection object (this will include validation for each Feature)
	validationErrors := validator.ValidateUpdateArea(&areaUpdate)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	// Call the service to update areas
	err := service.UpdateAreaService(r.Context(), &areaUpdate)
	if err != nil {
		logger.Error("Error modifying service area", "error", err)

		// Map service errors to appropriate HTTP responses
		if isDuplicateAreaNameError(err) {
			response.SendError(w, []response.ErrorDetail{{
				Field:   "areaName",
				Message: err.Error(),
			}}, http.StatusConflict)
		} else if isDuplicateKeyError(err) {
			// Return HTTP 409 Conflict if the error is due to a duplicate key
			response.SendError(w, []response.ErrorDetail{{
				Field:   "areaName, uuid",
				Message: err.Error(),
			}}, http.StatusConflict)
		} else if isNotFoundError(err) {
			// Return HTTP 404 Not Found if the error is due to a missing resource
			response.SendError(w, []response.ErrorDetail{{
				Field:   "areaName",
				Message: err.Error(),
			}}, http.StatusNotFound)
		} else {
			// Return HTTP 500 Internal Server Error for other issues
			response.SendInternalServerError(w, "Failed to modify service area")
		}
		return
	}

	// Respond with success
	response.SendSuccess(w, map[string]string{"message": "Service areas updated successfully"}, http.StatusOK)

}

func (sa *ServiceAreaHandler) ModifyServiceArea(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		logger.Error("Blog ID is required to fetch Blog")
		response.SendBadRequestError(w, "Blog ID is Required")
		return
	}

	// Check if the ID is a valid MongoDB ObjectID
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		logger.Error("Invalid Service Area ID", "error", err)
		response.SendBadRequestError(w, "Invalid Service Area ID")
		return
	}

	var oneAreaUpdate models.UpdateOneArea

	if err := json.NewDecoder(r.Body).Decode(&oneAreaUpdate); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	// Validate the FeatureCollection object (this will include validation for each Feature)
	validationErrors := validator.ValidateModifyArea(&oneAreaUpdate)
	if len(validationErrors) > 0 {
		logger.Error("Modify Area Validation failed", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	// Call the service to modify the service area
	err := service.ModifyServiceArea(r.Context(), id, &oneAreaUpdate)
	if err != nil {
		logger.Error("Error modifying service area", "error", err)

		// Map service errors to appropriate HTTP responses
		switch {
		case isDuplicateAreaNameError(err):
			response.SendError(w, []response.ErrorDetail{{
				Field:   "areaName",
				Message: err.Error(),
			}}, http.StatusConflict) // HTTP 409 for conflicts
		case isDuplicateKeyError(err):
			response.SendError(w, []response.ErrorDetail{{
				Field:   "areaName, uuid",
				Message: err.Error(),
			}}, http.StatusConflict) // HTTP 409 for conflicts
		case isNotFoundError(err):
			response.SendError(w, []response.ErrorDetail{{
				Field:   "id",
				Message: err.Error(),
			}}, http.StatusNotFound) // HTTP 404 for not found
		default:
			response.SendInternalServerError(w, "Failed to modify service area")
		}
		return
	}

	// Success response
	response.SendSuccess(w, map[string]string{"message": "Service area modified successfully"}, http.StatusOK)

}

func isDuplicateKeyError(err error) bool {
	if mongoErr, ok := err.(mongo.WriteException); ok {
		for _, writeErr := range mongoErr.WriteErrors {
			if writeErr.Code == 11000 { // Duplicate Key Error Code
				return true
			}
		}
	}
	return false
}

func isNotFoundError(err error) bool {
	return err.Error() == "document not found" // Customize based on service-layer error
}

// Helper function to check if a substring exists in an error message
func contains(fullText, subText string) bool {
	fmt.Printf("FullText: %s\n", fullText)
	fmt.Printf("SubText: %s\n", subText)
	return fullText != "" && subText != "" && strings.Contains(fullText, subText)
}

func isDuplicateAreaNameError(err error) bool {
	if err == nil {
		return false
	}

	fmt.Printf("Error text: %s\n", err.Error()) // Log the full error message for debugging
	return contains(err.Error(), "duplicate areaName")
}
