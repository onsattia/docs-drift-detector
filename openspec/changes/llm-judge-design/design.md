## Context

The judge labels each endpoint claim `consistent`, `drifted`, `missing`, or `unverifiable`. It needs evidence (the code search results) and a model to reason over that evidence. The model endpoint is undecided (epic #8 spike), so the design must not depend on a particular provider.

The Go codebase currently has no interfaces. The registry reads files directly, which is fine for it. The judge is the first part that calls an external service, so it is the first place a port is needed.

## Goals / Non-Goals

**Goals:**
- Verdict logic can be unit tested with a fake model, with no network access.
- Swapping the model provider changes one adapter and nothing in the core.
- The approved-endpoint rule can be enforced in one place.

**Non-Goals:**
- Choosing the provider or the model. That is the endpoint spike in #8.
- Prompt design and eval thresholds. Those belong to the judge specs when implementation starts.
- Code search and the draft-PR output. Those are separate epics.

## Decisions

**Ports live with the consumer.** The judge package defines the `Judge` interface it needs. The model adapter does not define the interface, and nothing outside the judge package imports the adapter directly. Go interfaces are satisfied implicitly, so the adapter needs no declaration. Alternative: one shared `ports` package. Rejected because it would collect interfaces that no single consumer owns.

**Port shape.** The core asks the port for one verdict at a time:

```go
type Judge interface {
    Label(ctx context.Context, claim Claim, evidence []Evidence) (Verdict, error)
}
```

`Claim` is an endpoint claim from the extractor (#6). `Evidence` is a code search hit from search (#7). `Verdict` is a typed string with the four allowed values, so an invalid label cannot be constructed by accident. Alternative: return a bare string. Rejected because a typo would pass the type checker.

**Adapters are separate packages.** Each provider gets its own package under `internal/judge/`, for example `internal/judge/bedrock`. The adapter translates `Claim` and `Evidence` into the provider's request and maps the response back to a `Verdict`. Alternative: a single client with switches for each provider. Rejected because every provider would then be compiled and tested together.

**The core owns the verdict rules.** Rules such as "no evidence means `missing`, not `drifted`" live in the judge package, not in the adapter. That logic is the part worth testing most, and it can be tested with a fake `Judge` that returns fixed verdicts.

**Test fake in the judge package.** A `fake` implementation returns recorded, synthetic verdicts keyed by claim ID. It is used by unit tests. The eval suite (CLAUDE.md) runs the real adapter against the golden set in CI.

**Context propagation.** `Label` takes a `context.Context` so that calls can be cancelled and time-limited. This is a standard Go convention for anything that does I/O.

## Risks / Trade-offs

- **Over-abstraction.** A port with one adapter can look like wasted code. → Mitigation: the port is justified by the fake for tests, which is needed from the first commit of #8 anyway.
- **Leaking provider details.** Request IDs or token counts could leak into `Verdict`. → Mitigation: `Verdict` carries only the label and a short reason. Provider metadata stays inside the adapter and goes to logs as counts only, per CLAUDE.md.
- **Approval drift.** A new adapter could call an unapproved endpoint. → Mitigation: each adapter takes its base URL from a config value checked against an allow-list, and that check has its own test.

## Migration Plan

No migration. This change adds no code. When #8 starts:
1. Write the judge specs and tests against the `Judge` port and the fake.
2. Implement the verdict rules in `internal/judge`.
3. Add one adapter after the endpoint spike is approved.

## Open Questions

None that change the port. The endpoint and model choice are decided in the spike, and they only affect which adapter is written.
