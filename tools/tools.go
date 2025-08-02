package tools

import (
	"ai-support-bot-service/services"

	"github.com/coder/aisdk-go"
)

var Tools = []aisdk.Tool{{
	Name:        "get_context",
	Description: "A test tool to get context on what the user is asking",
	Schema: aisdk.Schema{
		Required: []string{"query"},
		Properties: map[string]any{
			"query": map[string]any{
				"type": "string",
			},
		},
	},
}}

func HandleToolCall(toolCall aisdk.ToolCall, applicationID string) any {
	// Retrieve data from chroma for the message
	retrievedData, status, err := services.RetrieveData(toolCall.Args["query"].(string), applicationID)
	if err != nil {
		return map[string]any{
			"error":  err.Error(),
			"status": status,
		}
	}
	var message string
	for _, doc := range retrievedData {
		for _, doc := range doc {
			message += doc.ContentString()
		}
	}
	return map[string]any{
		"context": message,
	}
}
