with import <nixpkgs> {}; {
  fabEnv = stdenv.mkDerivation {
    name = "fabric";
    buildInputs = [ stdenv Fabric ];
    shellHook =
      ''
      '';
  };
}
