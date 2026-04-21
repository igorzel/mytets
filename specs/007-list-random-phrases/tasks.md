# Tasks: List Command - Random Phrase List

**Input**: Design documents from /specs/007-list-random-phrases/
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli-list-contract.md, quickstart.md

**Tests**: Included. Unit and integration tests are explicitly requested in the feature specification (FR-013, FR-014).

**Organization**: Tasks are grouped by user story so each story can be implemented and validated independently.

## Format: [ID] [P?] [Story] Description

- **[P]**: Can run in parallel (different files, no dependency on incomplete tasks)
- **[Story]**: User story label (US1, US2, US3)
- Every task includes an explicit file path

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create the reusable listing package scaffold that all user stories depend on

- [ ] T001 Create reusable listing package scaffold with package declaration and exported function stubs in internal/listing/listing.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Implement the injectable, reusable phrase-list domain logic that both US1 and US2 consume

**⚠️ CRITICAL**: No user story work starts before this phase is complete

- [ ] T002 Implement `Select(phrases []string, count int, rng *rand.Rand) []string` in internal/listing/listing.go — unique without-replacement sampling capped at available unique phrase count
- [ ] T003 [P] Implement unit tests for Select covering default count, custom count, oversized count cap, uniqueness guarantee, deduplication of duplicate source entries, and deterministic seam in internal/listing/listing_test.go

**Checkpoint**: Foundation complete. User stories can now proceed independently.

---

## Phase 3: User Story 1 - Random Phrase List as Plain Text (Priority: P1) 🎯 MVP

**Goal**: `mytets list` returns 5 unique phrases as plain text; `mytets list --count N` returns N unique phrases; exit code 0

**Independent Test**: Run `mytets list` and verify 5 non-empty unique lines on stdout, empty stderr, exit code 0. Run `mytets list --count 3` and verify 3 lines. Run `mytets list --count 999` and verify all available phrases returned without duplicates.

### Tests for User Story 1

- [ ] T004 [US1] Write integration tests for plain-text default (5 lines), `--count N`, oversized-count edge case, and invalid `--count` input (e.g., `--count abc`, `--count -1`) verifying non-zero exit code and stderr output in tests/integration/list_command_test.go

### Implementation for User Story 1

- [ ] T005 [US1] Implement list Cobra command with command-specific `--count` flag (default 5) calling `listing.Select` and printing plain-text output via `cmd.OutOrStdout()` in internal/commands/list/list.go
- [ ] T006 [P] [US1] Implement unit tests for list command construction, --count flag default, and plain-text output formatting in internal/commands/list/list_test.go
- [ ] T007 [US1] Register list command in the root command in internal/cli/root.go

**Checkpoint**: User Story 1 is independently functional and testable.

---

## Phase 4: User Story 2 - Random Phrase List as JSON (Priority: P1)

**Goal**: `mytets --output json list` returns a compact JSON array of `{"message":"..."}` objects; exit code 0

**Independent Test**: Run `mytets --output json list` and verify stdout is valid compact JSON array, each item has a `message` field with a phrase from the embedded set, stderr is empty, exit code 0. Run with `--count 2` and verify 2-item array.

### Tests for User Story 2

- [ ] T008 [US2] Add integration tests for JSON default output, `--output json list --count N`, and oversized-count JSON output in tests/integration/list_command_test.go

### Implementation for User Story 2

- [ ] T009 [US2] Add JSON output branch to the list command that marshals the selected phrases into a compact JSON array where each item is `{"message":"..."}` in internal/commands/list/list.go
- [ ] T010 [P] [US2] Add unit tests for JSON output format, compactness, array structure, and message field values in internal/commands/list/list_test.go

**Checkpoint**: User Story 2 is independently functional and testable.

---

## Phase 5: User Story 3 - Fail Fast When Phrase Data Is Unavailable (Priority: P2)

**Goal**: When the embedded phrase source is missing or empty, phrase-based commands (`one` and `list`) do not succeed and an error is printed to stderr

**Independent Test**: Confirm that the existing phrases-package initialization failure path propagates through `listing.Select` and the list command error handler such that the process exits with non-zero code and writes to stderr.

