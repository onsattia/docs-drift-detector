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

## Stack

- Python 3.12, managed with uv (`uv.lock` is committed and authoritative)
- Typer for the CLI, Pydantic v2 for data models, PyYAML for registry parsing
- pytest + pytest-cov, ruff (lint and format), mypy in strict mode
- GitHub Actions for CI, a multi-stage Dockerfile for the runtime image

## Layout

```
src/docs_drift/      package code (src layout)
tests/               pytest suite; fixtures under tests/fixtures/
examples/            sample sources.yaml files
openspec/            specs, changes, and project config
.claude/skills/      project skills (tdd)
```

## Commands

Planned until the `bootstrap-project` change is applied. Update this list when the commands change.

```bash
uv sync --locked                      # install from lockfile
uv run pytest                         # tests with coverage
uv run ruff check . && uv run ruff format --check .
uv run mypy src tests
uv run docs-drift --version
uv run docs-drift registry validate examples/sources.yaml
```

## Code conventions

- Match the surrounding code's style, naming, and comment density. Keep comments to the non-obvious "why".
- Type every function signature. No `Any` without a comment explaining why it is needed.
- Public functions have docstrings only when behavior is not obvious from the name and types.
- Errors are specific: include the field path, the file, or the claim that failed. Never swallow exceptions silently.
- Commit messages use conventional commits (`feat:`, `fix:`, `test:`, `chore:`, `docs:`).

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

- Do not edit `uv.lock` by hand. Use `uv add` or `uv lock`.
- Do not commit `.env` files, API keys, or tokens, even for testing.
- Do not disable a failing CI check to merge. Fix the cause or ask for review.
- Do not change the Python version or the base image digest without an explicit request.
