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
        # API サーバーパッケージ
        packages.api = pkgs.buildGoModule {
          pname = "k6-practice-api";
          version = "1.0.0";
          src = ./api;
          vendorHash = "sha256-deZ8L5aju1JraGTnjIW3vR1zm5Jc15F6D+goi1pVLpU=";

          meta = with pkgs.lib; {
            description = "k6 practice API server";
            license = licenses.mit;
          };
        };

        packages.default = self.packages.${system}.api;

        # アプリケーション定義
        apps.api = {
          type = "app";
          program = "${self.packages.${system}.api}/bin/api";
        };

        apps.default = self.apps.${system}.api;

        # API サーバー用環境 (Go)
        devShells.api = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
          ];

          shellHook = ''
            echo "=== API Server Environment ==="
            echo "Go version: $(go version)"
            echo ""
            echo "Commands:"
            echo "  cd api && go run .             - Start API server"
            echo "  cd api && go build -o api .    - Build API binary"
            echo ""
            echo "API will listen on http://localhost:8080"
          '';
        };

        # テスト実行用環境 (k6 + Bun)
        devShells.test = pkgs.mkShell {
          buildInputs = with pkgs; [
            k6
            bun
          ];

          shellHook = ''
            echo "=== Test Environment ==="
            echo "k6 version: $(k6 version)"
            echo "Bun version: $(bun --version)"
            echo ""
            echo "Target API: ''${API_URL:-http://localhost:8080}"
            echo ""
            echo "Setup: bun install"
            echo ""
            echo "Commands:"
            echo "  bun run build                  - Build TypeScript"
            echo "  bun run test:load              - Run load test"
            echo "  bun run test:stress            - Run stress test"
            echo "  bun run test:spike             - Run spike test"
            echo "  bun run test:all               - Run all tests"
            echo ""
            echo "Custom API URL:"
            echo "  API_URL=http://example.com bun run test:load"
          '';
        };

        # 統合環境 (開発用)
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            k6
            bun
            go
          ];

          shellHook = ''
            echo "=== Full Development Environment ==="
            echo "k6 version: $(k6 version)"
            echo "Bun version: $(bun --version)"
            echo "Go version: $(go version)"
            echo ""
            echo "Setup: bun install"
            echo ""
            echo "Available commands:"
            echo "  bun run api                    - Start API server"
            echo "  bun run build                  - Build TypeScript"
            echo "  bun run test:load              - Run load test"
            echo "  bun run test:stress            - Run stress test"
            echo "  bun run test:spike             - Run spike test"
            echo ""
            echo "Separate environments:"
            echo "  nix develop .#api              - API only (Go)"
            echo "  nix develop .#test             - Test only (k6 + Bun)"
          '';
        };
      }
    );
}
