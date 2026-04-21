<!--
## Sync Impact Report

**Version change**: (new) → 1.0.0
**Modified principles**: N/A — initial constitution fill
**Added sections**: Core Principles (I–VIII), Constraints & Technical Standards,
  Testing & Quality, Security, Governance
**Removed sections**: N/A
**Templates reviewed**:
  - .specify/templates/plan-template.md ✅ (Constitution Check gate is generic; no update required)
  - .specify/templates/spec-template.md ✅ (structure compatible; no update required)
  - .specify/templates/tasks-template.md ✅ (phase structure compatible; no update required)
**Deferred TODOs**: None
-->

# mytets Constitution

## Core Principles

### I. Clean, Self-Explained Code

Code MUST communicate intent without requiring comments to explain *what* it does.
Variable and function names MUST be descriptive enough to make logic obvious to
a new reader. Comments are reserved for *why*, never *what*. Code that requires
a lengthy comment to explain its logic MUST be refactored.

**Rationale**: The project is a lightweight single-binary tool; clarity reduces
maintenance cost and lowers the barrier for contributors.

### II. Simplicity

Commands MUST be intuitive and require minimal flags. The default invocation
(`mytets`) MUST produce useful output with zero arguments. Each flag MUST solve
a concrete, documented user need — no speculative options. When in doubt, leave
it out (YAGNI).

**Rationale**: A tool that is simple to invoke is simple to trust. Complexity
added prematurely is complexity that must be maintained forever.

### III. Reliability

Errors MUST be handled gracefully with a human-readable message printed to
`stderr` and a non-zero exit code returned. The tool MUST NOT panic in normal
operation. Every error path MUST be tested.

**Rationale**: Unreliable CLI tools break scripts and erode user trust.

### IV. Performance

The binary MUST start and produce output in under 100 ms on reference hardware
(modern laptop, cold filesystem cache). Memory allocations MUST be kept minimal;
prefer stack allocation and avoid unnecessary heap churn.

**Rationale**: A `fortune`-style tool is invoked frequently (e.g., in shell
`rc` files). Latency is a first-class concern.

### V. Extensibility

The internal package structure MUST allow new subcommands and future plugins
without modifying existing command logic. Public APIs MUST be minimal and
stable; breaking changes require a major-version bump.

**Rationale**: The tool may evolve (e.g., filters by category, multiple quote
sources). Architecture must accommodate growth without rewrites.

### VI. Documentation

Every command and flag MUST have concise help text accessible via `--help`.
The `README.md` MUST cover installation, usage, and a quick-start example.
Public Go symbols MUST have godoc-compatible comments.

**Rationale**: Documentation is part of the product. An undocumented flag is
an unusable flag.

### VII. Distribution

The project MUST produce a single static binary with no runtime dependencies.
Releases MUST target Linux (amd64/arm64), macOS (amd64/arm64), and Windows
(amd64) via cross-compilation. A reproducible build MUST be achievable with a
single `go build` invocation.

**Rationale**: Ease of installation directly affects adoption.

### VIII. Go Best Practices

The codebase MUST follow idiomatic Go conventions:

- Package names MUST be singular, short, and lowercase (e.g., `quote`, not
  `quotes` or `quotation`).
- Stuttering MUST be avoided (e.g., prefer `quote.Item`, not `quote.QuoteItem`).
- Packages MUST be small and focused — one clear responsibility per package.
- `main` MUST live in its own package under `cmd/mytets/`; all reusable logic
  MUST live in `internal/` packages.
- Non-public packages MUST use the `internal/` prefix to prevent misuse.
- Code MUST be organized by domain/functionality (e.g., `internal/quote`,
  `internal/output`), not by type.
- Unnecessary nesting MUST be avoided; flat structures are preferred.
- Only necessary symbols MUST be exported; APIs MUST be minimal and clean.

**Rationale**: Consistency with Go idioms makes the codebase instantly
navigable by any Go developer and prevents common anti-patterns.

## Constraints & Technical Standards

- The project MUST compile with **Go 1.25 or later**.
- External dependencies are PROHIBITED unless no reasonable standard-library
  alternative exists; any exception MUST be documented in `README.md` with
  justification.
- Default output MUST be human-readable plain text to `stdout`.
- Machine-readable output MUST be available via `--output json` (or `-o json`)
  for automation and scripting.
- The embedded quotation data MUST be stored as a Go source file (e.g., using
  `embed` or a generated slice) — no runtime file I/O required.
- Version and build metadata MUST be injected at link time via `-ldflags`; no
  hard-coded version strings in source.

## Testing & Quality

- Unit tests MUST cover all core functions in `internal/` packages.
- Integration tests MUST cover all CLI commands and flag combinations.
- Inputs and outputs MUST be mockable to enable deterministic testing of
  random-selection logic (e.g., injectable `rand.Source`).
- Test coverage MUST remain at or above **80 %** for `internal/` packages.
- `go vet`, `staticcheck`, and `golangci-lint` MUST pass with zero warnings
  before any merge.
- The CI pipeline MUST run tests on all three target platforms.

## Security

- All user-supplied inputs (flags, environment variables) MUST be validated
  before use; invalid inputs MUST produce a clear error and exit code 1.
- Sensitive data (tokens, keys, personal information) MUST NOT appear in log
  output at any verbosity level.
- Goroutine usage MUST follow Go memory-safety conventions; data races are a
  build-blocking defect (run `go test -race`).
- Dependencies (when permitted) MUST be reviewed for known CVEs before
  inclusion and on each release.

## Governance

This constitution supersedes all other project guidelines. Any practice that
conflicts with a principle stated here MUST be brought into compliance.

**Amendment procedure**:
1. Open a GitHub issue describing the proposed change and its rationale.
2. Discuss with maintainers; reach consensus (or majority vote for open-source
   contributors).
3. Update this file, increment the version per semantic versioning rules, and
   update `LAST_AMENDED_DATE`.
4. Document breaking changes in `CHANGELOG.md`.

**Versioning policy** (semantic versioning):
- MAJOR: Backward-incompatible removal or redefinition of a principle.
- MINOR: New principle or section added; materially expanded guidance.
- PATCH: Clarifications, wording fixes, non-semantic refinements.

**Compliance review**: All pull requests MUST include a self-checklist
confirming compliance with applicable principles. Reviewers MUST flag any
violation before approving.

**Version**: 1.0.0 | **Ratified**: 2026-04-21 | **Last Amended**: 2026-04-21
