package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/rhydianjenkins/rag-mcp-server/src/db"
)

func Status() {
	storage, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to storage: %v", err)
	}

	status, err := storage.GetStatus()
	if err != nil {
		log.Fatalf("Failed to get status: %v", err)
	}

	jsonOutput, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal status: %v", err)
	}

	fmt.Println(string(jsonOutput))
}
