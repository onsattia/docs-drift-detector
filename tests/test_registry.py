from pathlib import Path

import pytest
from typer.testing import CliRunner

from docs_drift.cli import app
from docs_drift.registry import RegistryError, load_registry

REPO_ROOT = Path(__file__).resolve().parent.parent
runner = CliRunner()


def write_registry(tmp_path: Path, text: str) -> Path:
    path = tmp_path / "sources.yaml"
    path.write_text(text, encoding="utf-8")
    return path


def test_valid_registry_exposes_projects_and_code_sources(tmp_path: Path) -> None:
    path = write_registry(
        tmp_path,
        """
projects:
  billing-docs:
    docs: docs/billing
    code:
      - repo: acme/billing-api
        ref: main
""",
    )

    registry = load_registry(path)

    project = registry.projects["billing-docs"]
    assert project.docs == "docs/billing"
    assert [(c.repo, c.ref) for c in project.code] == [("acme/billing-api", "main")]


def test_project_without_code_source_is_rejected(tmp_path: Path) -> None:
    path = write_registry(
        tmp_path,
        """
projects:
  billing-docs:
    docs: docs/billing
    code: []
""",
    )

    with pytest.raises(RegistryError) as excinfo:
        load_registry(path)

    assert any("billing-docs" in e and "at least one code source" in e for e in excinfo.value.errors)


def test_all_errors_are_reported_in_one_run(tmp_path: Path) -> None:
    path = write_registry(
        tmp_path,
        """
projects:
  billing-docs:
    docs: docs/billing
    code:
      - repo: not-a-repo
        ref: main
  search-docs:
    docs: docs/search
    code:
      - repo: acme/search-api
""",
    )

    with pytest.raises(RegistryError) as excinfo:
        load_registry(path)

    errors = excinfo.value.errors
    assert len(errors) == 2
    assert any("billing-docs" in e and "code[0].repo" in e for e in errors)
    assert any("search-docs" in e and "code[0].ref" in e for e in errors)


@pytest.mark.parametrize("docs_path", ["../other-repo/docs", "/etc/docs"])
def test_docs_path_must_be_repository_relative(tmp_path: Path, docs_path: str) -> None:
    path = write_registry(
        tmp_path,
        f"""
projects:
  billing-docs:
    docs: {docs_path}
    code:
      - repo: acme/billing-api
        ref: main
""",
    )

    with pytest.raises(RegistryError) as excinfo:
        load_registry(path)

    assert any(
        "billing-docs" in e and "docs" in e and "must stay within the repository" in e for e in excinfo.value.errors
    )


def test_openapi_path_must_be_repository_relative(tmp_path: Path) -> None:
    path = write_registry(
        tmp_path,
        """
projects:
  billing-docs:
    docs: docs/billing
    openapi: ../billing-api/openapi.yaml
    code:
      - repo: acme/billing-api
        ref: main
""",
    )

    with pytest.raises(RegistryError) as excinfo:
        load_registry(path)

    assert any("openapi" in e and "must stay within the repository" in e for e in excinfo.value.errors)


@pytest.mark.parametrize("field", ["token", "password", "secret", "api_key"])
def test_credential_fields_are_rejected_without_echoing_value(tmp_path: Path, field: str) -> None:
    secret_value = "s3cr3t-value-123"
    path = write_registry(
        tmp_path,
        f"""
projects:
  billing-docs:
    docs: docs/billing
    code:
      - repo: acme/billing-api
        ref: main
        {field}: {secret_value}
""",
    )

    with pytest.raises(RegistryError) as excinfo:
        load_registry(path)

    message = "\n".join(excinfo.value.errors)
    assert field in message
    assert secret_value not in message


def test_cli_validate_prints_project_summary_and_exits_zero() -> None:
    result = runner.invoke(app, ["registry", "validate", str(REPO_ROOT / "examples" / "sources.yaml")])

    assert result.exit_code == 0
    assert "billing-docs" in result.output
    assert "acme/billing-api@main" in result.output


def test_cli_validate_exits_nonzero_and_reports_errors(tmp_path: Path) -> None:
    path = write_registry(
        tmp_path,
        """
projects:
  billing-docs:
    docs: docs/billing
    code: []
""",
    )

    result = runner.invoke(app, ["registry", "validate", str(path)])

    assert result.exit_code != 0
    assert "billing-docs" in result.output
    assert "at least one code source" in result.output
