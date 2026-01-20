package handlers

import (
	"encoding/json"
	"testing"
)

func TestIndexResultJSON(t *testing.T) {
	tests := []struct {
		name     string
		result   IndexResult
		expected string
	}{
		{
			name: "successful index result",
			result: IndexResult{
				Success:      true,
				FilesIndexed: 5,
				TotalChunks:  42,
				Message:      "Successfully indexed",
			},
			expected: `{"success":true,"files_indexed":5,"total_chunks":42,"message":"Successfully indexed"}`,
		},
		{
			name: "failed index result with error",
			result: IndexResult{
				Success: false,
				Error:   "connection failed",
			},
			expected: `{"success":false,"files_indexed":0,"total_chunks":0,"message":"","error":"connection failed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.result)
			if err != nil {
				t.Fatalf("Failed to marshal IndexResult: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("JSON marshal = %s, want %s", string(data), tt.expected)
			}

			var unmarshaled IndexResult
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal IndexResult: %v", err)
			}

			if unmarshaled.Success != tt.result.Success {
				t.Errorf("Success = %v, want %v", unmarshaled.Success, tt.result.Success)
			}
		})
	}
}

func TestSearchResultJSON(t *testing.T) {
	result := SearchResult{
		Score:      0.95,
		Filename:   "test.txt",
		ChunkIndex: 3,
		Content:    "This is test content",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal SearchResult: %v", err)
	}

	var unmarshaled SearchResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal SearchResult: %v", err)
	}

	if unmarshaled.Score != result.Score {
		t.Errorf("Score = %v, want %v", unmarshaled.Score, result.Score)
	}
	if unmarshaled.Filename != result.Filename {
		t.Errorf("Filename = %v, want %v", unmarshaled.Filename, result.Filename)
	}
	if unmarshaled.ChunkIndex != result.ChunkIndex {
		t.Errorf("ChunkIndex = %v, want %v", unmarshaled.ChunkIndex, result.ChunkIndex)
	}
	if unmarshaled.Content != result.Content {
		t.Errorf("Content = %v, want %v", unmarshaled.Content, result.Content)
	}
}

func TestSearchResultsJSON(t *testing.T) {
	results := SearchResults{
		Success: true,
		Query:   "test query",
		Count:   2,
		Results: []SearchResult{
			{
				Score:      0.95,
				Filename:   "file1.txt",
				ChunkIndex: 0,
				Content:    "First result",
			},
			{
				Score:      0.87,
				Filename:   "file2.txt",
				ChunkIndex: 1,
				Content:    "Second result",
			},
		},
	}

	data, err := json.Marshal(results)
	if err != nil {
		t.Fatalf("Failed to marshal SearchResults: %v", err)
	}

	var unmarshaled SearchResults
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal SearchResults: %v", err)
	}

	if unmarshaled.Success != results.Success {
		t.Errorf("Success = %v, want %v", unmarshaled.Success, results.Success)
	}
	if unmarshaled.Query != results.Query {
		t.Errorf("Query = %v, want %v", unmarshaled.Query, results.Query)
	}
	if unmarshaled.Count != results.Count {
		t.Errorf("Count = %v, want %v", unmarshaled.Count, results.Count)
	}
	if len(unmarshaled.Results) != len(results.Results) {
		t.Errorf("Results length = %v, want %v", len(unmarshaled.Results), len(results.Results))
	}
}
