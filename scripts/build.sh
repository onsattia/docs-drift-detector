#!/bin/sh
# Builds docs-drift with the version from VERSION baked in.
# Usage: scripts/build.sh <output-path>
set -eu

out="${1:?usage: scripts/build.sh <output-path>}"
case "$out" in
/*) ;;
*) out="$PWD/$out" ;;
esac

root="$(cd "$(dirname "$0")/.." && pwd)"
version="$(tr -d '[:space:]' < "$root/VERSION")"

cd "$root"
go build -trimpath -ldflags "-X main.version=$version" -o "$out" ./cmd/docs-drift
