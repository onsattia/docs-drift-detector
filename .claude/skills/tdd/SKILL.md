---
name: tdd
description: Test-driven development workflow for docs-drift-detector. Use when implementing or changing behavior in src/docs_drift, including parsers, the registry, CLI commands, and verdict logic. Enforces red-green-refactor with pytest.
---

Implement behavior test-first. Write a failing test, watch it fail for the right reason, write the smallest code that passes, then refactor with the tests green.

## When to use

- Any new function, class, or CLI command in `src/docs_drift/`
- Any change to validation rules, parsing, or verdict logic
- Bug fixes: first write a test that reproduces the bug

Skip for pure configuration files (`pyproject.toml`, workflow YAML, Dockerfile). Verify those by running the tool they configure.

## Cycle

1. **Pick one behavior** from the change's `specs/` scenarios or `tasks.md`. Use one WHEN/THEN scenario as the unit of work.
2. **Red:** Write one test that expresses the scenario. Use a descriptive name such as `test_registry_rejects_path_escaping_repo_root`.
3. **Run it and confirm the failure is the expected one.** Run `uv run pytest path/to/test_file.py::test_name -q`. An import error or a typo is not a valid red; fix it until the assertion itself fails.
4. **Green:** Write the minimum code to pass. Hardcoding a return value is acceptable at this step if the next test forces generalization.
5. **Refactor:** Clean up names and duplication. Run the full suite after every refactor: `uv run pytest`.
6. Repeat with the next scenario.

## Rules

- **One failing test at a time.** Do not write several tests before making the first one pass.
- **Test behavior, not implementation.** Assert on outputs, errors, and exit codes. Do not assert on private helper calls.
- **Mock only at boundaries:** model API calls, network and git access, and the filesystem where practical. Use real Pydantic models and real YAML parsing in tests.
- **Use fixtures for input data.** Store sample registries and doc trees under `tests/fixtures/`. Keep them small and synthetic. Never copy real internal documentation into fixtures.
- **Every requirement scenario maps to at least one test.** Name the test after the scenario so the mapping is visible.
- **Bug fix order:** reproduce with a failing test, fix, then keep the test as a regression guard.
- **Coverage is a floor, not a goal.** Do not write tests that only raise the coverage number.

## Output expectations

When reporting back after a TDD session, state:
- Which scenarios now have tests
- The red-to-green result for each
- The final result of `uv run pytest` and `uv run mypy` with real output, not a summary of intent
