package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int
	Message string
}

func Error(w http.ResponseWriter, error string, code int) {
	errResponse := ErrorResponse{
		Code:    code,
		Message: error,
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errResponse)
}
