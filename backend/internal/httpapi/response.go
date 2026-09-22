package httpapi

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the common response envelope for every endpoint.
// Exactly one of Data and Error is non-nil.
type APIResponse[T any] struct {
	Data  *T        `json:"data"`
	Error *APIError `json:"error"`
}

// APIError contains a stable code for clients and a safe human-readable message.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeSuccess[T any](writer http.ResponseWriter, status int, data T) {
	writeJSON(writer, status, APIResponse[T]{Data: &data})
}

func writeError(writer http.ResponseWriter, status int, code, message string) {
	writeJSON(writer, status, APIResponse[any]{
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

func writeJSON(writer http.ResponseWriter, status int, response any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(response)
}
