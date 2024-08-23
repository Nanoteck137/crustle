{
  description = "Devshell and Building crustle";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";
    tagopus.url      = "github:nanoteck137/tagopus/v0.1.1";
  };

  outputs = { self, nixpkgs, flake-utils, tagopus, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        module = pkgs.buildGoModule {
          pname = "crustle";
          version = self.shortRev or "dirty";
          src = ./.;

          vendorHash = "sha256-+5XrSEgxv4YOtbiJkhfBOZtma7zw9bieMM4gCYK+IUo=";

          nativeBuildInputs = [pkgs.makeWrapper];

          postFixup = ''
            wrapProgram $out/bin/crustle \
            --set PATH ${pkgs.lib.makeBinPath [
              tagopus.packages.${system}.default
            ]}
          '';
        };
      in
      {
        packages.default = module;
        packages.crustle = module;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            tagopus.packages.${system}.default
          ];
        };
      }
    );
}
