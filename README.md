MCP Server for Retrieval-Augmented Generation (RAG)

# Getting Started

With [Nix](https://nixos.org/download), you can start the services by doing:

```sh
nix run # start services
```

Then, enter a dev shell with:

```sh
nix develop # enter dev shell with all dependencies
go run main.go # print help information
```

You will then need to index some knowledge bases.
For example, you can index the test emails in this repository:

```sh
go run . index --dataDir test-data/emails
```

Then you can run a search the knowledge base

```sh
go run . search "is there any deadlines coming up?"
```
