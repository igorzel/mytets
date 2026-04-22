# mytets Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-04-22

## Active Technologies
- Go 1.26.2 + `github.com/spf13/cobra` (already present) (speckit-try)
- N/A (stateless command) (speckit-try)
- Go 1.26.2 (with Go 1.25+ compatibility) + github.com/spf13/cobra (CLI command framework, already present) (006-one-random-phrases)
- Embedded JSON file (`phrases.json`) compiled into binary (006-one-random-phrases)
- Go 1.26.2 (with Go 1.25+ compatibility) + github.com/spf13/cobra (existing CLI framework), Go standard library (007-list-random-phrases)
- Embedded JSON file at `internal/phrases/phrases.json` compiled into the binary (007-list-random-phrases)
- Go 1.26.2 (with Go 1.25+ compatibility) + `github.com/spf13/cobra` (CLI framework, existing), `github.com/BurntSushi/toml` (TOML parsing — new dependency, justified by FR-006 requiring TOML localization files) (008-locale-messages)
- Embedded TOML files via `//go:embed` (no runtime file I/O) (008-locale-messages)
- Go 1.26.2 (with Go 1.25+ compatibility) + `snapcraft` (build-time tool, not a Go dependency), `make` (build orchestration) (009-snap-packaging)
- N/A (no data storage changes) (009-snap-packaging)

- Go 1.26.2 in repository; feature constrained to Go 1.25+ compatibility + `github.com/spf13/cobra` for command/flag parsing (new dependency), Go standard library (003-version-command)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.26.2 in repository; feature constrained to Go 1.25+ compatibility

## Code Style

Go 1.26.2 in repository; feature constrained to Go 1.25+ compatibility: Follow standard conventions

## Recent Changes
- 009-snap-packaging: Added Go 1.26.2 (with Go 1.25+ compatibility) + `snapcraft` (build-time tool, not a Go dependency), `make` (build orchestration)
- 008-locale-messages: Added Go 1.26.2 (with Go 1.25+ compatibility) + `github.com/spf13/cobra` (CLI framework, existing), `github.com/BurntSushi/toml` (TOML parsing — new dependency, justified by FR-006 requiring TOML localization files)
- 007-list-random-phrases: Added Go 1.26.2 (with Go 1.25+ compatibility) + github.com/spf13/cobra (existing CLI framework), Go standard library


<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
