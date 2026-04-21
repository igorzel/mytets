# Tasks: One Command - Random Phrase Output

**Input**: Design documents from /specs/006-one-random-phrases/
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/cli-one-contract.md, quickstart.md

**Tests**: Included. Unit and integration tests are explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story so each story can be implemented and validated independently.

## Format: [ID] [P?] [Story] Description

- [P]: Can run in parallel (different files, no dependency on incomplete tasks)
- [Story]: User story label (US1, US2)
- Every task includes explicit file path(s)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Prepare embedded phrase data and package scaffolding

- [X] T001 Create embedded phrase source file with required messages schema in internal/phrases/phrases.json
- [X] T002 [P] Create phrase package scaffold with embed declaration and package API stubs in internal/phrases/phrases.go
- [X] T003 [P] Create phrase package test scaffold in internal/phrases/phrases_test.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Build shared phrase loading and random selection infrastructure required by all stories

**CRITICAL**: No user story work starts before this phase is complete

- [X] T004 Implement embedded JSON decoding and schema validation for phrases in internal/phrases/phrases.go
- [X] T005 Implement random phrase selection using math/rand in internal/phrases/phrases.go
- [X] T006 Implement foundational unit tests for loader validation and selector behavior in internal/phrases/phrases_test.go
- [X] T007 Refactor command dependency wiring to consume phrase package APIs in internal/commands/one/one.go

**Checkpoint**: Foundation complete. User stories can now proceed.

---

## Phase 3: User Story 1 - Random Plain Text Output (Priority: P1) 🎯 MVP

**Goal**: Running mytets one returns one random phrase as plain text with exit code 0

**Independent Test**: Execute mytets one repeatedly and verify stdout is plain text from embedded phrase set, stderr is empty, exit code is 0

### Tests for User Story 1

- [X] T008 [US1] Update plain output unit tests to assert random phrase behavior in internal/commands/one/one_test.go
- [X] T009 [US1] Update plain output integration tests to verify phrase membership and success exit code in tests/integration/one_command_test.go

### Implementation for User Story 1

- [X] T010 [US1] Replace fixed message output with random phrase retrieval in plain output branch in internal/commands/one/one.go
- [X] T011 [US1] Update command descriptions to reflect random phrase behavior in internal/commands/one/one.go
- [X] T012 [US1] Ensure plain path returns descriptive errors to stderr with exit code 1 when phrase retrieval fails in internal/commands/one/one.go

**Checkpoint**: User Story 1 is independently functional and testable.

---

## Phase 4: User Story 2 - Random JSON Output (Priority: P1)

**Goal**: Running mytets --output json one returns compact JSON with random phrase and exit code 0

**Independent Test**: Execute mytets --output json one and verify compact JSON object with message field, valid phrase value, empty stderr, exit code 0

### Tests for User Story 2

- [X] T013 [US2] Update JSON unit tests to validate compact object format and message source in internal/commands/one/one_test.go
- [X] T014 [US2] Update JSON integration tests to use global flag order and validate output contract in tests/integration/one_command_test.go
- [X] T015 [P] [US2] Add CLI execution test for global flag before subcommand in internal/cli/run_test.go

### Implementation for User Story 2

- [X] T016 [US2] Return selected random phrase from JSON output branch in internal/commands/one/one.go
- [X] T017 [US2] Keep JSON response compact single-line marshaling with message field contract in internal/commands/one/one.go
- [X] T018 [US2] Ensure invalid output format remains descriptive and propagates as stderr error with exit code 1 in internal/flags/parser.go

**Checkpoint**: User Story 2 is independently functional and testable.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Finalize documentation and confidence checks across stories

- [X] T019 [P] Update command usage documentation for plain and JSON random output in README.md
- [X] T020 [P] Align quickstart verification commands with final behavior in specs/006-one-random-phrases/quickstart.md
- [X] T021 Add integration regression coverage for error and success matrix in tests/integration/one_command_test.go
- [X] T022 Add final command-level unit coverage for response formatting and error propagation in internal/commands/one/one_test.go

---

## Dependencies & Execution Order

### Phase Dependencies

- Phase 1 (Setup): No dependencies
- Phase 2 (Foundational): Depends on Phase 1 and blocks all user stories
- Phase 3 (US1): Depends on Phase 2
- Phase 4 (US2): Depends on Phase 2; can run in parallel with US1 after foundation if staffed
- Phase 5 (Polish): Depends on completion of selected user stories

### User Story Dependencies

- US1 (P1): Independent after foundation
- US2 (P1): Independent after foundation; integrates command output mode without requiring US1 completion

### Within Each User Story

- Tests first, then implementation
- Command behavior updates before documentation updates
- Story-level checks must pass before moving to polish

---

## Parallel Opportunities

- Setup parallel: T002 and T003
- US2 parallel: T015 can run alongside T013 and T014
- Polish parallel: T019 and T020

---

## Parallel Example: User Story 2

- Run T013 and T014 together (unit and integration updates in separate files)
- Run T015 in parallel (CLI execution test in internal/cli/run_test.go)
- After tests are in place, run T016 and T017 sequentially in internal/commands/one/one.go

---

## Implementation Strategy

### MVP First (US1)

1. Complete Phase 1
2. Complete Phase 2
3. Complete Phase 3 (US1)
4. Validate independently via tests/integration/one_command_test.go and internal/commands/one/one_test.go

### Incremental Delivery

1. Deliver US1 plain-text random output
2. Deliver US2 JSON random output
3. Execute polish/documentation updates

### Team Parallel Strategy

1. One developer handles phrase package foundation (Phase 2)
2. Another developer prepares US1 tests while foundation nears completion
3. Another developer prepares US2 CLI/run tests in internal/cli/run_test.go
4. Merge with story checkpoints before polish

---

## Notes

- Tasks are designed against current files already present in the repository.
- Existing one command and integration tests are updated, not recreated.
- All task descriptions use concrete repository paths for immediate execution.
