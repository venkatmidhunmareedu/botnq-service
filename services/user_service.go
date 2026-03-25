package services

import (
	"ai-support-bot-service/config"
	"ai-support-bot-service/models"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user *models.User) (*mongo.InsertOneResult, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Validate user
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return nil, http.StatusBadRequest, errors.New("invalid user details")
	}
	// Check if user already exists
	userExists, err := config.UserCollection.CountDocuments(ctx, bson.D{{Key: "email", Value: user.Email}})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to check if user already exists")
	}
	if userExists > 0 {
		return nil, http.StatusConflict, errors.New("user already exists")
	}
	user.ID = uuid.New().String()
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	newUser, err := config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return newUser, http.StatusOK, nil
}

func GenerateAPIKeyForUser(apiKey *models.APIKey) (result *mongo.InsertOneResult, status int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	apiKey.ID = uuid.New().String()
	apiKey.Key = uuid.New().String()
	result, err = config.APIKeyCollection.InsertOne(ctx, apiKey)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to generate api key")
	}
	return result, http.StatusOK, nil
}

func DeleteAPIKeyForUser(apiKeyID string) (result *mongo.DeleteResult, status int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err = config.APIKeyCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: apiKeyID}})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to delete api key")
	}
	return result, http.StatusOK, nil
}

func GetAPIKeysForUser(userID string) (result *mongo.Cursor, status int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err = config.APIKeyCollection.Find(ctx, bson.D{{Key: "user_id", Value: userID}})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve api keys")
	}
	return result, http.StatusOK, nil
}
