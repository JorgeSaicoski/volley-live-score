{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell{
  buildInputs = [
    pkgs.go
    pkgs.vscodium
    pkgs.docker
  ];
}