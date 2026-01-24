package readers

import (
	"log"
	"os"
)

func ReadPlainText(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file %s: %v", path, err)
		return ""
	}
	content := string(bytes)

	return content
}
