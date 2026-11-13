package response

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error"`
}

// JSON writes a serialized application/json response with the given HTTP status code.
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// Error writes a standardized {"error": "..."} JSON response.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorBody{Error: message})
}

// ErrorWithCode writes a standardized {"code": "...", "error": "..."} JSON response.
func ErrorWithCode(w http.ResponseWriter, status int, code string, message string) {
	JSON(w, status, ErrorBody{Code: code, Error: message})
}
