package config

import (
	"log"
	"net/http"
	"os"
	"time"
)

func PingOllama() {
	url := os.Getenv("OLLAMA_API_URL")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	log.Println("Pinging Ollama 🤙...")
	if err != nil {
		log.Fatalln("Ollama Ping Failed ❌, Err ", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalln("Ollama Ping Failed ❌ with Status Code ", resp.StatusCode)
	}
	log.Println("Ollama Ping Success ✅")
}
