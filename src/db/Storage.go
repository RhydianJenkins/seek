package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/qdrant/go-client/qdrant"
	"github.com/rhydianjenkins/seek/src/config"
)

func Connect() (*Storage, error) {
	cfg := config.Get()

	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   cfg.QdrantHost,
		Port:   cfg.QdrantPort,
		UseTLS: cfg.QdrantUseTLS,
	})

	if err != nil {
		return nil, err
	}

	storage := &Storage{
		client:         client,
		collectionName: cfg.CollectionName,
		ollamaURL:      cfg.OllamaURL,
		vectorSize:     cfg.VectorSize,
		embeddingModel: cfg.EmbeddingModel,
	}

	return storage, nil
}

func (storage *Storage) GetEmbedding(text string) ([]float32, error) {
	reqBody := ollamaEmbedRequest{
		Model:  storage.embeddingModel,
		Prompt: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		storage.ollamaURL+"/api/embeddings",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var embedResp ollamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return embedResp.Embedding, nil
}

func (storage *Storage) GenerateDb(points []*qdrant.PointStruct) error {
	exists, err := storage.client.CollectionExists(context.Background(), storage.collectionName)

	if err != nil {
		log.Fatal("Unable to check if collection exists")
		return err
	}

	if exists {
		storage.client.DeleteCollection(context.Background(), storage.collectionName)
	}

	storage.client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: storage.collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     storage.vectorSize,
			Distance: qdrant.Distance_Cosine,
		}),
	})

	operationInfo, err := storage.client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: storage.collectionName,
		Points:         points,
	})

	if err != nil {
		log.Fatal("Unable to upsert")
		return err
	}

	log.Println("qdrant upsert result", operationInfo)

	return nil
}

func (storage *Storage) Search(searchTerm string, limit int) ([]*qdrant.ScoredPoint, error) {
	embedding, err := storage.GetEmbedding(searchTerm)
	if err != nil {
		log.Printf("Failed to get embedding for search term: %v", err)
		return nil, fmt.Errorf("failed to get embedding: %w", err)
	}

	query := qdrant.NewQuery(embedding...)

	searchResult, err := storage.client.Query(
		context.Background(),
		&qdrant.QueryPoints{
			CollectionName: storage.collectionName,
			Query:          query,
			WithPayload:    qdrant.NewWithPayload(true),
			Limit:          qdrant.PtrOf(uint64(limit)),
		},
	)

	if err != nil {
		log.Printf("Unable to search for term: %v", err)
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return searchResult, nil
}

func (storage *Storage) GetStatus() (*CollectionStatus, error) {
	exists, err := storage.client.CollectionExists(context.Background(), storage.collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to check collection existence: %w", err)
	}

	if !exists {
		return &CollectionStatus{
			CollectionName: storage.collectionName,
			Exists:         false,
		}, nil
	}

	collectionInfo, err := storage.client.GetCollectionInfo(context.Background(), storage.collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection info: %w", err)
	}

	return &CollectionStatus{
		CollectionName: storage.collectionName,
		Exists:         true,
		VectorCount:    collectionInfo.GetPointsCount(),
		VectorSize:     storage.vectorSize,
	}, nil
}

func (storage *Storage) GetDocumentByFilename(filename string) ([]*qdrant.ScoredPoint, error) {
	filter := &qdrant.Filter{
		Must: []*qdrant.Condition{
			{
				ConditionOneOf: &qdrant.Condition_Field{
					Field: &qdrant.FieldCondition{
						Key: "filename",
						Match: &qdrant.Match{
							MatchValue: &qdrant.Match_Keyword{
								Keyword: filename,
							},
						},
					},
				},
			},
		},
	}

	scrollResult, err := storage.client.Scroll(
		context.Background(),
		&qdrant.ScrollPoints{
			CollectionName: storage.collectionName,
			Filter:         filter,
			WithPayload:    qdrant.NewWithPayload(true),
			Limit:          qdrant.PtrOf(uint32(1000)),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scroll documents: %w", err)
	}

	scoredPoints := make([]*qdrant.ScoredPoint, len(scrollResult))
	for i, point := range scrollResult {
		scoredPoints[i] = &qdrant.ScoredPoint{
			Id:      point.Id,
			Payload: point.Payload,
			Score:   1.0,
		}
	}

	return scoredPoints, nil
}

func (storage *Storage) ListDocuments() ([]string, error) {
	scrollResult, err := storage.client.Scroll(
		context.Background(),
		&qdrant.ScrollPoints{
			CollectionName: storage.collectionName,
			WithPayload:    qdrant.NewWithPayload(true),
			Limit:          qdrant.PtrOf(uint32(10000)),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scroll documents: %w", err)
	}

	filenameMap := make(map[string]bool)
	for _, point := range scrollResult {
		if point.Payload != nil {
			if filename, ok := point.Payload["filename"]; ok {
				filenameMap[filename.GetStringValue()] = true
			}
		}
	}

	filenames := make([]string, 0, len(filenameMap))
	for filename := range filenameMap {
		filenames = append(filenames, filename)
	}

	return filenames, nil
}
