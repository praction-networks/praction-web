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

type UserReferal struct{}

func (ur *UserReferal) ReferUser(w http.ResponseWriter, r *http.Request) {
	var userReferral models.UserRefrence

	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&userReferral); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload")
		return
	}

	logger.Info("Successfully parsed request body, Proceeding for Validation")

	// Validate user referral
	validationErrors := validator.ValidateUserRefrence(&userReferral)
	if len(validationErrors) > 0 {
		logger.Error("Validation failed for User Referral Request", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	logger.Info("User Referral attributes are valid")

	otp := utils.GenerateRandomOTP(6)
	userReferral.OTP = otp

	// Generate OTP and send it to the user

	logger.Info("OTP successfully sent to user", "Mobile", userReferral.ReferedBy.Mobile)

	// Call the service to create the referral
	err := service.CreateUserReferal(r.Context(), userReferral)
	if err != nil {
		logger.Error("Error Registering User Referral", "error", err)

		// Check for duplicate errors
		if service.ContainsSubstring(err.Error(), "duplicate email") {
			handleDuplicateFieldError(w, "Email", "This email is already associated with a referred person and cannot be referred again.")
			return
		} else if service.ContainsSubstring(err.Error(), "duplicate mobile") {
			handleDuplicateFieldError(w, "Mobile", "This mobile number is already associated with a referred person and cannot be referred again.")
			return
		}

		// Handle any other errors
		errorDetails := []response.ErrorDetail{
			{
				Field:   "User Referral",
				Message: "Failed to Register User Referral. Please try again later.",
			},
		}
		response.SendError(w, errorDetails, http.StatusInternalServerError)
		return
	}

	err = utils.SendOTP(userReferral.ReferedBy.Email, userReferral.ReferedBy.Mobile, otp, "text", false)
	if err != nil {
		logger.Error("Failed to send OTP to the Referral", "error", err)

		errorDetails := []response.ErrorDetail{
			{
				Field:   "OTP",
				Message: "Failed to send OTP. Please try again later.",
			},
		}
		response.SendError(w, errorDetails, http.StatusInternalServerError)
		return
	}

	logger.Info("Initial User Referral successfully created", "UUID", userReferral.UUID)

	// Respond with success
	response.SendSuccess(w, "Please enter the OTP received on Mobile and Email to complete the process.", http.StatusOK)
}

func (ur *UserReferal) VerifyUserOTP(w http.ResponseWriter, r *http.Request) {
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

	userData, result, err := service.UserReferrelVerifyOTPAndUpdate(*userOTP)

	if err != nil {
		logger.Error("Failed to Verify User OTP for Referrel Verification", "Error", err)

		response.SendInternalServerError(w, "Failed to VVerify  user, Internal Seerver error occured")
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

		err := utils.SendSuccessMailReferrel(userData)

		if err != nil {
			logger.Error("Fail to Submit User details over mail. please connect woth your administrator for further details")
		}

		response.SendSuccess(w, "OTP Verified sussfuly", http.StatusOK)
		return
	}
}

func (ur *UserReferal) ResendUserOTP(w http.ResponseWriter, r *http.Request) {

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

	user, err := service.GetUserDetailsFromDBRefeeral(userOTPResend)

	if err != nil {
		logger.Error(fmt.Sprintf("User Not found in databse with Mobile: %s or Email: %s", userOTPResend.Mobile, userOTPResend.Email))
		response.SendNotFoundError(w, "User Not Found")
		return
	}

	err = utils.SendOTP(user.ReferedBy.Email, user.ReferedBy.Mobile, user.OTP, userOTPResend.Resend, true)

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

	logger.Info("OTP successfully sent to user", "Mobile", user.ReferedBy.Mobile, "Email:", user.ReferedBy.Email)

	response.SendSuccess(w, "Please Enter OTP recived on Mobile and Email for complete the process, ", http.StatusOK)
}

func (ur *UserReferal) GetALl(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	params, err := utils.ParseQueryParams(r.URL.Query())
	if err != nil {
		logger.Error("Error parsing query parameters", "Error", err)
		response.SendBadRequestError(w, "Invalid query parameters")
		return
	}

	// Fetch data from the service
	users, err := service.GetAllIntrestUserService(ctx, params, "UserReferal")
	if err != nil {
		logger.Error("Error fetching user interest data", "Error", err)
		response.SendNotFoundError(w, "Failed to fetch user interest data")
		return
	}

	response.SendSuccess(w, users, http.StatusOK)
}

func handleDuplicateFieldError(w http.ResponseWriter, field, message string) {
	errorDetails := []response.ErrorDetail{
		{
			Field:   field,
			Message: message,
		},
	}
	response.SendError(w, errorDetails, http.StatusBadRequest)
}
