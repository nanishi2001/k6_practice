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
        packages = import ./nix/packages.nix { inherit pkgs self system; };
        apps = import ./nix/apps.nix { inherit self system; };
        devShells = import ./nix/shells.nix { inherit pkgs; };
      }
    );
}
