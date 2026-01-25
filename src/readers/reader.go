package readers

import (
	"path/filepath"
	"strings"
)

func NewReader() *Reader {
	reader := &Reader{
		readers:       make(map[string]FileReader),
		defaultReader: PlainTextReader{},
	}

	reader.register(".pdf", PDFReader{})
	reader.register(".xlsx", XLSXReader{})

	return reader
}

func (r *Reader) register(extension string, reader FileReader) {
	r.readers[strings.ToLower(extension)] = reader
}

func (r *Reader) getReader(path string) FileReader {
	ext := strings.ToLower(filepath.Ext(path))
	if reader, exists := r.readers[ext]; exists {
		return reader
	}
	return r.defaultReader
}

func (r *Reader) ReadFile(path string) string {
	reader := r.getReader(path)
	return reader.Read(path)
}
