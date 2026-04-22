# mytets

## Dependencies

- **[github.com/spf13/cobra](https://github.com/spf13/cobra)** — CLI command framework
- **[github.com/BurntSushi/toml](https://github.com/BurntSushi/toml)** — TOML parsing for localization files. Required because the Go standard library has no TOML parser. Used to load embedded per-language `.toml` translation files at startup.

## Installation

### Snap Store

```bash
sudo snap install mytets
```

### Local snap install

Build and install the snap locally:

```bash
make snap
sudo snap install packaging/snap/mytets_*.snap --dangerous
```

### Go install

```bash
go install github.com/igorzel/mytets/cmd/mytets@latest
```

Or build from source:

```bash
make build
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
# or global flag form
mytets --output json version
```

Output:

```json
{"version":"1.0.1"}
```

### Build with version injection (ldflags)

To embed a version at build time:

```bash
make build
./bin/mytets version  # prints version from git tags
```

Or manually:

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

### Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | Compile binary with version injection to `./bin/mytets` |
| `make snap` | Build a snap package via snapcraft |
| `make clean` | Remove the `./bin` build directory |

### CI / Script usage

```bash
# Capture version in a shell variable
VERSION=$(mytets version)

# Verify it matches semver
echo "$VERSION" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$'

# JSON capture
mytets version --output json | jq -r .version
```

### one

Print one random phrase from the embedded phrase set:

```bash
mytets one
```

Output (plain text, default; one of multiple embedded messages):

```text
Example message 2
```

JSON output:

```bash
mytets --output json one
# or
mytets one --output json
```

Output:

```json
{"message":"Example message 2"}
```

### list

Print a list of random phrases from the embedded phrase set:

```bash
mytets list
```

Output (plain text, default; 5 unique phrases, one per line):

```text
Example message 3
Fake message, tbd
Example message 1
Example message 2
```

Request a specific number of phrases:

```bash
mytets list --count 2
```

JSON output:

```bash
mytets --output json list
# or with a custom count
mytets --output json list --count 3
```

Output:

```json
[{"message":"Example message 3"},{"message":"Fake message, tbd"},{"message":"Example message 1"}]
```

## Dependency Rationale

This project introduces `github.com/spf13/cobra` for CLI command routing.
The dependency is used to provide stable subcommand and flag behavior,
consistent `--help` output, and clean extensibility for future commands.

The standard library flag parser was not selected for this feature because it
does not provide subcommand ergonomics and help output consistency at the same
level without additional custom command-routing scaffolding.
