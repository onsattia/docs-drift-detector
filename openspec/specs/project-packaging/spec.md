## Purpose

Defines how the tool is installed, versioned, and run, so that the same build works on a developer machine, in CI,
and in the production container.

## Requirements

### Requirement: Command-line entry point
The tool SHALL provide a `docs-drift` command. Running `docs-drift --version` SHALL print the installed version and exit
with status 0. Running `docs-drift --help` SHALL list available commands and exit with status 0.

#### Scenario: Version query
- **WHEN** a user runs `docs-drift --version`
- **THEN** the installed version string is printed and the process exits with status 0

#### Scenario: Unknown command
- **WHEN** a user runs `docs-drift` with an unrecognized subcommand
- **THEN** a usage message is printed to standard error and the process exits with a non-zero status

### Requirement: Reproducible install from a clean checkout
A fresh clone of the repository SHALL install and run the test suite using only the documented setup steps, with no
manual edits and no network access beyond the package index.

#### Scenario: Clean setup
- **WHEN** a developer clones the repository and runs the documented setup commands
- **THEN** the dependencies install from the lockfile and the test suite runs to completion

### Requirement: Container image builds from source
The repository SHALL provide a container build definition that produces a runnable image containing only the tool and
its runtime dependencies. The image SHALL run as a non-root user.

#### Scenario: Image build
- **WHEN** the container build runs on a clean checkout
- **THEN** an image is produced, and running `docs-drift --version` inside it prints the version

#### Scenario: Non-root execution
- **WHEN** the container starts with the default user
- **THEN** the process runs as a non-root user ID

### Requirement: Versioning
The package version SHALL follow semantic versioning and SHALL be the single source for the version reported by the
CLI and recorded in the container image metadata.

#### Scenario: Version consistency
- **WHEN** the version is read from the package metadata and from the CLI
- **THEN** both values are identical
