MCP Server for Retrieval-Augmented Generation (RAG).

Search through a given knowledge base for relevant information using embedded natural language.

This is an early prototype and breaking changes are likely.

## Features

- **Semantic Search**: Search your knowledge base using natural language queries with vector embeddings
- **Document Management**: Index and retrieve complete documents from your knowledge base
- **MCP Integration**: Full Model Context Protocol support for integration with MCP clients

# Getting Started

```sh
git clone git@github.com:rhydianjenkins/seek && cd seek
```

## then build from source

```sh
go build
./seek
```

## or with Nix

```sh
# start ollama and qdrant services
nix run .#start-services

# run the seek command
nix run . -- [flags] [command]
```

## Commands

### Embed Documents
Generate embeddings for all documents in a directory:
```sh
./seek embed --dataDir test-data --chunkSize 1000
```

### Search Knowledge Base
Search for documents using natural language:
```sh
./seek search "What is important for me to do this week?" --limit 3
```

### Get Document
Retrieve a complete document by filename:
```sh
./seek get "document.txt"
```

### Check Status
View the status of your knowledge base:
```sh
./seek status
```

### Run MCP Server
Start the MCP server for integration with MCP clients:
```sh
./seek mcp
```

## MCP Tools

When running as an MCP server, the following tools are available:

- `search` - Search the knowledge base using semantic similarity
- `embed` - Generate embeddings for documents in a directory
- `get_document` - Retrieve a full document by filename
- `status` - Get database status and statistics

## Configuration

Configure the server using command-line flags:

- `--ollamaHost` - Ollama server host (default: localhost)
- `--ollamaPort` - Ollama server port (default: 11434)
- `--qdrantHost` - Qdrant server host (default: localhost)
- `--qdrantPort` - Qdrant server port (default: 6334)
- `--collection` - Qdrant collection name (default: my_collection)
