package handlers

import (
	"log"

	"github.com/qdrant/go-client/qdrant"
	"github.com/rhydianjenkins/rag-mcp-server/src"
)

func Index() error {
	storage, err := src.Connect()

	if err != nil {
		log.Fatal("Unable to create storage");
		return err;
	}

	points := []*qdrant.PointStruct{
		{
			Id:      qdrant.NewIDNum(1),
			Vectors: qdrant.NewVectors(0.05, 0.61, 0.76, 0.74),
			Payload: qdrant.NewValueMap(map[string]any{"city": "London"}),
		},
		{
			Id:      qdrant.NewIDNum(2),
			Vectors: qdrant.NewVectors(0.19, 0.81, 0.75, 0.11),
			Payload: qdrant.NewValueMap(map[string]any{"age": 32}),
		},
		{
			Id:      qdrant.NewIDNum(3),
			Vectors: qdrant.NewVectors(0.36, 0.55, 0.47, 0.94),
			Payload: qdrant.NewValueMap(map[string]any{"vegan": true}),
		},
	}

	err = storage.Upsert(points);

	if err != nil {
		log.Fatal("Unable to upsert");
		return err;
	}

	return nil;
}
