# Tasks: System Locale Localized Messages

**Input**: Design documents from `/specs/008-locale-messages/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/cli-locale-contract.md

**Tests**: Included — spec explicitly requests unit tests for internal packages and integration tests for localized help.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add the TOML dependency and create the i18n package directory structure

- [x] T001 Add `github.com/BurntSushi/toml` dependency via `go get github.com/BurntSushi/toml`
- [x] T002 Create directory structure `internal/i18n/locales/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core i18n infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 Create English localization file `internal/i18n/locales/en.toml` with all translation keys from data-model.md (root.short, version.short, one.short, one.long, list.short, list.long, flag.output, flag.count, help.usage, help.available_commands, help.flags, help.global_flags, help.additional_help, help.aliases, help.help_for, error.no_phrases, error.invalid_output_format, error.failed_select_phrase, error.failed_encode_json, error.unknown_command, error.unknown_flag)
- [x] T004 Create Ukrainian localization file `internal/i18n/locales/uk.toml` with all the same keys as en.toml, translated to Ukrainian per contracts/cli-locale-contract.md
- [x] T005 Implement i18n package core in `internal/i18n/i18n.go`: embed `locales/*.toml` via `//go:embed`, implement `LoadBundle()` to auto-discover and parse all TOML files from embedded FS into `map[string]map[string]string`, implement `DetectLocale()` to read `LC_ALL` → `LC_MESSAGES` → `LANG` and extract language code, implement `Translate(key string) string` with two-tier fallback (active language → `en` → key name), and expose `SetLang(lang string)` for test injection
- [x] T006 Write unit tests in `internal/i18n/i18n_test.go`: test `DetectLocale()` with various env var combinations (`LC_ALL` set, `LC_MESSAGES` set, `LANG` set, none set, malformed values, `C`/`POSIX` values, bare `uk` without region), test `Translate()` returns correct value for `en` and `uk` keys, test fallback from `uk` to `en` for missing keys, test fallback to key name when key missing from all bundles, test `LoadBundle()` discovers both `en.toml` and `uk.toml`

**Checkpoint**: i18n package complete — locale detection works, all translations load, fallback chain verified

---

## Phase 3: User Story 1 — Help Output in Ukrainian (Priority: P1) 🎯 MVP

**Goal**: When system locale is `uk_UA.UTF-8`, all help text (command descriptions, flag descriptions, Cobra structural labels) is displayed in Ukrainian

**Independent Test**: Run `LANG=uk_UA.UTF-8 mytets help`, `mytets one --help`, `mytets list --help`, `mytets version --help` and verify all text is Ukrainian

### Tests for User Story 1

- [x] T007 [P] [US1] Write unit test in `internal/cli/run_test.go` to verify that when i18n is set to `uk`, `mytets help` output contains Ukrainian root description and structural labels from uk.toml
- [x] T008 [P] [US1] Write integration test in `tests/integration/locale_help_test.go` to verify `LANG=uk_UA.UTF-8 mytets help` prints Ukrainian root help, `mytets one --help` prints Ukrainian one-command help, `mytets list --help` prints Ukrainian list-command help, and `mytets version --help` prints Ukrainian version-command help

### Implementation for User Story 1

- [x] T009 [US1] Modify `internal/cli/run.go` to initialize i18n (call `i18n.LoadBundle()` and `i18n.DetectLocale()`) before building the command tree in both `Execute()` and `ExecuteArgs()`
- [x] T010 [US1] Modify `internal/cli/root.go` to use `i18n.Translate("root.short")` for root command `.Short`, set localized Cobra help/usage templates using `cmd.SetHelpTemplate()` and `cmd.SetUsageTemplate()` with translated structural labels (help.usage, help.available_commands, help.flags, help.global_flags, help.additional_help, help.aliases, help.help_for)
- [x] T011 [P] [US1] Modify `internal/cli/version_cmd.go` to use `i18n.Translate("version.short")` for `.Short` and `i18n.Translate("flag.output")` for the output flag description
- [x] T012 [P] [US1] Modify `internal/commands/one/one.go` to use `i18n.Translate("one.short")` for `.Short`, `i18n.Translate("one.long")` for `.Long`, and `i18n.Translate("flag.output")` for the output flag description
- [x] T013 [P] [US1] Modify `internal/commands/list/list.go` to use `i18n.Translate("list.short")` for `.Short`, `i18n.Translate("list.long")` for `.Long`, `i18n.Translate("flag.output")` for the output flag description, and `i18n.Translate("flag.count")` for the count flag description

**Checkpoint**: `LANG=uk_UA.UTF-8 mytets help` displays fully Ukrainian help. All subcommand `--help` displays Ukrainian descriptions and structural labels.

---

## Phase 4: User Story 2 — Error Messages in Ukrainian (Priority: P1)

**Goal**: When system locale is `uk_UA.UTF-8`, all application-defined and Cobra built-in error messages are displayed in Ukrainian