### Tests for User Story 3

- [ ] T011 [US3] Add unit test in internal/commands/list/list_test.go verifying that the list command returns an error and writes to stderr when the phrase source function is injected to return an empty slice (integration-level test is not feasible because `internal/phrases.init()` panics before command execution when the embedded source is invalid)
- [ ] T012 [P] [US3] Add unit test for Select when called with empty input slice, confirming it returns an empty result without panicking in internal/listing/listing_test.go

### Implementation for User Story 3

- [ ] T013 [US3] Confirm the list command propagates phrase-load errors from `phrases.Messages()` as a returned error (non-zero exit code, stderr output) in internal/commands/list/list.go

**Checkpoint**: All three user stories are independently functional and testable.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and final confidence checks across all stories

- [ ] T014 [P] Update README.md with list command usage examples covering plain-text default, --count, and --output json invocations
- [ ] T015 [P] Update quickstart verification commands to reflect final command behavior in specs/007-list-random-phrases/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1; **blocks all user stories**
- **Phase 3 (US1)**: Depends on Phase 2; no dependency on US2 or US3
- **Phase 4 (US2)**: Depends on Phase 2; no dependency on US1 (shares foundation only)
- **Phase 5 (US3)**: Depends on Phase 2; verifies error path already present in foundation
- **Phase 6 (Polish)**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1 (P1)**: Independent after foundation. Required for MVP.
- **US2 (P1)**: Independent after foundation. Shares `listing.Select` output but does not require US1.
- **US3 (P2)**: Independent after foundation. Verifies existing error propagation path; no new runtime code needed if error path is already correct.

### Within Each User Story

- Integration tests first, then implementation (test-first workflow per constitution)
- Command construction before output formatting
- Core implementation before error-path hardening
- Story-level checkpoint must pass before moving to polish

---

## Parallel Opportunities

- **Phase 2**: T002 (implementation) and T003 (tests in a separate file) can proceed in parallel if staffed
- **Phase 3**: T006 (unit tests) can run in parallel with T005 (implementation) across different files
- **Phase 4**: T010 (unit tests) can run in parallel with T009 (implementation) across different files
- **Phase 5**: T012 (unit test) can run in parallel with T011 (integration test)
- **Phase 6**: T014 and T015 are fully independent

---

## Parallel Example: User Story 1

```
Developer A                        Developer B
─────────────────────────────────────────────────────
T004 write integration tests
(tests/integration/list_command_test.go)
                                   T006 write unit tests
                                   (internal/commands/list/list_test.go)
     ↓ both done
T005 implement list command
(internal/commands/list/list.go)
T007 register in root.go
(internal/cli/root.go)
     ↓ run tests to confirm checkpoint
```

---

## Implementation Strategy

### MVP First (US1 only)

1. Complete Phase 1 (T001)
2. Complete Phase 2 (T002, T003)
3. Complete Phase 3 (T004 → T005 → T006 → T007)
4. Validate via `go test ./internal/listing/... ./internal/commands/list/... ./tests/integration/...`

### Incremental Delivery

1. Deliver US1 plain-text list output
2. Deliver US2 JSON list output (add JSON branch to existing list command)
3. Confirm US3 error propagation (verify or add single guard)
4. Execute polish (T014, T015)

### Suggested `go test` Commands

```bash
# Run after Phase 2 complete
go test ./internal/listing/...

# Run after Phase 3 complete
go test ./internal/listing/... ./internal/commands/list/... ./tests/integration/...

# Run full suite before merge
go test ./...
```

---

## Notes

- `internal/listing` must have no Cobra-specific imports; it must accept `[]string` and return `[]string` only.
- The injectable `*rand.Rand` seam in `listing.Select` is the unit-test handle; integration tests use the live random source.
- US3 may require no new runtime code if the existing `phrases.Messages()` error path already surfaces through the list command correctly — T013 is a verification task that adds a guard only if the propagation is missing.
- Task IDs are sequential in execution order. All [P]-marked tasks within a phase can be parallelized across separate files.
