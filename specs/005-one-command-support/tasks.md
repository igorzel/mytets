# Tasks: One Command Support

**Input**: Design documents from `specs/005-one-command-support/`
**Spec**: [spec.md](spec.md) | **Plan**: [plan.md](plan.md) | **Data Model**: [data-model.md](data-model.md) | **Contract**: [contracts/cli-one-contract.md](contracts/cli-one-contract.md)

**Tests**: Included — explicitly requested in feature specification (FR-007, FR-008)

**Organization**: Tasks grouped by user story (US1: plain text P1/MVP, US2: JSON output P2) to enable independent implementation and testing

## Format: `[ID] [P?] [Story] Description with file path`

- **[P]**: Parallelizable (different files, no blocking dependencies)
- **[Story]**: User story label (US1, US2, etc.)
- File paths relative to repository root

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project structure initialization for the new command

**Status**: ✅ Complete

- [X] T001 Create `internal/commands/` directory for subcommand implementations
- [X] T002 Create unit test file structure in `internal/commands/one/one_test.go` (stub)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that must be complete before user story work

**Status**: Complete — Cobra framework, global flags, and integration test infrastructure already in place

**Checkpoint**: Foundation ready; user story work can now proceed in parallel

---

## Phase 3: User Story 1 - Run `mytets one` for Plain Text Message (Priority: P1) 🎯 MVP

**Goal**: Implement basic `mytets one` command that outputs "Fake message, tbd" in plain text mode

**Independent Test**: Run `mytets one`; verify stdout equals `Fake message, tbd\n`, stderr is empty, exit code is 0

### Unit Tests for User Story 1 (Test-First Approach)

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation (TDD)**

- [X] T003 [P] [US1] Create unit test file `internal/commands/one/one_test.go` with test function stubs for plain text
- [X] T004 [P] [US1] Write unit test `TestOutputPlain()` in `internal/commands/one/one_test.go` that verifies plain text output behavior
- [X] T005 [P] [US1] Write unit test `TestNew()` in `internal/commands/one/one_test.go` that verifies command is properly constructed

### Integration Tests for User Story 1

- [X] T006 [P] [US1] Create integration test file `tests/integration/one_command_test.go` with test function stubs
- [X] T007 [P] [US1] Write integration test `TestOneCommandPlain()` in `tests/integration/one_command_test.go` that invokes `./main one` and verifies stdout/stderr/exit code

### Implementation for User Story 1

- [X] T008 [P] [US1] Implement command package skeleton in `internal/commands/one/one.go` with `New()` function that returns a Cobra command
- [X] T009 [P] [US1] Implement `outputPlain()` function in `internal/commands/one/one.go` that prints "Fake message, tbd" to stdout
- [X] T010 [US1] Implement `execute()` function in `internal/commands/one/one.go` that dispatches to output functions based on JSON mode flag (depends on T008, T009)
- [X] T011 [US1] Register one command in `internal/cli/root.go` by adding import and `root.AddCommand(one.New(cfg))` call
- [X] T012 [US1] Verify command registration by checking root.go imports and AddCommand call

**Checkpoint**: At this point, User Story 1 (plain text output) should be fully functional and pass all unit and integration tests

---

## Phase 4: User Story 2 - Run `mytets --json one` for Structured Output (Priority: P2)

**Goal**: Add JSON output support; respects global `--json` flag to emit compact JSON

**Independent Test**: Run `mytets --json one`; parse stdout as JSON; verify `{"message":"Fake message, tbd"}`; exit code 0

### Unit Tests for User Story 2

- [X] T013 [P] [US2] Write unit test `TestOutputJSON()` in `internal/commands/one/one_test.go` that verifies JSON output is valid and compact
- [X] T014 [P] [US2] Write unit test `TestResponseJSONFormat()` in `internal/commands/one/one_test.go` that marshals Response struct and verifies exact JSON format `{"message":"Fake message, tbd"}`
- [X] T015 [P] [US2] Write unit test `TestExecuteJSON()` in `internal/commands/one/one_test.go` that verifies execute() dispatches to outputJSON() when JSON mode is enabled

### Integration Tests for User Story 2

- [X] T016 [P] [US2] Write integration test `TestOneCommandJSON()` in `tests/integration/one_command_test.go` that invokes `./main --json one`, parses JSON response, and verifies message field
- [X] T017 [P] [US2] Write integration test `TestOneCommandJSONExitCode()` in `tests/integration/one_command_test.go` that verifies exit code is 0 in JSON mode

### Implementation for User Story 2

