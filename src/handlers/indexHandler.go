package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/qdrant/go-client/qdrant"
	"github.com/rhydianjenkins/rag-mcp-server/src/config"
	"github.com/rhydianjenkins/rag-mcp-server/src/storage"
)

func chunkText(text string, maxChunkSize int) []string {
	paragraphs := strings.Split(text, "\n\n")

	var chunks []string
	currentChunk := ""

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		if len(currentChunk) > 0 && len(currentChunk)+len(para) > maxChunkSize {
			chunks = append(chunks, currentChunk)
			currentChunk = para
		} else {
			if len(currentChunk) > 0 {
				currentChunk += "\n\n" + para
			} else {
				currentChunk = para
			}
		}
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

func readTextFiles(dataDir string) (map[string]string, error) {
	files := make(map[string]string)

	err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".txt") {
			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file %s: %v", path, err)
				return nil
			}

			relPath, _ := filepath.Rel(dataDir, path)
			files[relPath] = string(content)
		}

		return nil
	})

	return files, err
}

func IndexFiles(dataDir string, chunkSize int) (*IndexResult, error) {
	storage, err := storage.Connect(config.Get())
	if err != nil {
		return &IndexResult{
			Success: false,
			Error:   fmt.Sprintf("Unable to create storage: %v", err),
		}, err
	}

	files, err := readTextFiles(dataDir)
	if err != nil {
		return &IndexResult{
			Success: false,
			Error:   fmt.Sprintf("Unable to read files from directory %s: %v", dataDir, err),
		}, err
	}

	if len(files) == 0 {
		return &IndexResult{
			Success: false,
			Error:   fmt.Sprintf("No .txt files found in %s", dataDir),
		}, fmt.Errorf("no .txt files found in %s", dataDir)
	}

	var points []*qdrant.PointStruct
	pointID := uint64(1)

	for filename, content := range files {
		chunks := chunkText(content, chunkSize)

		for chunkIdx, chunk := range chunks {
			embedding, err := storage.GetEmbedding(chunk)
			if err != nil {
				log.Printf("Error generating embedding for %s chunk %d: %v", filename, chunkIdx, err)
				continue
			}

			payload := map[string]interface{}{
				"filename":    filename,
				"chunk_index": chunkIdx,
				"content":     chunk,
			}

			point := &qdrant.PointStruct{
				Id:      qdrant.NewIDNum(pointID),
				Vectors: qdrant.NewVectors(embedding...),
				Payload: qdrant.NewValueMap(payload),
			}

			points = append(points, point)
			pointID++
		}
	}

	if len(points) == 0 {
		return &IndexResult{
			Success: false,
			Error:   "No points to index",
		}, fmt.Errorf("no points to index")
	}

	err = storage.GenerateDb(points)
	if err != nil {
		return &IndexResult{
			Success: false,
			Error:   fmt.Sprintf("Unable to generate db: %v", err),
		}, err
	}

	return &IndexResult{
		Success:      true,
		FilesIndexed: len(files),
		TotalChunks:  len(points),
		Message:      fmt.Sprintf("Successfully indexed %d chunks from %d files", len(points), len(files)),
	}, nil
}

func Index(dataDir string, chunkSize int) error {
	log.Printf("Found files to index (chunk size: %d chars)", chunkSize)

	result, err := IndexFiles(dataDir, chunkSize)
	if err != nil {
		log.Fatal(result.Error)
		return err
	}

	log.Printf("%s", result.Message)
	return nil
}
