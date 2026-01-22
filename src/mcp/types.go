package mcp

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rhydianjenkins/seek/src/db"
)

type RAGServer struct {
	mcpServer *mcp.Server
	storage   *db.Storage
}

type SearchToolInput struct {
	Query string `json:"query" jsonschema:"required" jsonschema_description:"Search query text"`
	Limit int    `json:"limit" jsonschema_description:"Maximum number of results to return (default: 3)"`
}

type EmbedToolInput struct {
	DataDir   string `json:"dataDir" jsonschema:"required" jsonschema_description:"Directory containing .txt files to embed"`
	ChunkSize int    `json:"chunkSize" jsonschema_description:"Maximum chunk size in characters for splitting text (default: 1000)"`
}

type StatusToolInput struct{}

type GetDocumentToolInput struct {
	Filename string `json:"filename" jsonschema:"required" jsonschema_description:"The filename of the document to retrieve"`
}
