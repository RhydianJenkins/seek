package handlers

import (
	"fmt"

	"github.com/rhydianjenkins/seek/src/config"
	"github.com/rhydianjenkins/seek/src/ollama"
	"github.com/rhydianjenkins/seek/src/tools"
)

func AskQuestion(question string) error {
	cfg := config.Get()
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	client := ollama.NewClient(cfg.OllamaURL, cfg.ChatModel)

	messages := []ollama.Message{
		{
			Role: "system",
			Content: "You are a helpful assistant with access to a knowledge base. " +
				"When answering questions, use the search tool to find relevant information. " +
				"Provide accurate, concise answers based on the retrieved documents.",
		},
		{
			Role:    "user",
			Content: question,
		},
	}

	availableTools := tools.GetTools()

	maxIterations := 10

	for range maxIterations {
		response, err := client.Chat(messages, availableTools, func(chunk string) {
			fmt.Print(chunk)
		})
		if err != nil {
			return fmt.Errorf("chat request failed: %w", err)
		}

		messages = append(messages, *response)

		if len(response.ToolCalls) > 0 {
			for _, toolCall := range response.ToolCalls {
				fmt.Printf("[Calling tool: %s]\n", toolCall.Function.Name)

				result, err := tools.ExecuteTool(toolCall)
				if err != nil {
					return fmt.Errorf("tool execution failed: %w", err)
				}

				messages = append(messages, ollama.Message{
					Role:    "tool",
					Content: result,
				})
			}

			// Continue to next iteration to get LLM's response with tool results
			continue
		}

		// No tool calls means we got the final answer, so finish
		// TODO Rhydian allow the user to respond if they like?
		break
	}

	fmt.Println()
	return nil
}
