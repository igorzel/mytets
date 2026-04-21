# CLI Contract: mytets list command

**Feature**: `mytets list` command returning a random phrase list  
**Date**: 2026-04-21  
**Version**: 1.0  
**Status**: Design specification

## Command Invocation Syntax

### Default Plain-Text Output

```bash
mytets list
```

**Expected behavior**: prints 5 unique random phrases to stdout, one per line.

**Example**:
```text
Example message 3
Example message 1
Example message 4
Example message 2
Example message 5
```

**Exit code**: 0

### Plain-Text Output With Custom Count

```bash
mytets list --count 4
```

**Expected behavior**: prints 4 unique random phrases to stdout, one per line.

### JSON Output Via Existing Global Flag

```bash
mytets --output json list
```

**Expected behavior**: prints a compact JSON array of message objects to stdout.

**Example**:
```json
[{"message":"Example message 3"},{"message":"Example message 1"}]
```

### JSON Output With Custom Count

```bash
mytets --output json list --count 2
```

**Expected behavior**: prints a compact JSON array with exactly 2 unique message objects when at least 2 unique phrases are available.

## Flags

### Command-Specific `--count`

- Scope: `list` command only
- Type: integer
- Default: `5`
- Meaning: requested number of phrases to return

**Valid examples**:
```bash
mytets list --count 1
mytets list --count 5
mytets --output json list --count 3
```

**Required behavior**:
- Values larger than the number of unique available phrases return all available unique phrases.
- Invalid values fail with a CLI error and non-zero exit code.

### Global `--output` / `-o`

- Scope: global, already supported by the CLI
- Allowed values: `text`, `json`
- `list` reuses the existing leading global output parsing flow

## Output Specifications

### Plain Text

- One selected phrase per line
- No numbering, bullets, or extra labels
- Trailing newline allowed by standard CLI printing behavior

### JSON

- Root value is an array
- Each item is an object with exactly one field: `message`
- Output is compact, not pretty-printed

**Schema**:
```json
[
  {"message":"<phrase_text>"}
]
```

## Exit Codes

### Success

- Exit code `0`
- Stdout contains the selected list in the requested format
- Stderr is empty

### Error

- Non-zero exit code
- Stdout is empty
- Stderr contains a human-readable error message

## Error Scenarios

### Invalid Count

Examples:
```bash
mytets list --count 0
mytets list --count -1
mytets list --count invalid
```

**Expected behavior**: reject input using existing CLI validation conventions and return a non-zero exit code.

### Missing or Empty Phrase Source

**Expected behavior**: phrase-based commands fail consistently because shared phrase-source validation prevents successful execution.

### Unknown Flags

Examples:
```bash
mytets list --unknown
mytets --output xml list
```

**Expected behavior**: Cobra/global parser error with non-zero exit code.
