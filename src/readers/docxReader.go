package readers

import (
	"archive/zip"
	"encoding/xml"
	"io"
	"log"
	"strings"
)

type DOCXReader struct{}

// Ensure DOCXReader implements FileReader at compile time
var _ FileReader = DOCXReader{}

type wordDocument struct {
	XMLName xml.Name `xml:"document"`
	Body    wordBody `xml:"body"`
}

type wordBody struct {
	Paragraphs []wordParagraph `xml:"p"`
}

type wordParagraph struct {
	Runs []wordRun `xml:"r"`
}

type wordRun struct {
	Text []string `xml:"t"`
}

func (r DOCXReader) Read(path string) string {
	zipReader, err := zip.OpenReader(path)
	if err != nil {
		log.Printf("Error opening DOCX %s: %v", path, err)
		return ""
	}
	defer zipReader.Close()

	var documentXML string
	for _, file := range zipReader.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				log.Printf("Error opening document.xml in %s: %v", path, err)
				return ""
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				log.Printf("Error reading document.xml in %s: %v", path, err)
				return ""
			}
			documentXML = string(content)
			break
		}
	}

	if documentXML == "" {
		log.Printf("No document.xml found in %s", path)
		return ""
	}

	var doc wordDocument
	if err := xml.Unmarshal([]byte(documentXML), &doc); err != nil {
		log.Printf("Error parsing XML in %s: %v", path, err)
		return ""
	}

	var text strings.Builder
	for _, paragraph := range doc.Body.Paragraphs {
		for _, run := range paragraph.Runs {
			for _, t := range run.Text {
				text.WriteString(t)
			}
		}
		text.WriteString("\n")
	}

	return text.String()
}
