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

func (bh *BlogHandler) GeAlltBlogHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	params, err := utils.ParseQueryParams(r.URL.Query())
	if err != nil {
		logger.Error("Error parsing query parameters", "Error", err)
		response.SendBadRequestError(w, "Invalid query parameters")
		return
	}

	// Call the service to get all plans
	blogs, err := service.GetAdminAllBlogService(ctx, params)
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
		logger.Error("Blog identifier (ID or slug) is required")
		response.SendBadRequestError(w, "Blog identifier (ID or slug) is required")
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

func (bh *BlogHandler) AddView(w http.ResponseWriter, r *http.Request) {

	blogID := r.URL.Query().Get("blogID")

	if blogID == "" {
		logger.Error("blog ID is required to increment view to blog with params")
		response.SendBadRequestError(w, "blog ID is required to increment view")
		return
	}

	if !ValidateObjectID(blogID) {
		logger.Error("blogID must be a valid Mongo Object ID")

		response.SendBadRequestError(w, "Invalid blogID format")
		return
	}

	objectID, err := primitive.ObjectIDFromHex(blogID)

	if err != nil {
		logger.Warn("Invalid blogID format: %v", err)
	}

	err = service.CreateBlogView(r.Context(), objectID)

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
			logger.Error(fmt.Sprintf("Failed to increment view: %v", err))
			response.SendInternalServerError(w, "Failed to increment view")
			return
		}
	}

	response.SendCreated(w)

}

func (bh *BlogHandler) AddShare(w http.ResponseWriter, r *http.Request) {

	blogID := r.URL.Query().Get("blogID")

	if blogID == "" {
		logger.Error("blog ID is required to incremented share, Please pass blogID with parema in request")
		response.SendBadRequestError(w, "blog ID is required to incremented share")
		return
	}

	if !ValidateObjectID(blogID) {
		logger.Error("blogID must be a valid Mongo Object ID")

		response.SendBadRequestError(w, "Invalid blogID format")
		return
	}

	objectID, err := primitive.ObjectIDFromHex(blogID)

	if err != nil {
		logger.Warn("Invalid blogID format: %v", err)
	}

	err = service.CreateBlogShare(r.Context(), objectID)

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
			logger.Error(fmt.Sprintf("Failed to increment share: %v", err))
			response.SendInternalServerError(w, "Failed to increment share")
			return
		}
	}

	response.SendCreated(w)

}

func (bh *BlogHandler) DeleteBlogHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the Blog ID from the URL parameter
	id := chi.URLParam(r, "id")
	if id == "" {
		logger.Error("Blog ID is missing in the request")
		response.SendBadRequestError(w, "Blog ID is required to delete the blog")
		return
	}

	// Attempt to delete the blog using the DeleteOneBlog function
	err := service.DeleteOneBlog(r.Context(), id)
	if err != nil {
		// Check for specific error cases to provide meaningful responses
		if strings.Contains(err.Error(), "invalid blog ID format") {
			logger.Error("Invalid blog ID format", "blogID", id, "error", err)
			response.SendBadRequestError(w, "Invalid blog ID format")
			return
		}

		if strings.Contains(err.Error(), "no blog found") {
			logger.Warn("Blog not found for deletion", "blogID", id)
			response.SendNotFoundError(w, "Blog not found with the given ID")
			return
		}

		// Log and send a generic server error for unexpected issues
		logger.Error("Failed to delete blog", "blogID", id, "error", err)
		response.SendInternalServerError(w, "Failed to delete the blog")
		return
	}

	// Log success and send a success response
	logger.Info("Blog deleted successfully", "blogID", id)
	response.SendSuccess(w, "Blog deleted successfully", http.StatusOK)
}

func (bh *BlogHandler) UpdateBlogHandler(w http.ResponseWriter, r *http.Request) {
	// Extract Blog ID from the URL
	id := chi.URLParam(r, "id")
	if id == "" {
		logger.Error("Blog ID is missing in the request")
		response.SendBadRequestError(w, "Blog ID is required to update the blog")
		return
	}

	// Parse the request body into a BlogUpdate struct
	var blogUpdate models.BlogUpdate
	if err := json.NewDecoder(r.Body).Decode(&blogUpdate); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload for blog update")
		return
	}

	// Validate the parsed BlogUpdate struct
	validationErrors := validator.ValidateUpdateBlog(&blogUpdate)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed for Blog attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	// Call the service to update the blog
	err := service.UpdateOneBlog(r.Context(), id, &blogUpdate)
	if err != nil {
		// Handle specific error cases from UpdateOneBlog
		switch {
		case strings.Contains(err.Error(), "invalid blog ID format"):
			logger.Error("Invalid blog ID format", "blogID", id, "error", err)
			response.SendBadRequestError(w, "Invalid blog ID format")
			return
		case strings.Contains(err.Error(), "no blog found"):
			logger.Warn("Blog not found for update", "blogID", id)
			response.SendNotFoundError(w, "Blog not found with the given ID")
			return
		case strings.Contains(err.Error(), "did not modify any fields"):
			logger.Warn("No changes made to the blog during update", "blogID", id)
			response.SendBadRequestError(w, "No changes made to the blog")
			return
		default:
			// Generic error handler for unexpected issues
			logger.Error("Failed to update blog", "blogID", id, "error", err)
			response.SendInternalServerError(w, "Failed to update the blog")
			return
		}
	}

	// Success response
	logger.Info("Blog updated successfully", "blogID", id)
	response.SendSuccess(w, "Blog updated successfully", http.StatusOK)
}

