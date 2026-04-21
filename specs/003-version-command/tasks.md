# Tasks: Version Command

**Input**: Design documents from `/specs/003-version-command/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli-version-contract.md

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize dependencies and base CLI wiring entrypoint.

- [ ] T001 Add Cobra dependency declarations in go.mod and go.sum
- [ ] T002 Create CLI package skeleton files in internal/cli/run.go, internal/cli/root.go, and internal/cli/version_cmd.go
- [ ] T003 Create dedicated parser package scaffold in internal/flags/parser.go
- [ ] T004 Update process entrypoint to delegate parsing/execution to package code in cmd/mytets/main.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Build reusable command execution/parsing foundation required by all stories.

**⚠️ CRITICAL**: No user story implementation begins before this phase is complete.

- [ ] T005 Implement build-time version holder with fallback normalization in internal/version/version.go
- [ ] T006 Implement parser configuration type and constructor for CLI wiring in internal/flags/parser.go
- [ ] T007 Implement root command builder with injected parser config in internal/cli/root.go
- [ ] T008 Implement execution entrypoints `Execute` and `ExecuteArgs` with stdout/stderr capture in internal/cli/run.go
- [ ] T009 [P] Add parser-focused unit tests for configuration behavior in internal/flags/parser_test.go
- [ ] T010 [P] Add runner seam unit tests for stdout/stderr/exit-code capture in internal/cli/run_test.go

**Checkpoint**: Foundation complete; user stories can be implemented and tested independently.

---

## Phase 3: User Story 1 - Developer Queries Application Version (Priority: P1) 🎯 MVP

**Goal**: Running `mytets version` prints exactly one plain version string and exits successfully.

**Independent Test**: Build binary with `-ldflags` version injection, run `mytets version`, verify stdout is exactly one line `X.Y.Z`, stderr empty, and exit code 0.

### Tests for User Story 1

- [ ] T011 [P] [US1] Add unit tests for semantic version and fallback rendering in internal/version/version_test.go
- [ ] T012 [P] [US1] Add integration test for successful `mytets version` invocation output/exit code in tests/integration/version_command_test.go

### Implementation for User Story 1

- [ ] T013 [US1] Implement `version` Cobra subcommand handler to print single plain string in internal/cli/version_cmd.go
- [ ] T014 [US1] Register `version` subcommand on root command in internal/cli/root.go
- [ ] T015 [US1] Wire command execution path so `mytets version` reaches handler in internal/cli/run.go
- [ ] T016 [US1] Enforce no additional args for `version` command behavior in internal/cli/version_cmd.go

**Checkpoint**: `mytets version` is fully functional for interactive developer usage.

---

## Phase 4: User Story 2 - Scripted / CI Version Check (Priority: P2)

**Goal**: Scripts and CI can reliably capture stable plain-text output and success status.

**Independent Test**: Capture stdout in a script-like test, verify `^[0-9]+\.[0-9]+\.[0-9]+$` (or `dev` fallback), stderr empty on success, exit code 0.

### Tests for User Story 2

- [ ] T017 [P] [US2] Add integration test for regex-compatible scripting capture and empty stderr in tests/integration/version_command_test.go
- [ ] T018 [P] [US2] Add integration test for fallback `dev` output when version is not injected in tests/integration/version_command_test.go
- [ ] T019 [P] [US2] Add integration test for invalid-flag non-zero exit behavior in tests/integration/version_command_test.go

### Implementation for User Story 2

- [ ] T020 [US2] Finalize invocation contract behavior mapping parser errors to stderr/non-zero exits in internal/cli/run.go
- [ ] T021 [US2] Verify default plain-text output remains unchanged when no output flag is provided (covered by integration tests)
- [ ] T022 [US2] Document script/CI usage and ldflags injection examples in README.md
- [ ] T027 [US2] Implement output `json` flag parsing and version JSON rendering in internal/flags/parser.go and internal/cli/version_cmd.go
- [ ] T028 [P] [US2] Add integration tests for plain output versus JSON output contract in tests/integration/version_command_test.go

**Checkpoint**: `mytets version` is automation-safe and contract-compliant.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final quality, docs consistency, and validation against quickstart.

- [ ] T023 [P] Update command contract notes after implementation in specs/003-version-command/contracts/cli-version-contract.md
- [ ] T024 [P] Update quickstart validation commands to match final implementation in specs/003-version-command/quickstart.md
- [ ] T025 [P] Verify feature summary and constraints alignment in specs/003-version-command/plan.md
- [ ] T026 Run final test pass and record results in specs/003-version-command/quickstart.md
- [ ] T029 [P] Add performance validation test for `mytets version` (<100ms threshold) in tests/integration/version_performance_test.go
- [ ] T030 [P] Add CI step for race detector execution (`go test -race ./...`) in .github/workflows/ci.yml
- [ ] T031 [P] Add CI step for `go vet` and `staticcheck` in .github/workflows/ci.yml
- [ ] T032 [P] Add CI step for `golangci-lint` in .github/workflows/ci.yml
- [ ] T033 [P] Add CI matrix validation for Linux/macOS/Windows integration tests in .github/workflows/ci.yml
- [ ] T034 Add quickstart verification commands for race/lint/performance checks in specs/003-version-command/quickstart.md
- [ ] T035 Run full constitution gate validation and record results in specs/003-version-command/quickstart.md
- [ ] T036 [P] Add dependency exception justification for `github.com/spf13/cobra` in README.md, including why standard-library parsing was not selected
- [ ] T037 [P] Add integration test to verify concise help text for `version` command and `--output` flag in tests/integration/version_command_test.go

---

## Dependencies & Execution Order

### Phase Dependencies

- Phase 1 (Setup): No dependencies.
- Phase 2 (Foundational): Depends on Phase 1; blocks all user stories.
- Phase 3 (US1): Depends on Phase 2.
- Phase 4 (US2): Depends on Phase 2; can proceed after US1 command path exists.
- Phase 5 (Polish): Depends on completion of US1 and US2.

### User Story Dependencies

- US1 (P1): Starts after foundational phase; delivers MVP.
- US2 (P2): Builds on the same `version` command behavior and validates automation-focused guarantees.

### Parallel Opportunities

- Setup: none (shared files).
- Foundational: T009 and T010 can run in parallel.
- US1: T011 and T012 can run in parallel before implementation tasks.
- US2: T017, T018, T019, and T028 can run in parallel.
- Polish: T023, T024, T025, T029, T030, T031, T032, T033, T036, and T037 can run in parallel.

---

## Parallel Example: User Story 1

```bash
# Parallel test authoring
Task T011: internal/version/version_test.go
Task T012: tests/integration/version_command_test.go
```

## Parallel Example: User Story 2

```bash
# Parallel integration checks
Task T017: tests/integration/version_command_test.go
Task T018: tests/integration/version_command_test.go
Task T019: tests/integration/version_command_test.go
```

---

## Implementation Strategy

### MVP First (US1)

1. Complete Phase 1 and Phase 2.
2. Complete US1 tests and implementation (Phase 3).
3. Validate `mytets version` manually and via tests.
4. Demo/release MVP.

### Incremental Delivery

1. Build foundation once (Phases 1-2).
2. Deliver US1 (interactive usage) as MVP.
3. Deliver US2 (automation guarantees).
4. Finish with polish updates and final validation.
