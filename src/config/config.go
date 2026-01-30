package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	CollectionName string
	EmbeddingModel string
	ChatModel      string
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

func Initialize(cfg *Config) {
	once.Do(func() {
		instance = applyDefaults(cfg)
	})
}

func Get() *Config {
	if instance == nil {
		log.Println("Error: config not initialized")
		return nil
	}
	return instance
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	log.Println("Warning: environment variable not set:", key)
	return ""
}

func getEnvInt(key string) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	log.Println("Warning: environment variable not set:", key)
	return 0
}

func applyDefaults(cfg *Config) *Config {
	if cfg.CollectionName == "" {
		cfg.CollectionName = getEnv("COLLECTION_NAME")
	}
	if cfg.OllamaURL == "" {
		ollamaHost := getEnv("OLLAMA_HOST")
		ollamaPort := getEnvInt("OLLAMA_PORT")
		cfg.OllamaURL = fmt.Sprintf("http://%s:%d", ollamaHost, ollamaPort)
	}
	if cfg.QdrantHost == "" {
		cfg.QdrantHost = getEnv("QDRANT_HOST")
	}
	if cfg.QdrantPort == 0 {
		cfg.QdrantPort = getEnvInt("QDRANT_PORT")
	}
	if cfg.ServerVersion == "" {
		cfg.ServerVersion = "dev"
	}

	cfg.VectorSize = 768
	cfg.EmbeddingModel = "nomic-embed-text"
	cfg.ChatModel = "qwen2.5"

	return cfg
}
