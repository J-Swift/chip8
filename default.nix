with import <nixpkgs> {};

stdenv.mkDerivation rec {
  name = "chip8";
  env = buildEnv { name = name; paths = buildInputs; };
  buildInputs = [
    bash

    go
  ];
}
