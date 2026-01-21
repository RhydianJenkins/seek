package mcp

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rhydianjenkins/rag-mcp-server/src/config"
	"github.com/rhydianjenkins/rag-mcp-server/src/db"
	"github.com/rhydianjenkins/rag-mcp-server/src/handlers"
)

type RAGServer struct {
	mcpServer *mcp.Server
	storage   *db.Storage
}

func NewRAGServer() (*RAGServer, error) {
	cfg := config.Get()

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
}

type SearchToolInput struct {
	Query string `json:"query" jsonschema:"required" jsonschema_description:"Search query text"`
	Limit int    `json:"limit" jsonschema_description:"Maximum number of results to return (default: 3)"`
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

func (rs *RAGServer) Run(ctx context.Context) error {
	log.Println("Starting RAG MCP server on stdio transport...")

	if err := rs.mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return err
	}

	return nil
}
