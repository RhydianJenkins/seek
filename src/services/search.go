package services

import (
	"fmt"
	"strings"

	"github.com/rhydianjenkins/seek/src/db"
)

// SearchFiles performs a semantic search on the knowledge base
func SearchFiles(searchTerm string, limit int) (*SearchResults, error) {
	storage, err := db.Connect()
	if err != nil {
		return &SearchResults{
			Success: false,
			Error:   fmt.Sprintf("Unable to connect to storage: %v", err),
		}, err
	}

	// Normalize search term to lowercase for consistent embeddings
	normalizedTerm := strings.ToLower(searchTerm)

	searchResult, err := storage.Search(normalizedTerm, limit)
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
