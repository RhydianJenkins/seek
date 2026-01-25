package handlers

import (
	"fmt"
	"log"
	"sort"

	"github.com/rhydianjenkins/seek/src/db"
)

func List() {
	storage, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to storage: %v", err)
	}

	status, err := storage.GetStatus()
	if err != nil {
		log.Fatalf("Failed to get database status: %v", err)
	}

	if !status.Exists {
		return
	}

	filenames, err := storage.ListDocuments()
	if err != nil {
		log.Fatalf("Failed to list documents: %v", err)
	}

	sort.Strings(filenames)

	for _, filename := range filenames {
		fmt.Println(filename)
	}
}