**Independent Test**: Trigger `no phrases available`, invalid output format, and unknown command errors with `LANG=uk_UA.UTF-8` and verify Ukrainian text on stderr

### Tests for User Story 2

- [x] T014 [P] [US2] Write unit test in `internal/commands/one/one_test.go` to verify that when i18n is set to `uk`, the "no phrases available" error is in Ukrainian
- [x] T015 [P] [US2] Write unit test in `internal/commands/list/list_test.go` to verify that when i18n is set to `uk`, the "no phrases available" error is in Ukrainian
- [x] T016 [P] [US2] Write unit test in `internal/flags/parser_test.go` to verify that when i18n is set to `uk`, the "unsupported output format" error is in Ukrainian

### Implementation for User Story 2

- [x] T017 [US2] Modify `internal/flags/parser.go` to use `i18n.Translate("error.invalid_output_format")` (with `fmt.Sprintf` for the format argument) in the `ParseOutputFormat` error path
- [x] T018 [P] [US2] Modify `internal/commands/one/one.go` to use `i18n.Translate("error.no_phrases")` in the error returned when `randomMessage()` fails, and `i18n.Translate("error.failed_select_phrase")` for the wrap message, and `i18n.Translate("error.failed_encode_json")` for JSON encoding errors
- [x] T019 [P] [US2] Modify `internal/commands/list/list.go` to use `i18n.Translate("error.no_phrases")` in the error returned when `messageSource()` returns empty
- [x] T020 [P] [US2] Modify `internal/cli/version_cmd.go` to use `i18n.Translate("error.failed_encode_json")` for the JSON encoding error path
- [x] T021 [US2] Modify `internal/cli/root.go` to set `cmd.SetFlagErrorFunc()` on the root command to return localized flag error messages, and intercept unknown-command errors from `root.Execute()` in `internal/cli/run.go` to translate them using `i18n.Translate("error.unknown_command")`
- [x] T022 [P] [US2] Write integration test in `tests/integration/locale_error_test.go` to verify `LANG=uk_UA.UTF-8 mytets --output xml one` prints Ukrainian error on stderr, and `LANG=uk_UA.UTF-8 mytets foo` prints Ukrainian unknown-command error

**Checkpoint**: All error messages display in Ukrainian when `LANG=uk_UA.UTF-8`. Error exit codes remain unchanged.

---

## Phase 5: User Story 3 — Default English Fallback (Priority: P2)

**Goal**: When system locale is not Ukrainian (e.g., `en_US.UTF-8`, `de_DE.UTF-8`, or unset), all messages appear in English identical to current behavior

**Independent Test**: Run `LANG=en_US.UTF-8 mytets help` and verify English output matches pre-localization behavior

### Tests for User Story 3

- [x] T023 [P] [US3] Write unit test in `internal/cli/run_test.go` to verify that when i18n is set to `en`, `mytets help` output matches current English help text
- [x] T024 [P] [US3] Write unit test in `internal/i18n/i18n_test.go` to verify that when locale is `de_DE.UTF-8` (unsupported), `DetectLocale()` resolves to `en` fallback

### Implementation for User Story 3

- [x] T025 [US3] Verify in `internal/i18n/i18n.go` that `DetectLocale()` falls back to `en` when the detected language code has no matching TOML file, and when env vars are empty or set to `C`/`POSIX`
- [x] T026 [US3] Run all existing unit tests (`go test ./internal/...`) and all existing integration tests (`go test ./tests/integration/...`) to confirm backward compatibility — no test modifications should be needed

**Checkpoint**: English-locale users see identical behavior to pre-feature state. All existing tests pass.

---

## Phase 6: User Story 4 — Phrase Content Not Localized (Priority: P1)

**Goal**: Phrase content from `phrases.json` is displayed exactly as stored, unaffected by locale

**Independent Test**: Run `LANG=uk_UA.UTF-8 mytets one` and `mytets list` and verify phrase text matches `phrases.json` verbatim

### Tests for User Story 4

- [x] T027 [P] [US4] Write integration test in `tests/integration/one_command_test.go` (or extend existing) to verify that with `LANG=uk_UA.UTF-8`, `mytets one` output contains text from `phrases.json` unchanged
- [x] T028 [P] [US4] Write integration test in `tests/integration/list_command_test.go` (or extend existing) to verify that with `LANG=uk_UA.UTF-8`, `mytets list` output contains phrases from `phrases.json` unchanged

### Implementation for User Story 4

- [x] T029 [US4] Verify that `internal/phrases/phrases.go` has NO imports of the `i18n` package and NO calls to `Translate()` — the phrases package must remain completely independent of the localization system

**Checkpoint**: `phrases.json` content passes through unmodified regardless of active locale.

---

## Phase 7: User Story 5 — Adding New Language Without Code Changes (Priority: P2)

**Goal**: A contributor can add a new language by creating one TOML file and rebuilding — zero Go source changes

