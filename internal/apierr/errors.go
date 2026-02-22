package apierr

import (
	"encoding/json"
	"net/http"
)

type Code string

const (
	CodeInvalidInput Code = "invalid_input"
	CodeNotFound     Code = "not_found"
	CodeInternal     Code = "internal_error"
)

type APIError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func Write(w http.ResponseWriter, status int, code Code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIError{Code: code, Message: message})
}

func BadRequest(w http.ResponseWriter, message string) {
	Write(w, http.StatusBadRequest, CodeInvalidInput, message)
}

func NotFound(w http.ResponseWriter, message string) {
	Write(w, http.StatusNotFound, CodeNotFound, message)
}

func Internal(w http.ResponseWriter) {
	Write(w, http.StatusInternalServerError, CodeInternal, "internal server error")
}
