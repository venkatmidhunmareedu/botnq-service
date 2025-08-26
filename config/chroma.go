package config

import (
	"log"
	"os"
	"time"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	ollama "github.com/amikos-tech/chroma-go/pkg/embeddings/ollama"
)

type ChromaConfiguration struct {
	Client            *chroma.Client
	EmbeddingFunction ollama.OllamaEmbeddingFunction
}

var ChromaConfig ChromaConfiguration

func InitChroma() {
	client, err := chroma.NewHTTPClient(
		chroma.WithBaseURL(os.Getenv("CHROMA_DB_URL")),
		chroma.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatalln("Chroma Client Init Failed ❌, Err ", err)
	}
	ollamaef, err := ollama.NewOllamaEmbeddingFunction(ollama.WithBaseURL(os.Getenv("OLLAMA_API_URL")), ollama.WithModel("nomic-embed-text"))
	if err != nil {
		log.Fatalln("Chroma Client Init Failed ❌, Err ", err)
	}

	ChromaConfig.Client = &client
	ChromaConfig.EmbeddingFunction = *ollamaef
	log.Println("Chroma Client Intitialized ✅")
}
