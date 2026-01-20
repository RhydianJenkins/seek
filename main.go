package main

import (
	"context"
	"log"

	"github.com/rhydianjenkins/rag-mcp-server/src/config"
	"github.com/rhydianjenkins/rag-mcp-server/src/handlers"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	var (
		ollamaAddress  string
		qdrantHost     string
		qdrantPort     int
		collectionName string
	)

	var rootCmd = &cobra.Command{
		Short: "RAG MCP Server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.PersistentFlags().StringVar(&ollamaAddress, "ollamaAddress", "http://localhost:11434", "Ollama server address")
	rootCmd.PersistentFlags().StringVar(&qdrantHost, "qdrantHost", "localhost", "Qdrant server host")
	rootCmd.PersistentFlags().IntVar(&qdrantPort, "qdrantPort", 6334, "Qdrant server port")
	rootCmd.PersistentFlags().StringVar(&collectionName, "collection", "my_collection", "Qdrant collection name")

	var dataDir string
	var chunkSize int
	var indexCmd = &cobra.Command{
		Use: "index",
		Short: "Index the knowledge base",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			handlers.Index(ollamaAddress, dataDir, chunkSize)
		},
	}
	indexCmd.Flags().StringVar(&dataDir, "dataDir", "", "Directory containing .txt files to index (required)")
	indexCmd.Flags().IntVar(&chunkSize, "chunkSize", 1000, "Maximum chunk size in characters for splitting text")
	indexCmd.MarkFlagRequired("dataDir")
	rootCmd.AddCommand(indexCmd)

	var limit int
	var searchCmd = &cobra.Command{
		Use: "search",
		Short: "Search the knowledge base",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			searchTerm := args[0]
			handlers.Search(searchTerm, ollamaAddress, limit)
		},
	}
	searchCmd.Flags().IntVar(&limit, "limit", 3, "Maximum number of search results to return")
	rootCmd.AddCommand(searchCmd)

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the MCP server over stdio",
		Long:  "Starts the MCP server using stdio transport for integration with MCP clients",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &config.Config{
				QdrantHost:     qdrantHost,
				QdrantPort:     qdrantPort,
				QdrantUseTLS:   false,
				CollectionName: collectionName,
				VectorSize:     768,
				OllamaURL:      ollamaAddress,
				EmbeddingModel: "nomic-embed-text",
				ServerName:     "rag-mcp-server",
				ServerVersion:  "1.0.0",
			}

			ragServer, err := NewRAGServer(cfg)
			if err != nil {
				log.Fatalf("Failed to create RAG server: %v", err)
			}

			ctx := context.Background()
			if err := ragServer.Run(ctx); err != nil {
				log.Fatalf("MCP server error: %v", err)
			}
		},
	}
	rootCmd.AddCommand(runCmd)

	return rootCmd
}

func main() {
	initCmd().Execute()
}
