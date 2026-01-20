package main

import (
	"fmt"

	"github.com/rhydianjenkins/rag-mcp-server/src/handlers"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Short: "RAG MCP Server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var indexCmd = &cobra.Command{
		Use: "index",
		Short: "Index the knowledge base",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO Rhydian point to knowledge base source
			handlers.Index()
		},
	}
	rootCmd.AddCommand(indexCmd)

	var searchCmd = &cobra.Command{
		Use: "search",
		Short: "Search the knowledge base",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			searchTerm := args[0]
			handlers.Search(searchTerm)
		},
	}
	rootCmd.AddCommand(searchCmd)

	var runCmd = &cobra.Command{
		Use: "run",
		Short: "Run the mcp server",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO")
		},
	}
	rootCmd.AddCommand(runCmd)

	return rootCmd
}

func main() {
	initCmd().Execute()
}