**Independent Test**: Create a test `de.toml`, rebuild, set `LANG=de_DE.UTF-8`, verify German messages appear

### Tests for User Story 5

- [x] T030 [US5] Write unit test in `internal/i18n/i18n_test.go` to verify that `LoadBundle()` auto-discovers all `.toml` files in the embedded `locales/` directory (test that the count of loaded languages matches the number of `.toml` files)

### Implementation for User Story 5

- [x] T031 [US5] Verify in `internal/i18n/i18n.go` that `LoadBundle()` uses `fs.ReadDir()` to iterate over `embed.FS` entries and derives language codes from filenames (e.g., `uk.toml` → `uk`) — no hardcoded language list
- [x] T032 [US5] Verify that the `//go:embed locales/*.toml` directive uses a glob pattern so new `.toml` files are automatically included at build time without code changes

**Checkpoint**: Architecture verified — adding `de.toml` to `internal/i18n/locales/` and running `go build` is sufficient.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Final validation, backward compatibility confirmation, and cleanup

- [x] T033 Run full test suite with race detection: `go test -race ./...`
- [x] T034 Run `go vet ./...`, `staticcheck ./...`, and `golangci-lint run` and verify zero warnings
- [x] T035 Run `go test -coverprofile=cover.out ./internal/...` and verify test coverage is at or above 80% for `internal/` packages
- [x] T036 Update `README.md` with `github.com/BurntSushi/toml` dependency justification (required by constitution: external dependencies MUST be documented in README.md)
- [x] T037 Build the binary (`go build -o mytets ./cmd/mytets`) and run quickstart.md validation: `LANG=uk_UA.UTF-8 ./mytets help` produces Ukrainian output, `LANG=en_US.UTF-8 ./mytets help` produces English output, `LANG=uk_UA.UTF-8 ./mytets one` produces untranslated phrase content

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup (Phase 1) completion — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational (Phase 2) — help text localization
- **User Story 2 (Phase 4)**: Depends on Foundational (Phase 2) — error message localization; can run in parallel with US1
- **User Story 3 (Phase 5)**: Depends on US1 and US2 completion — validates English fallback works after localization is wired up
- **User Story 4 (Phase 6)**: Depends on Foundational (Phase 2) — can run in parallel with US1/US2
- **User Story 5 (Phase 7)**: Depends on Foundational (Phase 2) — verifies architecture; can run in parallel with US1/US2
- **Polish (Phase 8)**: Depends on all user stories being complete

### User Story Dependencies

- **US1 (P1)**: Depends only on Phase 2 — no dependencies on other stories
- **US2 (P1)**: Depends only on Phase 2 — no dependencies on other stories (can parallelize with US1)
- **US3 (P2)**: Depends on US1 + US2 (must validate fallback after localization is wired)
- **US4 (P1)**: Depends only on Phase 2 — verification task, no code changes
- **US5 (P2)**: Depends only on Phase 2 — architecture verification

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Template/structural changes before command-level changes
- Root command before subcommands
- Story complete before moving to next priority

### Parallel Opportunities

- T003 and T004 (TOML files) can be created in parallel
- T007 and T008 (US1 tests) can be written in parallel
- T011, T012, T013 (subcommand modifications) can run in parallel after T010
- T014, T015, T016 (US2 tests) can be written in parallel
- T018, T019, T020 (US2 subcommand error changes) can run in parallel
- T027 and T028 (US4 integration tests) can run in parallel
- US1, US2, US4, and US5 can all proceed in parallel after Phase 2

---

## Parallel Example: User Story 1

```bash
# Write tests first (parallel):
T007: "Unit test for Ukrainian help output in internal/cli/run_test.go"
T008: "Integration test for Ukrainian help in tests/integration/locale_help_test.go"

# Then wire up i18n init (sequential — root first):
T009: "Initialize i18n in internal/cli/run.go"
T010: "Localize root command and Cobra templates in internal/cli/root.go"

# Then localize subcommands (parallel — independent files):
T011: "Localize version command in internal/cli/version_cmd.go"
T012: "Localize one command in internal/commands/one/one.go"
T013: "Localize list command in internal/commands/list/list.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (add TOML dependency)
2. Complete Phase 2: Foundational (i18n package + TOML files + unit tests)
3. Complete Phase 3: User Story 1 (localized help text)
4. **STOP and VALIDATE**: `LANG=uk_UA.UTF-8 mytets help` shows Ukrainian
5. Then proceed to US2 (errors), US3 (fallback), US4 (phrases), US5 (extensibility)

### Incremental Delivery

| Increment | Stories | Value Delivered |
|-----------|---------|-----------------|
| MVP | US1 | Ukrainian help text — primary user-facing deliverable |
| +Errors | US1 + US2 | Complete Ukrainian UI (help + errors) |
| +Validation | US1–US4 | Full feature with backward compatibility confirmed |
| Complete | US1–US5 | Architecture validated for future language additions |
