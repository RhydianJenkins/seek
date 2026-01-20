package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rhydianjenkins/rag-mcp-server/src"
	"github.com/rhydianjenkins/rag-mcp-server/src/config"
	"github.com/rhydianjenkins/rag-mcp-server/src/handlers"
)

type RAGServer struct {
	config    *config.Config
	mcpServer *mcp.Server
	storage   *src.Storage
}

func NewRAGServer(cfg *config.Config) (*RAGServer, error) {
	storage, err := src.Connect(cfg)
	if err != nil {
		return nil, err
	}

	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    cfg.ServerName,
		Version: cfg.ServerVersion,
	}, nil)

	ragServer := &RAGServer{
		config:    cfg,
		mcpServer: mcpServer,
		storage:   storage,
	}

	ragServer.registerTools()

	return ragServer, nil
}

func (rs *RAGServer) registerTools() {
	mcp.AddTool(
		rs.mcpServer,
		&mcp.Tool{
			Name:        "index",
			Description: "Index text files from a directory into the RAG knowledge base. Splits files into chunks and generates embeddings.",
		},
		rs.handleIndexTool,
	)

	mcp.AddTool(
		rs.mcpServer,
		&mcp.Tool{
			Name:        "search",
			Description: "Search the RAG knowledge base for relevant content using semantic similarity.",
		},
		rs.handleSearchTool,
	)
}

type IndexToolInput struct {
	DataDir   string `json:"dataDir" jsonschema:"required" jsonschema_description:"Directory containing .txt files to index"`
	ChunkSize int    `json:"chunkSize" jsonschema_description:"Maximum chunk size in characters (default: 1000)"`
}

type SearchToolInput struct {
	Query string `json:"query" jsonschema:"required" jsonschema_description:"Search query text"`
	Limit int    `json:"limit" jsonschema_description:"Maximum number of results to return (default: 3)"`
}

func (rs *RAGServer) handleIndexTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input IndexToolInput,
) (*mcp.CallToolResult, *handlers.IndexResult, error) {
	// Apply defaults
	if input.ChunkSize == 0 {
		input.ChunkSize = 1000
	}

	log.Printf("Index tool called with dataDir=%s, chunkSize=%d", input.DataDir, input.ChunkSize)

	// Call refactored handler
	result, err := handlers.IndexFiles(rs.config, input.DataDir, input.ChunkSize)
	if err != nil {
		log.Printf("Index tool error: %v", err)
		return &mcp.CallToolResult{
			IsError: true,
		}, result, err
	}

	log.Printf("Index tool completed: %d files, %d chunks", result.FilesIndexed, result.TotalChunks)

	return &mcp.CallToolResult{
		IsError: false,
	}, result, nil
}

func (rs *RAGServer) handleSearchTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SearchToolInput,
) (*mcp.CallToolResult, *handlers.SearchResults, error) {
	// Apply defaults
	if input.Limit == 0 {
		input.Limit = 3
	}

	log.Printf("Search tool called with query=%s, limit=%d", input.Query, input.Limit)

	// Call refactored handler
	results, err := handlers.SearchFiles(rs.config, input.Query, input.Limit)
	if err != nil {
		log.Printf("Search tool error: %v", err)
		return &mcp.CallToolResult{
			IsError: true,
		}, results, err
	}

	log.Printf("Search tool completed: %d results found", results.Count)

	return &mcp.CallToolResult{
		IsError: false,
	}, results, nil
}

// Run starts the MCP server with stdio transport
func (rs *RAGServer) Run(ctx context.Context) error {
	log.Println("Starting RAG MCP server on stdio transport...")

	// Run server with stdio transport
	if err := rs.mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return err
	}

	return nil
}
