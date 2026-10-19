## Purpose

Defines the `sources.yaml` registry that maps documentation folders to the code repositories they describe. The
registry is the single place where the tool learns which docs belong to which code.

## Requirements

### Requirement: Registry file format
The registry SHALL be a YAML file with a top-level `projects` mapping. Each project entry SHALL declare a `docs` path
and a non-empty list of `code` sources. Each code source SHALL declare a `repo` in `owner/name` form and a `ref`.
An optional `openapi` path MAY be declared per project.

#### Scenario: Valid registry
- **WHEN** a registry declares one project with a docs path and one code source with `repo` and `ref`
- **THEN** validation succeeds and the project is available to later steps

#### Scenario: Project without code source
- **WHEN** a project entry declares no `code` sources
- **THEN** validation fails with an error naming the project and stating that at least one code source is required

### Requirement: Registry validation errors are specific
Validation SHALL report every error in the file in a single run, each with the project name and field path. It SHALL NOT
stop at the first error.

#### Scenario: Multiple errors
- **WHEN** a registry has an invalid `repo` format in one project and a missing `ref` in another
- **THEN** both errors are reported in one output and the process exits with a non-zero status

### Requirement: Registry paths are repository-relative
Documentation and OpenAPI paths SHALL be relative to the root of the docs repository. Absolute paths and paths that
escape the repository root through `..` segments SHALL be rejected.

#### Scenario: Path escapes repository
- **WHEN** a project declares a docs path of `../other-repo/docs`
- **THEN** validation fails with an error stating that paths must stay within the repository

### Requirement: Credentials are never stored in the registry
The registry SHALL NOT accept fields for tokens, passwords, or keys. A registry containing any such field SHALL fail
validation.

#### Scenario: Token field present
- **WHEN** a code source declares a `token` field
- **THEN** validation fails and the error message does not echo the value of the field

### Requirement: Validation command
The CLI SHALL provide a command that validates a registry file and prints a summary of projects and code sources on
success. Exit status SHALL be 0 on success and non-zero on any validation error.

#### Scenario: Successful validation summary
- **WHEN** a valid registry is validated
- **THEN** the output lists each project with its code repositories and the exit status is 0
