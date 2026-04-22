# Implementation Plan: System Locale Localized Messages

**Branch**: `008-locale-messages` | **Date**: 2026-04-21 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/008-locale-messages/spec.md`

## Summary

Add a localization layer to the `mytets` CLI that detects the system locale and loads translated UI strings (help text, error messages, Cobra structural labels) from embedded per-language TOML files. English (`en.toml`) is the default/fallback; Ukrainian (`uk.toml`) is the first non-English language. Phrase content from `phrases.json` is excluded from localization. Adding a new language requires only a new `.toml` file and a rebuild — no Go code changes.

## Technical Context

**Language/Version**: Go 1.26.2 (with Go 1.25+ compatibility)
**Primary Dependencies**: `github.com/spf13/cobra` (CLI framework, existing), `github.com/BurntSushi/toml` (TOML parsing — new dependency, justified by FR-006 requiring TOML localization files)
**Storage**: Embedded TOML files via `//go:embed` (no runtime file I/O)
**Testing**: `go test` (unit + integration), `go test -race`
**Target Platform**: Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64)
**Project Type**: CLI (single static binary)
**Performance Goals**: Binary start + output in under 100 ms (existing constraint maintained)
**Constraints**: Single binary, no runtime file I/O, Go 1.25+ compatible
**Scale/Scope**: 2 languages (en, uk) initially; architecture supports unlimited additions

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Clean, Self-Explained Code | PASS | Localization package will use descriptive names (`LoadBundle`, `Translate`, `DetectLocale`); no comments needed to explain *what* |
| II. Simplicity | PASS | No new CLI flags added. Locale detection is automatic from environment. Zero additional user cognitive load |
| III. Reliability | PASS | Missing translations fall back to `en.toml`; missing locale falls back to English. All error paths tested |
| IV. Performance | PASS | TOML files parsed once at init via `//go:embed`; no runtime I/O. Negligible impact on startup time |
| V. Extensibility | PASS | New language = new `.toml` file + rebuild. No existing code modified |
| VI. Documentation | PASS | Help text is the primary deliverable; `--help` output will be fully localized |
| VII. Distribution | PASS | Single static binary maintained; TOML files embedded at build time |
| VIII. Go Best Practices | PASS | New package `internal/i18n` — singular, short, lowercase. No stuttering. Domain-organized |
| External deps prohibition | JUSTIFIED | `github.com/BurntSushi/toml` is required because the Go standard library has no TOML parser. The spec mandates TOML format for localization files (FR-006). Documented in README |
| Go 1.25+ compatibility | PASS | `BurntSushi/toml` v1.x supports Go 1.21+; `//go:embed` available since Go 1.16 |

## Project Structure

### Documentation (this feature)

```text
specs/008-locale-messages/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── cli-locale-contract.md
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
internal/
├── i18n/                  # NEW — localization package
│   ├── i18n.go            # Bundle loading, locale detection, Translate() function
│   ├── i18n_test.go       # Unit tests
│   └── locales/           # Embedded TOML files (one per language)
│       ├── en.toml        # English translations (default/fallback, reference for translators)
│       └── uk.toml        # Ukrainian translations
├── cli/
│   ├── root.go            # MODIFIED — apply localized strings to Cobra commands
│   ├── run.go             # MODIFIED — initialize i18n before building command tree
│   ├── version_cmd.go     # MODIFIED — use localized descriptions/errors
│   └── run_test.go        # MODIFIED — tests for locale-aware execution
├── commands/
│   ├── one/
│   │   ├── one.go         # MODIFIED — use localized descriptions/errors
│   │   └── one_test.go    # MODIFIED — test localized error messages
│   └── list/
│       ├── list.go        # MODIFIED — use localized descriptions/errors
│       └── list_test.go   # MODIFIED — test localized error messages
├── flags/
│   └── parser.go          # MODIFIED — use localized error messages
└── phrases/               # UNCHANGED — phrases.json not affected by localization

tests/
└── integration/
    ├── locale_help_test.go     # NEW — integration tests for localized help output
    ├── version_command_test.go # EXISTING — verify backward compatibility
    ├── one_command_test.go     # EXISTING — verify backward compatibility
    └── list_command_test.go    # EXISTING — verify backward compatibility
```

**Structure Decision**: A single new package `internal/i18n` with embedded `locales/` directory. This follows the existing flat domain-organized pattern under `internal/`. The `i18n` package owns locale detection, TOML loading, and the `Translate()` API. Existing packages import `i18n` to resolve their user-facing strings.

## Complexity Tracking

> **No constitution violations requiring justification beyond the TOML dependency (documented in Constitution Check above).**
