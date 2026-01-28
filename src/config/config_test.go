package config

import (
	"sync"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	instance = nil
	once = sync.Once{}

	Initialize(&Config{})
	cfg := Get()

	if cfg == nil {
		t.Fatal("Get() returned nil")
	}

	tests := []struct {
		name     string
		got      any
		expected any
	}{
		{"QdrantHost", cfg.QdrantHost, "localhost"},
		{"QdrantPort", cfg.QdrantPort, 6334},
		{"QdrantUseTLS", cfg.QdrantUseTLS, false},
		{"CollectionName", cfg.CollectionName, "seek_collection"},
		{"VectorSize", cfg.VectorSize, uint64(768)},
		{"OllamaURL", cfg.OllamaURL, "http://localhost:11434"},
		{"EmbeddingModel", cfg.EmbeddingModel, "nomic-embed-text"},
		{"ServerName", cfg.ServerName, "seek"},
		{"ServerVersion", cfg.ServerVersion, "dev"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestCustomConfig(t *testing.T) {
	instance = nil
	once = sync.Once{}

	Initialize(&Config{
		QdrantHost:     "custom-host",
		QdrantPort:     9999,
		CollectionName: "custom_collection",
		OllamaURL:      "http://custom:8080",
	})

	cfg := Get()

	tests := []struct {
		name     string
		got      any
		expected any
	}{
		{"QdrantHost", cfg.QdrantHost, "custom-host"},
		{"QdrantPort", cfg.QdrantPort, 9999},
		{"CollectionName", cfg.CollectionName, "custom_collection"},
		{"OllamaURL", cfg.OllamaURL, "http://custom:8080"},
		{"EmbeddingModel", cfg.EmbeddingModel, "nomic-embed-text"}, // Should use default
		{"VectorSize", cfg.VectorSize, uint64(768)},                // Should use default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}
