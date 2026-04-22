# Research: System Locale Localized Messages

**Feature**: 008-locale-messages  
**Date**: 2026-04-21  
**Status**: Complete

## Research Task 1: Go TOML Parsing Library

**Decision**: Use `github.com/BurntSushi/toml`

**Rationale**: Industry-standard Go TOML parser. Requires Go 1.15+ (compatible with project's Go 1.25+ constraint). Natively supports flat `map[string]string` unmarshalling, which is exactly the localization file structure (flat key-value string pairs). Parses embedded files in <1ms — negligible startup cost.

**Alternatives considered**:
- `github.com/pelletier/go-toml/v2`: Newer, ~2× faster, requires Go 1.21+. Viable but BurntSushi/toml is more widely adopted and sufficient for this use case.
- Go standard library: No TOML parser exists in the standard library. JSON was rejected in the clarification phase (less ergonomic for translators). YAML would require another external dependency with more complex parsing.

## Research Task 2: Locale Detection from Environment Variables

**Decision**: Read `LC_ALL` → `LC_MESSAGES` → `LANG` in POSIX priority order; extract language code by splitting on `_` and `.`

**Rationale**: This is the standard POSIX convention. Works identically on Linux and macOS. On Windows, these variables may be unset (returns empty string), which correctly triggers the English fallback. No platform-specific system calls needed — environment variables are portable.

**Language code extraction**: Parse `uk_UA.UTF-8` → split on `_` or `.` → take first segment → `uk`. Handle edge cases: bare `uk` (no separator) returns `uk`; empty string returns fallback `en`.

**Alternatives considered**:
- `golang.org/x/text/language`: Provides BCP 47 tag matching. Heavyweight dependency for a simple 2-language lookup. Rejected per constitution's minimal-dependency principle.
- Windows `GetLocaleInfo()` API: Platform-specific, breaks single-binary portability. Env var approach works cross-platform.

## Research Task 3: Cobra Help Text Customization

**Decision**: Use a combination of `cmd.SetUsageTemplate()`, `cmd.SetHelpTemplate()`, `cmd.SetFlagErrorFunc()`, and direct `.Short`/`.Long` field assignment

**Rationale**: Cobra does NOT provide built-in i18n hooks. The customization surfaces are:

1. **Structural labels** (`Usage:`, `Flags:`, `Available Commands:`, etc.): Override via `cmd.SetUsageTemplate()` and `cmd.SetHelpTemplate()` with localized template strings.

2. **Command descriptions** (`.Short`, `.Long`): Set directly on each `cobra.Command` before registration. The `i18n` package provides translated values at command-tree construction time.

3. **Flag descriptions**: Set via `cmd.Flags().StringVarP()` — the description string is passed at construction time, so it uses the localized value.

4. **Built-in error messages** (`unknown command`, `unknown flag`): Override via `cmd.SetFlagErrorFunc()` for flag validation errors. For unknown-command errors, Cobra's `Execute()` returns the error which can be intercepted and translated before printing to stderr.

**Alternatives considered**:
- Monkey-patching Cobra template functions: Fragile, breaks on Cobra upgrades. Rejected.
- Forking Cobra: Massive maintenance burden for a small feature. Rejected.
- Post-processing output strings: Unreliable regex matching on natural language text. Rejected.

## Research Task 4: Embedding TOML Files with go:embed

**Decision**: Use `//go:embed locales/*.toml` with `embed.FS` and iterate via `fs.ReadDir()` for auto-discovery

**Rationale**: The `//go:embed` directive with a glob pattern captures all matching files at compile time. Adding a new `de.toml` file to the `locales/` directory and rebuilding automatically includes it — no Go code changes required (FR-007). The `fs.ReadDir()` + `fs.ReadFile()` APIs enable iterating over all embedded locale files at startup to auto-register languages.

**Alternatives considered**:
- Separate `//go:embed` per file: Requires a Go code change for each new language. Violates FR-007. Rejected.
- External file loading at runtime: Violates FR-005 (must be compiled into binary). Rejected.

## Research Task 5: Fallback Pattern

**Decision**: Two-tier fallback — language-specific TOML → English TOML → key name (graceful degradation)

**Rationale**: 
- Tier 1: Look up key in the active language's map (e.g., `uk`)
- Tier 2: If missing, look up key in `en` map (always present as reference)
- Tier 3: If both miss (should not happen in practice), return the key name itself — makes untranslated strings visible during development

Bundle is loaded once at startup into a `map[string]map[string]string` (lang → key → value). Read-only after initialization, so no mutex needed.

**Alternatives considered**:
- Per-call file reading: Violates performance constraint. Rejected.
- Interface-based strategy pattern: Over-engineering for 2 tiers. Simple map lookup is sufficient.

## Research Task 6: FR-007 Compliance (Zero-Code-Change Language Addition)

**Decision**: Auto-discover languages by iterating `embed.FS` directory entries; derive language code from filename (e.g., `uk.toml` → `uk`)

**Rationale**: The `//go:embed locales/*.toml` glob captures all TOML files at build time. At startup, the `i18n` package reads the embedded directory, extracts language codes from filenames, and loads each bundle. A contributor adds a new language by:
1. Creating `internal/i18n/locales/xx.toml` (copy `en.toml` as template)
2. Translating the values
3. Running `go build`

No Go source files are modified. The glob auto-includes the new file.

**Compliance verified**: ✅ FR-005 (embedded), ✅ FR-006 (single file per language), ✅ FR-007 (no code changes)
