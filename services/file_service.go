package services

import (
	"ai-support-bot-service/config"
	"ai-support-bot-service/models"
	"context"
	"time"

	"net/http"

	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UploadFile(fileBytes []byte, fileName string, contentType string, applicationID string) (message string, status int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Validate name
	if fileName == "" {
		return "File name is required", http.StatusBadRequest, errors.New("file name is required")
	}
	_, err = config.FileCollection.InsertOne(ctx, &models.File{
		ID:            uuid.New().String(),
		FileName:      fileName,
		ContentType:   contentType,
		FileBytes:     fileBytes,
		UploadTime:    time.Now(),
		ApplicationID: applicationID,
		Ingested:      false,
	})
	if err != nil {
		return "File upload failed", http.StatusInternalServerError, err
	}
	return "File uploaded successfully", http.StatusOK, nil
}

func GetFile(id string) (*models.File, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var file models.File
	err := config.FileCollection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, http.StatusNotFound, errors.New("no file found with the id " + id)
		}
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve file")
	}
	return &file, http.StatusOK, nil
}

func GetFilesByApplicationID(applicationID string) ([]models.FileResponse, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := config.FileCollection.Find(ctx, bson.D{{Key: "application_id", Value: applicationID}})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve files")
	}
	// files with all but file_bytes
	var files []models.FileResponse
	if err := cursor.All(ctx, &files); err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve files")
	}
	return files, http.StatusOK, nil
}

func GetFileCount(applicationID string) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := config.FileCollection.CountDocuments(ctx, bson.D{{Key: "application_id", Value: applicationID}})
	if err != nil {
		return 0, http.StatusInternalServerError, errors.New("failed to retrieve file count")
	}
	return int(count), http.StatusOK, nil
}
