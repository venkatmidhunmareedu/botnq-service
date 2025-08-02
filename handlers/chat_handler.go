package handlers

import (
	"ai-support-bot-service/config"
	"ai-support-bot-service/prompts"
	"ai-support-bot-service/services"
	"ai-support-bot-service/tools"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/coder/aisdk-go"
	"google.golang.org/genai"
)

type ChatRequest struct {
	*aisdk.Chat
	ApplicationID string `json:"appId"`
	APIKey        string `json:"apiKey"`
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tools_gemini, err := aisdk.ToolsToGoogle(tools.Tools)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	aisdk.WriteDataStreamHeaders(w)
	messages, err := aisdk.MessagesToGoogle(req.Messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch System prompt from MongoDB
	application, status, err := services.GetApplication(req.ApplicationID)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	stream := aisdk.GoogleToDataStream(config.GoogleClient.Models.GenerateContentStream(context.Background(), "gemini-2.5-flash", messages, &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Role: "system",
			Parts: []*genai.Part{
				{
					Text: prompts.SystemPrompt + application.SystemPrompt,
				},
			},
		},
		Tools: tools_gemini,
	}))

	var acc aisdk.DataStreamAccumulator
	stream = stream.WithToolCalling(func(toolCall aisdk.ToolCall) any {
		return tools.HandleToolCall(toolCall, req.ApplicationID)
	})
	stream = stream.WithAccumulator(&acc)
	err = stream.Pipe(w)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Messages = append(req.Messages, acc.Messages()...)
}
