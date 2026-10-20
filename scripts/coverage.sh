#!/bin/sh
# Runs the Go test suite with coverage and fails when total coverage is below the floor.
# The floor defaults to 90 and can be overridden with COVERAGE_MIN.
set -eu

min="${COVERAGE_MIN:-90}"
root="$(cd "$(dirname "$0")/.." && pwd)"
profile="$(mktemp)"
trap 'rm -f "$profile"' EXIT

cd "$root"
go test -coverprofile="$profile" ./...
total="$(go tool cover -func="$profile" | awk '/^total:/ { sub("%", "", $3); print $3 }')"

if ! awk -v t="$total" -v m="$min" 'BEGIN { exit !(t + 0 >= m + 0) }'; then
	echo "total coverage ${total}% is below the ${min}% floor" >&2
	exit 1
fi
echo "total coverage ${total}% (floor ${min}%)"
