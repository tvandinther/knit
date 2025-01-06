{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  packages = with pkgs.buildPackages; [
    git
    kubectl
    kustomize
    kubernetes-helm
    kcl
    jq
    yq-go
    go
    cobra-cli
  ];

  shellHook = ''
    export PATH=$(pwd)/scripts:$(pwd)/build:$PATH
  '';
}

