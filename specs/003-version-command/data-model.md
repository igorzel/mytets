# Data Model — Version Command

## Entity: VersionInfo

- Purpose: Represents version metadata rendered by the `version` command.
- Fields:
  - `value` (string, required): semantic version string (`X.Y.Z`) or fallback
    `dev` when not injected.
  - `source` (enum: `ldflags` | `fallback`, derived): indicates whether value
    came from linker injection or default.
- Validation rules:
  - If `value != "dev"`, it SHOULD match regex `^[0-9]+\.[0-9]+\.[0-9]+$`.
  - Empty string is invalid at render time and MUST be normalized to `dev`.
- State transitions:
  - `Uninitialized` -> `Resolved`: on command execution, load package variable.
  - `Resolved` -> `Rendered`: print value to stdout and exit 0.

## Entity: VersionCommandInvocation

- Purpose: Captures invocation constraints for `mytets version` behavior.
- Fields:
  - `command` (string, constant): `version`
  - `args` ([]string): trailing arguments (must be empty for success path)
  - `flags` ([]string): provided flags (none required for success path)
  - `exitCode` (int): expected process exit code (0 on success)
  - `stdout` (string): expected exact output (`<value>\n`)
  - `stderr` (string): expected empty string on success
- Validation rules:
  - Success path requires `args` length = 0.
  - Successful output must be exactly one line.

## Relationships

- `VersionCommandInvocation` reads one `VersionInfo` value.
- `VersionInfo` is produced by `internal/version` and consumed by
  `internal/cli/version_cmd.go`.
