## 1. Repository foundation

- [ ] 1.1 Add `.gitignore` for Python, virtual environments, coverage output, and editor files
- [ ] 1.2 Add `.editorconfig` (UTF-8, LF line endings, 4-space indent for Python, 2-space for YAML and JSON)
- [ ] 1.3 Add `CODEOWNERS` with a placeholder owner line that must be replaced before merge
- [ ] 1.4 Add `README.md` with purpose, prerequisites, setup steps, and the quality-gate commands

## 2. Package and tooling

- [ ] 2.1 Create `pyproject.toml` with project metadata, Python 3.12 requirement, and the `docs-drift` script entry point
- [ ] 2.2 Add runtime dependencies (typer, pydantic, PyYAML) and dev dependencies (ruff, mypy, pytest, pytest-cov, pre-commit) via uv
- [ ] 2.3 Generate and commit `uv.lock`
- [ ] 2.4 Configure ruff lint and format rules in `pyproject.toml`
- [ ] 2.5 Configure mypy with strict mode and the `src/` path
- [ ] 2.6 Configure pytest and pytest-cov with a coverage floor and document the value in the README
- [ ] 2.7 Create `src/docs_drift/__init__.py` exposing the version read from package metadata

## 3. CLI

- [ ] 3.1 Implement the Typer app in `src/docs_drift/cli.py` with `--version` and `--help`
- [ ] 3.2 Add a test that `--version` prints the package version and exits 0
- [ ] 3.3 Add a test that an unknown subcommand exits non-zero with a usage message on stderr

## 4. Source registry

- [ ] 4.1 Define Pydantic models for projects, code sources, and optional OpenAPI paths in `src/docs_drift/registry.py`
- [ ] 4.2 Implement a loader that parses YAML and collects all validation errors with project name and field path
- [ ] 4.3 Enforce repository-relative paths and reject paths containing `..` or starting with `/`
- [ ] 4.4 Reject `token`, `password`, `secret`, and `key` fields without echoing their values
- [ ] 4.5 Add `docs-drift registry validate <path>` with a summary on success and non-zero exit on error
- [ ] 4.6 Add an example `sources.yaml` under `examples/` and a test that it validates
- [ ] 4.7 Add tests for each validation rule, including the multiple-errors case

## 5. Quality gates

- [ ] 5.1 Add `.pre-commit-config.yaml` with ruff lint and format hooks using the same versions as `uv.lock`
- [ ] 5.2 Create `.github/workflows/ci.yml` triggered on pull requests to the default branch
- [ ] 5.3 In CI, install from the lockfile with `uv sync --locked` so an out-of-date lockfile fails the job
- [ ] 5.4 In CI, add separate jobs for ruff check, ruff format check, mypy, and pytest with coverage
- [ ] 5.5 In CI, add a dependency audit job using pip-audit against the locked dependencies
- [ ] 5.6 In CI, add a gitleaks job with full-history scanning on pull requests
- [ ] 5.7 Confirm each job fails on a deliberately broken commit, then revert the commit

## 6. Container

- [ ] 6.1 Write a multi-stage `Dockerfile` using a slim Python 3.12 base pinned by digest
- [ ] 6.2 Install locked dependencies in the build stage and copy only the environment and source into the runtime stage
- [ ] 6.3 Create and switch to a non-root user in the runtime stage
- [ ] 6.4 Add a CI job that builds the image without pushing and runs `docs-drift --version` inside it
- [ ] 6.5 Add a `.dockerignore` that excludes `.git`, virtual environments, caches, and `openspec/`

## 7. Verification

- [ ] 7.1 Run the full setup steps from the README on a fresh clone and record any corrections
- [ ] 7.2 Run ruff, mypy, pytest, and the container build locally and confirm all pass
- [ ] 7.3 Run `openspec validate bootstrap-project --strict` and resolve any findings
