package handler

import (
	"luxbiss/pkg/response"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "Server is healthy", nil)
}
