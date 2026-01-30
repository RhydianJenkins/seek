package handlers

import (
	"fmt"

	"github.com/rhydianjenkins/seek/src/services"
)

func Search(searchTerm string, limit int) error {
	results, err := services.SearchFiles(searchTerm, limit)
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
