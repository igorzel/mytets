# Quickstart: List Command - Random Phrase List

**Feature**: implement `mytets list` with reusable list-generation logic  
**Audience**: developers implementing and validating the feature

## 1. Work from the repository root

```bash
cd /home/igor/dev/workspace/mytets
```

## 2. Implement the reusable list package first

Create a dedicated internal package for domain logic, separate from Cobra wiring.

**Target package**: `internal/listing`

Implementation goals:
- accept the available phrase set and requested count
- select phrases uniquely without replacement
- expose a deterministic seam for unit tests by injecting randomness or a picker
- return neutral domain data that can be rendered by CLI today and reused by future adapters later

This package is the primary unit-test target for count handling, uniqueness, and oversized requests.

## 3. Wire the CLI adapter second

Create `internal/commands/list` to:
- define the `list` Cobra command
- add the command-specific `--count` flag
- call into `internal/listing`
- render plain-text lines or compact JSON arrays based on the existing output format selection

Also register the command in `internal/cli/root.go`.

## 4. Reuse the shared phrase source

Use `internal/phrases` as the single source of phrase loading and validation.

Implementation goals:
- do not duplicate phrase-source loading inside the new command
- preserve the existing startup failure behavior for missing or empty phrase data
- keep `list` and `one` aligned on phrase availability rules

## 5. Add focused tests

### Unit tests

Add or extend tests for:
- `internal/listing`: exact count handling, uniqueness, deduplication by text, and deterministic random selection
- `internal/commands/list`: plain-text and JSON adaptation behavior, plus command-specific flag validation
- `internal/phrases`: any small additions needed to support reusable phrase access

### Integration tests

Add `tests/integration/list_command_test.go` to verify:
- `mytets list` returns 5 unique phrases by default
- `mytets list --count N` returns the requested number when available
- `mytets --output json list` returns valid compact JSON
- success exit code is 0

## 6. Validate the feature

Run focused checks during implementation:

```bash
go test ./internal/listing/... ./internal/commands/list/... ./tests/integration/...
go test ./...
```

### Verify CLI behavior manually

```bash
go build -o ./bin/mytets ./cmd/mytets

# Plain text (default 5, capped to available phrases)
./bin/mytets list

# Custom count
./bin/mytets list --count 2

# JSON output
./bin/mytets --output json list

# JSON with custom count
./bin/mytets --output json list --count 3
```

If command registration or global output parsing changes are touched, rerun the broader test suite before merging.

## 7. Keep the design reusable

The future REST API mentioned in the specification is out of scope for this feature. The implementation should only preserve that path by ensuring `internal/listing` has no Cobra-specific types or stdout/stderr dependencies.