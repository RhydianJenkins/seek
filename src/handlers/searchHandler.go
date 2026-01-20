package handlers

import (
	"log"

	"github.com/rhydianjenkins/rag-mcp-server/src/storage"
)

func Search(knowledgeBasePath string) error {
	storage, err := storage.Connect()

	if err != nil {
		return err;
	}

	searchResult, err := storage.Search("Hello")

	if err != nil {
		return err;
	}

	log.Println("qdrant search result", searchResult)

	return nil;
}
