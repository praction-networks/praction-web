package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
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
