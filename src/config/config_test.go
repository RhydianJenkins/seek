package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"QdrantHost", cfg.QdrantHost, "localhost"},
		{"QdrantPort", cfg.QdrantPort, 6334},
		{"QdrantUseTLS", cfg.QdrantUseTLS, false},
		{"CollectionName", cfg.CollectionName, "my_collection"},
		{"VectorSize", cfg.VectorSize, uint64(768)},
		{"OllamaURL", cfg.OllamaURL, "http://localhost:11434"},
		{"EmbeddingModel", cfg.EmbeddingModel, "nomic-embed-text"},
		{"ServerName", cfg.ServerName, "rag-mcp-server"},
		{"ServerVersion", cfg.ServerVersion, "1.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}
