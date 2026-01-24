package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/qdrant/go-client/qdrant"
	"github.com/rhydianjenkins/seek/src/db"
	"github.com/rhydianjenkins/seek/src/readers"
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
	// TODO Rhydian we load all files and their content into memory
	// This might get too big for some systems
	files := make(map[string]string)

	err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		var content string

		switch ext {
		case ".txt":
			content = readers.ReadPlainText(path)
		case ".pdf":
			content = readers.ReadPDFFile(path)
		default:
			return nil
		}

		relPath, _ := filepath.Rel(dataDir, path)
		files[relPath] = content

		return nil
	})

	return files, err
}

func EmbedFilesWithProgress(dataDir string, chunkSize int, progressCallback ProgressCallback) (*EmbedResult, error) {
	storage, err := db.Connect()
	if err != nil {
		return &EmbedResult{
			Success: false,
			Error:   fmt.Sprintf("Unable to create storage: %v", err),
		}, err
	}

	files, err := readTextFiles(dataDir)
	if err != nil {
		return &EmbedResult{
			Success: false,
			Error:   fmt.Sprintf("Unable to read files from directory %s: %v", dataDir, err),
		}, err
	}

	if len(files) == 0 {
		return &EmbedResult{
			Success: false,
			Error:   fmt.Sprintf("No .txt or .pdf files found in %s", dataDir),
		}, fmt.Errorf("no .txt or .pdf files found in %s", dataDir)
	}

	var points []*qdrant.PointStruct
	pointID := uint64(1)

	totalFiles := len(files)
	currentFile := 0

	for filename, content := range files {
		currentFile++
		if progressCallback != nil {
			progressCallback(currentFile, totalFiles, filename)
		}

		chunks := chunkText(content, chunkSize)

		for chunkIdx, chunk := range chunks {
			embedding, err := storage.GetEmbedding(chunk)
			if err != nil {
				log.Printf("Error generating embedding for %s chunk %d: %v", filename, chunkIdx, err)
				continue
			}

			payload := map[string]any{
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
		return &EmbedResult{
			Success: false,
			Error:   "No points to index",
		}, fmt.Errorf("no points to index")
	}

	err = storage.GenerateDb(points)
	if err != nil {
		return &EmbedResult{
			Success: false,
			Error:   fmt.Sprintf("Unable to generate db: %v", err),
		}, err
	}

	return &EmbedResult{
		Success:      true,
		FilesIndexed: len(files),
		TotalChunks:  len(points),
		Message:      fmt.Sprintf("Successfully indexed %d chunks from %d files", len(points), len(files)),
	}, nil
}

// EmbedFiles is a wrapper for backwards compatibility (used by MCP server)
func EmbedFiles(dataDir string, chunkSize int) (*EmbedResult, error) {
	return EmbedFilesWithProgress(dataDir, chunkSize, nil)
}

func Embed(dataDir string, chunkSize int) error {
	fmt.Printf("Starting indexing (chunk size: %d chars)\n", chunkSize)

	startTime := time.Now()
	var lastUpdate time.Time

	progressCallback := func(current, total int, filename string) {
		now := time.Now()
		if now.Sub(lastUpdate) < 100*time.Millisecond && current != total {
			return
		}
		lastUpdate = now

		percent := float64(current) / float64(total) * 100
		barWidth := 40
		filled := int(float64(barWidth) * float64(current) / float64(total))

		var bar strings.Builder

		bar.WriteString("[")
		for i := range barWidth {
			if i < filled {
				bar.WriteString("=")
			} else if i == filled {
				bar.WriteString(">")
			} else {
				bar.WriteString(" ")
			}
		}
		bar.WriteString("]")

		displayName := filename
		if len(displayName) > 30 {
			displayName = "..." + displayName[len(displayName)-27:]
		}

		fmt.Printf("\r%s %3.0f%% (%d/%d) Processing: %-30s", bar.String(), percent, current, total, displayName)

		if current == total {
			fmt.Println()
		}
	}

	result, err := EmbedFilesWithProgress(dataDir, chunkSize, progressCallback)
	if err != nil {
		fmt.Printf("\nError: %s\n", result.Error)
		return err
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\n%s\n", result.Message)
	fmt.Printf("Completed in %.2f seconds\n", elapsed.Seconds())
	return nil
}
