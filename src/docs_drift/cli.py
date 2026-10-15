from pathlib import Path
from typing import Annotated

import typer

from docs_drift import __version__
from docs_drift.registry import RegistryError, load_registry

app = typer.Typer(help="Detect drift between documentation and code.", no_args_is_help=True)
registry_app = typer.Typer(help="Work with the sources.yaml registry.", no_args_is_help=True)
app.add_typer(registry_app, name="registry")


def _version_callback(value: bool) -> None:
    if value:
        typer.echo(__version__)
        raise typer.Exit()


@app.callback()
def main(
    version: Annotated[
        bool,
        typer.Option("--version", callback=_version_callback, is_eager=True, help="Show the version and exit."),
    ] = False,
) -> None:
    """Detect drift between documentation and code."""


@registry_app.command("validate")
def validate(path: Annotated[Path, typer.Argument(exists=True, dir_okay=False, readable=True)]) -> None:
    """Validate a sources.yaml registry and print a summary."""
    try:
        registry = load_registry(path)
    except RegistryError as exc:
        for error in exc.errors:
            typer.echo(f"error: {error}", err=True)
        raise typer.Exit(code=1) from exc

    typer.echo(f"Registry is valid: {len(registry.projects)} project(s)")
    for name, project in registry.projects.items():
        typer.echo(f"  {name}")
        for source in project.code:
            typer.echo(f"    - {source.repo}@{source.ref}")
