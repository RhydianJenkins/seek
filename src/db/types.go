package db

import "github.com/qdrant/go-client/qdrant"

type Storage struct {
	client         *qdrant.Client
	collectionName string
	ollamaURL      string
	vectorSize     uint64
	embeddingModel string
}

type ollamaEmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaEmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}
