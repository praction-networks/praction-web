package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Status string      `json:"status"`
	Count  *int        `json:"count,omitempty"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status string        `json:"status"`
	Errors []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// WriteSuccess writes a successful JSON response
// SendSuccess sends a JSON success response
// SendSuccess sends a JSON success response
func SendSuccess(w http.ResponseWriter, data interface{}, statusCode int) {
	var count *int

	// Calculate count if the data is a slice or map
	switch v := data.(type) {
	case []interface{}: // Slice of interfaces (equivalent to []any in modern Go)
		countValue := len(v)
		count = &countValue
	case map[string]interface{}: // Map
		countValue := len(v)
		count = &countValue
	default:
		count = nil // No count for unsupported types
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Prepare the response structure
	response := SuccessResponse{
		Status: "success",
		Count:  count, // Include count if available
		Data:   data,
	}

	// Encode and send the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// SendBadRequestError sends a JSON error response for bad requests
func SendBadRequestError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := ErrorResponse{
		Status: "error",
		Errors: []ErrorDetail{
			{
				Field:   "request",
				Message: message,
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

// SendError sends a JSON error response
func SendError(w http.ResponseWriter, errors []ErrorDetail, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := ErrorResponse{
		Status: "error",
		Errors: errors,
	}
	json.NewEncoder(w).Encode(response)
}

// SendInternalServerError sends a JSON error response for server errors
func SendInternalServerError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	response := ErrorResponse{
		Status: "error",
		Errors: []ErrorDetail{
			{
				Field:   "server",
				Message: message,
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

// SendInternalServerError sends a JSON error response for server errors
func SendServiceUnavailableError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)

	response := ErrorResponse{
		Status: "error",
		Errors: []ErrorDetail{
			{
				Field:   "server",
				Message: message,
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func SendNotFoundError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	response := ErrorResponse{
		Status: "error",
		Errors: []ErrorDetail{
			{
				Field:   "Not Found",
				Message: message,
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func SendConflictError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)

	response := ErrorResponse{
		Status: "error",
		Errors: []ErrorDetail{
			{
				Field:   "conflict",
				Message: message,
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

// SendUnsupportedMediaTypeError sends a JSON error response for unsupported media types (HTTP 415)
func SendUnsupportedMediaTypeError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnsupportedMediaType)

	response := ErrorResponse{
		Status: "error",
		Errors: []ErrorDetail{
			{
				Field:   "media_type",
				Message: message,
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func SendUnauthorizedError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	response := ErrorResponse{
		Status: "error",
		Errors: []ErrorDetail{
			{
				Field:   "unauthenticated",
				Message: message,
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

// SendCreated sends a JSON success response for resource creation without requiring data
func SendCreated(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Set the status code to 201

	// Prepare the response structure
	response := struct {
		Status string `json:"status"`
	}{
		Status: "success",
	}

	// Encode and send the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
