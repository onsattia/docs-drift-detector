from typing import Annotated

import typer

from docs_drift import __version__

app = typer.Typer(help="Detect drift between documentation and code.", no_args_is_help=True)


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
