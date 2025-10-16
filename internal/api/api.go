package api

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the standardized response wrapper for all API endpoints
type APIResponse struct {
	Error *APIError `json:"error"`
	Data  any       `json:"data"`
}

// APIError represents an API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// BuildRouter creates and configures the main HTTP router
func BuildRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Delegate to resource-specific routers
	buildWorkflowsRouter(mux)
	buildRunsRouter(mux)

	return mux
}

// writeError writes an error response with the given status code
func writeError(
	w http.ResponseWriter,
	code int,
	message string,
) {
	w.WriteHeader(code)
	writeJSON(w, APIResponse{
		Error: &APIError{code, message},
		Data:  nil,
	})
}

// writeData writes a success response with the given data
func writeData(
	w http.ResponseWriter,
	code int,
	data any,
) {
	if data != nil {
		w.WriteHeader(code)
		writeJSON(w, APIResponse{
			Error: nil,
			Data:  data,
		})
	} else {
		w.WriteHeader(code)
	}
}

// writeJSON writes a JSON response
func writeJSON(
	w http.ResponseWriter,
	data any,
) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
