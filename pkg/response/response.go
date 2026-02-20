package response

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func JSON(w http.ResponseWriter, status int, success bool, message string, data interface{}, errors interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := JSONResponse{
		Success: success,
		Message: message,
		Data:    data,
		Errors:  errors,
	}

	json.NewEncoder(w).Encode(resp)
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	JSON(w, status, true, message, data, nil)
}

func Error(w http.ResponseWriter, status int, message string, errors interface{}) {
	JSON(w, status, false, message, nil, errors)
}
