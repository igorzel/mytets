# Feature Specification: Version Command

**Feature Branch**: `003-version-command`
**Created**: 2026-04-21
**Status**: Draft
**Input**: User description: "Feature: version command support"

## Clarifications

### Session 2026-04-21

- Q: Where should the version variable live for ldflags injection? → A: `internal/version` package (`internal/version.Version`)
- Q: Should `mytets version` support `--output json`? → A: No, this command is plain-text only; constitution exception is documented for this command.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Developer Queries Application Version (Priority: P1)

A developer or operator runs `mytets version` to confirm which version of the
tool is installed on their machine. The command prints the version string and
exits immediately with a success code.

**Why this priority**: This is the sole purpose of the feature and directly
maps to the primary acceptance criterion. All other stories depend on this
working first.

**Independent Test**: Run `mytets version` from a terminal; verify exactly one
line is printed matching `X.Y.Z` and the process exits with code 0.

**Acceptance Scenarios**:

1. **Given** the binary is built with version `1.0.1` injected at link time,
   **When** the user runs `mytets version`,
   **Then** the tool prints `1.0.1` to stdout and exits with code 0.

2. **Given** the binary is built with version `0.9.0` injected at link time,
   **When** the user runs `mytets version`,
   **Then** the tool prints `0.9.0` to stdout and exits with code 0.

3. **Given** the user runs `mytets version` with no additional arguments,
   **Then** no error message is produced and stderr remains empty.

---

### User Story 2 — Scripted / CI Version Check (Priority: P2)

An automated script or CI pipeline captures the output of `mytets version` to
compare it against an expected release tag or embed it in a release artifact.

**Why this priority**: Machine-readable, predictable output is critical for
automation. The format must be stable and unambiguous.

**Independent Test**: Capture stdout of `mytets version` in a shell variable;
assert the string matches the regex `^[0-9]+\.[0-9]+\.[0-9]+$` and that the
exit code is 0.

**Acceptance Scenarios**:

1. **Given** a script runs `mytets version` and captures stdout,
   **When** the captured string is compared against `^[0-9]+\.[0-9]+\.[0-9]+$`,
   **Then** the match succeeds with no surrounding whitespace (beyond a trailing
   newline).

2. **Given** the version is captured in a CI pipeline,
   **When** the pipeline checks the exit code,
   **Then** the exit code is 0 allowing the pipeline to continue without error.

---

### Edge Cases

- What happens when `mytets version` is invoked with unexpected flags (e.g.,
  `mytets version --foo`)? → The tool MUST print an appropriate error to stderr
  and exit with a non-zero code (standard Cobra behaviour; no extra
  implementation required).
- What happens when a user expects JSON output for `mytets version`? → This
  command is explicitly plain-text-only and MUST NOT emit JSON; this is a
  documented command-level exception to the broader project automation-output
  guidance.
- What happens when the version string was not injected at build time (empty
  `ldflags`)? → A fallback sentinel value (e.g., `dev`) MUST be displayed so
  the command never exits with an error due to a missing version.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The tool MUST expose a `version` subcommand that is invocable as
  `mytets version`.
- **FR-002**: `mytets version` MUST print the application version to stdout in
  `X.Y.Z` format (e.g., `1.0.1`).
- **FR-003**: The output MUST be a single plain string — no labels, no
  decoration, no extra lines beyond the standard trailing newline.
- **FR-004**: The `version` subcommand MUST NOT require any additional arguments
  or flags to function.
- **FR-005**: The version value MUST be embedded in the binary at build time via
  linker flags; no runtime file access or environment variable is required for
  normal operation.
- **FR-006**: `mytets version` MUST exit with code 0 on every successful
  invocation.
- **FR-007**: A fallback version value (e.g., `dev`) MUST be used when no
  version is injected at build time, so the command always succeeds.
- **FR-008**: Unit tests MUST verify that the version output matches the
  injected value.
- **FR-009**: Integration tests MUST verify the full end-to-end CLI invocation,
  including exit code and stdout content.
- **FR-010**: `mytets version` MUST remain plain-text-only and MUST NOT emit
  JSON output for this feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: `mytets version` completes and exits in under 100 ms on any
  supported platform.
- **SC-002**: The output is exactly one line whose content matches the
  `X.Y.Z` pattern (or the `dev` fallback) — verified by automated tests.
- **SC-003**: Exit code is 0 for every invocation of `mytets version` without
  unexpected arguments — verified by integration tests on Linux, macOS, and
  Windows.
- **SC-004**: The displayed version matches the value injected at build time —
  verified by a unit test that controls the injected value.
- **SC-005**: No output is written to stderr during a successful `mytets version`
  invocation.

## Assumptions

- The existing project already uses `cobra` for command routing; the `version`
  subcommand will be registered as a new `cobra.Command` within the same CLI
  structure.
- The version string MUST be injected via linker flags targeting the variable
  `Version` in the `internal/version` package; the canonical ldflags value is
  `-X github.com/igorzel/mytets/internal/version.Version=<semver>`.
- The build system is responsible for supplying the correct `-ldflags` argument;
  the spec does not mandate a specific CI/CD toolchain.
- `mytets version` is an explicit exception to constitution-level optional JSON
  output guidance; this command returns only a plain version string.
- The `dev` fallback is sufficient for local development builds where no version
  is injected; no warning or error is needed for the fallback case.
- Mobile/browser output support is out of scope; only terminal stdout is
  targeted.
- The version command does not need to support `--help` beyond the default Cobra
  help that is generated automatically.
