#!/usr/bin/env bash
set -euo pipefail

self_dir="$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")"
root_dir=$(readlink -f "$self_dir/..")

version=$1
short_hash=$(git rev-parse --short HEAD)

go build -o build/knit -ldflags "-X 'knit/pkg/util.Version=$version' -X 'knit/pkg/util.ShortHash=$short_hash'" "$root_dir/main.go" 
