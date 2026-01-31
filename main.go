package main

import (
	_ "embed"
	"fmt"
	"log"
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

func initCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Short: "RAG MCP Server",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.Initialize(&config.Config{
				ServerVersion: strings.TrimSpace(version),
			})
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var dataDir string
	var chunkSize int
	var embedCmd = &cobra.Command{
		Use:   "embed --dataDir <directory>",
		Short: "Generate embeddings for the knowledge base",
		Long:  "Process documents in a directory, split them into chunks, generate embeddings using Ollama, and store them in the Qdrant vector database for semantic search.",
		Example: `  seek embed --dataDir ./documents
  seek embed --dataDir ./docs --chunkSize 500`,
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			handlers.Embed(dataDir, chunkSize)
		},
	}
	embedCmd.Flags().StringVar(&dataDir, "dataDir", "", "Directory containing .txt files to embed (required)")
	embedCmd.Flags().IntVar(&chunkSize, "chunkSize", 1000, "Maximum chunk size in characters for splitting text")
	embedCmd.MarkFlagRequired("dataDir")
	rootCmd.AddCommand(embedCmd)

	var askCmd = &cobra.Command{
		Use:   "ask <question>",
		Short: "Ask a question about the knowledge base",
		Long:  "Ask a natural language question and get answers based on your indexed documents. The AI will search the knowledge base and provide relevant information.",
		Example: `  seek ask "What is the company culture?"
  seek ask "How does authentication work?"`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			question := args[0]
			err := handlers.AskQuestion(question)

			if err != nil {
				log.Println("Error:", err)
			}
		},
	}
	rootCmd.AddCommand(askCmd)

	var limit int
	var searchCmd = &cobra.Command{
		Use:   "search <query>",
		Short: "Search the knowledge base",
		Long:  "Perform semantic search on the indexed documents using embeddings. Returns the most relevant chunks of text based on similarity.",
		Example: `  seek search "authentication"
  seek search "company culture" --limit 5`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			searchTerm := args[0]
			handlers.Search(searchTerm, limit)
		},
	}
	searchCmd.Flags().IntVar(&limit, "limit", 3, "Maximum number of search results to return")
	rootCmd.AddCommand(searchCmd)

	var getCmd = &cobra.Command{
		Use:   "get <filename>",
		Short: "Get a full document by filename",
		Long:  "Retrieve and display the complete contents of a document from the knowledge base by its filename. All chunks are reassembled in order.",
		Example: `  seek get "README.md"
  seek get "docs/architecture.txt"`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]
			handlers.GetDocument(filename)
		},
	}
	rootCmd.AddCommand(getCmd)

	var statusCmd = &cobra.Command{
		Use:     "status",
		Short:   "Show the status of the database",
		Long:    "Display information about the Qdrant vector database including whether the collection exists, how many vectors are stored, and collection configuration.",
		Example: `  seek status`,
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			handlers.Status()
		},
	}
	rootCmd.AddCommand(statusCmd)

	var listLimit int
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all document names in the database",
		Long:  "Display all unique document filenames that have been indexed in the knowledge base.",
		Example: `  seek list
  seek list --limit 50`,
		Args: cobra.ExactArgs(0),
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

	rootCmd.AddCommand(mcp.NewCommand(logfile))

	return rootCmd
}

func main() {
	godotenv.Overload(".env.default", ".env")
	initCmd().Execute()
}
