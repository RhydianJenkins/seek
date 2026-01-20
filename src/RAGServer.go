package src

type RAGServer struct {
    // vectorStore *VectorStore
    // embedder    *OllamaClient
    // kbPath      string
}

func NewRAGServer() (*RAGServer, error) {
	return &RAGServer{}, nil;
}
