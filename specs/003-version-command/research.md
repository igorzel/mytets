# Phase 0 Research — Version Command

## Decision 1: Use Cobra for CLI command and flag parsing

- Decision: Add `github.com/spf13/cobra` and route all command parsing through Cobra.
- Rationale: The feature request explicitly requires Cobra and requires parsing
  to live outside `main`; Cobra provides tested command dispatch and standard
  help/error behavior.
- Alternatives considered:
  - Standard library `flag`: rejected because it does not model subcommands
    cleanly without custom parsing glue.
  - Manual `os.Args` parsing: rejected by requirement and hard to test reliably.

## Decision 2: Keep `main` as thin process entrypoint

- Decision: `cmd/mytets/main.go` only calls exported package API (e.g.
  `cli.Execute()`) and handles process exit.
- Rationale: Enforces separation of concerns and allows parser/command behavior
  to be unit tested independently of process startup.
- Alternatives considered:
  - Parse commands in `main`: rejected by explicit feature requirement.

## Decision 3: Create dedicated parsing package

- Decision: Introduce `internal/flags` package for parser configuration and
  invocation behavior consumed by `internal/cli`.
- Rationale: Satisfies requirement for a dedicated package and provides a seam
  for focused unit tests.
- Alternatives considered:
  - Put parser logic only in `internal/cli`: rejected because the user asked for
    a dedicated parsing package with clear exported API.

## Decision 4: Build-time version metadata location

- Decision: Store the version variable in `internal/version` as exported
  `Version` with default `dev`.
- Rationale: Clarified in spec; allows `-ldflags -X` injection and enables
  unit tests without binary rebuilding.
- Alternatives considered:
  - `main.Version`: rejected due to weaker reuse/testability.
  - Runtime file/env lookup: rejected by requirement for binary-embedded value.

## Decision 5: Output contract for `mytets version`

- Decision: `mytets version` prints a single plain string and trailing newline;
  no labels and no JSON mode.
- Rationale: Explicit FR-003 and clarification decision B.
- Alternatives considered:
  - JSON output via `--output json`: rejected for this command by clarification.

## Decision 6: Testing strategy

- Decision: Add unit tests for version rendering/parsing seams and integration
  tests for full invocation (`stdout`, `stderr`, `exit code`).
- Rationale: Required by FR-008/FR-009 and constitution testing standards.
- Alternatives considered:
  - Integration-only tests: rejected due to reduced fault localization.
