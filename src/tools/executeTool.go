package tools

import (
	"encoding/json"
	"fmt"

	"github.com/rhydianjenkins/seek/src/ollama"
	"github.com/rhydianjenkins/seek/src/services"
)

func ExecuteTool(toolCall ollama.ToolCall) (string, error) {
	switch toolCall.Function.Name {
	case "search":
		var input struct {
			Query string `json:"query"`
			Limit int    `json:"limit"`
		}
		if err := json.Unmarshal(toolCall.Function.Arguments, &input); err != nil {
			return "", fmt.Errorf("failed to parse search arguments: %w", err)
		}
		if input.Limit == 0 {
			input.Limit = 3
		}

		results, err := services.SearchFiles(input.Query, input.Limit)
		if err != nil {
			return "", fmt.Errorf("search failed: %w", err)
		}

		resultJSON, err := json.Marshal(results)
		if err != nil {
			return "", fmt.Errorf("failed to marshal search results: %w", err)
		}
		return string(resultJSON), nil

	case "get_document":
		var input struct {
			Filename string `json:"filename"`
		}
		if err := json.Unmarshal(toolCall.Function.Arguments, &input); err != nil {
			return "", fmt.Errorf("failed to parse get_document arguments: %w", err)
		}

		result, err := services.GetDocumentByFilename(input.Filename)
		if err != nil {
			return "", fmt.Errorf("get_document failed: %w", err)
		}

		resultJSON, err := json.Marshal(result)
		if err != nil {
			return "", fmt.Errorf("failed to marshal document result: %w", err)
		}
		return string(resultJSON), nil

	default:
		return "", fmt.Errorf("unknown tool: %s", toolCall.Function.Name)
	}
}
