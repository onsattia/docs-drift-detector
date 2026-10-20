# docs-drift-detector

Detects drift between documentation and the code it describes. It extracts checkable claims from docs (function
names, env vars, endpoints, CLI flags, config keys), finds matching code with hybrid search, and asks an LLM judge
to label each claim `consistent`, `drifted`, `missing`, or `unverifiable`. It runs as a scheduled CI job, not a
web service. It never applies doc changes automatically.

## Working agreement

- Specs are the source of truth. Behavior changes start in `openspec/changes/<name>/specs/`, not in code.
- Use OpenSpec for any non-trivial change: propose, then apply. Do not implement a change before its artifacts exist.
- Implement behavior with the `tdd` skill (`.claude/skills/tdd/SKILL.md`). Write the failing test first.
- Keep changes small. One scenario or task per commit where practical.
- Ask before adding a dependency, changing the model endpoint, or touching CI permissions.

## Workflow

- Never commit directly to `main`. Create a branch per change and open a pull request.
- Lifecycle for a change: `openspec propose` → review the specs → implement tasks with `tdd` → verify → archive.
  Sync delta specs into `openspec/specs/` before archiving so the main specs stay current.
- Before every commit, run `make check`. Do not commit with any of its gates failing.
- Commit each finished task or scenario separately, with a conventional commit message. Keep the message about the
  why, not a list of files.
- Mark the task checkboxes in `tasks.md` as they are completed, so an interrupted session can resume from the change
  directory alone.
- Keep the Commands section above in sync when a command or flag changes.
- Keep specs compact and easy to read. Use short requirements with one scenario each, plain language, and no
  repeated background. If a spec needs a long explanation, move the rationale to `design.md`.
- Before opening a PR, rebase on the latest `main` and confirm CI is green locally. Fix the cause of any failing
  check rather than skipping it.

## Stack

- Go 1.26.9, pinned in `go.mod` (`go.sum` is committed and authoritative)
- cobra for the CLI, `gopkg.in/yaml.v3` for registry parsing
- `go test` with a coverage floor, gofmt, golangci-lint, govulncheck
- GitHub Actions for CI, a multi-stage Dockerfile that produces a static binary on a distroless image

## Layout

```
cmd/docs-drift/      binary entry point
internal/cli/        cobra command tree
internal/registry/   sources.yaml parsing and validation
scripts/             build.sh (version injection) and coverage.sh (floor check)
examples/            sample sources.yaml files
openspec/            specs, changes, and project config
.claude/skills/      project skills (tdd)
```

## Commands

Run `make help` for the list. Update the Makefile and this section together when a command changes.

```bash
make check                          # every quality gate CI runs
make build                          # bin/docs-drift with the version from VERSION
go mod download                       # fetch modules
sh scripts/build.sh bin/docs-drift    # build with the version from VERSION
sh scripts/coverage.sh                # tests with coverage floor (90%)
gofmt -l .                            # must print nothing
golangci-lint run ./...
go mod tidy -diff && go mod verify
govulncheck ./...
bin/docs-drift --version
bin/docs-drift registry validate examples/sources.yaml
```

## Code conventions

- Match the surrounding code's style and naming. Add a code comment only when it is very important, such as a
  non-obvious constraint or a workaround. Do not comment what the code already says.
- Type every function signature. No `any` without a comment explaining why it is needed.
- Public functions have docstrings only when behavior is not obvious from the name and types.
- Errors are specific: include the field path, the file, or the claim that failed. Never swallow exceptions silently.
- Commit messages use conventional commits (`feat:`, `fix:`, `test:`, `chore:`, `docs:`).

### Naming

- **Name a thing for what it holds or does, not its shape.** `errs` for an error collector, `messages` for the list
  inside it. Avoid one-letter names outside short loops and receivers.
- **Keep receiver names short and the same across one type.** Every method on `errorList` uses `errs`. Do not name a
  receiver after one of its fields.
- **Use the domain's words from the specs.** `docs`, `code`, `repo`, `ref`, `openapi`, `claim`, and `verdict` mean the
  same thing everywhere. Do not add synonyms.
- **Do not repeat the package name in exported identifiers.** Write `registry.Load`, not `registry.LoadRegistry`.
- **Name files for the concern they hold,** one concern per file. Two names that differ by a few letters or sound alike,
  such as `validate.go` and `validator.go`, make the code hard to navigate, so pick names that cannot be confused.
- **Use the Go conventions for common roles:** `err` for errors, `ok` or `found` for a second boolean result, `ctx` for a
  context, and `i` or `k` for loop indexes.
- **Name tests after the behavior or scenario,** such as `TestRegistryRejectsPathEscapingRepo`, so the mapping to specs
  is visible.
- **Prefer full words to abbreviations,** except the common Go roles above.

## Security and data handling

- **No secrets in the repo, logs, fixtures, or images.** Tokens come from the CI secret store or the runtime
  environment only. The `sources.yaml` schema rejects credential fields by design; keep it that way.
- **Read-only access to code and docs repos.** The tool reads source; it never writes to repos except to open an
  issue or a draft proposal, and only when that feature is specified.
- **Redact logs.** Log claim IDs, file paths, verdicts, and counts. Do not log full document text, code snippets from
  internal repos, or model prompts and responses in production runs.
- **Approved endpoints only.** Internal docs and code go only to model endpoints approved for the organization's data.
  Do not add a new provider or base URL without asking.
- **Fixtures are synthetic.** Never commit real internal documentation or code as test data.
- Data is internal documentation and falls under GDPR review. Treat any personal data found in docs as sensitive.

## Testing

- Unit tests for parsers, validators, and verdict logic run without network or model access.
- Model-dependent checks belong in a separate eval suite that runs in CI against the golden set. Do not put live
  model calls in the unit tests.
- Mock model calls at the client boundary with recorded, synthetic responses.

## Do not

- Do not edit `go.sum` by hand. Use `go get` or `go mod tidy`.
- Do not commit `.env` files, API keys, or tokens, even for testing.
- Do not disable a failing CI check to merge. Fix the cause or ask for review.
- Do not change the Go version or the base image digests without an explicit request.
