package readers

type FileReader interface {
	Read(path string) string
}

type Reader struct {
	readers       map[string]FileReader
	defaultReader FileReader
}
