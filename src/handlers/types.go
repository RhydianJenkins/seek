package handlers

type IndexResult struct {
	Success      bool   `json:"success"`
	FilesIndexed int    `json:"files_indexed"`
	TotalChunks  int    `json:"total_chunks"`
	Message      string `json:"message"`
	Error        string `json:"error,omitempty"`
}

type SearchResult struct {
	Score      float32 `json:"score"`
	Filename   string  `json:"filename"`
	ChunkIndex int64   `json:"chunk_index"`
	Content    string  `json:"content"`
}

type SearchResults struct {
	Success bool           `json:"success"`
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
	Count   int            `json:"count"`
	Error   string         `json:"error,omitempty"`
}

type DocumentChunk struct {
	ChunkIndex int64  `json:"chunk_index"`
	Content    string `json:"content"`
}

type DocumentResult struct {
	Success    bool            `json:"success"`
	Filename   string          `json:"filename"`
	ChunkCount int             `json:"chunk_count"`
	Chunks     []DocumentChunk `json:"chunks"`
	FullText   string          `json:"full_text"`
	Error      string          `json:"error,omitempty"`
}
