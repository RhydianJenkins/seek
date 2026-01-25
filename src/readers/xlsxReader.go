package readers

import (
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

type XLSXReader struct{}

// Ensure XLSXReader implements FileReader at compile time
var _ FileReader = XLSXReader{}

func (r XLSXReader) Read(path string) string {
	f, err := excelize.OpenFile(path)
	if err != nil {
		log.Printf("Error opening XLSX %s: %v", path, err)
		return ""
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing XLSX %s: %v", path, err)
		}
	}()

	var text strings.Builder
	sheets := f.GetSheetList()

	for _, sheetName := range sheets {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			log.Printf("Error reading sheet %s in %s: %v", sheetName, path, err)
			continue
		}

		text.WriteString("Sheet: ")
		text.WriteString(sheetName)
		text.WriteString("\n\n")

		for _, row := range rows {
			text.WriteString(strings.Join(row, "\t"))
			text.WriteString("\n")
		}
		text.WriteString("\n")
	}

	return text.String()
}
