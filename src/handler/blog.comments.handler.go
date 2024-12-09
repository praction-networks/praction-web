package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
	"github.com/praction-networks/quantum-ISP365/webapp/src/validator"
)

type BlogCommentsHandler struct{}

func (bh *BlogCommentsHandler) CreateBlogCommentsHandler(w http.ResponseWriter, r *http.Request) {
	var blogComments *models.Comments

	// Decode the incoming request body into the blogComments struct
	if err := json.NewDecoder(r.Body).Decode(&blogComments); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	// Validate the blog comments
	validationErrors := validator.ValidateBlogComments(blogComments)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed for Blog attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	logger.Info("Blog comments attributes are valid, proceeding to create")

	// Call service layer to create the blog comment
	err := service.CreateBlogComments(r.Context(), *blogComments)

	// Handle errors returned by the service layer
	if err != nil {
		// Check specific errors for more detailed responses
		if strings.Contains(err.Error(), "does not exist") {
			logger.Error(fmt.Sprintf("Blog not found: %v", err))
			response.SendNotFoundError(w, err.Error())
			return
		} else if strings.Contains(err.Error(), "internal database error") {
			logger.Error(fmt.Sprintf("Internal database error: %v", err))
			response.SendInternalServerError(w, "Internal server error, please try again later or connect with Web admin")
			return
		} else {
			// General error
			logger.Error(fmt.Sprintf("Failed to create comment: %v", err))
			response.SendInternalServerError(w, "Failed to create comment")
			return
		}
	}

	// If no error occurred, return a successful response
	responsePayload := map[string]interface{}{
		"message": "Comment added successfully",
	}

	// Send the success response
	response.SendSuccess(w, responsePayload, http.StatusCreated)
}
