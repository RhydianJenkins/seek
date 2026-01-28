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
<summary>Use with Docker</summary>

Build the Docker image using Nix (requires Nix with flakes enabled):
```sh
# Build the Docker image
nix build .#docker

# Load it into Docker
docker load < result
```

This creates a `seek` image tagged with the version from the VERSION file (e.g., `seek:1.0.0`).

Run the required services and seek:
```sh
# Start Qdrant
docker run -d --name qdrant \
  -p 6333:6333 -p 6334:6334 \
  -v qdrant_storage:/qdrant/storage \
  qdrant/qdrant:latest

# Start Ollama
docker run -d --name ollama \
  -p 11434:11434 \
  -v ollama_data:/root/.ollama \
  ollama/ollama:latest

# Pull the embedding model
docker exec ollama ollama pull nomic-embed-text

# Run seek (adjust tag to match your version)
docker run --rm \
  -e QDRANT_HOST=host.docker.internal \
  -e QDRANT_PORT=6333 \
  -e OLLAMA_HOST=host.docker.internal \
  -e OLLAMA_PORT=11434 \
  -v $(pwd)/data:/data \
  seek:1.0.0 --help
```

Note: Use `host.docker.internal` on Mac/Windows or `172.17.0.1` on Linux to connect to services running on the host.

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
- [x] Docker image support
- [ ] `seek ask` command that uses ollama for answers
