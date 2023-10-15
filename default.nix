{ pkgs, ... }:

pkgs.buildGoApplication rec {
  pname = "abots";
  version = "0.0.2";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;

  ldflags = [
    "-X github.com/adamgoose/abots/cmd.Version=${version}"
  ];
}
