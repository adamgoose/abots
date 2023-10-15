{ pkgs, pkgsLinux, ... }:
let
  abots = pkgsLinux.callPackage ./default.nix { };
in
pkgs.dockerTools.buildImage {
  name = "abots";
  config = {
    Cmd = [ "${abots}/bin/abots" ];
  };
}
