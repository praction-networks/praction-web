package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	middlewares "github.com/praction-networks/quantum-ISP365/webapp/src/middleware"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
	"github.com/praction-networks/quantum-ISP365/webapp/src/validator"
)

type User struct{}

// Login authenticates the user and returns a JWT
func (a *User) Login(w http.ResponseWriter, r *http.Request) {

	var user models.LoginUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	logger.Info("Successfully parsed request body, Proceeding for Validation")

	validationErrors := validator.ValidateLoginUser(&user)

	if len(validationErrors) > 0 {
		logger.Error("Validation failed for admin user attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
	}

	logger.Info("User attributes are valid")

	// Validate user credentials
	// Authenticate user credentials
	authenticatedUser, err := service.AuthenticateUser(r.Context(), user.Username, user.Password)
	if err != nil {
		logger.Error(fmt.Sprintf("Authentication failed: %v", err))
		// Send an error response using the custom response format
		response.SendError(w, []response.ErrorDetail{
			{
				Field:   "authentication",
				Message: "Invalid username or password",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Generate JWT token for authenticated user
	token, err := service.GenerateJWT(authenticatedUser)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to generate JWT: %v", err))
		// Send an error response for internal server error
		response.SendError(w, []response.ErrorDetail{
			{
				Field:   "server",
				Message: "Internal server error",
			},
		}, http.StatusInternalServerError)
		return
	}

	// Send the success response with the generated JWT token
	resp := map[string]string{"token": token}
	response.SendSuccess(w, resp, http.StatusOK)

}

func (a *User) Create(w http.ResponseWriter, r *http.Request) {

	var user models.Admin

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	logger.Info("Successfully parsed request body, Proceeding for Validation")

	validationErrors := validator.ValidateCreateAdmin(&user)

	if len(validationErrors) > 0 {
		logger.Error("Validation failed for admin user createa attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	logger.Info("Admin attributes are valid")

	err := service.CreateUser(r.Context(), user)
	if err != nil {

		if service.IsDuplicateKeyError(err) {
			logger.Warn("User Already Registred with request user details", err)
			response.SendError(w, []response.ErrorDetail{
				{
					Field:   "duplicate user",
					Message: "User is already created with user requeest details",
				},
			}, http.StatusConflict)
			return

		}

		logger.Error(fmt.Sprintf("Failed to create admin: %v", err))
		// Send an error response using the custom response format
		response.SendError(w, []response.ErrorDetail{
			{
				Field:   "user createa",
				Message: "Failed to create Admin user",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Authenticate user credentials
	authenticatedUser, err := service.AuthenticateUser(r.Context(), user.Username, user.Password)
	if err != nil {
		logger.Error(fmt.Sprintf("Authentication failed: %v", err))
		// Send an error response using the custom response format
		response.SendError(w, []response.ErrorDetail{
			{
				Field:   "authentication",
				Message: "Invalid username or password",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Generate JWT token for authenticated user
	token, err := service.GenerateJWT(authenticatedUser)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to generate JWT: %v", err))
		// Send an error response for internal server error
		response.SendError(w, []response.ErrorDetail{
			{
				Field:   "server",
				Message: "Internal server error",
			},
		}, http.StatusInternalServerError)
		return
	}

	// Send the success response with the generated JWT token
	resp := map[string]string{"token": token}
	response.SendSuccess(w, resp, http.StatusOK)

}

func (a *User) Logout(w http.ResponseWriter, r *http.Request) {
	// Send the success response with the generated JWT token
	resp := map[string]string{"token": ""}
	response.SendSuccess(w, resp, http.StatusOK)

}

// GetOne retrieves a single user based on the provided user ID or username
func (a *User) GetOne(w http.ResponseWriter, r *http.Request) {
	// Get the user ID or username from the URL path
	vars := chi.URLParam(r, "id") // assuming the URL is /user/{id}

	logger.Info("Finding Admin USer With ", "ID", vars)

	// Call the service to fetch the user
	user, err := service.GetOneUser(r.Context(), vars)
	if err != nil {
		// Handle error: user not found or internal error
		if err.Error() == "user not found" {
			// Return a 404 not found error if the user does not exist
			response.SendError(w, []response.ErrorDetail{
				{
					Field:   "user",
					Message: "User not found",
				},
			}, http.StatusNotFound)
		} else {
			// Internal server error for other issues
			response.SendError(w, []response.ErrorDetail{
				{
					Field:   "user",
					Message: "Internal server error",
				},
			}, http.StatusInternalServerError)
		}
		return
	}

	// Return the user data in the response
	response.SendSuccess(w, user, http.StatusOK)
}

// GetAll retrieves all users from the database
func (a *User) GetAll(w http.ResponseWriter, r *http.Request) {
	// Call the service to fetch all users
	users, err := service.GetAllUsers(r.Context())
	if err != nil {
		// Handle error: internal server error
		logger.Error("Failed to fetch users", "error", err)
		response.SendError(w, []response.ErrorDetail{
			{
				Field:   "user",
				Message: "Failed to fetch users due to internal error",
			},
		}, http.StatusInternalServerError)
		return
	}

	// Send the success response with the list of users
	response.SendSuccess(w, users, http.StatusOK)
}

// Update updates the user details (admin user)
func (a *User) Update(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL
	userID := chi.URLParam(r, "id") // assuming the URL is /user/{id}

	// Parse the request body to get the updated user data
	var updatedAdmin models.UpdateAdmin
	if err := json.NewDecoder(r.Body).Decode(&updatedAdmin); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	// Validate the updated data
	validationErrors := validator.ValidateUpdateAdmin(&updatedAdmin)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed for admin user attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	// Call the service to update the user
	updatedUser, err := service.UpdateUser(r.Context(), userID, &updatedAdmin)
	if err != nil {
		// Handle error: user not found or internal error
		if err.Error() == "user not found" {
			response.SendError(w, []response.ErrorDetail{
				{
					Field:   "user",
					Message: "User not found",
				},
			}, http.StatusNotFound)
		} else {
			// Internal server error for other issues
			response.SendError(w, []response.ErrorDetail{
				{
					Field:   "user",
					Message: "Internal server error",
				},
			}, http.StatusInternalServerError)
		}
		return
	}

	// Send the success response with the updated user data
	response.SendSuccess(w, updatedUser, http.StatusOK)
}

// SetPassword allows a user to update their password
func (a *User) SetPassword(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(middlewares.UserContextKey).(*models.Admin)

	if !ok || user == nil {
		// Handle the case where the user is not found in context
		response.SendError(w, []response.ErrorDetail{
			{
				Field:   "user",
				Message: "User not authenticated",
			},
		}, http.StatusUnauthorized)
		return
	}
	// Decode the request body to get the current and new passwords
	var passwordData struct {
		CurrentPassword    string `json:"current_password"`
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password"`
	}

	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&passwordData); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	// Validate that the new password matches the confirmation
	if passwordData.NewPassword != passwordData.ConfirmNewPassword {
		response.SendError(w, []response.ErrorDetail{
			{
				Field:   "new_password",
				Message: "New password and confirmation do not match",
			},
		}, http.StatusBadRequest)
		return
	}

	// Validate that the new password meets any necessary complexity requirements (e.g., minimum length)
	if len(passwordData.NewPassword) < 8 {
		response.SendError(w, []response.ErrorDetail{
			{
				Field:   "new_password",
				Message: "New password must be at least 8 characters long",
			},
		}, http.StatusBadRequest)
		return
	}

	err := service.SetUserPassword(r.Context(), user.Username, passwordData.CurrentPassword, passwordData.NewPassword)
	if err != nil {
		if err.Error() == "incorrect current password" {
			response.SendError(w, []response.ErrorDetail{
				{
					Field:   "current_password",
					Message: "Current password is incorrect",
				},
			}, http.StatusUnauthorized)
		} else {
			// Internal server error for any other issues
			response.SendError(w, []response.ErrorDetail{
				{
					Field:   "password",
					Message: "Failed to update password due to internal error",
				},
			}, http.StatusInternalServerError)
		}
		return
	}

	// Send a success response
	response.SendSuccess(w, map[string]string{"message": "Password updated successfully"}, http.StatusOK)
}

func (a *User) SetPasswordByAdmin(w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "id")

	var passwordData struct {
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password"`
	}
	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&passwordData); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}
	// Validate that the new password matches the confirmation
	// Validate that the new password matches the confirmation
	if passwordData.NewPassword != passwordData.ConfirmNewPassword {
		response.SendError(w, []response.ErrorDetail{
			{
				Field:   "new_password",
				Message: "New password and confirmation do not match",
			},
		}, http.StatusBadRequest)
		return
	}

	// Ensure that the new password meets the minimum length requirement
	if len(passwordData.NewPassword) < 8 {
		response.SendError(w, []response.ErrorDetail{
			{
				Field:   "new_password",
				Message: "New password must be at least 8 characters long",
			},
		}, http.StatusBadRequest)
		return
	}

	// Call the service to update the password for the user
	err := service.SetUserPasswordByAdmin(r.Context(), userID, passwordData.NewPassword)
	if err != nil {
		// Handle specific errors (e.g., user not found)
		if err.Error() == "user not found" {
			response.SendError(w, []response.ErrorDetail{
				{
					Field:   "user",
					Message: "User not found",
				},
			}, http.StatusNotFound)
		} else {
			// Internal server error for any other issues
			response.SendError(w, []response.ErrorDetail{
				{
					Field:   "password",
					Message: "Failed to update password due to internal error",
				},
			}, http.StatusInternalServerError)
		}
		return
	}

	// Send a success response
	response.SendSuccess(w, map[string]string{"message": "Password updated successfully"}, http.StatusOK)
}
