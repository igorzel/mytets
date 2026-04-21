# mytets

## Installation

```bash
go install github.com/igorzel/mytets/cmd/mytets@latest
```

Or build from source:

```bash
go build -o ./bin/mytets ./cmd/mytets
```

## Usage

### version

Print the application version and exit:

```bash
mytets version
```

Output (plain text, default):

```text
1.0.1
```

JSON output for automation and CI:

```bash
mytets version --output json
# or
mytets version -o json
```

Output:

```json
{"version":"1.0.1"}
```

### Build with version injection (ldflags)

To embed a version at build time:

```bash
go build \
  -ldflags "-X github.com/igorzel/mytets/internal/version.Version=1.0.1" \
  -o ./bin/mytets ./cmd/mytets
```

When built without ldflags the version falls back to `dev`:

```bash
go build -o ./bin/mytets ./cmd/mytets
./bin/mytets version  # prints: dev
```

### CI / Script usage

```bash
# Capture version in a shell variable
VERSION=$(mytets version)

# Verify it matches semver
echo "$VERSION" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$'

# JSON capture
mytets version --output json | jq -r .version
```

## Dependency Rationale

This project introduces `github.com/spf13/cobra` for CLI command routing.
The dependency is used to provide stable subcommand and flag behavior,
consistent `--help` output, and clean extensibility for future commands.

The standard library flag parser was not selected for this feature because it
does not provide subcommand ergonomics and help output consistency at the same
level without additional custom command-routing scaffolding.