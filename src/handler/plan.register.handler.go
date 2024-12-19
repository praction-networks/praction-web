package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"
	"github.com/praction-networks/quantum-ISP365/webapp/src/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlanRegister struct{}

func (p *PlanRegister) RegisterUserIntrestForPlan(w http.ResponseWriter, r *http.Request) {
	var userRequestRegister models.AvailableUserRequest

	if err := json.NewDecoder(r.Body).Decode(&userRequestRegister); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload, valid JSON required for User Creation")
		return
	}

	logger.Info("Successfully parsed request body of User Request for Plan, Proceeding for Validation")

	validationErrors := validator.ValidateAvailableUser(&userRequestRegister)

	if len(validationErrors) > 0 {
		logger.Error("Validation failed for user intrest for plan", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	logger.Info("User attributes are verified and procedding for User register")

	userData, result, err := service.CheckUserIntrestForPlanDuplicate(&userRequestRegister)

	if err != nil {
		logger.Error("Internal server error occurred While Connecting MongoDB", "Error", err)
		response.SendInternalServerError(w, "An unexpected error occurred. Please try again later.")
		return
	}

	switch result {
	case "Verified":
		logger.Info(fmt.Sprintf("User %s is already registered and verified with mobile %s and email %s.", userRequestRegister.FirstName, userRequestRegister.Mobile, userRequestRegister.Email))
		response.SendSuccess(w, "Thank you for your interest in our services. We have already received your request, and our representative will contact you shortly.", http.StatusAlreadyReported)

	case "NotVerified":
		logger.Info(fmt.Sprintf("User %s is already registered but not verified with mobile %s and email %s. Reuesting User again to submit OTP for sussfull registration", userRequestRegister.FirstName, userRequestRegister.Mobile, userRequestRegister.Email)) // Assume a utility function for OTP generation

		err = utils.SendOTP(userData.Email, userData.Mobile, userData.OTP, "text", true)

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

		logger.Info("OTP successfully sent to user", "Mobile", userData.Mobile, "Email:", userData.Email)

		response.SendSuccess(w, "Please Enter OTP recived on Mobile and Email for complete the process, ", http.StatusOK)
		return

	case "NotFound":
		// Send OTP to the user
		otp := utils.GenerateRandomOTP(6) // Assume a utility function for OTP generation
		userRequestRegister.OTP = otp

		err = utils.SendOTP(userRequestRegister.Email, userRequestRegister.Mobile, otp, "text", false)

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

		logger.Info("OTP successfully sent to user", "Mobile", userRequestRegister.Mobile)

		//Updaing user in database

		// Call service layer to create the plan
		err = service.CreateUserIntrestForPlan(r.Context(), userRequestRegister)

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

		logger.Info("Initial User interest successfully created", "UUID", userRequestRegister.UUID)

		// Respond with success
		response.SendSuccess(w, "Please Enter OTP recived on MObile and Email for complete the process, ", http.StatusOK)
		return
	}

}

func (p *PlanRegister) VerifyUserOTPForPlan(w http.ResponseWriter, r *http.Request) {
	var userOTP *models.UserOTPPlanVerify

	if err := json.NewDecoder(r.Body).Decode(&userOTP); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload for OTP Verification, valid JSON required for user registration for plan")
		return
	}

	logger.Info("Successfully parsed request body of user intrest for plan, Proceeding for Validation")

	validationErrors := validator.ValidateShowUserIntrestOTPPlan(userOTP)

	if len(validationErrors) > 0 {
		logger.Error("Validation failed for user Intrest attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	userData, result, err := service.UserInterestVerifyOTPAndUpdateForPlan(*userOTP)

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
		logger.Info("OTP Verified, Send Deatils to Team and User")

		planDetails, err := service.GetOnePlanByUUID(r.Context(), userData.PlanID)

		if err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				logger.Warn("No plans found in the database.")
			case strings.Contains(err.Error(), "decoding plan"):
				logger.Warn("Decoding error", "Error:", err)
			case strings.Contains(err.Error(), "error fetching plans"):
				logger.Warn("Database fetch error:", "Error:", err)
			default:
				logger.Error("An unexpected error occurred:", "Error:", err)
			}
		} else {
			logger.Info("Retrieved PlanSpecific:", "Plan", planDetails)
		}

		err = utils.SendSuccessMailIntrestForPlan(userData, planDetails)

		if err != nil {
			logger.Error("Fail to Submit User details over mail. please connect woth your administrator for further details")
		}

		response.SendSuccess(w, "OTP Verified sussfuly", http.StatusOK)
		return
	}
}

func (p *PlanRegister) ResendUserOTPForPlan(w http.ResponseWriter, r *http.Request) {

	var userOTPResend *models.UserOTPPlanResend

	if err := json.NewDecoder(r.Body).Decode(&userOTPResend); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload for OTP Verification, valid JSON required for Plan Create")
		return
	}

	logger.Info("Successfully parsed request body of plan, Proceeding for Validation")

	validationErrors := validator.ValidateShowUserIntrestOTPResndPlan(userOTPResend)

	if len(validationErrors) > 0 {
		logger.Error("Validation failed for user Intrest attributes", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	user, err := service.GetUserDetailsFromDBForPlan(userOTPResend)

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

	response.SendSuccess(w, "Please Enter OTP recived on Mobile and Email for complete the Registraion Process, ", http.StatusOK)
}

func (p *PlanRegister) GetALl(w http.ResponseWriter, r *http.Request) {

	// Parse query parameters
	query := r.URL.Query()

	// Parse pagination
	page, err := strconv.Atoi(query.Get("Page"))
	if err != nil || page < 1 {
		page = 1 // Default to page 1
	}

	pageSize, err := strconv.Atoi(query.Get("PageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // Default to 10 items per page
	}

	// Parse sorting
	sortField := query.Get("sortField")
	if sortField == "" {
		sortField = "createdAt" // Default sorting field
	}

	sortOrder := 1 // Default to ascending
	if query.Get("sortOrder") == "desc" {
		sortOrder = -1
	}

	// Parse filters
	filters := make(map[string]string)
	for key, values := range query {
		// Skip reserved parameters (pagination and sorting)
		if key == "Page" || key == "PageSize" || key == "sortField" || key == "sortOrder" {
			continue
		}
		// Use the first value for simplicity
		if len(values) > 0 {
			filters[key] = values[0]
		}
	}

	// Build the pagination parameters
	params := utils.PaginationParams{
		Page:      page,
		PageSize:  pageSize,
		SortField: sortField,
		SortOrder: sortOrder,
		Filters:   filters,
	}

	// Fetch data from the service
	users, err := service.GetAllIntrestUserForPlanService(r.Context(), params, "UserPlanIntrest")
	if err != nil {
		logger.Error("Error fetching user interest data", "Error", err)
		response.SendNotFoundError(w, "Failed to fetch user interest data")
		return
	}

	response.SendSuccess(w, users, http.StatusOK)
}
