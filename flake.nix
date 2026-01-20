{
  description = "RAG MCP Server - A RAG server with MCP support";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    process-compose-flake.url = "github:Platonic-Systems/process-compose-flake";
  };

  outputs = { self, nixpkgs, flake-utils, process-compose-flake }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        # Define the process composition
        packages.default = pkgs.writeShellScriptBin "services" ''
          ${pkgs.process-compose}/bin/process-compose -f ${pkgs.writeText "process-compose.yaml" ''
            version: "0.5"
            processes:
              qdrant:
                command: ${pkgs.qdrant}/bin/qdrant
                availability:
                  restart: on_failure
              ollama:
                command: ${pkgs.ollama}/bin/ollama serve
                availability:
                  restart: on_failure
          ''}
        '';

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/services";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            golangci-lint
            qdrant
            ollama
            process-compose
          ];

          shellHook = ''
            echo "RAG MCP Server development environment"
            echo "Go version: $(go version)"
            echo "Qdrant version: $(qdrant --version)"
            echo "Ollama version: $(ollama --version)"
            echo ""
            echo "Services available:"
            echo "  - Qdrant (vector database)"
            echo "  - Ollama (local embeddings)"
            echo ""
            echo "To start all services: nix run"
            echo "Or from dev shell: process-compose up"
          '';
        };
      }
    );
}
