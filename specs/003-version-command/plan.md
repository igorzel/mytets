# Implementation Plan: Version Command

**Branch**: `003-version-command` | **Date**: 2026-04-21 | **Spec**: `/home/igor/dev/workspace/mytets/specs/003-version-command/spec.md`
**Input**: Feature specification from `/home/igor/dev/workspace/mytets/specs/003-version-command/spec.md`

## Summary

Add `mytets version` so it prints exactly one plain version string and exits 0.
The implementation will use Cobra-based command routing in a dedicated CLI/flag
package, keeping `main` thin and delegating argument parsing outside `main` for
independent unit testing. Version value will be embedded at build time via
ldflags into `internal/version.Version`, with `dev` fallback when not provided.

## Technical Context

**Language/Version**: Go 1.26.2 in repository; feature constrained to Go 1.25+ compatibility  
**Primary Dependencies**: `github.com/spf13/cobra` for command/flag parsing (new dependency), Go standard library  
**Storage**: N/A (no persistent storage; in-memory version string only)  
**Testing**: `go test ./...`, table-driven unit tests for parser/handlers, integration CLI tests for end-to-end invocation  
**Target Platform**: Linux, macOS, Windows CLI environments
**Project Type**: Single-binary CLI application  
**Performance Goals**: `mytets version` completion under 100 ms; minimal allocations  
**Constraints**: Plain-text-only output for `version`; no custom parsing in `main`; build-time ldflags injection; clear non-zero exits on invalid invocation  
**Scale/Scope**: Single command addition (`version`) plus foundational CLI package separation (`main` -> `internal/cli` and `internal/flags`)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Principle I (Clean Code): PASS. `main` will delegate to explicit package APIs.
- Principle II (Simplicity): PASS. `mytets version` accepts no additional args.
- Principle III (Reliability): PASS. deterministic stdout/stderr and exit code.
- Principle IV (Performance): PASS. constant-time command path, no I/O.
- Principle V (Extensibility): PASS. command registration in `internal/cli`.
- Principle VI (Documentation): PASS. command help and README usage updates.
- Principle VII (Distribution): PASS. version injected by `-ldflags` in build.
- Principle VIII (Go Practices): PASS. parsing and reusable logic outside `main`.

Constitution exception acknowledged from spec clarification:
- `mytets version` is plain-text only and does not provide `--output json`.

## Project Structure

### Documentation (this feature)

```text
/home/igor/dev/workspace/mytets/specs/003-version-command/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── cli-version-contract.md
└── tasks.md
```

### Source Code (repository root)

```text
/home/igor/dev/workspace/mytets/
├── cmd/
│   └── mytets/
│       └── main.go
├── internal/
│   ├── cli/
│   │   ├── run.go
│   │   ├── root.go
│   │   └── version_cmd.go
│   ├── flags/
│   │   └── parser.go
│   └── version/
│       └── version.go
├── tests/
│   └── integration/
│       └── version_command_test.go
└── go.mod
```

**Structure Decision**: Single-project CLI layout with command-line parsing and
command wiring implemented in `internal/cli` and `internal/flags` rather than
`cmd/mytets/main.go`. This directly supports unit testing and aligns with the
constitution requirement for reusable logic in `internal/`.

## Post-Design Constitution Check

- Principle I (Clean Code): PASS. clear package boundaries and small APIs.
- Principle II (Simplicity): PASS. version command remains no-arg/no-flag.
- Principle III (Reliability): PASS. explicit success/error output contracts.
- Principle IV (Performance): PASS. O(1) command path and no runtime I/O.
- Principle V (Extensibility): PASS. additional subcommands can be registered in
	`internal/cli` without growing `main`.
- Principle VI (Documentation): PASS. quickstart and CLI contract authored.
- Principle VII (Distribution): PASS. ldflags-based embedding retained.
- Principle VIII (Go Practices): PASS. reusable logic in `internal/*` packages.

Result: PASS. No constitution violations introduced by design artifacts.

## Complexity Tracking

No constitution violations requiring justification.
