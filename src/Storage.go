package src

import (
	"context"
	"log"

	"github.com/qdrant/go-client/qdrant"
)

type Storage struct {
	client *qdrant.Client
	collectionName string
}

func Connect() (*Storage, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
		UseTLS: false,
		// APIKey: "<your-api-key>",
		// PoolSize: 3,
		// KeepAliveTime: 10,
		// KeepAliveTimeout: 2,
		// TLSConfig: &tls.Config{...},
		// GrpcOptions: []grpc.DialOption{},
	})

	if err != nil {
		return nil, err
	}

	storage := &Storage {
		client: client,
		collectionName: "my_collection",
	}

	return storage, nil
}

func (storage *Storage) Upsert(points []*qdrant.PointStruct) (error) {
	storage.client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: storage.collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     4,
			Distance: qdrant.Distance_Cosine,
		}),
	})

	operationInfo, err := storage.client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: storage.collectionName,
		Points: points,
	})

	if err != nil {
		log.Fatal("Unable to upsert");
		return err;
	}

	log.Println("qdrant upsert result", operationInfo)

	return nil;
}

func (storage *Storage) Search(searchTerm string) ([]*qdrant.ScoredPoint, error) {
	// TODO Rhydian convert search term into vector to query
	query := qdrant.NewQuery(0.5, 0.2, 0.1, 0.2)

	searchResult, err := storage.client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: storage.collectionName,
		Query:          query,
	})

	if err != nil {
		log.Fatal("Unable to search for term")
		return nil, err
	}

	return searchResult, nil
}
