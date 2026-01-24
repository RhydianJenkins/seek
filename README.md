# MCP Server for Retrieval-Augmented Generation (RAG).

Search through a given knowledge base for relevant information using embedded natural language.

This is an early prototype and breaking changes are likely.

## Features

- **Semantic Search**: Search your knowledge base using natural language queries with vector embeddings
- **Document Management**: Index and retrieve complete documents from your knowledge base
- **MCP Integration**: Full Model Context Protocol support for integration with MCP clients

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
nix run github:rhydianjenkins/seek#start-services

# Or you can bring your own services
seek --ollamaHost your.ollama.host \
    --ollamaPort 11434 \
    --qdrantHost your.qdrant.host \
    --qdrantPort 6334 \
    --dataDir test-data \
    [command]
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

## Configuration

Configure the server using command-line flags:

- `--ollamaHost` - Ollama server host (default: localhost)
- `--ollamaPort` - Ollama server port (default: 11434)
- `--qdrantHost` - Qdrant server host (default: localhost)
- `--qdrantPort` - Qdrant server port (default: 6334)
- `--collection` - Qdrant collection name (default: my_collection)

## Supported File Types

- `.txt` - Plain text files
- `.pdf` - PDF documents

## TODO

- [ ] Incremental updates - you lose everything each time you re-embed
- [ ] Backup/restore functionality
- [ ] No auth/TLS support
- [ ] Better file support (.md, .docx, .pptx, .xlsx, etc)
- [ ] Allow embedding model to be specified (currently hardcoded to `nomic-embed-text`)
- [ ] Image/OCR support
