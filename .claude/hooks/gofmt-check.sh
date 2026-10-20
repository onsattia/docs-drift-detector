#!/bin/sh
# Reports Go files that gofmt would change, right after an edit.
file=$(jq -r '.tool_input.file_path // ""')
case "$file" in
  *.go)
    [ -f "$file" ] || exit 0
    if [ -n "$(gofmt -l "$file")" ]; then
      echo "gofmt would change $file. Run: gofmt -w $file" >&2
      exit 2
    fi
    ;;
esac
exit 0
