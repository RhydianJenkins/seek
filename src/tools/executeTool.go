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
		// Parse flexibly - Ollama may send limit as string or int
		var rawInput map[string]any
		if err := json.Unmarshal(toolCall.Function.Arguments, &rawInput); err != nil {
			return "", fmt.Errorf("failed to parse search arguments: %w", err)
		}

		query, _ := rawInput["query"].(string)

		// Handle limit as either int or string
		limit := 3 // default
		if limitVal, ok := rawInput["limit"]; ok && limitVal != nil {
			switch v := limitVal.(type) {
			case float64: // JSON numbers are float64
				limit = int(v)
			case string:
				if parsedLimit, err := json.Number(v).Int64(); err == nil {
					limit = int(parsedLimit)
				}
			case int:
				limit = v
			}
		}

		results, err := services.SearchFiles(query, limit)
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
