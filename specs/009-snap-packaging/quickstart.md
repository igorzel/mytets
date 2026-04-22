# Quickstart: Snap Packaging

**Feature**: 009-snap-packaging  
**Date**: 2026-04-22

## Prerequisites

- Linux system (snap packaging is Linux-only)
- `snapcraft` installed: `sudo snap install snapcraft --classic`
- `make` installed: `sudo apt install make` (usually pre-installed)
- Git repository with at least one tag for version derivation (optional — falls back to `dev`)

## Build the Go Binary

```bash
make build
```

This compiles the binary to `./bin/mytets` with the version injected from git tags:

```bash
./bin/mytets version
# Output: 1.0.1  (or "dev" if no git tags)
```

## Build the Snap Package

```bash
make snap
```

This runs snapcraft from `packaging/snap/`, which:
1. Clones the repository source (preserving git tags)
2. Installs Go toolchain via snap
3. Compiles the binary with version-injected ldflags
4. Packages everything into a `.snap` file

Output: `packaging/snap/mytets_<version>_amd64.snap`

## Install and Test Locally

```bash
# Install the locally built snap (--dangerous for unsigned local snaps)
sudo snap install packaging/snap/mytets_*.snap --dangerous

# Verify all commands work
mytets version
mytets one
mytets list
mytets --output json version

# Remove when done
sudo snap remove mytets
```

## Upload to Snap Store

```bash
# Login to your Snap Store account (one-time)
snapcraft login

# Upload and release to stable channel
snapcraft upload packaging/snap/mytets_*.snap --release=stable
```

## Clean Build Artifacts

```bash
make clean
```

## File Structure

```text
Makefile                          # make build, make snap, make clean
packaging/
└── snap/
    └── snapcraft.yaml            # Snap build configuration
```

## Troubleshooting

### `snapcraft: command not found`
Install snapcraft: `sudo snap install snapcraft --classic`

### Version shows "dev"
Ensure you have at least one git tag: `git tag v1.0.0 && git push --tags`

### Snap build fails with Go errors
Ensure your Go code compiles locally first: `make build`

### Permission denied on snap install
Local snaps require `sudo` and `--dangerous` flag: `sudo snap install <file>.snap --dangerous`
