package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client
var MongoClient *mongo.Client

// Database
var MongoDatabase *mongo.Database

// Collections
var UserCollection *mongo.Collection
var ApplicationCollection *mongo.Collection
var FileCollection *mongo.Collection
var JobCollection *mongo.Collection

func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_DB_URL"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connection error ❌: %v", err)
	}
	log.Println("MongoDB Client Intitialized ✅")
	log.Println("Checking for Connection...")

	// Check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB ping error ❌: %v", err)
	}

	log.Println("MongoDB Ping Success ✅")
	MongoClient = client
	MongoDatabase = client.Database("ai-support-bot")
	UserCollection = client.Database("ai-support-bot").Collection("users")
	ApplicationCollection = client.Database("ai-support-bot").Collection("applications")
	FileCollection = client.Database("ai-support-bot").Collection("files")
	JobCollection = client.Database("ai-support-bot").Collection("jobs")
}
