package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"
	"github.com/praction-networks/quantum-ISP365/webapp/src/validator"
)

type UnAvailableAreaRegister struct{}

func (un *UnAvailableAreaRegister) UnavailavleAreaUserIntrest(w http.ResponseWriter, r *http.Request) {
	var unAvailableArea models.UnAvailableArea

	if err := json.NewDecoder(r.Body).Decode(&unAvailableArea); err != nil {
		logger.Error("Error parsing request body", "error", err)
		response.SendBadRequestError(w, "Invalid request payload, valid JSON required for User Creation")
		return
	}

	logger.Info("Successfully parsed request body of User Request for Plan, Proceeding for Validation")

	validationErrors := validator.ValidateUnAvailableArea(&unAvailableArea)

	if len(validationErrors) > 0 {
		logger.Error("Validation failed for unavailableArea", "validationErrors", validationErrors)
		response.SendError(w, validationErrors, http.StatusBadRequest)
		return
	}

	logger.Info("User attributes are verified and procedding for User register")

	err := service.CreateUnAvailableAreaRequest(r.Context(), unAvailableArea)

	if err != nil {
		logger.Error("Internal server error occurred While Connecting MongoDB", "Error", err)
		response.SendInternalServerError(w, "An unexpected error occurred. Please try again later.")
		return
	}

	logger.Info("User Request Register successfully")
	response.SendCreated(w)

}

func (un *UnAvailableAreaRegister) GetALl(w http.ResponseWriter, r *http.Request) {

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
	users, err := service.GetAllUnavailableAreaUserService(r.Context(), params, "unavailableAreaRequest")
	if err != nil {
		logger.Error("Error fetching user interest data", "Error", err)
		response.SendNotFoundError(w, "Failed to fetch user interest data")
		return
	}

	response.SendSuccess(w, users, http.StatusOK)
}
