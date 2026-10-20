#!/bin/sh
# Blocks hand edits to go.sum. Use go get or go mod tidy instead (CLAUDE.md "Do not").
file=$(jq -r '.tool_input.file_path // ""')
case "$file" in
  */go.sum|go.sum)
    echo "go.sum is managed by go get and go mod tidy. Do not edit it by hand." >&2
    exit 2
    ;;
esac
exit 0
