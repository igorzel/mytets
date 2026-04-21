# CLI Contract: mytets one command

**Feature**: `mytets one` command returning random phrases  
**Date**: 2026-04-21  
**Version**: 1.0  
**Status**: Design specification  

## Command Invocation Syntax

### Basic Invocation (Plain Text Output)

```bash
mytets one
```

**Expected behavior**: Prints single phrase as plain text to stdout

**Output example**:
```
Example message 1
```

**Exit code**: 0 (success)

---

### With Global --output json Flag

```bash
mytets --output json one
```

**Alternative flag syntax**:
```bash
mytets -o json one
```

**Expected behavior**: Prints phrase wrapped in JSON to stdout

**Output example**:
```json
{"message":"Example message 1"}
```

**Exit code**: 0 (success)

---

## Output Format Specifications

### Format: Plain Text (default)

**Trigger**: No `--output` flag specified, or `--output text`

**Specification**:
- Single line of text
- Phrase content only (no markup)
- Encoded: UTF-8
- Trailing newline: Yes (standard for terminal output)
- Max length: TBD (practical: terminal width - 1, ~80-120 chars)

**Example**:
```
Example message 1
```

**Schema**: Plain text string
```
<phrase_text>\n
```

---

### Format: JSON

**Trigger**: `--output json` or `-o json` global flag

**Specification**:
- Single JSON object (not array, not multiple objects)
- Key: `"message"` (literal string)
- Value: phrase text (JSON string)
- Compact format (no pretty-printing, no whitespace inside JSON)
- Encoded: UTF-8
- Trailing newline: Yes (standard for terminal output)

**Structure**:
```json
{"message":"<phrase_text>"}
```

**Example**:
```json
{"message":"Example message 1"}
```

**Validation**:
- ✅ Must be valid JSON
- ✅ Must parse as object (not array)
- ✅ Must contain `"message"` key
- ✅ Value must be non-empty string
- ✅ No pretty-printing (single line)

**Why this structure?**
- Single object (not array) provides cleaner JSON contract
- `"message"` key is self-documenting
- Value is always a string (consistent type)
- Compact format is standard for CLI tools (for scripting, piping)

---

## Exit Codes

### Exit Code 0: Success

**When**: 
- Command executes successfully
- Phrase selected and output without errors

**Output**:
- Stdout: Phrase (plain text or JSON per format)
- Stderr: Empty

---

### Exit Code 1: Error

**When**:
- Invalid flag syntax
- Unsupported `--output` format value
- Any error during phrase selection/output formatting
- (Runtime errors should be rare if init validation is correct)

**Output**:
- Stdout: Empty (no partial output on error)
- Stderr: Human-readable error message

**Error Message Format**:
- Single line to stderr
- Starts with lowercase letter (convention)
- Begins with context or error type: "invalid output format", "failed to select phrase", etc.
- Example: `invalid output format: expected 'text' or 'json', got 'xml'`

---

## Flag Specifications

### Global --output / -o Flag

**Scope**: Global (applies to all commands)  
**Name**: `--output` (long form)  
**Short**: `-o` (short form)  
**Type**: String  
**Values**:
- `text` — Plain text format (default)
- `json` — JSON object format

**Default**: `text` (if flag omitted, plain text is output)

**Validation**:
- ✅ Must be one of the allowed values
- ✅ Case-sensitive (lowercase: `json`, not `JSON` or `Json`)
- ✅ Rejects unknown formats with error message

**Examples**:
```bash
mytets one                    # default, same as --output text
mytets -o json one           # short flag
mytets --output json one     # long flag
mytets --output text one     # explicit plain text
```

**Invalid usage** (should error):
```bash
mytets --output xml one      # unsupported format → error
mytets --output JSON one     # wrong case → error
mytets one --json            # --json not valid (only --output flag)
```

---

## Flag Combinations

### Positional Arguments

