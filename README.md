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

Pull the Docker image:
```sh
docker pull ghcr.io/rhydianjenkins/seek:latest
```

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

# Run seek
docker run --rm \
  -e QDRANT_HOST=host.docker.internal \
  -e QDRANT_PORT=6333 \
  -e OLLAMA_HOST=host.docker.internal \
  -e OLLAMA_PORT=11434 \
  -v $(pwd)/data:/data \
  ghcr.io/rhydianjenkins/seek:latest --help
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

# Development

## Publishing Docker Images

To publish a new version to GitHub Container Registry:

```sh
# Build the image with Nix
nix build .#docker

# Load it into Docker
docker load < result

# Tag for GitHub Container Registry
docker tag seek:1.0.0 ghcr.io/rhydianjenkins/seek:1.0.0
docker tag seek:1.0.0 ghcr.io/rhydianjenkins/seek:latest

# Login to GitHub Container Registry (requires a PAT with write:packages scope)
echo $CR_PAT | docker login ghcr.io -u rhydianjenkins --password-stdin

# Push both tags
docker push ghcr.io/rhydianjenkins/seek:1.0.0
docker push ghcr.io/rhydianjenkins/seek:latest
```

Make sure the repository visibility is set to public in GitHub Container Registry settings.

## TODO

- [ ] Add auth/TLS support
- [ ] Image/OCR support
- [ ] `seek ask` command that uses ollama for answers
