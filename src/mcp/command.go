package mcp

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func NewCommand(logfile string) *cobra.Command {
	var httpMode bool
	var httpPort int

	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Run the MCP server over stdio or HTTP",
		Long:  "Starts the MCP server using stdio transport (default) or HTTP transport with --http flag",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// Only set up log file for stdio mode
			if !httpMode {
				logFile, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					log.Fatalf("Failed to open log file: %v", err)
				}
				defer logFile.Close()
				log.SetOutput(logFile)
			}

			ragServer, err := NewRAGServer()
			if err != nil {
				log.Fatalf("Failed to create RAG server: %v", err)
			}

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			if httpMode {
				if err := ragServer.RunHTTP(ctx, httpPort); err != nil {
					log.Fatalf("MCP server error: %v", err)
				}
			} else {
				if err := ragServer.Run(ctx); err != nil {
					log.Fatalf("MCP server error: %v", err)
				}
			}
		},
	}

	cmd.Flags().BoolVar(&httpMode, "http", false, "Run server in HTTP mode instead of stdio")
	cmd.Flags().IntVar(&httpPort, "port", 8080, "Port to listen on when using --http mode")

	return cmd
}
