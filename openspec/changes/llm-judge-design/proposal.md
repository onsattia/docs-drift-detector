## Why

The LLM judge (epic #8) will call a model endpoint that is not chosen yet. Without a boundary between the verdict logic and the model call, the logic can't be unit tested without a live model, and switching endpoints would touch the core. CLAUDE.md requires that unit tests run without network or model access and that model calls are mocked at the client boundary. That boundary needs to be designed before the judge is built.

## What Changes

- Record the design for the judge as ports and adapters: the core verdict logic depends on a small `Judge` interface, and each model provider is an adapter behind it.
- Decide where the judge's types live and which package owns each rule, so the judge can be built and tested in slices.
- No judge code is written in this change. Implementation starts with epic #8.

## Capabilities

### New Capabilities

None yet. The judge's requirements will be written as specs when implementation starts. This change records the design only.

### Modified Capabilities

None.

## Impact

- **Code:** none in this change. The design affects the future `internal/judge` package and its adapters.
- **Epic #8:** implementation will follow this design. The endpoint spike (#8) still decides which adapter gets built first.
- **Security:** model calls go only to an approved endpoint, per CLAUDE.md. The design keeps adapters separate so the endpoint can be checked in one place.
