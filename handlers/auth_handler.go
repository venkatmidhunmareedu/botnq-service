package handlers

import (
	"ai-support-bot-service/models"
	"ai-support-bot-service/services"
	"ai-support-bot-service/utils"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user *models.LoginRequest
	mySecret := os.Getenv("JWT_SECRET")
	exp := time.Now().Add(time.Hour * 1).Unix()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.HandleErrorResponse("Invalid request", http.StatusBadRequest, w)
		return
	}
	temp_user, status, err := services.Login(user)
	if err != nil {
		utils.HandleErrorResponse(err.Error(), status, w)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": temp_user.ID,
		"exp":     exp,
	})
	tokenString, err := token.SignedString([]byte(mySecret))
	if err != nil {
		utils.HandleErrorResponse("Failed to generate access token", http.StatusInternalServerError, w)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		HttpOnly: true,
		Secure:   false, // Use true in production
		Path:     "/",
	})
	loginResponse := &models.LoginResponse{
		User: &models.UserResponse{
			Name:      temp_user.Name,
			Email:     temp_user.Email,
			ID:        temp_user.ID,
			CreatedAt: temp_user.CreatedAt,
			UpdatedAt: temp_user.UpdatedAt,
		},
		Exp: exp,
	}
	utils.HandleSuccessResponse("Login successful", loginResponse, w)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		utils.HandleErrorResponse("No access token found", http.StatusUnauthorized, w)
		return
	}
	cookie.Expires = time.Now().Add(-time.Hour * 24)
	cookie.Value = ""
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteLaxMode
	w.Header().Add("Set-Cookie", cookie.String())
	utils.HandleSuccessResponse("Logout successful", nil, w)
}
