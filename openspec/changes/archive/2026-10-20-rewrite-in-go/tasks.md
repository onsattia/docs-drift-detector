## 1. Approval and baseline

- [x] 1.1 Get approval for the Go dependencies listed in design.md (cobra, gopkg.in/yaml.v3, golangci-lint, govulncheck)
- [x] 1.2 Turn each scenario in `openspec/specs/project-packaging` and `openspec/specs/quality-gates` into a checklist in the PR description, so the Go version is checked against the same scenarios
- [x] 1.3 Confirm the Python suite passes on `main` to establish the baseline

## 2. Go module skeleton

- [x] 2.1 Create `go.mod` pinned to a current stable Go release (confirm the version before writing it), and commit `go.sum` (pinned to Go 1.26.9; see design.md)
- [x] 2.2 Add `cmd/docs-drift/main.go` that wires cobra and prints usage for `--help`
- [x] 2.3 Add `--version` that reads the version injected with `-ldflags` and exits 0 (the `-ldflags` injection itself is task 4.1)
- [x] 2.4 Verify unknown subcommands print usage to stderr and exit non-zero

## 3. Registry parsing and validation

- [x] 3.1 Write failing tests for the registry file format, using the existing `examples/sources.yaml` and `tests/fixtures/`
- [x] 3.2 Implement typed registry structs in `internal/registry`. Unknown keys are checked by hand rather than with `KnownFields(true)`, so each error names its field path
- [x] 3.3 Write failing tests for each validation error, asserting the message includes the field path
- [x] 3.4 Implement validation rules to pass 3.3
- [x] 3.5 Write failing tests that credential fields are rejected by name
- [x] 3.6 Implement the credential-field rejection to pass 3.5
- [x] 3.7 Write failing tests that registry paths are repository-relative, and implement the check
- [x] 3.8 Add `docs-drift registry validate <file>` and test it against `examples/sources.yaml`

## 4. Packaging and version

- [x] 4.1 Add a `VERSION` file and a build script that injects it with `-ldflags`
- [x] 4.2 Add a test that the version in the binary matches the `VERSION` file
- [x] 4.3 Verify a clean clone builds with the documented commands and no network beyond the module proxy

## 5. Container image

- [x] 5.1 Add a multi-stage Dockerfile: Go builder with `CGO_ENABLED=0`, then a minimal runtime image, with base images pinned by digest
- [x] 5.2 Run the image as a non-root UID
- [x] 5.3 Verify `docker run <image> --version` prints the version and exits 0

## 6. Quality gates

- [x] 6.1 Add test coverage reporting with a floor check that fails below the threshold, and document the threshold in the README
- [x] 6.2 Add golangci-lint configuration with the strict linters chosen in design.md
- [x] 6.3 Add `go mod tidy -diff` and `go mod verify` to CI so lockfile drift fails the build
- [x] 6.4 Add `govulncheck` to CI
- [x] 6.5 Update `.pre-commit-config.yaml` to run `gofmt` and golangci-lint with the same versions as CI

## 7. Switch CI and tooling

- [x] 7.1 Update `.github/workflows/ci.yml` to run the Go gate, keeping the secret-scan job
- [x] 7.2 Update `CLAUDE.md` (stack, layout, commands) and the `tdd` skill (`go test` instead of pytest)
- [x] 7.3 Update the README setup and test instructions
- [x] 7.4 Add a Makefile so every quality gate and build step has one named target that CI and contributors share

## 8. Remove the Python implementation

- [x] 8.1 Confirm every Python-tested behavior has a passing Go test
- [x] 8.2 Delete `src/docs_drift/`, the Python tests, `pyproject.toml`, `uv.lock`, and the Python-only CI jobs in one commit
- [x] 8.3 Verify the Python-era spec scenarios still pass on the Go build

## 9. Archive

- [x] 9.1 Run `openspec validate rewrite-in-go` and fix any issues
- [x] 9.2 Archive the change with `openspec-archive-change` once merged
