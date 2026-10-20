---
name: search-first
description: Research-before-coding workflow for docs-drift-detector. Use before adding a Go module dependency, a new integration, or a helper that may already exist in the Go ecosystem, the repo, or OpenSpec specs. Decides whether to adopt, extend, or build.
---

Search for an existing solution before writing new code. Report what you searched and what you skipped.

## When to use

- Starting a change in `openspec/changes/` where a library might cover the behavior (YAML parsing, CLI parsing, HTTP clients, retries, hybrid search, LLM client wrappers)
- Adding or upgrading a Go module, or any new dependency (the project asks before adding one; this skill gives you the evidence for that request)
- About to write a utility, helper, or abstraction under `internal/`

Skip for pure edits to existing behavior that the repo already implements. Search the repo first, then stop.

## Workflow

0. **Preflight.** Check which channels you can use. Report skipped channels rather than claiming coverage you do not have.
1. **Need analysis.** State the behavior in one sentence. Note the constraints: Go 1.26.9 (pinned in `go.mod`), the approved model endpoints only, and the stack already in use (cobra, `gopkg.in/yaml.v3`).
2. **Search the repo.** `rg` for the function, type, or concept across `internal/`, `cmd/`, `scripts/`, and `openspec/specs/`. A spec may already define the behavior.
3. **Search the ecosystem.** Check the standard library first. Then search the Go module ecosystem (pkg.go.dev, `go list -m -versions`), and GitHub for maintained implementations.
4. **Check the project's own skills and MCP tools.** Look in `.claude/skills/` and the MCP servers configured for the session. Use docs MCP or official docs for library APIs.
5. **Evaluate candidates.** Score each on fit, maintenance (recent release, open issues), license (must be compatible with the repo), transitive dependencies, and whether it adds a network call or sends data outside the approved endpoints.
6. **Decide.**

| Signal | Action |
|--------|--------|
| Standard library or existing repo code covers it | **Reuse** — no new code or dependency |
| Maintained module fits exactly, compatible license, small dependency tree | **Adopt** — `go get`, then use directly |
| Module covers most of it | **Extend** — `go get` plus a thin wrapper in `internal/` |
| Nothing suitable | **Build** — write it under `internal/` with `tdd`, noting what was rejected and why |

7. **Record the decision.** For a non-trivial change, put the rejected options and reasons in the change's `design.md`. For a dependency, state the module, version, and license in the PR description.

## Rules for this repo

- **Ask before adding a dependency.** Present the search results and the recommendation, then wait for approval. Do not run `go get` on your own judgment.
- **Never edit `go.sum` by hand.** Use `go get` or `go mod tidy`, then `go mod verify`.
- **Ask before any new model provider or base URL.** A library that calls an LLM endpoint is a data-handling decision, not only a code decision. Internal docs and code only go to approved endpoints.
- **Watch for dependency bloat.** One small feature does not justify a large framework. Prefer packages with small dependency trees.
- **Check licenses.** Record the license of any candidate that is not MIT, Apache-2.0, or BSD.
- **Check the security scan.** A candidate must pass `govulncheck ./...` after it is added.

## Search shortcuts for this project

| Need | Start with |
|------|-----------|
| CLI commands and flags | cobra (already in use); check `internal/cli/` first |
| YAML parsing and validation | `gopkg.in/yaml.v3` (already in use); check `internal/registry/` first |
| Schema or field validation | Standard library or hand-written checks; prefer a small validator over a framework |
| Retries and timeouts for HTTP | `net/http` with a context deadline; a retry library only if the need is real |
| Markdown or doc parsing | `goldmark` or the standard library; check whether the claim extractor already covers it |
| Test helpers and fixtures | `testing`, `testdata/` directories, and existing helpers under `internal/` |
| Lint and vulnerability checks | Already set up: golangci-lint, govulncheck; do not add a second tool for the same job |

## Examples

### Example 1: "Parse sources.yaml with strict field checks"
```
Need: YAML parsing that rejects unknown fields and credential fields
Repo: internal/registry/ already parses sources.yaml
Action: REUSE — extend the existing parser's validation, no new dependency
```

### Example 2: "Add retries to the model client"
```
Need: Retry transient 429 and 5xx errors from the approved model endpoint
Search: standard library net/http, pkg.go.dev retry modules
Found: a small retry helper module (license MIT, one transitive dependency)
Action: ASK first — dependency approval needed; the stdlib version is about 30 lines
Result: Waiting for approval before running go get
```

## Anti-patterns

- **Jumping to code.** Writing a helper without checking the repo, the stdlib, and the specs.
- **Silent skipping.** Saying "nothing found" when a channel was unavailable or never searched.
- **Adding a dependency to avoid a small function.** Count the lines and the transitive dependencies before adding a module.
- **Ignoring the approved-endpoint rule.** A library that sends data to a new host is a blocker, not a detail.
