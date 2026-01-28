package mcp

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rhydianjenkins/seek/src/config"
	"github.com/rhydianjenkins/seek/src/db"
	"github.com/rhydianjenkins/seek/src/handlers"
)

func NewRAGServer() (*RAGServer, error) {
	cfg := config.Get()
	if cfg == nil {
		return nil, fmt.Errorf("config not initialized: call config.Initialize() before creating server")
	}

	storage, err := db.Connect()
	if err != nil {
		return nil, err
	}

	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    cfg.ServerName,
		Version: cfg.ServerVersion,
	}, nil)

	ragServer := &RAGServer{
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
			Name:        "search",
			Description: "Search the RAG knowledge base for relevant content using semantic similarity.",
		},
		rs.handleSearchTool,
	)

	mcp.AddTool(
		rs.mcpServer,
		&mcp.Tool{
			Name:        "embed",
			Description: "Generate embeddings for documents in a directory and store them in the knowledge base.",
		},
		rs.handleEmbedTool,
	)

	mcp.AddTool(
		rs.mcpServer,
		&mcp.Tool{
			Name:        "status",
			Description: "Get the status of the knowledge base database, including collection information and vector count.",
		},
		rs.handleStatusTool,
	)

	mcp.AddTool(
		rs.mcpServer,
		&mcp.Tool{
			Name:        "get_document",
			Description: "Retrieve a full document by filename, returning all chunks in order.",
		},
		rs.handleGetDocumentTool,
	)
}

func (rs *RAGServer) handleSearchTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SearchToolInput,
) (*mcp.CallToolResult, *handlers.SearchResults, error) {
	if input.Limit == 0 {
		input.Limit = 3
	}

	log.Printf("Search tool called with query=%s, limit=%d", input.Query, input.Limit)

	results, err := handlers.SearchFiles(input.Query, input.Limit)
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

func (rs *RAGServer) handleEmbedTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input EmbedToolInput,
) (*mcp.CallToolResult, *handlers.EmbedResult, error) {
	if input.ChunkSize == 0 {
		input.ChunkSize = 1000
	}

	log.Printf("Embed tool called with dataDir=%s, chunkSize=%d", input.DataDir, input.ChunkSize)

	results, err := handlers.EmbedFiles(input.DataDir, input.ChunkSize)
	if err != nil {
		log.Printf("Embed tool error: %v", err)
		return &mcp.CallToolResult{
			IsError: true,
		}, results, err
	}

	log.Printf("Embed tool completed: %d files processed, %d chunks created", results.FilesIndexed, results.TotalChunks)

	return &mcp.CallToolResult{
		IsError: false,
	}, results, nil
}

func (rs *RAGServer) handleStatusTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StatusToolInput,
) (*mcp.CallToolResult, *db.CollectionStatus, error) {
	log.Println("Status tool called")

	status, err := rs.storage.GetStatus()
	if err != nil {
		log.Printf("Status tool error: %v", err)
		return &mcp.CallToolResult{
			IsError: true,
		}, status, err
	}

	log.Printf("Status tool completed: collection=%s, exists=%v, count=%d", status.CollectionName, status.Exists, status.VectorCount)

	return &mcp.CallToolResult{
		IsError: false,
	}, status, nil
}

func (rs *RAGServer) handleGetDocumentTool(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetDocumentToolInput,
) (*mcp.CallToolResult, *handlers.DocumentResult, error) {
	log.Printf("Get document tool called with filename=%s", input.Filename)

	result, err := handlers.GetDocumentByFilename(input.Filename)
	if err != nil {
		log.Printf("Get document tool error: %v", err)
		return &mcp.CallToolResult{
			IsError: true,
		}, result, err
	}

	log.Printf("Get document tool completed: filename=%s, chunks=%d", result.Filename, result.ChunkCount)

	return &mcp.CallToolResult{
		IsError: false,
	}, result, nil
}

func (rs *RAGServer) Run(ctx context.Context) error {
	log.Println("Starting RAG MCP server on stdio transport...")

	if err := rs.mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return err
	}

	return nil
}
