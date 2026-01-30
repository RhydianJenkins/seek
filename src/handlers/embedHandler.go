package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/rhydianjenkins/seek/src/services"
)

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

	result, err := services.EmbedFilesWithProgress(dataDir, chunkSize, progressCallback)
	if err != nil {
		fmt.Printf("\nError: %s\n", result.Error)
		return err
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\n%s\n", result.Message)
	fmt.Printf("Completed in %.2f seconds\n", elapsed.Seconds())
	return nil
}
