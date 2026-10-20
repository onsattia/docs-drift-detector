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

## Testing

Run the automated checks first. They cover the same gates CI runs:

```bash
make check
```

Then check the behavior by hand. Run these from the project root.

Build the binary:

```bash
make build
```

Check the version and help:

```bash
bin/docs-drift --version
bin/docs-drift --help
```

Validate the example registry. You should see a summary and exit code 0:

```bash
bin/docs-drift registry validate examples/sources.yaml
echo $?
```

Try a bad registry. Save this as `/tmp/bad.yaml`:

```yaml
projects:
  billing-docs:
    docs: ../escape
    token: abc123
    code:
      - repo: not-a-repo
```

Then run:

```bash
bin/docs-drift registry validate /tmp/bad.yaml
echo $?
```

You should see four errors, one per line: the `token` field, the `..` in the docs path, the `repo` format, and the missing `ref`. The exit code should be 1. The token value `abc123` should not appear anywhere in the output.

Check the unknown-command behavior. Usage goes to stderr and the exit code is 1:

```bash
bin/docs-drift nope
echo $?
```

Test the container the same way:

```bash
make image
docker run --rm docs-drift:local --version
docker run --rm -v "$PWD/examples:/data:ro" docs-drift:local registry validate /data/sources.yaml
```

Run `make help` to list the other targets.

## Source registry

`sources.yaml` maps documentation folders to the code repositories they describe. See
[examples/sources.yaml](examples/sources.yaml) for the format.

```bash
go run ./cmd/docs-drift registry validate path/to/sources.yaml
```

The command reports every problem in the file in one run and exits non-zero if any are found.
