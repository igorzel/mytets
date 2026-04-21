# Feature Specification: One Command Support

**Feature Branch**: `005-one-command-support`
**Created**: 2026-04-21
**Status**: Draft
**Input**: User description: "Feature: add the \"one\" command support"

## Clarifications

### Session 2026-04-21

- Q: Does the CLI framework currently support a global `--json` flag? → A: Yes, it is an existing feature; the `one` command will respect it automatically.
- Q: Does the CLI parser support both `mytets --json one` and `mytets one --json`? → A: No, only `mytets --json one` (flag before subcommand) is supported.
- Q: How should unsupported flags be handled? → A: Per standard Cobra behavior, reject with error message and non-zero exit code.

### Implementation Notes

- Plain-text mode output consists of only the message value as a single line, excluding labels or prefixes.
- JSON mode output consists of a compact object with exactly one field: `message`.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Run `mytets one` for Plain Message (Priority: P1)

A user runs the `one` command and immediately sees the expected message in plain text so they can confirm the command behavior quickly.

**Why this priority**: Plain command execution is the core capability and the default expected behavior.

**Independent Test**: Run `mytets one`; verify stdout contains exactly `Fake message, tbd` (with standard trailing newline), stderr is empty, and process exits with code 0.

**Acceptance Scenarios**:

1. **Given** the user invokes `mytets one` without global output modifiers, **When** the command executes, **Then** stdout contains `Fake message, tbd` as plain text and the process exits with code 0.
2. **Given** the command executes successfully, **When** completion is observed, **Then** no error output is written to stderr.

---

### User Story 2 - Run `mytets --json one` for Structured Output (Priority: P2)

A user or script runs the `one` command with the global JSON flag and gets structured output to support automation.

**Why this priority**: Structured output is needed for scriptability but depends on the command existing first.

**Independent Test**: Run `mytets --json one`, parse stdout as JSON, verify object equals `{"message":"Fake message, tbd"}`, and confirm exit code 0.

**Acceptance Scenarios**:

1. **Given** the user invokes the command with global JSON mode enabled, **When** the command executes, **Then** stdout is valid compact JSON containing exactly one field `message` with value `Fake message, tbd`.
2. **Given** a script consumes the command output in JSON mode, **When** it parses stdout as JSON, **Then** parsing succeeds without requiring pretty-print handling.

---

### Edge Cases

- What happens if the command is invoked repeatedly in the same session? The output MUST remain deterministic and identical across runs for both plain and JSON modes.
- What happens if global JSON mode is enabled for other commands in the same invocation pattern? The `one` command MUST still return the same semantic message value in JSON format.
- What happens if consumers compare output exactly? The JSON response MUST be compact (single-line, non-pretty format) to avoid whitespace-dependent parsing issues.
- What happens if unsupported flags are passed to the `one` command (e.g., `mytets --json one --invalid-flag`)? The command MUST reject the flag, print an error message to stderr, and exit with non-zero code (standard Cobra behavior).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST expose a `one` subcommand invocable as `mytets one`.
- **FR-002**: By default, `mytets one` MUST print exactly `Fake message, tbd` to stdout as plain text (single message line, no labels or prefixes).
- **FR-003**: When global JSON output mode is enabled, the command MUST print valid JSON equivalent to `{\"message\":\"Fake message, tbd\"}`.
- **FR-004**: JSON output for this command MUST be compact (not pretty printed) and MUST include no fields other than `message`.
- **FR-005**: Successful execution in both plain and JSON modes (invocations: `mytets one`, `mytets --json one`) MUST exit with code 0.
- **FR-006**: Successful execution of `mytets one` MUST NOT write any output to stderr.
- **FR-007**: Unit tests MUST verify internal command behavior for plain and JSON output paths.
- **FR-008**: Integration tests MUST verify end-to-end CLI behavior for `mytets one` and JSON-enabled invocation, including stdout content and exit code.

### Key Entities *(include if feature involves data)*

- **One Command Response**: The user-visible message payload emitted by the `one` command, represented either as plain text or as a JSON object with a `message` field.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of successful `mytets one` invocations in automated integration tests return exit code 0.
- **SC-002**: 100% of default-mode test runs for `mytets one` produce exactly the message `Fake message, tbd` in stdout.
- **SC-003**: 100% of JSON-mode test runs for `mytets one` produce valid compact JSON that parses to a single `message` value of `Fake message, tbd`.
- **SC-004**: Unit and integration test suites include explicit coverage for both plain output and JSON output paths.

## Assumptions

- A global `--json` flag already exists in the CLI framework and is placed before the subcommand (e.g., `mytets --json one`). The `one` command will automatically respect this flag without implementing standalone flag parsing.
- The CLI framework is Cobra-based; unsupported flags are rejected with error messages and non-zero exit codes.
- The command has no required arguments; all error handling and unknown flag rejection follows Cobra conventions.
- Message text is fixed to `Fake message, tbd` for this feature version and is not user-configurable.
- Standard CLI behavior includes a trailing newline on stdout output in both plain and JSON modes.
