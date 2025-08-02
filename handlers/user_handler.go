package handlers

import (
	"ai-support-bot-service/models"
	"ai-support-bot-service/services"
	"ai-support-bot-service/utils"
	"encoding/json"
	"net/http"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.HandleErrorResponse("Invalid request", http.StatusBadRequest, w)
		return
	}
	_, status, err := services.CreateUser(user)
	if err != nil {
		utils.HandleErrorResponse(err.Error(), status, w)
		return
	}
	utils.HandleSuccessResponse("User created successfully", nil, w)
}
