package handlers

import (
	"log"

	"github.com/rhydianjenkins/rag-mcp-server/src"
)

func Search(searchTerm string) error {
	storage, err := src.Connect()

	if err != nil {
		return err;
	}

	searchResult, err := storage.Search(searchTerm)

	if err != nil {
		return err;
	}

	log.Println("qdrant search result", searchResult)

	return nil;
}