**Expected**: None  
**Behavior**: Reject any positional arguments with error

**Valid**:
```bash
mytets one
```

**Invalid** (should error):
```bash
mytets one "extra argument"   # extra positional arg
mytets one extra              # extra positional arg
```

---

### Multiple --output Flags

**Expected**: Only one `--output` value per invocation  
**Behavior**: If multiple specified, Cobra typically uses the last one (standard behavior)

**Example** (last wins):
```bash
mytets --output text --output json one   # Uses 'json'
```

---

### Unknown Flags

**Expected**: None  
**Behavior**: Cobra rejects unknown flags with "unknown flag" error

**Example**:
```bash
mytets --unknown one          # Error: unknown flag '--unknown'
mytets one --unknown          # Error: unknown flag '--unknown'
```

---

## Stderr Output Specifications

### Success Case
- **Content**: Empty
- **Lines**: 0

### Error Case
- **Content**: Human-readable error message
- **Lines**: 1 (single line)
- **Format**: Lowercase start, clear context
- **Examples**:
  - `invalid output format: expected 'text' or 'json', got 'xml'`
  - `failed to parse output format flag: <details>`

### Important Notes
- ✅ Errors ALWAYS go to stderr, never stdout
- ✅ Error messages do NOT include usage/help text (just the error itself)
- ✅ One error message per error (no multi-line errors unless necessary)

---

## Examples: Complete Scenarios

### Scenario 1: Plain Text Success

```bash
$ mytets one
Example message 1
```

**Stdout**: `Example message 1\n`  
**Stderr**: (empty)  
**Exit code**: 0  

---

### Scenario 2: JSON Success

```bash
$ mytets --output json one
{"message":"Example message 1"}
```

**Stdout**: `{"message":"Example message 1"}\n`  
**Stderr**: (empty)  
**Exit code**: 0  

---

### Scenario 3: Invalid Output Format

```bash
$ mytets --output xml one
```

**Stdout**: (empty)  
**Stderr**: `invalid output format: expected 'text' or 'json', got 'xml'\n`  
**Exit code**: 1  

---

### Scenario 4: Multiple Plain Text Runs (Randomness)

```bash
$ mytets one
Example message 1

$ mytets one
Example message 2

$ mytets one
Example message 1
```

**Behavior**: Different phrases returned on each invocation (demonstrating random selection)  
**All exit codes**: 0  

---

### Scenario 5: Multiple JSON Runs

```bash
$ mytets --output json one
{"message":"Example message 1"}

$ mytets --output json one
{"message":"Example message 3"}
```

**JSON output**: Always compact, always valid  
**All exit codes**: 0  

---

## Contract Validation Checklist

Implementation must satisfy:

- [ ] `mytets one` returns plain text phrase + newline, exit 0
- [ ] `mytets one` returns different phrases on repeated runs (randomness)
- [ ] `mytets --output json one` returns valid, compact JSON, exit 0
- [ ] JSON output has `{"message":"..."}` structure
- [ ] `mytets --output text one` returns plain text, exit 0
- [ ] `mytets -o json one` returns JSON, exit 0 (short flag works)
- [ ] `mytets --output xml one` returns error to stderr, exit 1
- [ ] `mytets one extra-arg` returns error, exit 1
- [ ] Error messages are single line on stderr
- [ ] Exit code 0 only on success, 1 only on error
- [ ] No stdout output on error
- [ ] Performance: <100ms total (init + random + output)

---

## Future Extensibility

While not implemented in v1, the contract supports:
- **Additional output formats** (e.g., `--output yaml`): Add to flag validation
- **Multiple phrases per run**: Change to return array in JSON or multi-line plain text
- **Phrase filtering**: Add future flags like `--filter <category>` (filter before random selection)
- **Phrase metadata**: Extend `Phrase` entity with optional fields (author, category, etc.)

Current contract is minimal and stable; future extensions layer on top without breaking existing behavior.
