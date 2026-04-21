# Quickstart — Version Command Feature

## Prerequisites

- Go 1.25+ installed
- Workspace root: `/home/igor/dev/workspace/mytets`

## Build With Version Injection

```bash
cd /home/igor/dev/workspace/mytets
go build -o ./bin/mytets -ldflags "-X github.com/igorzel/mytets/internal/version.Version=1.0.1" ./cmd/mytets
```

## Run — Plain Text (default)

```bash
./bin/mytets version
```

Expected output:

```text
1.0.1
```

## Run — JSON Output

```bash
./bin/mytets version --output json
# or
./bin/mytets version -o json
```

Expected output:

```json
{"version":"1.0.1"}
```

## Build Without Injection (Fallback)

```bash
cd /home/igor/dev/workspace/mytets
go build -o ./bin/mytets ./cmd/mytets
./bin/mytets version
```

Expected output:

```text
dev
```

## Test

```bash
cd /home/igor/dev/workspace/mytets

# Unit tests (with race detector)
go test -race ./internal/...

# Coverage gate (≥80%)
go test -coverprofile=cover.out ./internal/...
go tool cover -func=cover.out

# Integration tests
go test -timeout 120s ./tests/integration/...

# All tests
go test ./...
```

Integration validation checks:

- `mytets version` prints exactly one line to stdout.
- stderr is empty on success.
- process exit code is 0.
- ldflags-injected version matches output exactly.
- `--output json` / `-o json` returns a valid JSON object with a `version` field.
- unsupported format (e.g. `--output yaml`) returns non-zero exit and non-empty stderr.
- execution completes in under 100 ms.

## Race / Lint / Performance Checks

```bash
# Race detector
go test -race ./...

# Vet
go vet ./...

# Performance (covered by TestVersionCommandPerformance in integration tests)
go test -run TestVersionCommandPerformance -v ./tests/integration/...
```

## Final Test Results (2026-04-21)

All tests pass on Linux:

- Unit tests (internal/cli, internal/flags, internal/version): PASS
- Integration tests (13 scenarios): PASS
- Coverage (internal/): 82.1% — above 80% threshold

## Suggested Package API Shape

- `internal/cli.Execute() int`: process entrypoint; returns exit code for `os.Exit`.
- `internal/cli.ExecuteArgs(args []string) (stdout, stderr string, exitCode int)`:
  test seam for invocation behavior.
- `internal/flags.NewParserConfig() ParserConfig`: centralized parser defaults.
- `internal/flags.ParseOutputFormat(raw string) (OutputFormat, error)`: validated
  format parsing consumed by the version command.
