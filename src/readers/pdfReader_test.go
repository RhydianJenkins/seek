package readers

import (
	"testing"
)

func TestReadPDFFile(t *testing.T) {
	content := PDFReader{}.Read("../../test-data/pdfs/pdf_test.pdf")

	if content == "" {
		t.Errorf("PDFReader{}.Read(...) returned an empty string")
	}
}
