package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rhydianjenkins/seek/src/config"
	"github.com/rhydianjenkins/seek/src/handlers"
	"github.com/rhydianjenkins/seek/src/mcp"
	"github.com/spf13/cobra"
)

//go:embed VERSION
var version string
var logfile = "seek.log"

func loadEnvFiles() {
	godotenv.Overload(".env.default", ".env")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func initCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Short: "RAG MCP Server",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ollamaHost := getEnv("OLLAMA_HOST", "localhost")
			ollamaPort := getEnvInt("OLLAMA_PORT", 11434)
			ollamaURL := fmt.Sprintf("http://%s:%d", ollamaHost, ollamaPort)
			config.Initialize(&config.Config{
				CollectionName: getEnv("COLLECTION_NAME", "seek_collection"),
				OllamaURL:      ollamaURL,
				QdrantHost:     getEnv("QDRANT_HOST", "localhost"),
				QdrantPort:     getEnvInt("QDRANT_PORT", 6334),
				ServerVersion:  strings.TrimSpace(version),
			})
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var dataDir string
	var chunkSize int
	var embedCmd = &cobra.Command{
		Use:   "embed",
		Short: "Generate embeddings for the knowledge base",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			handlers.Embed(dataDir, chunkSize)
		},
	}
	embedCmd.Flags().StringVar(&dataDir, "dataDir", "", "Directory containing .txt files to embed (required)")
	embedCmd.Flags().IntVar(&chunkSize, "chunkSize", 1000, "Maximum chunk size in characters for splitting text")
	embedCmd.MarkFlagRequired("dataDir")
	rootCmd.AddCommand(embedCmd)

	var limit int
	var searchCmd = &cobra.Command{
		Use:   "search",
		Short: "Search the knowledge base",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			searchTerm := args[0]
			handlers.Search(searchTerm, limit)
		},
	}
	searchCmd.Flags().IntVar(&limit, "limit", 3, "Maximum number of search results to return")
	rootCmd.AddCommand(searchCmd)

	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a full document by filename",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]
			handlers.GetDocument(filename)
		},
	}
	rootCmd.AddCommand(getCmd)

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show the status of the database",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			handlers.Status()
		},
	}
	rootCmd.AddCommand(statusCmd)

	var listLimit int
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all document names in the database",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			handlers.List(listLimit)
		},
	}
	listCmd.Flags().IntVar(&listLimit, "limit", 100, "Maximum number of documents to scan")
	rootCmd.AddCommand(listCmd)

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(strings.TrimSpace(version))
		},
	}
	rootCmd.AddCommand(versionCmd)

	var mcpCmd = &cobra.Command{
		Use:   "mcp",
		Short: "Run the MCP server over stdio",
		Long:  "Starts the MCP server using stdio transport for integration with MCP clients",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			logFile, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("Failed to open log file: %v", err)
			}
			defer logFile.Close()
			log.SetOutput(logFile)

			ragServer, err := mcp.NewRAGServer()
			if err != nil {
				log.Fatalf("Failed to create RAG server: %v", err)
			}

			ctx := context.Background()
			if err := ragServer.Run(ctx); err != nil {
				log.Fatalf("MCP server error: %v", err)
			}
		},
	}
	rootCmd.AddCommand(mcpCmd)

	return rootCmd
}

func main() {
	loadEnvFiles()
	initCmd().Execute()
}
