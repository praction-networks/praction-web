package handler

import (
	"encoding/json"
	"net/http"

	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/validator"
)

type BlogHandler struct{}

// CreateJobHandler handles the POST request for creating a job posting.
func (bh *BlogHandler) CreateBlogHandler(w http.ResponseWriter, r *http.Request) {

	var blog *models.Blog

	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	validationErrors := validator.ValidateBlog(blog)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed for Blog attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	logger.Info("Blog attributes are valid Blog Create")

	// // Call service layer to create the plan
	// err := service.CreateJob(r.Context(), job)

}
