package config

type Config struct {
	// Qdrant Configuration
	QdrantHost   string
	QdrantPort   int
	QdrantUseTLS bool

	// Collection Configuration
	CollectionName string
	VectorSize     uint64

	// Ollama Configuration
	OllamaURL      string
	EmbeddingModel string

	// MCP Server Configuration
	ServerName    string
	ServerVersion string
}

func DefaultConfig() *Config {
	return &Config{
		QdrantHost:     "localhost",
		QdrantPort:     6334,
		QdrantUseTLS:   false,
		CollectionName: "my_collection",
		VectorSize:     768,
		OllamaURL:      "http://localhost:11434",
		EmbeddingModel: "nomic-embed-text",
		ServerName:     "rag-mcp-server",
		ServerVersion:  "1.0.0",
	}
}
