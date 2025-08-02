package handlers

import (
	"ai-support-bot-service/jobs"
	"ai-support-bot-service/models"
	"ai-support-bot-service/services"
	"ai-support-bot-service/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	var application *models.Application
	err := json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		utils.HandleErrorResponse("Invalid request", http.StatusBadRequest, w)
		return
	}
	newApplication, status, err := services.CreateApplication(application)
	if err != nil {
		utils.HandleErrorResponse(err.Error(), status, w)
		return
	}
	utils.HandleSuccessResponse("Application created successfully", newApplication, w)
}

func GetApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	applications, status, err := services.GetApplications()
	if err != nil {
		utils.HandleErrorResponse(err.Error(), status, w)
		return
	}
	utils.HandleSuccessResponse("Applications retrieved successfully", applications, w)
}

func GetApplicationHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	application, status, err := services.GetApplication(id)
	if err != nil {
		utils.HandleErrorResponse(err.Error(), status, w)
		return
	}
	utils.HandleSuccessResponse("Application Retrived Successfully", application, w)
}

func UpdateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	var application *models.Application
	err := json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		utils.HandleErrorResponse("Invalid request", http.StatusBadRequest, w)
		return
	}
	updatedApplication, status, err := services.UpdateApplication(*application)
	if err != nil {
		utils.HandleErrorResponse(err.Error(), status, w)
		return
	}
	utils.HandleSuccessResponse("Application updated successfully", updatedApplication, w)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file data", http.StatusInternalServerError)
		return
	}
	message, status, err := services.UploadFile(data, header.Filename, header.Header.Get("Content-Type"), chi.URLParam(r, "id"))
	if err != nil {
		utils.HandleErrorResponse(err.Error(), status, w)
		return
	}
	utils.HandleSuccessResponse(message, nil, w)
}

func IngestDataHandler(w http.ResponseWriter, r *http.Request) {
	FileID := chi.URLParam(r, "fileId")
	ApplicationID := chi.URLParam(r, "id")
	data, status, err := services.IngestData(FileID, ApplicationID)
	if err != nil {
		utils.HandleErrorResponse(err.Error(), status, w)
		return
	}
	utils.HandleSuccessResponse("Ingestion started successfully", data, w)
}

func IngestJobStatusSSEHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	jobID := chi.URLParam(r, "job_id")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	lastStatus := ""
	for {
		status := jobs.GetJobStatus(jobID)
		if status != lastStatus {
			fmt.Fprintf(w, "data: %s\n\n", status)
			flusher.Flush()
			lastStatus = status
		}

		if status == string(jobs.JobStatusSuccess) || strings.HasPrefix(status, "Error") {
			break
		}

		time.Sleep(1 * time.Second)
	}
}
