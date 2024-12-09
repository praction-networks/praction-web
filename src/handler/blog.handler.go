package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"
	"github.com/praction-networks/quantum-ISP365/webapp/src/validator"
)

type BlogHandler struct{}

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

	// Call service layer to create the blog
	err := service.CreateBlog(r.Context(), *blog)
	if err != nil {
		// Error handling based on the returned error from the service layer
		if err.Error() == "blog exists, cannot create new one with same blog title" {
			// Blog with the same title already exists
			logger.Error("Blog creation failed", "error", err)
			response.SendError(w, []response.ErrorDetail{
				{Field: "blogTitle", Message: "Blog title already exists"},
			}, http.StatusConflict)
			return
		} else if err.Error() == "one or more blog categories not found" {
			// Some of the categories were not found
			logger.Error("Category validation failed", "error", err)
			response.SendError(w, []response.ErrorDetail{
				{Field: "category", Message: "One or more categories not found"},
			}, http.StatusNotFound)
			return
		} else if err.Error() == "one or more blog tag not found" {
			// Some of the tags were not found
			logger.Error("Tag validation failed", "error", err)
			response.SendError(w, []response.ErrorDetail{
				{Field: "tag", Message: "One or more tags not found"},
			}, http.StatusNotFound)
			return
		} else if err.Error() == "blog image not found" {
			// Blog image not found error
			logger.Error("Blog image not found", "error", err)
			response.SendError(w, []response.ErrorDetail{
				{Field: "blogImage", Message: "Blog image not found"},
			}, http.StatusNotFound)
			return
		} else if err.Error() == "feature image not found" {
			// Feature image not found error
			logger.Error("Feature image not found", "error", err)
			response.SendError(w, []response.ErrorDetail{
				{Field: "featureImage", Message: "Feature image not found"},
			}, http.StatusNotFound)
			return
		} else {
			// Generic server error
			logger.Error("Blog creation failed", "error", err)
			response.SendError(w, []response.ErrorDetail{
				{Field: "general", Message: "An unexpected error occurred while creating the blog"},
			}, http.StatusInternalServerError)
			return
		}
	}

	// Return success response if no error
	response.SendSuccess(w, "Blog created successfully", http.StatusCreated)
}

func (bh *BlogHandler) GetBlogHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	params, err := utils.ParseQueryParams(r.URL.Query())
	if err != nil {
		logger.Error("Error parsing query parameters", "Error", err)
		response.SendBadRequestError(w, "Invalid query parameters")
		return
	}

	// Call the service to get all plans
	blogs, err := service.GetAllBlogService(ctx, params)
	if err != nil {
		logger.Error("Failed to retrieve blogs: " + err.Error())
		http.Error(w, "Failed to retrieve blogs", http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w, blogs, http.StatusOK)

}

func (bh *BlogHandler) GetOneBlogHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	id := chi.URLParam(r, "id")
	if id == "" {
		logger.Error("Blog ID is required to fetch Blog")

		response.SendBadRequestError(w, "Blog ID is Required")
		return
	}

	blog, err := service.GetOneBlog(r.Context(), id)

	if err != nil {
		logger.Error(err.Error())
		response.SendInternalServerError(w, "Error while retriving Blog please connect with you administrator")
		return
	}

	if blog == nil {
		logger.Info("No blog found with the ID")
		response.SendNotFoundError(w, "No Blog Found")
		return
	}

	logger.Info("Blog found with ID")
	response.SendSuccess(w, blog, http.StatusOK)
}
