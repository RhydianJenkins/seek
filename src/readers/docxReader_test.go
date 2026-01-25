package readers

import (
	"testing"
)

func TestDOCXReader(t *testing.T) {
	content := DOCXReader{}.Read("../../test-data/docs/DOCX_TestPage.docx")

	if content == "" {
		t.Errorf("DOCXReader{}.Read(...) returned an empty string")
	}
}
