# CLI Contract: `mytets version`

## Command

- Invocation: `mytets version`
- Purpose: Print the application version and exit.

## Input Contract

- Positional args: none required.
- Supported flags for this feature: none.
- Unexpected/invalid flags: standard CLI parser error behavior, non-zero exit.

## Output Contract (Success)

- Stream: stdout only.
- Format: single plain string and trailing newline.
- Value:
  - Semantic version `X.Y.Z` when injected at build time.
  - `dev` when no build-time version was injected.
- Example:

```text
1.0.1
```

## Error Contract

- For invalid invocation (unknown flags/args), stderr contains parser error
  details and exit code is non-zero.

## Exit Codes

- `0`: success (`mytets version` valid invocation)
- non-zero: parser/invocation error

## Build Contract

- Version must be injectable with ldflags:

```bash
go build -ldflags "-X github.com/igorzel/mytets/internal/version.Version=1.0.1" ./cmd/mytets
```

## Notes

- This command is intentionally plain-text-only and does not implement JSON
  output.
