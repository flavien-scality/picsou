with import <nixpkgs> { };

let
  buildDeps = with pkgs; [ stdenv buildGoPackage fetchgit gnumake ];
in

buildGoPackage rec {
  name = "picsou-${version}";
  version = "20170120";

  goPackagePath = "github.com/scality/picsou";

  goDeps = with buildDeps; ./deps.nix;

  extraCmds = ''
    export GOPATH=~/go:$GOPATH
  '';
}
