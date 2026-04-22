# CLI Snap Contract: Snap Packaging

**Feature**: 009-snap-packaging  
**Date**: 2026-04-22

## Overview

This contract defines the interfaces exposed by the snap packaging feature: Makefile targets (developer-facing) and the snap package metadata (Snap Store-facing).

---

## Contract 1: Makefile Targets

**Interface type**: CLI build commands  
**Consumer**: Developer (local machine)

### Target: `make build`

**Purpose**: Compile the Go binary with version injection.

**Input**: None (reads git tags from repository)

**Output**: Binary at `./bin/mytets`

**Behavior**:
```
$ make build
go build -ldflags "-X github.com/igorzel/mytets/internal/version.Version=1.0.1" -o ./bin/mytets ./cmd/mytets

$ ./bin/mytets version
1.0.1
```

**Error cases**:
- Go toolchain not installed → standard `make` error: `go: command not found`
- No git tags → version falls back to `dev`

---

### Target: `make snap`

**Purpose**: Build a snap package containing the mytets application.

**Input**: None (requires snapcraft installed)

**Output**: `.snap` file in `packaging/snap/` directory (produced by snapcraft)

**Behavior**:
```
$ make snap
cd packaging/snap && snapcraft
...
Created snap package mytets_1.0.1_amd64.snap
```

**Error cases**:
- snapcraft not installed → `make` error: `snapcraft: command not found`
- No git tags → snap version falls back to `dev`

---

### Target: `make clean`

**Purpose**: Remove build artifacts.

**Behavior**:
```
$ make clean
rm -rf ./bin
```

---

## Contract 2: Snap Package Metadata

**Interface type**: Snap Store metadata  
**Consumer**: Snap Store, end users

### Required Fields

| Field | Value | Snap Store Requirement |
|-------|-------|----------------------|
| `name` | `mytets` | Must be registered on Snap Store |
| `version` | Derived from git tags | Semantic version string |
| `summary` | Short description (≤79 chars) | Required for store listing |
| `description` | Multi-line description | Required for store listing |
| `license` | `GPL-3.0` | SPDX identifier |
| `confinement` | `strict` | Required for stable channel |
| `base` | `core24` | Ubuntu 24.04 LTS |

### Application Declaration

| App Name | Command | Snap Path |
|----------|---------|-----------|
| `mytets` | `bin/mytets` | `/snap/mytets/current/bin/mytets` |

After installation, the following commands are available:

```
$ mytets version          # Print version
$ mytets one              # Print one random phrase
$ mytets list             # List all phrases
$ mytets --help           # Show help
$ mytets one --help       # Show command help
$ mytets --output json version   # JSON output
```

### Installation Contract

```
# Install from local build
$ sudo snap install mytets_1.0.1_amd64.snap --dangerous

# Install from Snap Store (after upload)
$ sudo snap install mytets

# Remove
$ sudo snap remove mytets
```

---

## Contract 3: Snapcraft Configuration Schema

**Interface type**: Build configuration  
**Consumer**: Snapcraft build system

### Expected `snapcraft.yaml` structure

```yaml
name: mytets
adopt-info: mytets
base: core24
summary: <short description>
description: |
  <multi-line description>
confinement: strict
license: GPL-3.0

apps:
  mytets:
    command: bin/mytets

parts:
  mytets:
    plugin: nil
    source: ../..
    source-type: git
    build-snaps:
      - go/latest/stable
    override-build: |
      <version derivation from git>
      <craftctl set version>
      <go build with ldflags>
```

### Build Environment

| Requirement | Source |
|-------------|--------|
| Go toolchain | `build-snaps: [go/latest/stable]` |
| Git history | `source-type: git` (clones repo with tags) |
| Build isolation | Snapcraft managed (LXD or Multipass) |

---

## Versioning

This contract follows the application's existing versioning scheme:
- Version is derived from git tags via `git describe --tags --always`
- Falls back to `dev` when no tags exist
- Version appears in both the snap metadata and the binary's `version` command output
