# MCP Server for Retrieval-Augmented Generation (RAG).

Search through a given knowledge base for relevant information using embedded natural language.

This is an early prototype and breaking changes are likely.

## Features

- **Semantic Search**: Search your knowledge base using natural language queries with vector embeddings
- **Document Management**: Index and retrieve complete documents from your knowledge base
- **MCP Integration**: Full Model Context Protocol support for integration with MCP clients

## Supported File Formats

- **Plain Text** - `.txt`, `.md`, and other text-based files
- **PDF** - `.pdf` documents with text extraction
- **Word** - `.docx` documents with text extraction
- **Excel** - `.xlsx` spreadsheets (all sheets and cells)

# Getting Started

<details>
<summary>Install with Nix</summary>

Install the binary to your system:
```sh
nix profile install github:rhydianjenkins/seek
seek --help
```

Or use without installing:
```sh
# Run directly without installing
nix run github:rhydianjenkins/seek -- --help

# Temporary shell with seek available
nix shell github:rhydianjenkins/seek
seek --help
```

Start the required services ([Ollama](https://ollama.com) and [Qdrant](https://qdrant.tech)):
```sh
# Start Qdrant (in one terminal)
nix run github:rhydianjenkins/seek#qdrant

# Start Ollama (in another terminal - auto-pulls nomic-embed-text model)
nix run github:rhydianjenkins/seek#ollama

# Or bring your own services by creating a .env file
cp .env.default .env
```

</details>

<details>
<summary>Install Globally</summary>

Install the binary globally on your system using Go:
```sh
go install github.com/rhydianjenkins/seek@latest
seek --help
```

Make sure `$GOPATH/bin` (usually `~/go/bin`) is in your PATH.

</details>

<details>
<summary>Build from Source</summary>

```sh
git clone git@github.com:rhydianjenkins/seek
cd seek
go build
./seek --help
```

</details>

## Commands

### Embed Documents
Generate embeddings for all documents in a directory:
```sh
seek embed --dataDir test-data --chunkSize 1000
```

### Search Knowledge Base
Search for documents using natural language:
```sh
seek search "What is important for me to do this week?" --limit 3
```

### Get Document
Retrieve a complete document by filename:
```sh
# show all documents in current database
seek list

seek get "document.txt"
```

### Check Status
View the status of your knowledge base:
```sh
seek status
```

### Run MCP Server
Start the MCP server for integration with MCP clients:
```sh
seek mcp
```

## MCP Tools

When running as an MCP server, the following tools are available:

- `search` - Search the knowledge base using semantic similarity
- `embed` - Generate embeddings for documents in a directory
- `get_document` - Retrieve a full document by filename
- `status` - Get database status and statistics

## TODO

- [ ] Add auth/TLS support
- [ ] Image/OCR support
- [ ] Docker image support
- [ ] `seek ask` command that uses ollama for answers
