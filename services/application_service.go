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

func CreateApplication(application *models.Application) (*mongo.InsertOneResult, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Check if the name already exists
	applicationExists, err := config.ApplicationCollection.CountDocuments(ctx, bson.D{{Key: "name", Value: application.Name}})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to check if application already exists")
	}
	if applicationExists > 0 {
		return nil, http.StatusConflict, errors.New("application name already exists")
	}
	// Validate application
	if application.Name == "" || application.UserID == "" {
		return nil, http.StatusBadRequest, errors.New("invalid application details")
	}
	// Check if user exists
	userExists, err := config.UserCollection.CountDocuments(ctx, bson.D{{Key: "_id", Value: application.UserID}})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to check if user exists")
	}
	if userExists == 0 {
		return nil, http.StatusNotFound, errors.New("user not found")
	}
	application.ID = uuid.New().String()

	newApplication, err := config.ApplicationCollection.InsertOne(ctx, application)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to create application")
	}
	return newApplication, http.StatusOK, nil
}

func GetApplications() ([]models.Application, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := config.ApplicationCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve applications")
	}
	var applications []models.Application
	if err := cursor.All(ctx, &applications); err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve applications")
	}
	return applications, http.StatusOK, nil
}

func GetApplication(id string) (*models.Application, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var application models.Application
	err := config.ApplicationCollection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&application)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, http.StatusNotFound, errors.New("no application found with the id " + id)
		}
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve application")
	}
	return &application, http.StatusOK, nil
}

func UpdateApplication(application models.Application) (*mongo.UpdateResult, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Check if application exists
	applicationExists, err := config.ApplicationCollection.CountDocuments(ctx, bson.D{{Key: "_id", Value: application.ID}})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to check if application exists")
	}
	if applicationExists == 0 {
		return nil, http.StatusNotFound, errors.New("application not found")
	}
	// Validation
	if application.Name == "" || application.SystemPrompt == "" {
		return nil, http.StatusBadRequest, errors.New("invalid application details")
	}
	result, err := config.ApplicationCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: application.ID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: application.Name}, {Key: "system_prompt", Value: application.SystemPrompt}}}})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return result, http.StatusOK, nil
}
