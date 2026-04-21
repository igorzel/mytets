# Implementation Plan: One Command Support

**Branch**: `005-one-command-support` | **Date**: 2026-04-21 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `specs/005-one-command-support/spec.md`

## Summary

The `one` command feature adds a new subcommand (`mytets one`) that prints a fixed message ("Fake message, tbd") to stdout, with automatic JSON formatting support via the global `--json` flag. The implementation will create a dedicated internal package for the command to enable isolated unit testing and future reusability, following Go best practices and the project's Cobra-based CLI architecture.

## Technical Context

**Language/Version**: Go 1.26.2  
**Primary Dependencies**: `github.com/spf13/cobra` (already present)  
**Storage**: N/A (stateless command)  
**Testing**: Go's standard `testing` package + integration tests  
**Target Platform**: Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64)  
**Project Type**: CLI tool (single static binary)  
**Performance Goals**: Execution under 100 ms (per constitution)  
**Constraints**: No dynamic allocation for core message output; clean error handling  
**Scale/Scope**: Single new subcommand with two output modes (plain text, JSON)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle Compliance Matrix

| Principle | Requirement | Feature Status | Justification |
|-----------|-------------|---|---|
| I. Clean, Self-Explained Code | Variable/function names explain intent; comments reserved for *why* | ✅ PASS | Command logic and output formatting will be explicit; no ambiguous names |
| II. Simplicity | Intuitive commands, minimal flags, YAGNI | ✅ PASS | `mytets one` has no required arguments or subcommand-specific flags; respects global `--json` only |
| III. Reliability | Graceful error handling, stderr for errors, non-zero exit codes | ✅ PASS | Cobra handles flag errors automatically; command always exits 0 on success, unit/integration tests verify paths |
| IV. Performance | Binary starts and outputs in under 100 ms | ✅ PASS | Command is pure function (message is hardcoded); no I/O or external lookups; well within 100 ms budget |
| V. Extensibility | Internal packages allow new commands without modifying existing logic | ✅ PASS | Dedicated `internal/commands/one` package isolates implementation; root command registration is minimal |
| VI. Documentation | All commands have help text; README covered; godoc comments | ✅ PASS | Cobra auto-generates help; command will have Short/Long description; godoc-compatible comments on exported functions |
| VII. Distribution | Single static binary, no runtime dependencies | ✅ PASS | No new dependencies; command is pure Go standard library logic |
| VIII. Go Best Practices | Idiomatic Go: singular package names, avoid stuttering, internal/ prefix, domain organization | ✅ PASS | Package structure: `internal/commands/one`; no stuttering; command tests co-located with implementation |

**Gate Status**: ✅ **PASS** — No principle violations; feature aligns with all constitution requirements.

## Project Structure

### Documentation (this feature)

```text
specs/005-one-command-support/
├── plan.md              # This file
├── research.md          # Phase 0 output (minimal)
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (if needed)
└── spec.md              # Feature specification
```

### Source Code (repository root)

Current structure (relevant sections):

```text
cmd/
└── mytets/
    └── main.go          # Entry point; calls cli.Execute()

internal/
├── cli/
│   ├── root.go          # Root command; subcommand registration
│   └── version_cmd.go   # Existing version command (reference)
├── commands/            # NEW: Subcommand implementations
│   └── one.go           # NEW: One command implementation
├── flags/
│   └── parser.go        # Global flag parsing (--json, etc)
└── version/
    └── version.go       # Version management

tests/
└── integration/
    └── one_command_test.go  # NEW: Integration tests for `mytets one`
```

**Structure Decision**: Implement the `one` command in a dedicated `internal/commands/one` package (or `internal/commands/one/one.go` if command-specific utilities are added later). This mirrors the version command pattern but provides better isolation for unit testing and future copying to other projects.

## Phase 0: Research

**Scope**: Verify no ambiguities in specification or dependencies.  
**Status**: ✅ **COMPLETE** — Specification was clarified in `/speckit.clarify` session. All 3 clarification questions resolved:

1. ✅ Global `--json` flag is pre-existing; command respects it automatically
2. ✅ Only `mytets --json one` flag placement is supported (flag before subcommand)
3. ✅ Unsupported flags follow Cobra convention: error message to stderr, non-zero exit

**No new research required.** Proceed to Phase 1 design.

---

## Phase 1: Design & Contracts

### 1.1 Data Model

See: [data-model.md](data-model.md)

### 1.2 CLI Command Contract

See: [contracts/cli-one-contract.md](contracts/cli-one-contract.md)

### 1.3 Quick Start

See: [quickstart.md](quickstart.md)

---

## Implementation Approach (User Guidance)

Based on user input, this feature will:

1. **Use Cobra library** for command-line flag parsing (already present in the project as `github.com/spf13/cobra`)
2. **Create a separate package** (`internal/commands/one`) for easy unit testing and isolation
3. **Follow existing patterns** from `internal/cli/version_cmd.go` for Cobra command registration and output formatting
4. **Respect global `--json` flag** via the existing parser configuration passed to all commands
5. **Output message exactly as specified**: `Fake message, tbd` in plain mode; `{"message":"Fake message, tbd"}` in JSON mode

---

## Next Steps

1. Complete Phase 1 design artifacts (data-model.md, contracts/, quickstart.md)
2. Run `/speckit.tasks` to generate actionable task list
3. Begin implementation and testing
