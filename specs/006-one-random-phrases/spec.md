# Feature Specification: One Command - Random Phrase Output

**Feature Branch**: `006-one-random-phrases`  
**Created**: 2026-04-21  
**Status**: Draft  
**Input**: User description: "the 'one' command should return a random phrase from the file embedded in the application"

## Clarifications

### Session 2026-04-21

- Q: Should the JSON output flag be `mytets one --json` or `mytets --output json one`? → A: Use global `--output json` flag before the command (Option B)
- Q: What random selection strategy should be used for phrase selection? → A: Standard random selection using Go's math/rand (Option A)
- Q: How should errors be handled and reported? → A: Show clear error message to stderr and exit with code 1 (Option A)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Get Random Phrase as Plain Text (Priority: P1)

A user runs the `mytets one` command without any flags and expects to receive a single random phrase displayed as plain text in the terminal.

**Why this priority**: This is the core functionality of the feature - the primary user interaction. It delivers the basic value of selecting and displaying a random phrase.

**Independent Test**: Can be fully tested by running `mytets one` from the command line and verifying that a valid message from the embedded phrases.json file is printed as plain text to stdout, and the exit code is 0.

**Acceptance Scenarios**:

1. **Given** the application has an embedded phrases.json file with multiple messages, **When** the user runs `mytets one`, **Then** a single phrase (plain text) is printed to stdout
2. **Given** the user runs `mytets one` multiple times, **When** each execution completes with exit code 0, **Then** the phrase may be the same or different (randomness is working)
3. **Given** the user runs `mytets one`, **When** the phrase is printed, **Then** no JSON formatting or extra markup is included - just the text

---

### User Story 2 - Get Random Phrase as JSON (Priority: P1)

A user runs the `mytets --output json one` command and expects to receive the random phrase wrapped in a JSON object with a "message" field.

**Why this priority**: JSON output is explicitly required as an alternative output format. It's equally critical as plain text for supporting different use cases (programmatic consumption, integration with other tools).

**Independent Test**: Can be fully tested by running `mytets --output json one` and verifying the output is valid, compact JSON containing a "message" field with the phrase text, with exit code 0.

**Acceptance Scenarios**:

1. **Given** the application has an embedded phrases.json file with multiple messages, **When** the user runs `mytets --output json one`, **Then** a JSON object `{"message":"<phrase text>"}` is printed to stdout
2. **Given** the output is in JSON format, **When** the JSON is parsed, **Then** it is valid, well-formed JSON
3. **Given** the `--output json` global flag is used before the command, **When** output is produced, **Then** it is compact (not pretty-printed with extra whitespace or newlines)
4. **Given** a valid phrase is selected, **When** it is included in the JSON output, **Then** the exit code is 0

---

### Edge Cases

- **Empty or malformed phrases.json**: System MUST output an error message to stderr and exit with code 1
- **File not found**: System MUST output an error message to stderr and exit with code 1
- **Unknown flags**: Command should reject unknown flags and display usage error message to stderr with exit code 1
- **Empty messages array**: System MUST output an error message to stderr and exit with code 1
- **Random distribution**: All phrases MUST have equal probability of selection across multiple runs (uniform random distribution)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The application MUST embed a phrases.json file during the build process
- **FR-002**: The phrases.json file MUST contain a "messages" array with objects containing a "text" field for each phrase
- **FR-003**: The `mytets one` command MUST read the embedded phrases.json file
- **FR-004**: The `mytets one` command MUST select a phrase randomly from the messages array
- **FR-005**: By default (without flags), the command MUST output the selected phrase as plain text to stdout with no additional formatting
- **FR-006**: When invoked with the global `--output json` flag, the command MUST output a JSON object with the structure `{"message":"<phrase text>"}` to stdout
- **FR-007**: The JSON output MUST be compact (single line, no pretty-printing)
- **FR-008**: The command MUST exit with code 0 on successful execution
- **FR-009**: The command MUST output a descriptive error message to stderr and exit with code 1 if the embedded file is missing or cannot be read
- **FR-010**: The command MUST output a descriptive error message to stderr and exit with code 1 if no valid phrases are available

### Key Entities

- **Phrase**: A text message with a "text" field, stored as part of the "messages" array in phrases.json
- **phrases.json**: An embedded JSON file containing the structure `{ "messages": [{ "text": "..." }, ...] }`

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can run `mytets one` and receive a random phrase from the embedded file in under 100ms
- **SC-002**: Users can run `mytets --output json one` and receive valid, compact JSON output in under 100ms
- **SC-003**: 100% of valid phrases in phrases.json are selectable (uniform random distribution)
- **SC-004**: All unit tests for internal package functions pass without errors
- **SC-005**: All integration tests pass, confirming correct behavior for both plain text and JSON output modes
- **SC-006**: Command returns exit code 0 on success and non-zero on failure in 100% of test cases

## Assumptions

- The embedded phrases.json file will always have at least one valid message in the "messages" array
- The phrase selection uses Go's standard math/rand library for random selection
- The `--output json` global flag is already available (from existing flag parsing infrastructure)
- The application follows existing CLI conventions and patterns established by other commands (version, help)
- No special characters or escape sequences are required for phrase text in plain text output
- JSON output should use standard JSON encoding without custom serialization
