# docs-drift-detector

Detects drift between documentation and the code it describes. It extracts checkable claims from docs, finds the
matching code, and labels each claim `consistent`, `drifted`, `missing`, or `unverifiable`. It runs as a scheduled
CI job and never edits documentation.

## Prerequisites

- [Go 1.26.9](https://go.dev/dl/) (the version in `go.mod`; Go downloads it automatically if your local toolchain is older)
- Docker, only to build the container image locally

## Setup

```bash
git clone <repository-url> docs-drift-detector
cd docs-drift-detector
go mod download
```

Install the pre-commit hooks with [pre-commit](https://pre-commit.com/), then the Go linter that the hooks call:

```bash
pre-commit install
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
```

## Build

```bash
make build
bin/docs-drift --version
```

`scripts/build.sh` reads the version from [VERSION](VERSION) and bakes it into the binary.

## Quality gates

Run `make check` locally before pushing. It runs the same gates CI runs on every pull request to `main`. Run `make help` to see every target.

```bash
test -z "$(gofmt -l .)"
golangci-lint run ./...
go mod tidy -diff
go mod verify
sh scripts/coverage.sh
govulncheck ./...
docker build -t docs-drift:local . && docker run --rm docs-drift:local --version
```

Install `govulncheck` with `go install golang.org/x/vuln/cmd/govulncheck@v1.1.4`.

## Test coverage

`sh scripts/coverage.sh` fails if total statement coverage drops below **90%**. Override the floor for a local run
with `COVERAGE_MIN=85 sh scripts/coverage.sh`; CI always uses the default.

## Source registry

`sources.yaml` maps documentation folders to the code repositories they describe. See
[examples/sources.yaml](examples/sources.yaml) for the format.

```bash
go run ./cmd/docs-drift registry validate path/to/sources.yaml
```

The command reports every problem in the file in one run and exits non-zero if any are found.
