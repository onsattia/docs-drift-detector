from importlib.metadata import version

from typer.testing import CliRunner

from docs_drift.cli import app

runner = CliRunner()


def test_version_flag_prints_package_version_and_exits_zero() -> None:
    result = runner.invoke(app, ["--version"])

    assert result.exit_code == 0
    assert result.output.strip() == version("docs-drift-detector")


def test_unknown_subcommand_exits_nonzero_with_usage_on_stderr() -> None:
    result = runner.invoke(app, ["no-such-command"])

    assert result.exit_code != 0
    assert "Usage" in result.stderr
