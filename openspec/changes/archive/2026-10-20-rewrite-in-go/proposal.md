## Why

The tool is a scheduled batch job that ships as a container image. Python needs a runtime, a virtualenv, and a multi-stage build to package it, and the team wants one static binary that is simpler to build, scan, and deploy to AWS (see epic #5). Rewriting now, before the endpoint extraction and judge work (#6 to #9) grows, keeps the cost of the change low.

## What Changes

- Reimplement the `docs-drift` CLI and its packages in Go, with the same user-facing behavior and the same `sources.yaml` format.
- Replace the Python toolchain with Go tooling: Go modules with `go.sum` as the lockfile, `go test` with coverage, `gofmt` and `go vet` with golangci-lint for lint and types.
- Replace the Python container build with a Go build stage and a minimal static runtime image.
- Update `CLAUDE.md`, the `tdd` skill, `ci.yml`, and the Dockerfile to match.
- Remove the Python package, its tests, and `uv.lock` in the final commit of this change, after the Go version passes the same checks.
- **BREAKING** for contributors: the `uv`, `pytest`, `ruff`, and `mypy` commands in `CLAUDE.md` are replaced by Go equivalents. The CLI name and flags stay the same.

## Capabilities

### New Capabilities

None. The rewrite keeps existing behavior.

### Modified Capabilities

None. The existing specs (`project-packaging`, `quality-gates`, `source-registry`) describe behavior in tool-neutral terms, such as "lockfile", "formatting and lint", and "container image". The rewrite keeps every requirement true, so the scenarios in those specs become acceptance checks for the Go version instead of changing. The change is marked `skip_specs: true`.

## Impact

- **Code:** all of `src/docs_drift/` and `tests/` are replaced. Fixtures under `tests/fixtures/` carry over if they are format-neutral.
- **Dependencies:** adds Go modules for CLI parsing (cobra), YAML parsing (`gopkg.in/yaml.v3`), and test tooling. Adding dependencies needs approval under CLAUDE.md.
- **CI and images:** `ci.yml`, the Dockerfile, and the deploy epic (#5) change. The OIDC role and deploy workflow from #5 are unaffected by the language choice.
- **Docs:** `CLAUDE.md` (stack, layout, commands) and the `tdd` skill must be updated in the same change.
- **Out of scope:** changing the model endpoint or any CI permissions. Both still need separate approval.
