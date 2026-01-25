package readers

import (
	"log"
	"strings"

	"github.com/ledongthuc/pdf"
)

type PDFReader struct{}

// Ensure PDFReader implements FileReader at compile time
var _ FileReader = PDFReader{}

func (r PDFReader) Read(path string) string {
	f, pdfReader, err := pdf.Open(path)
	if err != nil {
		log.Printf("Error opening PDF %s: %v", path, err)
		return ""
	}
	defer f.Close()

	var text strings.Builder
	totalPages := pdfReader.NumPage()

	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		p := pdfReader.Page(pageNum)
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
