package services

import (
	"ai-support-bot-service/config"
	"ai-support-bot-service/models"
	"context"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Login(user *models.LoginRequest) (*models.User, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userExists, err := config.UserCollection.CountDocuments(ctx, bson.D{{Key: "email", Value: user.Email}})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to check if user already exists")
	}
	if userExists == 0 {
		return nil, http.StatusNotFound, errors.New("user does not exist")
	}
	var temp_user models.User
	err = config.UserCollection.FindOne(ctx, bson.D{{Key: "email", Value: user.Email}}).Decode(&temp_user)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve user")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(temp_user.Password), []byte(user.Password)); err != nil {
		return nil, http.StatusUnauthorized, errors.New("invalid password")
	}
	return &temp_user, http.StatusOK, nil
}
