package readers

import (
	"testing"
)

func TestReadPDFFile(t *testing.T) {
	pdfContent := ReadPDFFile("../../test-data/pdfs/pdf_test.pdf")

	if pdfContent == "" {
		t.Errorf("ReadPDFFile() returned an empty string")
	}
}
