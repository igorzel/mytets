# Implementation Plan: List Command - Random Phrase List

**Branch**: `007-list-random-phrases` | **Date**: 2026-04-21 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/007-list-random-phrases/spec.md`

**Note**: This plan is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Implement a new `mytets list` command that returns a unique random list of phrases from the existing embedded phrase source, with plain-text output by default and compact JSON output via the existing global `--output json` flow. The core list-generation logic will live in a reusable internal package decoupled from Cobra so it can be unit-tested independently now and reused later by non-CLI adapters such as an HTTP endpoint.

## Technical Context

**Language/Version**: Go 1.26.2 (with Go 1.25+ compatibility)  
**Primary Dependencies**: github.com/spf13/cobra (existing CLI framework), Go standard library  
**Storage**: Embedded JSON file at `internal/phrases/phrases.json` compiled into the binary  
**Testing**: Go `testing` package, package-level unit tests, CLI integration tests in `tests/integration`  
**Target Platform**: Linux, macOS, Windows single-binary CLI  
**Project Type**: CLI application with reusable internal domain packages  
**Performance Goals**: Command startup and response remain under 100 ms on reference hardware  
**Constraints**: No runtime file I/O, no new external dependencies, count flag remains command-specific, list generation must be reusable outside CLI wiring  
**Scale/Scope**: Single binary, embedded phrase catalog of tens to low hundreds of phrases, one new subcommand plus shared internal package refactoring

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Pre-Research Gate

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Clean, Self-Explained Code | ✅ PASS | Plan separates CLI wiring, phrase source access, and reusable list generation into focused packages |
| II. Simplicity | ✅ PASS | User-facing surface stays small: `mytets list` and `mytets list --count N` |
| III. Reliability | ✅ PASS | Shared phrase-source errors remain explicit; invalid counts and format errors return stderr output and non-zero exit codes |
| IV. Performance | ✅ PASS | Embedded data plus in-memory unique sampling keeps execution well below the 100 ms target |
| V. Extensibility | ✅ PASS | Reusable list package is intentionally decoupled from Cobra for future adapters such as REST handlers |
| VI. Documentation | ✅ PASS | Plan includes command contract and quickstart updates for help, testing, and usage |
| VII. Distribution | ✅ PASS | No runtime dependencies added; feature remains part of the existing single binary |
| VIII. Go Best Practices | ✅ PASS | Reusable logic stays in `internal/`; packages remain small and domain-focused |

**Gate Result**: ✅ PASS

### Post-Design Re-Check

The Phase 1 design artifacts preserve the same alignment: the reusable list package returns domain data rather than CLI-specific output, testing remains deterministic via injectable randomness, and no constitutional violations were introduced.

**Post-Design Result**: ✅ PASS

## Project Structure

### Documentation (this feature)

```text
specs/007-list-random-phrases/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── spec.md
├── contracts/
│   └── cli-list-contract.md
├── checklists/
│   └── requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/mytets/
└── main.go

internal/
├── cli/
│   ├── root.go
│   ├── run.go
│   ├── run_test.go
│   └── version_cmd.go
├── commands/
│   ├── one/
│   │   ├── one.go
│   │   └── one_test.go
│   └── list/
│       ├── list.go
│       └── list_test.go
├── flags/
│   ├── parser.go
│   └── parser_test.go
├── listing/
│   ├── listing.go
│   └── listing_test.go
├── phrases/
│   ├── phrases.go
│   ├── phrases.json
│   └── phrases_test.go
└── version/
    ├── version.go
    └── version_test.go

tests/
└── integration/
    ├── list_command_test.go
    ├── one_command_test.go
    └── version_command_test.go
```

**Structure Decision**: Keep CLI registration and output adaptation in `internal/commands/list`, and place reusable list-generation behavior in a separate `internal/listing` package. `internal/listing` will depend on phrase retrieval abstractions rather than Cobra, making it suitable for unit tests and future reuse by non-CLI delivery layers without forcing a later refactor.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|

