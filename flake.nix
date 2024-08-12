{
  description = "GoSkew is a program for post-processing g-code to account for axis skew.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      perSystem = { config, pkgs, ... }: {
        packages = rec {
          goskew = pkgs.callPackage ./package.nix { };
          default = goskew;
        };
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [ git go_1_22 ];
        };

        overlayAttrs = {
          inherit (config.packages) goskew;
        };
      };
      systems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin" ];
      imports = [ inputs.flake-parts.flakeModules.easyOverlay ];
    };
}
