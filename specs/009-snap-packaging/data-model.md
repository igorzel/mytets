# Data Model: Snap Packaging

**Feature**: 009-snap-packaging  
**Date**: 2026-04-22

## Overview

This feature introduces no new Go entities or runtime data structures. It adds build-time configuration files (Makefile, snapcraft.yaml) that define how the application is compiled and packaged. The "data model" here describes the structure and relationships of these configuration artifacts.

## Entities

### 1. Makefile

**What it represents**: The build orchestration file at the repository root that provides standardized targets for building and packaging.

**Key attributes**:
- `VERSION`: Derived from `git describe --tags --always`, fallback `dev`
- `LDFLAGS`: Go linker flags for version injection
- `BINARY_NAME`: `mytets`
- `BUILD_DIR`: `./bin`

**Targets**:
| Target | Description | Dependencies |
|--------|-------------|--------------|
| `build` | Compile Go binary with version injection | Go toolchain |
| `snap` | Build snap package via snapcraft | snapcraft |
| `clean` | Remove build artifacts | None |

**Relationships**: The `snap` target invokes `snapcraft` which reads `packaging/snap/snapcraft.yaml`.

### 2. Snapcraft Configuration (`snapcraft.yaml`)

**What it represents**: The snap build definition that tells snapcraft how to compile the Go source, what metadata to attach, and how to expose the application.

**Key attributes**:

| Field | Value | Purpose |
|-------|-------|---------|
| `name` | `mytets` | Snap Store identifier |
| `base` | `core24` | Ubuntu 24.04 LTS runtime |
| `confinement` | `strict` | Required for public Snap Store |
| `license` | `GPL-3.0` | Matches project LICENSE |
| `adopt-info` | `mytets` | Version derived from build part |
| `summary` | Short description | Snap Store listing |
| `description` | Extended description | Snap Store listing |

**Parts**:

| Part | Plugin | Source | Purpose |
|------|--------|--------|---------|
| `mytets` | `nil` | `../..` (repo root via git) | Compile Go binary with ldflags |

**Apps**:

| App | Command | Description |
|-----|---------|-------------|
| `mytets` | `bin/mytets` | Main CLI entry point |

**Relationships**: References the Go source at the repository root. Produces a `.snap` file containing the compiled binary.

### 3. Snap Package (`.snap` file)

**What it represents**: The distributable artifact produced by `snapcraft`, ready for local installation or Snap Store upload.

**Key attributes**:
- File name pattern: `mytets_<version>_<arch>.snap`
- Contains: compiled `mytets` binary at `/snap/mytets/current/bin/mytets`
- Metadata: name, version, summary, description, license, confinement
- Architecture: host architecture (typically `amd64`)

**Relationships**: Produced by snapcraft from the snapcraft.yaml configuration. Installed via `snap install`.

## State Transitions

N/A — No runtime state changes. This feature is purely build-time.

## Validation Rules

- Makefile `VERSION` must be a non-empty string (falls back to `dev` if no git tags)
- `snapcraft.yaml` must pass `snapcraft lint` validation
- The produced `.snap` file must contain the `mytets` binary at the expected path
- The binary inside the snap must report the correct version (not `dev` when built from a tagged commit)

## File Layout

```text
Repository root
├── Makefile                          # Build orchestration
├── packaging/
│   └── snap/
│       └── snapcraft.yaml            # Snap build configuration
├── cmd/mytets/main.go                # (existing) Go entry point
├── internal/                         # (existing) Go packages
└── bin/                              # Build output (gitignored)
    └── mytets                        # Compiled binary from `make build`
```
