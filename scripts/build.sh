#!/usr/bin/env bash
set -euo pipefail

self_dir="$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")"
root_dir=$(readlink -f "$self_dir/..")

version=${1:-$(git rev-parse --abbrev-ref HEAD)}
short_hash=$(git rev-parse --short HEAD)

go build -o "$root_dir/build/knit" -ldflags \
    "-w -s -X 'knit/pkg/util.Version=$version' -X 'knit/pkg/util.ShortHash=$short_hash' -X 'knit/pkg/util.Architecture=${GOARCH:-amd64}'" \
    "$root_dir/main.go" 
