# CLI Contract: `mytets version`

## Command

- Invocation: `mytets version`
- Purpose: Print the application version and exit.

## Input Contract

- Positional args: none accepted (extra positional args yield non-zero exit).
- Supported flags:
  - `--output <format>` / `-o <format>`: output format, one of `text` (default) or `json`.
- Unexpected/invalid flags: standard CLI parser error behavior, non-zero exit.

## Output Contract (Success — plain text, default)

- Stream: stdout only.
- Format: single plain string and trailing newline.
- Value:
  - Semantic version `X.Y.Z` when injected at build time.
  - `dev` when no build-time version was injected.
- Example:

```text
1.0.1
```

## Output Contract (Success — JSON, `--output json` / `-o json`)

- Stream: stdout only.
- Format: JSON object, single line, trailing newline.
- Fields: `version` (string).
- Example:

```json
{"version":"1.0.1"}
```

## Error Contract

- For invalid invocation (unknown flags/args, unsupported output format), stderr
  contains a human-readable error message and exit code is non-zero.

## Exit Codes

- `0`: success (`mytets version` valid invocation)
- non-zero: parser/invocation error or unsupported output format

## Build Contract

- Version must be injectable with ldflags:

```bash
go build -ldflags "-X github.com/igorzel/mytets/internal/version.Version=1.0.1" ./cmd/mytets
```

## Notes

- Plain-text output is the default; `--output json` / `-o json` enables the
  JSON envelope for automation and CI usage (FR-010).
- Implementation verified: all integration tests pass including JSON, unsupported
  format, ldflags injection, and performance (<100ms) scenarios.
