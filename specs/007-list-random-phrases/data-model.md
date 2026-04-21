# Data Model: List Command - Random Phrase List

**Feature**: `mytets list` command with reusable list generation  
**Date**: 2026-04-21  
**Domain**: phrase selection, bounded list requests, CLI output adaptation

## Entities

### Phrase

**Definition**: A single message text available from the embedded phrase source.

**Fields**:
- `text` (string, required): the user-visible message value

**Rules**:
- Must be non-empty after trimming whitespace.
- Phrase text is the uniqueness key for one command result.

### Phrase Collection

**Definition**: The full set of phrases loaded from `internal/phrases/phrases.json`.

**Fields**:
- `messages` (ordered collection of Phrase): source data loaded at initialization

**Rules**:
- Must be available before phrase-based commands succeed.
- May contain duplicate source entries, but downstream list results must not repeat the same phrase text within one invocation.

### List Request

**Definition**: The reusable input passed to the list-generation package.

**Fields**:
- `count` (integer, required): requested number of phrases

**Rules**:
- Default value is 5 when the CLI flag is omitted.
- If `count` exceeds the number of unique available phrases, the result size is capped at the available unique count.
- Invalid count values are rejected by CLI validation before domain execution.

### List Result

**Definition**: The reusable domain output produced by the list-generation package.

**Fields**:
- `messages` (ordered collection of string): selected unique phrases in display order
- `requestedCount` (integer, optional metadata): original requested count used for traceable testing
- `returnedCount` (integer, derived): actual number of phrases returned

**Rules**:
- Contains only unique phrase texts.
- Contains at most one entry for each unique phrase text in the source.
- Contains up to the requested count, but may be smaller when fewer unique phrases exist.

### CLI JSON Item

**Definition**: The adapter-level representation of one selected phrase in JSON output.

**Fields**:
- `message` (string, required): one selected phrase text

**Rules**:
- Produced only by the CLI adapter when `--output json` is requested.
- Serialized as a compact JSON array item.

## Relationships

```text
Phrase Collection
    └── provides source phrases to List Request
List Request
    └── produces List Result
List Result
    ├── renders as plain-text lines in CLI output
    └── maps to CLI JSON Item[] for JSON output
```

## Validation Rules

### Phrase Source Validation
- Phrase collection load must fail if the embedded source is missing, malformed, or empty.
- Empty or whitespace-only phrase text is invalid.

### List Generation Validation
- Requested count is normalized by CLI defaults and validated before list generation.
- Duplicate source phrase texts are coalesced for per-invocation uniqueness.
- The result length is `min(requestedCount, uniqueAvailableCount)`.

### Output Validation
- Plain-text output prints one selected phrase per line.
- JSON output is a compact array of objects shaped as `{ "message": "..." }`.
