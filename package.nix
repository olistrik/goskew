{ pkgs, ... }:
pkgs.buildGoModule {
  name = "goskew";
  version = "1.6.0";
  vendorHash = "sha256-83r+2dK0gHF/f2iyBfa0N7KlemOjdngFx1c1Yx00+as=";
  src = ./.;
}
