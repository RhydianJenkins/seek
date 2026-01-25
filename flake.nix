{
  description = "RAG MCP Server - A RAG server with MCP support";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        version = pkgs.lib.trim (builtins.readFile ./VERSION);
      in
      {
        packages = {
          seek = pkgs.buildGoModule {
            pname = "seek";
            inherit version;
            src = ./.;
            vendorHash = "sha256-3YeLHsu0kISSROHeyWPPISOQLbtw61KCkCy9Z/6RN6s=";

            meta = with pkgs.lib; {
              description = "RAG MCP Server";
              homepage = "https://github.com/rhydianjenkins/seek";
              license = licenses.mit;
            };
          };

          qdrant = pkgs.writeShellScriptBin "qdrant-service" ''
            exec ${pkgs.qdrant}/bin/qdrant
          '';

          ollama = pkgs.writeShellScriptBin "ollama-service" ''
            export PATH="${pkgs.lib.makeBinPath [ pkgs.ollama pkgs.curl ]}:$PATH"

            # Start Ollama in the background
            ${pkgs.ollama}/bin/ollama serve &
            OLLAMA_PID=$!

            # Wait for Ollama to be ready with timeout
            echo "Starting Ollama..."
            MAX_RETRIES=30
            RETRY_COUNT=0
            until curl -s http://localhost:11434/api/tags > /dev/null 2>&1; do
              RETRY_COUNT=$((RETRY_COUNT + 1))
              echo "Checking if Ollama is ready... ($RETRY_COUNT/$MAX_RETRIES)"
              if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
                echo "ERROR: Ollama failed to start after $MAX_RETRIES seconds" >&2
                kill $OLLAMA_PID 2>/dev/null
                exit 1
              fi
              sleep 1
            done
            echo "Ollama is ready!"

            # Check for and pull the model if needed
            if ollama list | grep -q "nomic-embed-text"; then
              echo "Model nomic-embed-text already exists"
            else
              echo "Pulling nomic-embed-text model (this may take a few minutes on first run)..."
              ollama pull nomic-embed-text
              echo "Model pulled successfully!"
            fi

            # Wait for the Ollama process
            wait $OLLAMA_PID
          '';

          default = self.packages.${system}.seek;
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
            curl
          ];

          shellHook = ''
            echo "RAG MCP Server development environment"
            echo "Go version: $(go version)"
            echo "Qdrant version: $(qdrant --version)"
            echo "Ollama version: $(ollama --version)"
            echo ""
            echo "Services:"
            echo "  nix run .#qdrant - Start Qdrant vector database"
            echo "  nix run .#ollama - Start Ollama server (auto-pulls nomic-embed-text)"
          '';
        };
      }
    );
}
