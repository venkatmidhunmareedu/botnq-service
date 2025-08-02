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
)

func GenerateAPIKeyForUser(user *models.User) (result *mongo.UpdateResult, status int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err = config.UserCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: user.ID}}, bson.D{{Key: "api_keys", Value: bson.D{{Key: "$push", Value: uuid.New().String()}}}})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to generate api key")
	}
	return result, http.StatusOK, nil
}
