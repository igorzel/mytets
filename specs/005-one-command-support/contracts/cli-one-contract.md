# CLI Contract: One Command

## Command Signature

```
mytets one [flags]
```

## Global Flags

Inherits all global flags from the root command. The relevant flag for this feature:

- `--json` (boolean, optional): Enables JSON output mode. Default: disabled (plain text mode).

## Subcommand-Specific Flags

None. The `one` command accepts no subcommand-specific flags.

## Invocation Examples

### Plain Text Mode (Default)

```bash
$ mytets one
Fake message, tbd
```

**Exit Code**: 0  
**Stdout**: `Fake message, tbd\n`  
**Stderr**: (empty)

### JSON Mode

```bash
$ mytets --json one
{"message":"Fake message, tbd"}
```

**Exit Code**: 0  
**Stdout**: `{"message":"Fake message, tbd"}\n`  
**Stderr**: (empty)

## Error Cases

### Invalid/Unsupported Flags

```bash
$ mytets --json one --invalid-flag
Error: unknown flag: --invalid-flag
```

**Exit Code**: Non-zero (Cobra default)  
**Stdout**: (empty)  
**Stderr**: Error message from Cobra

### Help Request

```bash
$ mytets one --help
Usage:
  mytets one [flags]

Short description...

Flags:
  -h, --help   help for one

Global Flags:
  --json   output in JSON format
```

**Exit Code**: 0  
**Stdout**: Help text  
**Stderr**: (empty)

## Output Format Specification

### Plain Text Mode

- Single line of output: the message text
- Trailing newline included
- No labels, prefixes, or decoration
- No color codes or ANSI escape sequences

**Example**:
```
Fake message, tbd
```

### JSON Mode

- Valid JSON object in compact format (no whitespace)
- Single field: `message`
- Trailing newline included
- Field order: `message` only

**Example**:
```
{"message":"Fake message, tbd"}
```

## Behavioral Guarantees

1. **Determinism**: Output is identical on every invocation with the same mode (plain or JSON)
2. **No Side Effects**: Command does not modify state, create files, or access external resources
3. **Exit Codes**:
   - **0** (success): All valid invocations (plain mode, JSON mode)
   - **Non-zero**: Invalid invocation (unsupported flags, syntax errors)
4. **Performance**: Command executes in under 100 ms (per constitution)
5. **Error Handling**: Invalid input produces an error message to stderr and a non-zero exit code (Cobra standard)

## Dependencies on Global Configuration

- **JSON Mode Flag**: The command respects the global `--json` flag passed before the subcommand (e.g., `mytets --json one`)
- **Parser Configuration**: The command receives a `flags.ParserConfig` object from the CLI that includes the JSON mode setting

## Integration Points

- **Root Command Registration**: The `one` command is registered in `internal/cli/root.go` using Cobra's `root.AddCommand()` method
- **Execution Path**: `cmd/mytets/main.go` → `internal/cli.Execute()` → root command → `one` subcommand
