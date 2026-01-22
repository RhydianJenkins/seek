package handlers

import (
	"fmt"
	"log"
	"sort"

	"github.com/rhydianjenkins/rag-mcp-server/src/db"
)

func GetDocumentByFilename(filename string) (*DocumentResult, error) {
	storage, err := db.Connect()
	if err != nil {
		return &DocumentResult{
			Success: false,
			Error:   fmt.Sprintf("Unable to connect to storage: %v", err),
		}, err
	}

	points, err := storage.GetDocumentByFilename(filename)
	if err != nil {
		return &DocumentResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to retrieve document: %v", err),
		}, err
	}

	if len(points) == 0 {
		return &DocumentResult{
			Success: false,
			Error:   fmt.Sprintf("No document found with filename: %s", filename),
		}, fmt.Errorf("no document found with filename: %s", filename)
	}

	chunks := make([]DocumentChunk, 0, len(points))
	for _, point := range points {
		chunk := DocumentChunk{}
		if point.Payload != nil {
			if content, ok := point.Payload["content"]; ok {
				chunk.Content = content.GetStringValue()
			}
			if chunkIdx, ok := point.Payload["chunk_index"]; ok {
				chunk.ChunkIndex = chunkIdx.GetIntegerValue()
			}
		}
		chunks = append(chunks, chunk)
	}

	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].ChunkIndex < chunks[j].ChunkIndex
	})

	fullContent := ""
	for _, chunk := range chunks {
		if fullContent != "" {
			fullContent += "\n\n"
		}
		fullContent += chunk.Content
	}

	return &DocumentResult{
		Success:    true,
		Filename:   filename,
		ChunkCount: len(chunks),
		Chunks:     chunks,
		FullText:   fullContent,
	}, nil
}

func GetDocument(filename string) error {
	result, err := GetDocumentByFilename(filename)
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
		return err
	}

	fmt.Printf("\nDocument: %s\n", result.Filename)
	fmt.Printf("Total chunks: %d\n\n", result.ChunkCount)
	fmt.Println(result.FullText)

	return nil
}
