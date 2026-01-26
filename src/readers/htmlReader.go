package readers

import (
	"log"
	"os"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

type HTMLReader struct{}

// Ensure HTMLReader implements FileReader at compile time
var _ FileReader = HTMLReader{}

func (r HTMLReader) Read(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading HTML file %s: %v", path, err)
		return ""
	}

	converter := md.NewConverter("", true, nil)

	converter.Remove("script")
	converter.Remove("style")
	converter.Remove("nav")
	converter.Remove("footer")
	converter.Remove("header")

	markdown, err := converter.ConvertString(string(bytes))
	if err != nil {
		log.Printf("Error converting HTML to markdown %s: %v", path, err)
		return ""
	}

	return markdown
}
