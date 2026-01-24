package readers

import (
	"log"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ReadPDFFile(path string) string {
	f, r, err := pdf.Open(path)
	if err != nil {
		log.Printf("Error opening PDF %s: %v", path, err)
		return ""
	}
	defer f.Close()

	var text strings.Builder
	totalPages := r.NumPage()

	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		p := r.Page(pageNum)
		if p.V.IsNull() {
			continue
		}

		content, err := p.GetPlainText(nil)
		if err != nil {
			log.Printf("Error reading page %d of %s: %v", pageNum, path, err)
			continue
		}
		text.WriteString(content)
		text.WriteString("\n\n")
	}

	return text.String()
}
