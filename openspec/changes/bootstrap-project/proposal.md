## Why

The docs drift detector has no project structure yet. Before any detection logic is written, it needs a
reproducible Python package, quality tooling, and a CI pipeline, so every later change is built on the
same standards and can be deployed as a scheduled job without rework.

## What Changes

- Create the `docs_drift` Python package (src layout) managed with uv and a committed lockfile.
- Configure ruff (lint and format) and mypy in strict mode, plus pytest with a coverage floor.
- Add a CI workflow that runs lint, type check, tests, and dependency and secret scanning on every pull request.
- Add a placeholder CLI entry point (`docs-drift --version`) to prove the packaging and install path work end to end.
- Add a container build definition that builds the image in CI without publishing it.
- Add a `sources.yaml` schema placeholder and a loader that validates the registry format. No code is read or
  scanned yet.
- Add repository hygiene: `.gitignore`, `.editorconfig`, pre-commit hooks, CODEOWNERS, and a README with setup steps.
- **Adopt test-driven development for all code in this change.** Each behavior is written as a failing test first,
  seen to fail for the expected reason, then made to pass with the minimum code, then refactored. Tests map one-to-one
  to the WHEN/THEN scenarios in the specs. The `tdd` project skill (`.claude/skills/tdd/SKILL.md`) defines the
  workflow, and the test suite is the first verification step in the tasks.
- Add a `CLAUDE.md` at the repository root that records the project's stack, layout, commands, code conventions,
  security rules, and testing approach, so every later session starts from the same standards.

No detection logic, model calls, or network access are included in this change.

## Capabilities

### New Capabilities

- `quality-gates`: Rules that every change must pass before merge: lint, format, strict type check, tests with a
  coverage floor, dependency scanning, and secret scanning. Covers the CI workflow and local pre-commit hooks.
- `project-packaging`: Requirements for installing and running the tool: a locked, reproducible dependency set, a
  CLI entry point, and a container image that builds from a clean checkout.
- `source-registry`: Format and validation rules for `sources.yaml`, which maps documentation folders to the code
  repositories they describe.

### Modified Capabilities

None. This is the first change in the project.

## Impact

- **Code:** new `src/docs_drift/`, `tests/`, and root configuration files (`pyproject.toml`, `uv.lock`, `Dockerfile`,
  `.pre-commit-config.yaml`).
- **Process:** TDD is required for every behavior-changing task. Test files are committed in the same change as the
  code they cover, and a test added without a prior failing run is a review finding.
- **Documentation for agents:** `CLAUDE.md` and `.claude/skills/tdd/` are new. They apply to all future changes.
- **CI:** new GitHub Actions workflow. No deploy step and no scheduled job yet.
- **Dependencies:** runtime is minimal (a YAML parser, a validation library, and a CLI library). Dev tools are
  ruff, mypy, pytest, pytest-cov, and pre-commit.
- **Systems:** none. No secrets, tokens, or model endpoints are configured in this change.
- **Assumption:** Python 3.12 is the target. This was chosen as a default and can be changed before apply.
