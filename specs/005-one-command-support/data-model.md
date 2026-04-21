# Data Model: One Command

## Entities

### OneCommandResponse

Represents the output payload of the `one` command.

**Fields**:

- `message` (string): The fixed message text: "Fake message, tbd"

**Representation**:

- **Plain Text Mode**: String only: `"Fake message, tbd"`
- **JSON Mode**: Compact JSON object: `{"message":"Fake message, tbd"}`

**Constraints**:

- Message is immutable (hardcoded)
- Message is always returned with exit code 0 on success
- No additional fields or metadata is included
- JSON output must be compact (single line, no whitespace)

## State Transitions

The `one` command is stateless. On every invocation:

1. Command is invoked (`mytets one` or `mytets --json one`)
2. Check if global JSON mode is enabled (parser configuration)
3. Output message in appropriate format (plain text or JSON)
4. Exit with code 0

No state persistence, caching, or side effects occur.

## Storage & Persistence

N/A — Command is pure function; message is hardcoded; no state.

## Validation Rules

- **Input**: None (command accepts no arguments or subcommand-specific flags)
- **Output (plain text mode)**: Must be exactly `Fake message, tbd\n` (with trailing newline)
- **Output (JSON mode)**: Must be valid compact JSON: `{"message":"Fake message, tbd"}\n` (with trailing newline)

## Relationships to Other Entities

- **Depends on**: `flags.ParserConfig` (to determine if JSON mode is enabled)
- **Used by**: Root command registration in `internal/cli/root.go`
- **Tested by**: Unit tests in `internal/commands/one/one_test.go` and integration tests in `tests/integration/one_command_test.go`
