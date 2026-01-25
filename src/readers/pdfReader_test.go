package readers

import (
	"testing"
)

func TestReadPDFFile(t *testing.T) {
	pdfContent := PDFReader{}.Read("../../test-data/pdfs/pdf_test.pdf")

	if pdfContent == "" {
		t.Errorf("PDFReader{}.Read(...) returned an empty string")
	}
}
