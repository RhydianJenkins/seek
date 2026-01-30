package tools

import "github.com/rhydianjenkins/seek/src/ollama"

func GetTools() []ollama.Tool {
	return []ollama.Tool{
		{
			Type: "function",
			Function: ollama.FunctionDef{
				Name:        "search",
				Description: "Search the RAG knowledge base for relevant content using semantic similarity. Use this to find information related to the user's question.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "The search query to find relevant documents",
						},
						"limit": map[string]any{
							"type":        "integer",
							"description": "Maximum number of search results to return (default: 3)",
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: ollama.FunctionDef{
				Name:        "get_document",
				Description: "Retrieve a full document by filename, returning all chunks in order. Use this when you need the complete content of a specific document.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"filename": map[string]any{
							"type":        "string",
							"description": "The name of the file to retrieve",
						},
					},
					"required": []string{"filename"},
				},
			},
		},
	}
}
