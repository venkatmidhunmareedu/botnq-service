package utils

import (
	"ai-support-bot-service/models"
	"encoding/json"
	"net/http"
)

func HandleSuccessResponse(message string, data any, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Message: message,
		Success: true,
		Data:    data,
	})
}

func HandleErrorResponse(message string, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.Response{
		Message: message,
		Success: false,
	})
}
