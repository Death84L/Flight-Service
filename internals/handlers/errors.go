package handlers

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequest(msg string) *AppError {
	return &AppError{Message: msg, Code: http.StatusBadRequest}
}

func NewInternalError(msg string) *AppError {
	return &AppError{Message: msg, Code: http.StatusInternalServerError}
}

func NewNotFound(msg string) *AppError {
	return &AppError{Message: msg, Code: http.StatusNotFound}
}

// WriteError writes an error response to the client
func WriteError(w http.ResponseWriter, err *AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}
