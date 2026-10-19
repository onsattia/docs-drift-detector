## Context

The repository is empty apart from the OpenSpec scaffolding. See proposal.md - Why. The tool will run as a scheduled
CI job on internal documentation, so the foundation must support reproducible installs, strict typing, automated
security checks, and a minimal container image from the first commit.

## Goals / Non-Goals

**Goals:**
- A package layout and toolchain that every later change uses without re-deciding.
- Quality gates enforced by CI and mirrored locally (see specs/quality-gates).
- A CLI that proves install, entry point, and container build work end to end.
- A registry validator that defines the `sources.yaml` contract before any code reads source repositories.

**Non-Goals:**
- Claim extraction, retrieval, model calls, or report generation.
- Deployment to a scheduler or cluster. That follows the platform decision in a later change.
- Authentication to any external repository or model endpoint.

## Decisions

**Python 3.12 with uv.** uv provides fast installs and a single lockfile (`uv.lock`) for reproducible CI. Alternative:
Poetry, which works but is slower and has a heavier lock format. Alternative: pip-tools, which needs more manual
wiring.

**src layout (`src/docs_drift/`).** Prevents tests from importing the working tree instead of the installed package,
which hides packaging bugs. Alternative: flat layout, rejected for that reason.

**Typer for the CLI.** Gives typed arguments and help generation with little code. Alternative: argparse, which is
built in but more verbose. The choice has no effect on behavior, only on maintenance cost.

**Pydantic v2 for the registry model.** Validates the YAML into typed objects and produces field-path error messages.
Using it here means later steps receive validated data without re-checking it. Alternative: hand-written validation
with `dict` access, rejected because error reporting would be inconsistent.

**ruff for lint and format; mypy in strict mode.** One fast tool replaces flake8, isort, and black. Strict mypy
catches the interface mistakes that matter most when a later LLM-judge step consumes structured output.

**pytest with pytest-cov.** Standard, and coverage is enforced through a configured floor rather than a convention.

**Container: multi-stage build on a slim Python base, non-root user.** The build stage installs locked dependencies;
the runtime stage copies only the virtual environment and source. Base image is pinned by digest in the Dockerfile.

**Secret scanning with gitleaks in CI and pre-commit.** It scans history, not only the working tree, which matches the
quality-gates requirement.

**Registry rejects credential fields outright.** Credentials belong in the platform's secret store. Rejecting the
fields at validation time prevents a common mistake from reaching the repository.

## Risks / Trade-offs

- **Pinning the base image by digest** → Digests go stale and must be updated. Mitigation: Dependabot updates the
  Dockerfile digest.
- **Strict mypy slows early development** → Accept the cost now; relaxing it later is harder than tightening it.
- **Pydantic adds a runtime dependency** → Small and widely used; it removes a class of validation bugs.
- **Coverage floor may block small changes** → Start the floor at a level the bootstrap code meets, then raise it as
  the detector is added.

## Migration Plan

Not applicable. The repository has no existing code or deployments to migrate. Rollback is reverting the merge commit.
