package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rhydianjenkins/rag-mcp-server/src/config"
	"github.com/rhydianjenkins/rag-mcp-server/src/handlers"
	"github.com/rhydianjenkins/rag-mcp-server/src/mcp"
	"github.com/spf13/cobra"
)

//go:embed VERSION
var version string
var logfile = "rag-mcp-server.log"

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
		Use:   "index",
		Short: "Index the knowledge base",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize(&config.Config{
				OllamaURL:      ollamaAddress,
				QdrantHost:     qdrantHost,
				QdrantPort:     qdrantPort,
				CollectionName: collectionName,
			})
			handlers.Index(dataDir, chunkSize)
		},
	}
	indexCmd.Flags().StringVar(&dataDir, "dataDir", "", "Directory containing .txt files to index (required)")
	indexCmd.Flags().IntVar(&chunkSize, "chunkSize", 1000, "Maximum chunk size in characters for splitting text")
	indexCmd.MarkFlagRequired("dataDir")
	rootCmd.AddCommand(indexCmd)

	var limit int
	var searchCmd = &cobra.Command{
		Use:   "search",
		Short: "Search the knowledge base",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize(&config.Config{
				OllamaURL:      ollamaAddress,
				QdrantHost:     qdrantHost,
				QdrantPort:     qdrantPort,
				CollectionName: collectionName,
			})
			searchTerm := args[0]
			handlers.Search(searchTerm, limit)
		},
	}
	searchCmd.Flags().IntVar(&limit, "limit", 3, "Maximum number of search results to return")
	rootCmd.AddCommand(searchCmd)

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(strings.TrimSpace(version))
		},
	}
	rootCmd.AddCommand(versionCmd)

	var runCmd = &cobra.Command{
		Use:   "run",
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

			config.Initialize(&config.Config{
				QdrantHost:     qdrantHost,
				QdrantPort:     qdrantPort,
				CollectionName: collectionName,
				OllamaURL:      ollamaAddress,
				ServerVersion:  strings.TrimSpace(version),
			})

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
	rootCmd.AddCommand(runCmd)

	return rootCmd
}

func main() {
	initCmd().Execute()
}
