package handlers

import (
	"fmt"

	"github.com/rhydianjenkins/rag-mcp-server/src/config"
	"github.com/rhydianjenkins/rag-mcp-server/src/storage"
)

func SearchFiles(searchTerm string, limit int) (*SearchResults, error) {
	storage, err := storage.Connect(config.Get())
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

func Search(searchTerm string, limit int) error {
	results, err := SearchFiles(searchTerm, limit)
	if err != nil {
		return err
	}

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
