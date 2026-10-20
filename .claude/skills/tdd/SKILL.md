---
name: tdd
description: Test-driven development workflow for docs-drift-detector. Use when implementing or changing behavior in internal/ or cmd/, including parsers, the registry, CLI commands, and verdict logic. Enforces red-green-refactor with go test.
---

Implement behavior test-first. Write a failing test, watch it fail for the right reason, write the smallest code that passes, then refactor with the tests green.

## When to use

- Any new function, type, or CLI command in `internal/` or `cmd/`
- Any change to validation rules, parsing, or verdict logic
- Bug fixes: first write a test that reproduces the bug

Skip for pure configuration files (`.golangci.yml`, workflow YAML, Dockerfile, `.pre-commit-config.yaml`). Verify those by running the tool they configure.

## Cycle

1. **Pick one behavior** from the change's `specs/` scenarios or `tasks.md`. Use one WHEN/THEN scenario as the unit of work.
2. **Red:** Write one test that expresses the scenario. Use a descriptive name such as `TestRegistryRejectsPathEscapingRepoRoot`.
3. **Run it and confirm the failure is the expected one.** Run `go test ./internal/registry/ -run TestName -v`. A build error or a typo is not a valid red; fix it until the assertion itself fails.
4. **Green:** Write the minimum code to pass. Hardcoding a return value is acceptable at this step if the next test forces generalization.
5. **Refactor:** Clean up names and duplication. Run the full suite after every refactor: `go test ./...`.
6. Repeat with the next scenario.

## Rules

- **One failing test at a time.** Do not write several tests before making the first one pass.
- **Test behavior, not implementation.** Assert on outputs, returned errors, and exit codes. Do not assert on unexported helper calls.
- **Mock only at boundaries:** model API calls, network and git access, and the filesystem where practical. Use real YAML decoding and real structs in tests. Judge tests use a client interface with a fake implementation, never a live model.
- **Use fixtures for input data.** Store sample registries and doc trees under `testdata/` or `tests/fixtures/`. Keep them small and synthetic. Never copy real internal documentation into fixtures.
- **Table-driven tests** for a family of inputs that share one behavior. Name each case after the scenario.
- **Every requirement scenario maps to at least one test.** Name the test after the scenario so the mapping is visible.
- **Bug fix order:** reproduce with a failing test, fix, then keep the test as a regression guard.
- **Coverage is a floor, not a goal.** Do not write tests that only raise the coverage number.

## Output expectations

When reporting back after a TDD session, state:
- Which scenarios now have tests
- The red-to-green result for each
- The final result of `go test ./...` and `scripts/coverage.sh` with real output, not a summary of intent
