{
  description = "RAG MCP Server - A RAG server with MCP support";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };
        version = pkgs.lib.trim (builtins.readFile ./VERSION);
      in
      {
        packages = {
          seek = pkgs.buildGoModule {
            pname = "seek";
            inherit version;
            src = ./.;
            vendorHash = "sha256-LQ2dzPqxTBWTlWgX48ylAIqa5pSKAO4s8DejAYxEcvo=";

            meta = with pkgs.lib; {
              description = "RAG MCP Server";
              homepage = "https://github.com/rhydianjenkins/seek";
              license = licenses.mit;
            };
          };

          qdrant = pkgs.writeShellScriptBin "qdrant-service" ''
            # Load environment variables from .env if it exists
            set -a
            source .env.default
            if [ -f .env ]; then
              source .env
            fi
            set +a

            # Set Qdrant-specific environment variables
            export QDRANT__SERVICE__HOST=''${QDRANT_HOST:-localhost}
            export QDRANT__SERVICE__HTTP_PORT=''${QDRANT_PORT:-6334}

            exec ${pkgs.qdrant}/bin/qdrant
          '';

          ollama = pkgs.writeShellScriptBin "ollama-service" ''
            # Load environment variables from .env if it exists
            set -a
            source .env.default
            if [ -f .env ]; then
              source .env
            fi
            set +a

            export PATH="${pkgs.lib.makeBinPath [ pkgs.ollama pkgs.curl ]}:$PATH"

            # Use environment variables with defaults
            OLLAMA_HOST=''${OLLAMA_HOST:-localhost}
            OLLAMA_PORT=''${OLLAMA_PORT:-11434}

            # Start Ollama in the background
            ${pkgs.ollama}/bin/ollama serve &
            OLLAMA_PID=$!

            # Wait for Ollama to be ready with timeout
            echo "Starting Ollama..."
            MAX_RETRIES=30
            RETRY_COUNT=0
            until curl -s http://$OLLAMA_HOST:$OLLAMA_PORT/api/tags > /dev/null 2>&1; do
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

          open-webui = pkgs.writeShellScriptBin "open-webui-service" ''
            # Load environment variables from .env if it exists
            set -a
            source .env.default
            if [ -f .env ]; then
              source .env
            fi
            set +a

            # Set writable data directory
            export DATA_DIR="''${XDG_DATA_HOME:-$HOME/.local/share}/open-webui"
            mkdir -p "$DATA_DIR"

            # Use environment variables with defaults
            OLLAMA_HOST=''${OLLAMA_HOST:-localhost}
            OLLAMA_PORT=''${OLLAMA_PORT:-11434}
            WEBUI_PORT=''${WEBUI_PORT:-3000}

            export OLLAMA_BASE_URL=http://$OLLAMA_HOST:$OLLAMA_PORT
            exec ${pkgs.open-webui}/bin/open-webui serve --port $WEBUI_PORT
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
            echo ""
            echo "Or start individually:"
            echo "  nix run .#qdrant      - Start Qdrant vector database"
            echo "  nix run .#ollama      - Start Ollama server (auto-pulls nomic-embed-text)"
            echo "  nix run .#open-webui  - Start Open WebUI (Ollama chat interface)"
          '';
        };
      }
    );
}
