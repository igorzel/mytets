# Data Model: One Command - Random Phrase Feature

**Feature**: `mytets one` command returning random phrases  
**Date**: 2026-04-21  
**Domain**: CLI phrase display with multi-format output  

## Entities

### 1. Phrase

**Definition**: A single text message from the embedded phrases file  
**Purpose**: The atomic unit of data returned by the `one` command  
**Lifecycle**: Immutable; loaded once at initialization from embedded file

**Attributes**:
- `text` (string, required): The message text to display
  - Constraints: Non-empty, max length TBD (reasonable for terminal display, e.g., <500 chars)
  - Encoding: UTF-8
  - Special characters: No escaping required; standard JSON string format

**Relationships**:
- Belongs to `PhrasesFile` (container of all phrases)
- Selected by `RandomSelector` (randomization service)
- Output by `OutputFormatter` (display service)

**Example**:
```json
{
    "text": "Example message 1"
}
```

---

### 2. PhrasesFile (Embedded)

**Definition**: The embedded JSON file containing all available phrases  
**Purpose**: Single source of truth for phrase data; compiled into binary  
**Storage**: `embedded/phrases.json` (file path at build time; becomes string in binary at runtime)  
**Lifecycle**: Immutable after build; loaded once at package initialization

**Structure**:
```json
{
    "messages": [
        {
            "text": "Example message 1"
        },
        {
            "text": "Example message 2"
        },
        {
            "text": "Example message 3"
        }
    ]
}
```

**Schema**:
- `messages` (array, required): Collection of phrase objects
  - Constraints: Non-empty array (enforced at init; panic if empty)
  - Min length: 1 phrase
  - Max length: No hard limit (practical: 1000s of phrases)

**Attributes**:
- Each element in `messages` is a `Phrase` object
- All phrases must have non-empty `text` field
- JSON must be valid and parseable at init time

**Validation Rules**:
- ✅ File must exist (embedded at compile-time)
- ✅ Must be valid JSON
- ✅ Must have `messages` array
- ✅ Array must not be empty
- ✅ Each message must have `text` field

**Error Handling**:
- If validation fails at init → panic (fails fast, prevents silent errors)
- If validation fails at command execution → error returned to stderr with exit code 1

**Example**:
```json
{
    "messages": [
        {"text": "Hello, World!"},
        {"text": "Welcome to mytets"},
        {"text": "Random phrase selection"}
    ]
}
```

---

### 3. CommandOutput (Polymorphic)

**Definition**: The formatted result of the `one` command  
**Purpose**: Represents output in different formats depending on user flags  
**Variants**: `PlainTextOutput`, `JSONOutput`

#### 3.1 PlainTextOutput

**Representation**: Single phrase as plain text  
**Format**: Raw text string, no JSON markup  
**Encoding**: UTF-8, single line  
**Example**:
```
Example message 1
```

**Constraints**:
- ✅ No JSON formatting
- ✅ No extra whitespace or markup
- ✅ Single line (phrase + newline)
- ✅ Exit code: 0 on success, 1 on error

#### 3.2 JSONOutput

**Representation**: Phrase wrapped in JSON object  
**Format**: Compact JSON (no pretty-printing)  
**Encoding**: UTF-8, single line  
**Structure**:
```json
{"message":"Example message 1"}
```

**Schema**:
- Root: object (single JSON object, not an array)
- Key: `"message"` (string literal, required)
- Value: phrase text (string)

**Constraints**:
- ✅ Valid, well-formed JSON
- ✅ Compact (no extra whitespace, no newlines inside JSON)
- ✅ Must include `"message"` key
- ✅ Single line in output
- ✅ Exit code: 0 on success, 1 on error

**Example**:
```json
{"message":"Welcome to mytets"}
```

---

## Entity Relationships

```
PhrasesFile (embedded JSON)
    ├── contains
    └── Phrase[] (array of 1..n phrases)
        └── selected by RandomSelector
            └── returns single Phrase
                └── formatted by OutputFormatter
                    ├── PlainTextOutput (default)
                    └── JSONOutput (via --output json)
```

---

## Validation Rules

### At Build Time
- `embedded/phrases.json` exists and is valid JSON
- Array is not empty
- Each phrase has non-empty `text` field

### At Runtime (Command Execution)
1. **Package Initialization** (`init()` block):
   - Parse embedded JSON (panic if invalid — fail fast)
   - Verify non-empty messages array (panic if empty)
   - Cache in package-level variable

2. **Command Execution** (`RunE`):
   - Retrieve all valid phrases from cache
   - Select one randomly using `math/rand.Intn()`
   - Format output based on `--output` flag
   - Print to stdout
   - Return nil (exit code 0) on success
   - Return error (exit code 1) on flag parsing failure

### Error Cases
- Invalid `--output` format value → error to stderr, exit code 1
- Missing/malformed embedded file → panic at startup (caught by init validation)
- Empty messages array → panic at startup (caught by init validation)

---

## Storage & Serialization

### Embedded File Path
**Location**: `embedded/phrases.json`  
**Format**: UTF-8 JSON  
**Built into**: Single binary via `//go:embed` directive  
**Access**: Package-level string variable loaded at `package init()`

### JSON Unmarshaling
**Library**: `encoding/json` (Go standard library)  
**Structure**: Typed struct matching JSON schema  
**Timing**: Once at package initialization, cached  
**Error Handling**: Panic on unmarshal error (fail-fast prevents silent corruption)

---

## Constraints & Limits

| Constraint | Value | Reason |
|-----------|-------|--------|
| Min phrases | 1 | Enforced; command needs ≥1 phrase to return |
| Max phrases | TBD (practical: 1000s) | File size limit; JSON parse time must stay <1ms |
| Max phrase length | TBD (reasonable: <500 chars) | Terminal display; long text is unusual for fortune-style tool |
| Execution time | <100ms | Per constitution principle IV |
| JSON output lines | 1 | Compact format; no pretty-printing |
| Plain text output lines | 1 | Single phrase + newline |
| File encoding | UTF-8 | Standard for JSON |
| Exit code success | 0 | Unix convention |
| Exit code error | 1 | Unix convention |

---

## Notes for Implementation

1. **Package Organization**: Consider `internal/phrases/` package for all phrase-related logic:
   - Loading/unmarshaling
   - Random selection
   - Message retrieval
   - Validation

2. **Error Messages** (stderr on exit code 1):
   - "invalid output format: expected 'text' or 'json', got '<value>'"
   - "failed to select phrase: <detailed error>"
   - Any panic message (should be rare/impossible at runtime)

3. **Performance**: All operations expected <10ms:
   - Package init: <1ms
   - Random selection: <1μs
   - JSON marshaling: <1ms
   - Output formatting: <1ms

4. **Testing Strategy**:
   - Unit: Verify phrase loading, random selection, formatting
   - Integration: Verify stdout/stderr/exit codes via CLI
   - Randomness: Verify uniform distribution across multiple runs

---

## Glossary

- **Phrase**: Individual text message from the embedded file
- **PhrasesFile**: The JSON file containing all phrases
- **OutputFormatter**: Logic that converts a phrase into display format
- **PlainTextOutput**: Phrase displayed as raw text
- **JSONOutput**: Phrase wrapped in `{"message":"..."}` JSON object
- **RandomSelector**: Logic that picks one phrase uniformly from all phrases
- **Exit Code**: Integer returned by command (0 = success, 1 = error)
