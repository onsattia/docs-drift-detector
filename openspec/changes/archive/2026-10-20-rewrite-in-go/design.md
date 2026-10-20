## Context

The current implementation is a Python 3.12 package under `src/docs_drift/` with tests in `tests/`, built with uv, checked with pytest, ruff, and mypy strict, and packaged by a multi-stage Dockerfile. See proposal.md for why we are changing it.

The behavior that must survive the rewrite is defined in `openspec/specs/project-packaging`, `openspec/specs/quality-gates`, and `openspec/specs/source-registry`. Those specs are unchanged by this change, so their scenarios are the acceptance checks for the Go version.

## Goals / Non-Goals

**Goals:**
- The `docs-drift` binary passes every existing spec scenario, including `--version`, `--help`, unknown-command errors, and `registry validate`.
- The same `sources.yaml` files parse and validate with the same error messages, including the field path in each error.
- The container image is a single static binary on a minimal runtime image, running as a non-root user.
- CI runs the same gate as before: formatting, lint, vet, tests with a coverage floor, dependency scanning, and secret scanning.

**Non-Goals:**
- No new features. Endpoint extraction, search, and the judge are built in their own epics (#6 to #9), in Go, after this change lands.
- No change to the model endpoint or to CI permissions. Both need separate approval.
- No change to the output format of verdicts or logs, so the redaction rules in CLAUDE.md still apply unchanged.

## Decisions

**Go version: 1.26.9, not 1.27.2.** Go 1.27.2 is the latest stable release, but golangci-lint v2.12.2 cannot read Go 1.27's standard-library export data, so its typecheck fails on every file. Pin the latest 1.26 patch release until a golangci-lint release supports 1.27, then move `go.mod` up.

**Module layout.** Use `cmd/docs-drift/` for the entry point and `internal/` for packages (`internal/registry`, `internal/cli`). Use `internal/` so nothing is exposed as a public API by accident. Alternative: a flat `pkg/` layout. Rejected because the project is a CLI, not a library.

**CLI framework: cobra.** It matches the Typer command structure (subcommands, `--version`, generated `--help`) and is the standard in Go CLIs. Alternative: the standard `flag` package. Rejected because subcommands and help text would need hand-written code, and the spec requires a usage message on unknown commands.

**Registry parsing: `gopkg.in/yaml.v3` into typed structs, then explicit validation.** Pydantic v2 gave field-level errors. In Go, validation is explicit: each rule returns an error that names the field path, file, and claim, as CLAUDE.md requires. The credential-field rejection is a list of forbidden keys checked after decoding with `KnownFields(true)`, so unknown keys fail the same way the Pydantic schema did.

**Version: set at build time with `-ldflags`.** The version lives in one place (a `VERSION` file read by the build) and is injected into the binary and the image label. This meets the "single source for the version" requirement in `project-packaging`.

**Tooling: `go test -coverprofile` with a floor check, `gofmt`, `go vet`, golangci-lint, and `govulncheck` for dependency scanning.** This covers what ruff, mypy, and the dependency scanner did before. Alternative: keep a Python-only check for lint. Rejected because it would keep Python in the toolchain.

**Lockfile: `go.sum`, committed, and CI runs `go mod verify` plus `go mod tidy -diff`.** The second check fails when `go.mod` and `go.sum` drift from the source, which is the behavior the "Dependency versions are locked" requirement asks for.

**Container: multi-stage build, Go builder, then `gcr.io/distroless/static` or `scratch` for runtime.** The binary is built with `CGO_ENABLED=0`. The runtime runs as a non-root UID. Base images are pinned by digest, as the current Dockerfile does. Image digests are pinned, so the change does not alter the base-image policy in CLAUDE.md.

**Test strategy.** Port the Python tests by behavior, not line by line. Unit tests use table-driven `go test` with no network. The judge, once added, uses a mocked client interface at the boundary, as CLAUDE.md requires. Fixtures under `tests/fixtures/` (or `testdata/` by Go convention) are kept if they are format-neutral.

## Risks / Trade-offs

- **Behavioral drift in validation messages.** Go's YAML decoding and error text differ from PyYAML and Pydantic. → Mitigation: keep the existing spec scenarios as tests before porting, so each one fails until the Go version matches.
- **Loss of Python-only checks.** Mypy strict catches some issues that `go vet` does not. → Mitigation: enable golangci-lint with the stricter linters (for example `errcheck`, `staticcheck`, `nilnil`), and document the choice in CLAUDE.md.
- **Two implementations at once.** Keeping Python around while the Go version grows could confuse contributors. → Mitigation: the Python tree is removed in one commit at the end, after the Go version passes all scenarios. Until then, `CLAUDE.md` states which one is authoritative.
- **Adding Go dependencies needs approval** under CLAUDE.md. → Mitigation: list the dependencies (cobra, yaml.v3, golangci-lint, govulncheck) in the first task and get approval before adding them.

## Migration Plan

1. Land the spec-acceptance tests for the Python version first, so the behavior is pinned down in tests.
2. Build the Go version beside the Python one in the same repo, under `cmd/` and `internal/`.
3. Switch CI and the Dockerfile to Go once the Go version passes every scenario. Keep the Python jobs until then.
4. Delete the Python package, its tests, `pyproject.toml`, and `uv.lock` in one commit.
5. Rollback: revert the switch commit (step 3). The Python jobs and Dockerfile remain in history until step 4.

## Open Questions

None that change the specs or task breakdown. The judge endpoint (#8) and deploy design (#5) are tracked separately and do not depend on the language.
