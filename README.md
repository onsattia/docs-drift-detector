# docs-drift-detector

Detects drift between documentation and the code it describes. It extracts checkable claims from docs, finds the
matching code, and labels each claim `consistent`, `drifted`, `missing`, or `unverifiable`. It runs as a scheduled
CI job and never edits documentation.

## Prerequisites

- [uv](https://docs.astral.sh/uv/) (installs Python 3.12 for you)
- Docker, only to build the container image locally

## Setup

```bash
git clone <repository-url> docs-drift-detector
cd docs-drift-detector
uv sync --locked
uv run pre-commit install
```

## Quality gates

Run these locally before pushing. CI runs the same checks on every pull request to `main`.

```bash
uv run ruff check .
uv run ruff format --check .
uv run mypy src tests
uv run pytest
uv run docs-drift --version
uv run docs-drift registry validate examples/sources.yaml
docker build -t docs-drift:local . && docker run --rm docs-drift:local --version
```

## Test coverage

`uv run pytest` fails if line coverage drops below **90%**. The floor is set by `--cov-fail-under` in
`pyproject.toml`.

## Type checking

mypy runs in strict mode. PyYAML ships no type stubs, so `registry.py` carries a targeted
`# type: ignore[import-untyped]`. Adding `types-PyYAML` as a dev dependency would remove it.

## Source registry

`sources.yaml` maps documentation folders to the code repositories they describe. See
[examples/sources.yaml](examples/sources.yaml) for the format.

```bash
uv run docs-drift registry validate path/to/sources.yaml
```
