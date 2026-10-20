#!/bin/sh
# Blocks git commit while HEAD is main. Create a branch first (CLAUDE.md Workflow).
command=$(jq -r '.tool_input.command // ""')
case "$command" in
  *"git commit"*)
    branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
    if [ "$branch" = "main" ]; then
      echo "Never commit directly to main. Create a branch and open a pull request." >&2
      exit 2
    fi
    ;;
esac
exit 0
