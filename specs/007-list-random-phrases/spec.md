# Feature Specification: List Command - Random Phrase List

**Feature Branch**: `007-list-random-phrases`  
**Created**: 2026-04-21  
**Status**: Draft  
**Input**: User description: "the \"list\" command should print a list of random phrases from the file embedded in the application"

## Clarifications

### Session 2026-04-21

- Q: Which flag order defines JSON output for the `list` command? → A: Use the existing global form `mytets --output json list`.
- Q: How should random phrase selection enforce uniqueness? → A: Select phrases without replacement for each command invocation so the same phrase text does not appear twice in one result.
- Q: What does phrase-source validation "on startup" mean for this feature? → A: Application initialization validates that the shared phrase source is available and non-empty before phrase-based commands run; if validation fails, `one` and `list` both report an error and do not succeed.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Get a Random Phrase List as Plain Text (Priority: P1)

A user runs the `mytets list` command and expects a plain-text list of random phrases printed to the terminal, one phrase per line.

**Why this priority**: This is the primary feature behavior. It delivers the main user value by returning multiple random phrases in a terminal-friendly format.

**Independent Test**: Can be fully tested by running `mytets list` and `mytets list --count N`, then verifying that stdout contains one phrase per line, the default list size is 5, the requested count is respected when possible, each phrase is unique within the result, and the exit code is 0.

**Acceptance Scenarios**:

1. **Given** the application has an embedded phrases.json file with at least 5 phrases, **When** the user runs `mytets list`, **Then** 5 unique phrases are printed to stdout as plain text, one per line
2. **Given** the application has an embedded phrases.json file with multiple phrases, **When** the user runs `mytets list --count 4`, **Then** 4 unique phrases are printed to stdout as plain text, one per line
3. **Given** the requested count is larger than the number of available phrases, **When** the user runs `mytets list --count N`, **Then** all available phrases are printed once each and no phrase is repeated
4. **Given** the command finishes successfully, **When** output is produced, **Then** the process exits with code 0

---

### User Story 2 - Get a Random Phrase List as JSON (Priority: P1)

A user requests machine-readable output and expects the `list` command to return the selected phrases as a compact JSON array of message objects.

**Why this priority**: JSON output is a stated requirement and enables scripting and downstream automation without changing the underlying command behavior.

**Independent Test**: Can be fully tested by running the command with the global `--output json` flag and verifying that stdout contains valid compact JSON, each array item has a `message` field, the selected phrases remain unique, and the exit code is 0.

**Acceptance Scenarios**:

1. **Given** the application has an embedded phrases.json file with multiple phrases, **When** the user runs `mytets --output json list`, **Then** stdout contains a compact JSON array of 5 objects with a `message` field
2. **Given** the user requests a specific count in JSON mode, **When** the user runs `mytets --output json list --count 2`, **Then** stdout contains a compact JSON array with 2 unique message objects
3. **Given** the requested count is larger than the number of available phrases, **When** the user runs the command in JSON mode, **Then** the JSON array contains every available phrase at most once and may contain fewer items than requested

---

### User Story 3 - Fail Fast When Phrase Data Is Unavailable (Priority: P2)

A user starts the application when the embedded phrase source is missing or empty and expects the application to report the problem clearly so phrase-based commands do not run with invalid data.

**Why this priority**: Phrase selection depends entirely on the embedded data source. Detecting invalid startup state prevents confusing empty output and keeps command behavior predictable.

**Independent Test**: Can be fully tested by starting the application with a missing or empty phrase source and verifying that an error is written to stderr during startup and that phrase-based commands such as `one` and `list` do not succeed.

**Acceptance Scenarios**:

1. **Given** the embedded phrases.json file is missing, **When** the application starts, **Then** an error is printed to stderr and phrase-based commands do not run successfully
2. **Given** the embedded phrases.json file is present but contains no phrases, **When** the application starts, **Then** an error is printed to stderr and phrase-based commands do not run successfully

---

### Edge Cases

- The requested count is greater than the number of available phrases
- The requested count is exactly equal to the number of available phrases
- The phrase source contains duplicate entries; the returned list still must not repeat the same phrase text within a single command result
- The user provides an invalid `--count` value; the command should fail consistently with existing CLI argument validation rules
- The phrase source is missing or contains zero usable phrases at startup

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The application MUST reuse the existing embedded phrases.json file as the source for the `list` functionality
- **FR-002**: The application MUST provide a `list` command that selects random phrases from the embedded phrase source and writes them to stdout
- **FR-003**: When the user runs `mytets list` without a count flag, the command MUST return 5 phrases by default
- **FR-004**: The `list` command MUST support a command-specific `--count N` flag that changes the requested number of phrases returned
- **FR-005**: The `--count` flag MUST apply only to the `list` command and MUST NOT become a global flag
- **FR-006**: Each command result produced by `list` MUST contain unique phrases only; phrase selection for a single invocation MUST be performed without replacement so the same phrase text does not appear more than once in a single result
- **FR-007**: When the requested count is greater than the number of available phrases, the command MUST return all available phrases without duplication, even if that produces fewer items than requested
- **FR-008**: By default, the `list` command MUST output the selected phrases as plain text with one phrase per line and no additional markup
- **FR-009**: When the global `--output json` flag is set to `json`, the `list` command MUST output a compact JSON array using the existing invocation form `mytets --output json list`, where each item has the structure `{ "message": "<phrase text>" }`
- **FR-010**: Successful execution of the `list` command MUST exit with code 0
- **FR-011**: During application initialization, if the shared phrase source is missing or cannot be loaded, the application MUST print an error to stderr and phrase-based commands such as `one` and `list` MUST NOT work
- **FR-012**: During application initialization, if the shared phrase source contains zero usable phrases, the application MUST print an error to stderr and phrase-based commands such as `one` and `list` MUST NOT work
- **FR-013**: The feature MUST include unit tests that verify internal phrase-loading, selection, counting, uniqueness, and output-formatting behavior
- **FR-014**: The feature MUST include integration tests that verify `mytets list` produces the expected plain-text and JSON outputs

### Key Entities

- **Phrase**: A single message available for selection and display
- **Phrase Collection**: The full set of phrases loaded from the embedded phrases.json source
- **List Result Item**: A returned phrase represented either as a plain-text line or as an object with a `message` field in JSON output

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can run the `list` command with default settings and receive a 5-item phrase list in under 1 second
- **SC-002**: Users can request a custom list size and receive up to the requested number of phrases in under 1 second
- **SC-003**: 100% of phrases within a single command result are unique
- **SC-004**: 100% of successful plain-text and JSON invocations of the `list` command exit with code 0
- **SC-005**: 100% of startup failures caused by a missing or empty phrase source display an error on stderr and prevent phrase-based commands from succeeding
- **SC-006**: All defined unit and integration acceptance tests for the `list` feature pass

## Assumptions

- The existing global output-format option already supports `json` and remains the mechanism for switching `list` output between plain text and JSON
- Phrase-based commands share the same phrase source, so startup validation for missing or empty phrase data applies consistently to both `one` and `list`
- The phrase source is expected to contain enough valid entries for normal operation, but the command may return fewer items than requested when the source contains fewer unique phrases
- Invalid `--count` inputs continue to follow the application’s existing command-line validation behavior