func (bh *BlogHandler) ApproveBlogHandler(w http.ResponseWriter, r *http.Request) {
	// Extract Blog ID from the URL
	id := chi.URLParam(r, "id")
	if id == "" {
		logger.Error("Blog ID is missing in the request")
		response.SendBadRequestError(w, "Blog ID is required to approve the blog")
		return
	}

	// Parse the request body into a BlogApproval struct
	var blogApproval models.BlogApproval
	if err := json.NewDecoder(r.Body).Decode(&blogApproval); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload for blog approval")
		return
	}

	// Validate the parsed BlogApproval struct
	validationErrors := validator.ValidateApproveBlog(&blogApproval)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed for BlogApproval input", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	// Call the service to approve the blog
	err := service.ApproveBlog(r.Context(), id, blogApproval.Approve)
	if err != nil {
		// Handle specific error cases from ApproveBlog
		switch {
		case strings.Contains(err.Error(), "invalid blog ID format"):
			logger.Error("Invalid blog ID format", "blogID", id, "error", err)
			response.SendBadRequestError(w, "Invalid blog ID format")
			return
		case strings.Contains(err.Error(), "no blog found"):
			logger.Warn("Blog not found for approval", "blogID", id)
			response.SendNotFoundError(w, "Blog not found with the given ID")
			return
		case strings.Contains(err.Error(), "did not modify any fields"):
			logger.Warn("No changes made to the blog during approval", "blogID", id)
			response.SendBadRequestError(w, "No changes made to the blog")
			return
		default:
			// Generic error handler for unexpected issues
			logger.Error("Failed to approve blog", "blogID", id, "error", err)
			response.SendInternalServerError(w, "Failed to approve the blog")
			return
		}
	}

	// Success response
	logger.Info("Blog approved successfully", "blogID", id)
	response.SendSuccess(w, "Blog approved successfully", http.StatusOK)
}

func (bh *BlogHandler) PublishBlogHandler(w http.ResponseWriter, r *http.Request) {
	// Extract Blog ID from the URL
	id := chi.URLParam(r, "id")
	if id == "" {
		logger.Error("Blog ID is missing in the request")
		response.SendBadRequestError(w, "Blog ID is required to publish the blog")
		return
	}

	// Parse the request body into a BlogPublish struct
	var blogPublish models.BlogPublish
	if err := json.NewDecoder(r.Body).Decode(&blogPublish); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload for blog publishing")
		return
	}

	// Validate the parsed BlogPublish struct
	validationErrors := validator.ValidatePublishBlog(&blogPublish)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed for BlogPublish input", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	// Call the service to publish the blog
	err := service.PublishBlog(r.Context(), id, blogPublish.Publish)
	if err != nil {
		// Handle specific error cases from PublishBlog
		switch {
		case strings.Contains(err.Error(), "invalid blog ID format"):
			logger.Error("Invalid blog ID format", "blogID", id, "error", err)
			response.SendBadRequestError(w, "Invalid blog ID format")
			return
		case strings.Contains(err.Error(), "no blog found"):
			logger.Warn("Blog not found for publishing", "blogID", id)
			response.SendNotFoundError(w, "Blog not found with the given ID")
			return
		case strings.Contains(err.Error(), "did not modify any fields"):
			logger.Warn("No changes made to the blog during publishing", "blogID", id)
			response.SendBadRequestError(w, "No changes made to the blog")
			return
		default:
			// Generic error handler for unexpected issues
			logger.Error("Failed to publish blog", "blogID", id, "error", err)
			response.SendInternalServerError(w, "Failed to publish the blog")
			return
		}
	}

	// Success response
	logger.Info("Blog published successfully", "blogID", id)
	response.SendSuccess(w, "Blog published successfully", http.StatusOK)
}
