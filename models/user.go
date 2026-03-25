package models

import "time"

type User struct {
	ID        string    `json:"id,omitempty" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type APIKey struct {
	ID     string `json:"id,omitempty" bson:"_id"`
	UserID string `json:"user_id" bson:"user_id"`
	Name   string `json:"name" bson:"name"`
	Key    string `json:"key" bson:"key"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type LoginResponse struct {
	User *UserResponse `json:"user"`
	Exp  int64         `json:"exp"`
}
