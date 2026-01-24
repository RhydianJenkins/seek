package readers

import (
	"log"
	"os"
	"unicode/utf8"
)

func ReadPlainText(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file %s: %v", path, err)
		return ""
	}
	content := string(bytes)

	if utf8.ValidString(content) {
		return content
	}

	log.Println("Warning: invalid UTF-8 content in file", path)
	return ""
}
