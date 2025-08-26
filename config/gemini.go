package config

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"
)

var GoogleClient *genai.Client

func InitGemini(ctx context.Context) {
	var err error
	GoogleClient, err = genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  os.Getenv("GOOGLE_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalln("Google Client Init Failed ❌, Err ", err)
	}
	log.Println("Google Client Intitialized ✅")
}
