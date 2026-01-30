package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func NewClient(baseURL, model string) *Client {
	return &Client{
		baseURL: baseURL,
		model:   model,
	}
}

func (c *Client) Chat(messages []Message, tools []Tool, onChunk func(string)) (*Message, error) {
	request := chatRequest{
		Model:    c.model,
		Messages: messages,
		Tools:    tools,
		Stream:   true,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		c.baseURL+"/api/chat",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)

	var finalMessage Message
	var fullContent strings.Builder

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var response chatResponse
		if err := json.Unmarshal(line, &response); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		if response.Message.Content != "" {
			fullContent.WriteString(response.Message.Content)
			if onChunk != nil {
				onChunk(response.Message.Content)
			}
		}

		if len(response.Message.ToolCalls) > 0 {
			finalMessage.ToolCalls = response.Message.ToolCalls
		}

		if response.Message.Role != "" {
			finalMessage.Role = response.Message.Role
		}

		if response.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	finalMessage.Content = fullContent.String()
	return &finalMessage, nil
}
