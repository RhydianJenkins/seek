package handlers

import (
	"fmt"

	"github.com/rhydianjenkins/rag-mcp-server/src"
	"github.com/rhydianjenkins/rag-mcp-server/src/config"
)

// SearchFiles performs a search and returns structured results
// This function is used by both CLI and MCP modes
func SearchFiles(cfg *config.Config, searchTerm string, limit int) (*SearchResults, error) {
	storage, err := src.Connect(cfg)
	if err != nil {
		return &SearchResults{
			Success: false,
			Error:   fmt.Sprintf("Unable to connect to storage: %v", err),
		}, err
	}

	searchResult, err := storage.Search(searchTerm, limit)
	if err != nil {
		return &SearchResults{
			Success: false,
			Error:   fmt.Sprintf("Search failed: %v", err),
		}, err
	}

	results := &SearchResults{
		Success: true,
		Query:   searchTerm,
		Count:   len(searchResult),
		Results: make([]SearchResult, 0, len(searchResult)),
	}

	for _, result := range searchResult {
		sr := SearchResult{
			Score: result.Score,
		}

		if result.Payload != nil {
			if filename, ok := result.Payload["filename"]; ok {
				sr.Filename = filename.GetStringValue()
			}
			if chunkIdx, ok := result.Payload["chunk_index"]; ok {
				sr.ChunkIndex = chunkIdx.GetIntegerValue()
			}
			if content, ok := result.Payload["content"]; ok {
				sr.Content = content.GetStringValue()
			}
		}

		results.Results = append(results.Results, sr)
	}

	return results, nil
}

// Search is the CLI wrapper for backward compatibility
func Search(searchTerm string, ollamaURL string, limit int) error {
	cfg := config.DefaultConfig()
	cfg.OllamaURL = ollamaURL

	results, err := SearchFiles(cfg, searchTerm, limit)
	if err != nil {
		return err
	}

	// Pretty print for CLI users
	fmt.Printf("\nSearch results for: '%s'\n", results.Query)
	fmt.Printf("Found %d results:\n", results.Count)

	for i, result := range results.Results {
		fmt.Printf("\n--- Result %d (Score: %.4f) ---\n", i+1, result.Score)
		fmt.Printf("File: %s\n", result.Filename)
		fmt.Printf("Chunk: %d\n", result.ChunkIndex)
		fmt.Println()
		fmt.Println(result.Content)
		fmt.Println()
	}

	return nil
}
