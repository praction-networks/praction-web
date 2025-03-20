package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"
	"github.com/praction-networks/quantum-ISP365/webapp/src/validator"
)

type UserIntrest struct{}

func (h *UserIntrest) ShowIntresrtHandler(w http.ResponseWriter, r *http.Request) {

	var userIntrest *models.UserInterest

	if err := json.NewDecoder(r.Body).Decode(&userIntrest); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload, valid JSON required for Plan Create")
		return
	}

	logger.Info("Successfully parsed request body of plan, Proceeding for Validation")

	validationErrors := validator.ValidateShowUserIntrest(userIntrest)

	if len(validationErrors) > 0 {
		logger.Error("Validation failed for user Intrest attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	//Checking if User already shown intrest

	userData, result, err := service.CheckUserIntrestDuplicate(*userIntrest)

	if err != nil {
		logger.Error("Internal server error occurred While Connecting MongoDB", "Error", err)
		response.SendInternalServerError(w, "An unexpected error occurred. Please try again later.")
		return
	}

	switch result {
	case "Verified":
		logger.Info(fmt.Sprintf("User %s is already registered and verified with mobile %s and email %s.", userIntrest.Name, userIntrest.Mobile, userIntrest.Email))
		response.SendSuccess(w, "Thank you for your interest in our services. We have already received your request, and our representative will contact you shortly.", http.StatusAlreadyReported)

	case "NotVerified":
		logger.Info(fmt.Sprintf("User %s is already registered but not verified with mobile %s and email %s. Reuesting User again to submit OTP for sussfull registration", userIntrest.Name, userIntrest.Mobile, userIntrest.Email)) // Assume a utility function for OTP generation

		err = utils.SendOTP(userIntrest.Email, userIntrest.Mobile, userData.OTP, "text", true)

		if err != nil {
			logger.Error("Failed to send OTP to the user", "error", err)

			errorDetails := []response.ErrorDetail{
				{
					Field:   "OTP",                                         // The field that caused the error
					Message: "Failed to send OTP. Please try again later.", // The error message
				},
			}
			response.SendError(w, errorDetails, http.StatusInternalServerError)
			return
		}

		logger.Info("OTP successfully sent to user", "Mobile", userIntrest.Mobile, "Email:", userIntrest.Email)

		response.SendSuccess(w, "Please Enter OTP recived on Mobile and Email for complete the process, ", http.StatusOK)
		return

	case "NotFound":
		// Send OTP to the user
		otp := utils.GenerateRandomOTP(6) // Assume a utility function for OTP generation
		userIntrest.OTP = otp

		err = utils.SendOTP(userIntrest.Email, userIntrest.Mobile, otp, "text", false)

		if err != nil {
			logger.Error("Failed to send OTP to the user", "error", err)

			errorDetails := []response.ErrorDetail{
				{
					Field:   "OTP",                                         // The field that caused the error
					Message: "Failed to send OTP. Please try again later.", // The error message
				},
			}
			response.SendError(w, errorDetails, http.StatusInternalServerError)
			return
		}

		logger.Info("OTP successfully sent to user", "Mobile", userIntrest.Mobile)

		//Updaing user in database

		// Call service layer to create the plan
		err = service.CreateUserIntrest(r.Context(), *userIntrest)

		// If there was an error creating the plan, handle it appropriately
		if err != nil {
			logger.Error("Error Registering User Intial Intrest", "error", err)

			// Create an array of ErrorDetail with relevant fields
			errorDetails := []response.ErrorDetail{
				{
					Field:   "User Intrest",                                             // The field that caused the error
					Message: "Failed to Register User INtrest. Please try again later.", // The error message
				},
			}

			// Pass the error details and the HTTP status code to SendError
			response.SendError(w, errorDetails, http.StatusInternalServerError)
			return
		}

		logger.Info("Initial User interest successfully created", "UUID", userIntrest.UUID)

		// Respond with success
		response.SendSuccess(w, "Please Enter OTP recived on MObile and Email for complete the process, ", http.StatusOK)
		return
	}

}

func (h *UserIntrest) VerifyUserOTP(w http.ResponseWriter, r *http.Request) {
	var userOTP *models.UserOTPVerify

	if err := json.NewDecoder(r.Body).Decode(&userOTP); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload for OTP Verification, valid JSON required for Plan Create")
		return
	}

	logger.Info("Successfully parsed request body of plan, Proceeding for Validation")

	validationErrors := validator.ValidateShowUserIntrestOTP(userOTP)

	if len(validationErrors) > 0 {
		logger.Error("Validation failed for user Intrest attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	userData, result, err := service.UserInterestVerifyOTPAndUpdate(*userOTP)

	if err != nil {
		logger.Error("Failed to Verify User OTP for Intrest Verification", "Error", err)

		response.SendInternalServerError(w, "Failed to Very Uuser, Internal Seerver error occured")
		return
	}

	if result == "User-not-found" {
		logger.Error(fmt.Sprintf("User Not found in databse with Mobile: %s or Email: %s", userOTP.Mobile, userOTP.Email))
		response.SendNotFoundError(w, "User Not Found")
		return
	}

	if result == "User-already-verified" {

		logger.Info(fmt.Sprintf("User is already registered and verified with mobile %s and email %s.", userOTP.Mobile, userOTP.Email))
		response.SendSuccess(w, "Thank you for your interest in our services. We have already received your request, and our representative will contact you shortly.", http.StatusAlreadyReported)
		return
	}

	if result == "OTP-expired" {

		logger.Info(fmt.Sprintf("OTP expired for mobile %s and email %s.", userOTP.Mobile, userOTP.Email))
		response.SendSuccess(w, "OTP is already expired. please genrate a new OTP.", http.StatusAlreadyReported)
		return
	}

	if result == "OTP-mismatch" {
		logger.Warn("User OTP is not Matched with Request OTP")

		response.SendUnauthorizedError(w, "OTP Not Mached")
		return
	}

	if result == "Verification-successful" {
		logger.Info("OTP Verified, Send Deatils to Team")

		err := utils.SendSuccessMailIntrest(userData)

		if err != nil {
			logger.Error("Fail to Submit User details over mail. please connect woth your administrator for further details")
		}

		response.SendSuccess(w, "OTP Verified sussfuly", http.StatusOK)
		return
	}
}

func (h *UserIntrest) ResendUserOTP(w http.ResponseWriter, r *http.Request) {

	var userOTPResend *models.UserOTPResend

	if err := json.NewDecoder(r.Body).Decode(&userOTPResend); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload for OTP Verification, valid JSON required for Plan Create")
		return
	}

	logger.Info("Successfully parsed request body of plan, Proceeding for Validation")

	validationErrors := validator.ValidateShowUserIntrestOTPResnd(userOTPResend)

	if len(validationErrors) > 0 {
		logger.Error("Validation failed for user Intrest attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	user, err := service.GetUserDetailsFromDB(userOTPResend)

	if err != nil {
		logger.Error(fmt.Sprintf("User Not found in databse with Mobile: %s or Email: %s", userOTPResend.Mobile, userOTPResend.Email))
		response.SendNotFoundError(w, "User Not Found")
		return
	}

	err = utils.SendOTP(user.Email, user.Mobile, user.OTP, userOTPResend.Resend, true)

	if err != nil {
		logger.Error("Failed to send OTP to the user", "error", err)

		errorDetails := []response.ErrorDetail{
			{
				Field:   "OTP",                                         // The field that caused the error
				Message: "Failed to send OTP. Please try again later.", // The error message
			},
		}
		response.SendError(w, errorDetails, http.StatusInternalServerError)
		return
	}

	logger.Info("OTP successfully sent to user", "Mobile", user.Mobile, "Email:", user.Email)

	response.SendSuccess(w, "Please Enter OTP recived on Mobile and Email for complete the process, ", http.StatusOK)
}

func (h *UserIntrest) GetALl(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	params, err := utils.ParseQueryParams(r.URL.Query())
	if err != nil {
		logger.Error("Error parsing query parameters", "Error", err)
		response.SendBadRequestError(w, "Invalid query parameters")
		return
	}

	// Fetch data from the service
	users, err := service.GetAllIntrestUserService(ctx, params, "UserIntrest")
	if err != nil {
		logger.Error("Error fetching user interest data", "Error", err)
		response.SendNotFoundError(w, "Failed to fetch user interest data")
		return
	}

	response.SendSuccess(w, users, http.StatusOK)
}

func (h *UserIntrest) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from URL parameters
	id := chi.URLParam(r, "id")
	if id == "" {
		logger.Error("User ID is missing in the request")
		response.SendBadRequestError(w, "User ID is required to update the User Interest")
		return
	}

	// Parse the request body into a UserInterestUpdate struct
	var userUpdate models.UserInterestUpdate
	if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload for User Interest update")
		return
	}

	// Validate the parsed UserInterestUpdate struct
	validationErrors := validator.ValidateUserInterestUpdate(&userUpdate)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed for User Interest Update attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	// Call the UpdateUserInterest function to update the database
	err := service.UpdateUserInterest(ctx, id, &userUpdate)
	if err != nil {
		logger.Error("Failed to update User Interest", "userID", id, "error", err)
		if err.Error() == "no user interest found with the given ID: "+id {
			response.SendNotFoundError(w, "User Interest not found with the given ID")
		} else if err.Error() == "no valid fields provided for update" {
			response.SendBadRequestError(w, "No valid fields provided for update")
		} else {
			response.SendInternalServerError(w, "Failed to update User Interest")
		}
		return
	}

	// Send success response
	response.SendSuccess(w, "User Interest updated successfully", http.StatusOK)
	logger.Info("User Interest updated successfully", "userID", id)
}
