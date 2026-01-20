package handlers

// IndexResult represents the result of an indexing operation
type IndexResult struct {
	Success      bool   `json:"success"`
	FilesIndexed int    `json:"files_indexed"`
	TotalChunks  int    `json:"total_chunks"`
	Message      string `json:"message"`
	Error        string `json:"error,omitempty"`
}

// SearchResult represents a single search result
type SearchResult struct {
	Score      float32 `json:"score"`
	Filename   string  `json:"filename"`
	ChunkIndex int64   `json:"chunk_index"`
	Content    string  `json:"content"`
}

// SearchResults represents the complete search response
type SearchResults struct {
	Success bool           `json:"success"`
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
	Count   int            `json:"count"`
	Error   string         `json:"error,omitempty"`
}
