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

type BlogTagHandler struct{}

// CreateJobHandler handles the POST request for creating a job posting.
func (bt *BlogTagHandler) CreateBlogTagHandler(w http.ResponseWriter, r *http.Request) {

	var blogTag models.BlogTag

	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&blogTag); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	// Validate the job object
	validationErrors := validator.ValidateBlogTag(&blogTag)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed for Blog Tag attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	logger.Info("User attributes are valid Create blog tag now")

	err := service.CreateBlogTagService(r.Context(), blogTag)

	// If there was an error creating the plan, handle it appropriately
	if err != nil {
		logger.Error("Error creating blog tag", "error", err)

		// Create an array of ErrorDetail with relevant fields
		errorDetails := []response.ErrorDetail{
			{
				Field:   "blog",                                                         // The field that caused the error
				Message: "Failed to create blog. Please try again later." + err.Error(), // The error message
			},
		}

		// Pass the error details and the HTTP status code to SendError
		response.SendError(w, errorDetails, http.StatusInternalServerError)
		return
	}

	// Send a success response (you might want to return the created plan or a success message)
	response.SendSuccess(w, "blog tag created successfully", http.StatusCreated)

}
func (bt *BlogTagHandler) GetAllBlogTagHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters using the common function
	params, err := utils.ParseQueryParams(r.URL.Query())
	if err != nil {
		logger.Error("Error parsing query parameters", "Error", err)
		response.SendBadRequestError(w, "Invalid query parameters")
		return
	}

	// Fetch blog categories using the parsed parameters
	blogTag, err := service.GetAllBlogTag(r.Context(), params)
	if err != nil {
		logger.Error("Error fetching blog tag", "Error", err)
		response.SendNotFoundError(w, "Failed to fetch blog tag")
		return
	}

	// Return the fetched categories
	response.SendSuccess(w, blogTag, http.StatusOK)
}

func (bt *BlogTagHandler) DeleteBlogTag(w http.ResponseWriter, r *http.Request) {
	logger.Info("Initiating Blog Tag deletion process")

	// Extract tag ID from the request
	tagID := chi.URLParam(r, "id")
	if tagID == "" {
		logger.Error("Missing tag ID in request")
		response.SendBadRequestError(w, "Missing tag ID in request")
		return
	}

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(tagID)
	if err != nil {
		logger.Error("Invalid tag ID format", "tagID", tagID, "error", err)
		response.SendBadRequestError(w, "Invalid tag ID format, should be a valid MongoDB ObjectID")
		return
	}

	ctx := r.Context()

	// Call service to delete blog tag
	err = service.DeleteBlogTagByID(ctx, objID)
	if err != nil {
		if strings.Contains(err.Error(), "no blog tag found") {
			logger.Info("Blog tag not found", "tagID", tagID)
			response.SendNotFoundError(w, fmt.Sprintf("No blog tag found with ID: %s", tagID))
			return
		}

		logger.Error("Failed to delete blog tag", "tagID", tagID, "error", err)
		response.SendInternalServerError(w, "Internal Server Error: Failed to delete blog tag")
		return
	}

	logger.Info("Blog tag successfully deleted", "tagID", tagID)
	response.SendSuccess(w, "Blog tag deleted successfully", http.StatusOK)
}
