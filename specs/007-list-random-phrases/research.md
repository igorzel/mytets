# Research: List Command - Random Phrase List

**Feature**: `mytets list` command with reusable list-generation logic  
**Date**: 2026-04-21  
**Scope**: reusable package boundaries, unique random selection, CLI adapter design, test seams

## 1. Reusable Package Boundary

**Decision**: Implement the phrase-list generation logic in a dedicated `internal/listing` package and keep Cobra-specific concerns in `internal/commands/list`.

**Rationale**:
- The user explicitly requested the `list` functionality be decoupled from the CLI for easier unit testing.
- A reusable package that returns domain values instead of formatted bytes can be reused later by alternative delivery layers such as an HTTP handler.
- This preserves a clean separation of concerns: command parsing and output formatting in one place, list generation in another.

**Alternatives considered**:
- Put all logic in `internal/commands/list`: rejected because it couples business logic to Cobra and makes future reuse harder.
- Expand `internal/phrases` to own the list use case: rejected because phrase storage/loading and list-request behavior are separate responsibilities.

## 2. Unique Random Selection Strategy

**Decision**: Build each list by sampling phrases without replacement from the available phrase set, capping the result length at the number of unique available phrases.

**Rationale**:
- The specification requires no repeated phrase within one result.
- Sampling without replacement is a direct fit for the requirement and remains efficient for the current data size.
- Deduplicating by phrase text before final selection handles the edge case where the source contains duplicate entries.

**Alternatives considered**:
- Repeated random picks with retry-on-duplicate: rejected because it is less predictable and gets inefficient as the requested count approaches the available set.
- Shuffle the entire list every time and slice: acceptable, but the selected approach can still use a partial shuffle internally while keeping the decision framed around the required behavior rather than one concrete algorithm.

## 3. Deterministic Testing Seam

**Decision**: Design `internal/listing` to accept an injected randomness source or picker function so unit tests can assert exact outcomes deterministically.

**Rationale**:
- The constitution requires mockable inputs and outputs for deterministic testing of random-selection logic.
- Injecting the randomness seam makes unit tests stable without relying on timing or probabilistic assertions.
- CLI integration tests can still verify broad behavioral properties while unit tests verify exact selection and count rules.

**Alternatives considered**:
- Use the package-global RNG directly in all tests: rejected because it produces brittle tests.
- Mock the entire phrase source at the command layer: rejected because it tests the wrong abstraction boundary.

## 4. Output Formatting Boundary

**Decision**: Keep plain-text and JSON formatting in the CLI-facing package and have the reusable `internal/listing` package return a neutral list result such as `[]string` or a domain result struct.

**Rationale**:
- Output format is adapter-specific, while phrase selection is reusable domain logic.
- Returning neutral data keeps the reusable package suitable for future HTTP or other integrations.
- The current CLI can still produce compact JSON arrays and plain-text lines without leaking Cobra into domain code.

**Alternatives considered**:
- Have `internal/listing` return already formatted JSON/plain text: rejected because it would hard-code a CLI delivery concern into reusable logic.

## 5. Shared Phrase-Source Validation

**Decision**: Reuse the existing `internal/phrases` initialization and validation path as the source of truth for phrase availability, and propagate errors through `listing` and command handlers.

**Rationale**:
- The spec requires missing or empty phrase data to stop both `one` and `list` from working.
- The repository already centralizes embedded phrase loading in `internal/phrases`.
- Keeping validation in one source package avoids divergent behavior between commands.

**Alternatives considered**:
- Add separate startup checks inside the new `list` command: rejected because it duplicates logic and can drift from `one`.
