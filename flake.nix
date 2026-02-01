{
  description = "k6 performance testing environment with TypeScript";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            k6
            nodejs_22
            pnpm
            go
          ];

          shellHook = ''
            echo "k6 + TypeScript performance testing environment"
            echo "k6 version: $(k6 version)"
            echo "Node.js version: $(node --version)"
            echo "pnpm version: $(pnpm --version)"
            echo "Go version: $(go version)"
            echo ""
            echo "Setup: pnpm install"
            echo ""
            echo "Available commands:"
            echo "  pnpm build                     - Build TypeScript"
            echo "  pnpm test:load                 - Run load test"
            echo "  pnpm test:stress               - Run stress test"
            echo "  pnpm test:spike                - Run spike test"
            echo "  pnpm api                       - Start API server"
          '';
        };
      }
    );
}
