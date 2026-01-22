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

type CollectionStatus struct {
	CollectionName string `json:"collection_name"`
	Exists         bool   `json:"exists"`
	VectorCount    uint64 `json:"vector_count,omitempty"`
	VectorSize     uint64 `json:"vector_size,omitempty"`
}
