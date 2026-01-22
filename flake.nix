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
        version = builtins.readFile ./VERSION;
      in
      {
        packages = {
          seek = pkgs.buildGoModule {
            pname = "seek";
            inherit version;
            src = ./.;
            vendorHash = "sha256-Q7JcnsdRuHA0bUJ711UG4KqbrnUfWND+yY/PW/2x5Ig=";

            meta = with pkgs.lib; {
              description = "RAG MCP Server";
              homepage = "https://github.com/rhydianjenkins/seek";
              license = licenses.mit;
            };
          };

          start-services = pkgs.writeShellScriptBin "start-services" ''
            export PATH="${pkgs.lib.makeBinPath [ pkgs.qdrant pkgs.ollama pkgs.curl ]}:$PATH"
            exec ${pkgs.process-compose}/bin/process-compose -f ${./process-compose.yaml}
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
            process-compose
            curl
          ];

          shellHook = ''
            echo "RAG MCP Server development environment"
            echo "Go version: $(go version)"
            echo "Qdrant version: $(qdrant --version)"
            echo "Ollama version: $(ollama --version)"
            echo ""
            echo "To start all services, do 'nix run .#start-services'"
          '';
        };
      }
    );
}
