## Purpose

Defines the quality checks every change must pass before it can be merged, so that code standards are enforced by
the pipeline rather than by review alone.

## Requirements

### Requirement: Pull requests run the full quality gate
The CI pipeline SHALL run formatting, lint, strict type checking, tests, dependency scanning, and secret scanning on
every pull request targeting the default branch. A pull request SHALL NOT be mergeable while any of these checks fail.

#### Scenario: All checks pass
- **WHEN** a pull request has no formatting, lint, type, test, dependency, or secret findings
- **THEN** every required check reports success and the pull request is mergeable

#### Scenario: A check fails
- **WHEN** a pull request introduces a type error
- **THEN** the type-check job fails, the pull request is blocked from merging, and the failure output identifies the file and line

### Requirement: Test coverage floor
The test suite SHALL enforce a minimum line coverage threshold, configured in the project and documented in the
README. Coverage below the threshold SHALL fail the test job.

#### Scenario: Coverage drops below the floor
- **WHEN** a change reduces total line coverage below the configured threshold
- **THEN** the test job fails and reports the measured coverage and the threshold

### Requirement: Secret scanning blocks committed credentials
The CI pipeline SHALL scan the full commit history of each pull request for credential patterns. A detected secret
SHALL fail the job.

#### Scenario: Credential committed
- **WHEN** a commit contains a value matching a known credential pattern
- **THEN** the secret-scan job fails and reports the file path and commit, without printing the secret value

### Requirement: Local hooks mirror CI
The repository SHALL provide pre-commit hooks that run formatting and lint checks on staged files. Hook configuration
SHALL use the same tool versions as CI.

#### Scenario: Developer commits unformatted code
- **WHEN** a developer commits a file that the formatter would change
- **THEN** the commit is rejected until the file is formatted

### Requirement: Dependency versions are locked
The project SHALL commit a lockfile. CI SHALL install dependencies from the lockfile and SHALL fail if the lockfile is
out of date with the project's dependency declarations.

#### Scenario: Lockfile out of date
- **WHEN** a dependency is added to the project declaration without updating the lockfile
- **THEN** the CI install step fails with a message instructing the developer to update the lockfile
