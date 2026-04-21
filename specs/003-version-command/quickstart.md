# Quickstart — Version Command Feature

## Prerequisites

- Go 1.25+ installed
- Workspace root: `/home/igor/dev/workspace/mytets`

## Build With Version Injection

```bash
cd /home/igor/dev/workspace/mytets
go build -o ./bin/mytets -ldflags "-X github.com/igorzel/mytets/internal/version.Version=1.0.1" ./cmd/mytets
```

## Run

```bash
./bin/mytets version
```

Expected output:

```text
1.0.1
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
go test ./...
```

Integration focus checks:

- `mytets version` prints exactly one line to stdout.
- stderr is empty on success.
- process exit code is 0.

## Suggested Package API Shape

This feature expects `main` to call into parser/CLI packages rather than parse
arguments directly:

- `internal/cli.Execute() error`: process entrypoint for runtime execution.
- `internal/cli.ExecuteArgs(args []string) (stdout, stderr string, exitCode int, err error)`:
  test seam for invocation behavior.
- `internal/flags.NewParserConfig() ParserConfig`: centralized parser behavior.
