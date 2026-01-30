package handlers

import (
	"fmt"
	"log"

	"github.com/rhydianjenkins/seek/src/services"
)

func GetDocument(filename string) error {
	result, err := services.GetDocumentByFilename(filename)
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
		return err
	}

	fmt.Printf("\nDocument: %s\n", result.Filename)
	fmt.Printf("Total chunks: %d\n\n", result.ChunkCount)
	fmt.Println(result.FullText)

	return nil
}
