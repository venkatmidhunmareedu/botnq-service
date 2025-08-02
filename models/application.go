package models

import "time"

type Application struct {
	ID           string `json:"id,omitempty" bson:"_id"`
	UserID       string `json:"user_id" bson:"user_id"`
	Name         string `json:"name" bson:"name"`
	SystemPrompt string `json:"system_prompt,omitempty" bson:"system_prompt"`
}

type File struct {
	ID            string    `json:"id,omitempty" bson:"_id"`
	ApplicationID string    `json:"application_id" bson:"application_id"`
	FileName      string    `json:"file_name" bson:"file_name"`
	ContentType   string    `json:"content_type" bson:"content_type"`
	FileBytes     []byte    `json:"file_bytes" bson:"file_bytes"`
	UploadTime    time.Time `json:"upload_time" bson:"upload_time"`
	Ingested      bool      `json:"ingested" bson:"ingested"`
}

type FileResponse struct {
	ID            string    `json:"id,omitempty" bson:"_id"`
	ApplicationID string    `json:"application_id" bson:"application_id"`
	FileName      string    `json:"file_name" bson:"file_name"`
	ContentType   string    `json:"content_type" bson:"content_type"`
	UploadTime    time.Time `json:"upload_time" bson:"upload_time"`
	Ingested      bool      `json:"ingested" bson:"ingested"`
}

type Job struct {
	ID            string   `json:"id,omitempty" bson:"_id"`
	FileID        string   `json:"file_id" bson:"file_id"`
	ApplicationID string   `json:"application_id" bson:"application_id"`
	Status        string   `json:"status" bson:"status"`
	Log           []string `json:"log" bson:"log"`
}