- [X] T018 [P] [US2] Define `Response` struct in `internal/commands/one/one.go` with json-tagged `Message` field
- [X] T019 [P] [US2] Implement `outputJSON()` function in `internal/commands/one/one.go` that marshals Response struct and prints compact JSON to stdout (use `json.Marshal`, not `json.MarshalIndent`)
- [X] T020 [US2] Update `execute()` function in `internal/commands/one/one.go` to check `cfg.OutputJSON` and dispatch to `outputJSON()` when true (depends on T018, T019)
- [X] T021 [US2] Verify JSON output format is compact by confirming no whitespace between fields

**Checkpoint**: At this point, User Story 2 (JSON output) should be fully functional; both plain text and JSON modes work correctly

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final validation, documentation, and test coverage verification

✅ **Status**: Complete

- [X] T022 [P] Run all unit tests: `go test ./internal/commands/one/...` and verify all tests pass
- [X] T023 [P] Run integration tests: `go test ./tests/integration/...` and verify all tests pass
- [X] T024 [P] Build the project: `go build -o main ./cmd/mytets/` and verify no errors
- [X] T025 Manually verify `./main one` outputs "Fake message, tbd" with exit code 0
- [X] T026 Manually verify `./main --json one` outputs valid JSON with exit code 0
- [X] T027 Verify help text displays: `./main one --help` shows command description and global flags
- [X] T028 Verify error handling: `./main --json one --invalid-flag` exits non-zero with error message to stderr
- [X] T029 Run coverage check: `go test -cover ./internal/commands/one/...` and document coverage percentage

---

## Implementation Dependencies & Parallel Execution

### Execution Order (MVP Path: US1 Only)

For a minimal viable product (Phase 3 only):

1. **Setup** (T001-T002): 2 tasks, quick directory creation
2. **US1 Tests** (T003-T007): 5 parallel-able test tasks
3. **US1 Implementation** (T008-T012): 5 tasks, T008-T009 parallel, T010 depends on both, T011-T012 depend on T010

**Estimated US1 time**: ~3-4 hours for experienced Go developer

### Full Feature (MVP + US2)

Add Phase 4 (US2 tests + implementation):
- US2 tests (T013-T017): 5 tasks, mostly parallelizable
- US2 implementation (T018-T021): 4 tasks with minimal dependencies

**Estimated total time**: ~5-6 hours including testing and verification

### Parallel Execution Example (Full Feature)

```
Round 1 (concurrent):  T003, T004, T005, T006, T007, T008, T009
Round 2 (depends on T008, T009): T010
Round 3 (concurrent):  T011, T012, T013, T014, T015, T016, T017, T018, T019
Round 4 (depends on T018, T019): T020, T021
Round 5 (verification): T022-T029
```

---

## Task Summary

| Phase | Count | Scope | MVP? |
|-------|-------|-------|------|
| Phase 1: Setup | 2 | Directory structure | ✅ |
| Phase 2: Foundational | 0 | (Already complete) | ✅ |
| Phase 3: US1 (Plain Text) | 10 | Core feature implementation & tests | ✅ |
| Phase 4: US2 (JSON Output) | 9 | Enhanced feature & tests | ⭕ |
| Phase 5: Polish | 8 | Verification & coverage | ✅ |
| **TOTAL** | **29** | **Complete feature with full test coverage** | **23 for MVP** |

---

## Success Criteria (from Spec)

Each task supports one or more success criteria:

| SC-001 | 100% exit code 0 on success | Tasks: T025, T026, T022-T023 |
| SC-002 | Plain text output exact match | Tasks: T004, T007, T025 |
| SC-003 | JSON output valid & parseable | Tasks: T013-T014, T016, T026 |
| SC-004 | Full test coverage both paths | Tasks: T003-T007 (US1), T013-T017 (US2) |

---

## Acceptance Checklist

✅ **All tasks complete!**

- [X] All unit tests pass (`go test ./internal/commands/one/...`) — 5/5 tests passing
- [X] All integration tests pass (`go test ./tests/integration/...`) — 5/5 tests passing
- [X] Project builds cleanly (`go build -o main ./cmd/mytets/`) — Binary created successfully
- [X] Plain text mode works: `./main one` outputs exactly "Fake message, tbd\n"
- [X] JSON mode works: `./main --json one` outputs valid compact JSON `{"message":"Fake message, tbd"}`
- [X] Exit codes correct: 0 on success, non-zero on invalid flags
- [X] No output to stderr on success
- [X] Help text displays via `--help` — Shows "The one command outputs a fixed message in plain text or JSON format."
- [X] Code review ready and tested
