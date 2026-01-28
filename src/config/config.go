package config

import (
	"os"
	"sync"
)

type Config struct {
	CollectionName string
	EmbeddingModel string
	OllamaURL      string
	QdrantHost     string
	QdrantPort     int
	QdrantUseTLS   bool
	ServerName     string
	ServerVersion  string
	VectorSize     uint64
}

var (
	instance *Config
	once     sync.Once
)

// Initialize sets up the global configuration with the provided values.
// Any zero values will be replaced with defaults.
func Initialize(cfg *Config) {
	once.Do(func() {
		instance = applyDefaults(cfg)
	})
}

// Get returns the global configuration instance.
// If Initialize hasn't been called, it returns a default configuration.
func Get() *Config {
	if instance == nil {
		Initialize(&Config{})
	}
	return instance
}

// getEnv returns the value of an environment variable or the default value if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// applyDefaults fills in default values for any zero-value fields
// Reads from environment variables for supported fields, then falls back to hardcoded defaults
func applyDefaults(cfg *Config) *Config {
	if cfg.CollectionName == "" {
		cfg.CollectionName = getEnv("COLLECTION_NAME", "seek_collection")
	}
	if cfg.EmbeddingModel == "" {
		cfg.EmbeddingModel = "nomic-embed-text"
	}
	if cfg.OllamaURL == "" {
		cfg.OllamaURL = "http://localhost:11434"
	}
	if cfg.QdrantHost == "" {
		cfg.QdrantHost = getEnv("QDRANT_HOST", "localhost")
	}
	if cfg.QdrantPort == 0 {
		cfg.QdrantPort = 6334
	}
	if cfg.ServerName == "" {
		cfg.ServerName = "seek"
	}
	if cfg.ServerVersion == "" {
		cfg.ServerVersion = "dev"
	}
	if cfg.VectorSize == 0 {
		cfg.VectorSize = 768
	}
	return cfg
}
