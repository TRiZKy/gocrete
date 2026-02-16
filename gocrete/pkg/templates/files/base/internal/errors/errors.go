package errors

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

func BadRequest(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, message)
}

func NotFound(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusNotFound, message)
}

func InternalServerError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusInternalServerError, message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusUnauthorized, message)
}
