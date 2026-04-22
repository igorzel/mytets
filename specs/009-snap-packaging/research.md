# Research: Snap Packaging

**Feature**: 009-snap-packaging  
**Date**: 2026-04-22

## Research Task 1: Snapcraft Go Plugin Capabilities (core24)

**Decision**: Use `nil` plugin with `override-build` instead of the built-in `go` plugin

**Rationale**: The current snapcraft Go plugin (via `craft-parts`) uses `go install ./...` with **no ldflags support**. It supports `go-buildtags` and `go-generate`, but there is no mechanism to pass `-ldflags` for version injection (FR-006). The legacy plugin used `$SNAPCRAFT_GO_LDFLAGS`, but this was removed in the craft-parts rewrite. Using the `nil` plugin with `override-build` provides full control over the `go build` command, including ldflags injection.

**Alternatives considered**:
- **Go plugin with `override-build`**: Combines Go plugin setup (dependency download, GOBIN) with custom build commands. Workable but redundant — the plugin's `go mod download` and GOBIN setup add little value when we override the build entirely.
- **Go plugin with `GOFLAGS` env var**: Set `GOFLAGS=-ldflags=...` via `build-environment`. Fragile — flag quoting issues with nested strings, and couples environment config to build logic.
- **Go plugin as-is**: Would build with `go install ./...` but cannot inject version. Rejected because FR-006 requires version injection.

## Research Task 2: Version Injection in Snap Builds

**Decision**: Use `source-type: git` in snapcraft.yaml so the cloned source retains git history; derive version from `git describe --tags --always` inside `override-build`; use `adopt-info` + `craftctl set version` to propagate version to snap metadata.

**Rationale**: Snap builds run in an isolated environment. With `source-type: local`, snapcraft copies files without `.git`, making `git describe` unavailable. Using `source-type: git` with a local path triggers a `git clone`, preserving tags and history. The `adopt-info` mechanism lets the build scriptlet set the snap's display version dynamically, eliminating hardcoded version strings.

**Alternatives considered**:
- **Makefile passes version as env var**: The Makefile computes `VERSION` and passes it to snapcraft. Problem: snapcraft's build environment is isolated; host env vars don't propagate into the build container without explicit `build-environment` mapping. This creates a fragile coupling.
- **VERSION file in repo**: Read a static file at build time. Problem: requires manual maintenance and diverges from the git-tags-as-truth-source decision.
- **`snapcraftctl set-version` (legacy)**: Replaced by `craftctl set version` in core24. Not available.

## Research Task 3: Custom snapcraft.yaml Location

**Decision**: Place `snapcraft.yaml` at `packaging/snap/snapcraft.yaml`; invoke `snapcraft` from the `packaging/snap/` directory via the Makefile; use `source: ../..` with `source-type: git` to reference the repo root.

**Rationale**: The user chose `packaging/snap/` as the packaging directory (Clarification Session 2026-04-22). Snapcraft resolves `source` paths relative to the directory containing `snapcraft.yaml`. With `source: ../..` and `source-type: git`, snapcraft clones the repo root, preserving git history for version derivation. The Makefile `cd`s into the directory before invoking `snapcraft`, keeping the command simple.

**Alternatives considered**:
- **`snap/snapcraft.yaml` (default location)**: Snapcraft convention. Rejected because the user chose a `packaging/` directory structure to accommodate future packaging formats.
- **`snapcraft --project-dir packaging/snap` from repo root**: Snapcraft's `--project-dir` flag. Rejected due to inconsistent behavior with relative source paths across snapcraft versions.
- **Symlink**: `snap/snapcraft.yaml` → `packaging/snap/snapcraft.yaml`. Rejected — adds indirection and can confuse tooling.

## Research Task 4: Snap Base and Confinement for CLI Tools

**Decision**: Use `base: core24` with `confinement: strict`

**Rationale**: `core24` is based on Ubuntu 24.04 LTS (supported until 2029) and is the current recommended base for new snaps. `strict` confinement is required for public Snap Store distribution and is appropriate for a CLI tool like mytets that needs no special system access (no network, no filesystem, no hardware). The application reads only embedded data and writes to stdout/stderr.

**Alternatives considered**:
- **`core22`** (Ubuntu 22.04 LTS): Older but still supported. Rejected per user choice (Clarification Session 2026-04-22).
- **`classic` confinement**: Provides unrestricted access. Rejected — unnecessary for this application, and requires Snap Store review approval. `strict` is simpler.
- **`devmode` confinement**: For development only. Not publishable to stable channel. Rejected.

## Research Task 5: Makefile Conventions for Go Projects

**Decision**: Standard GNU Make with `.PHONY` targets, variable-based version derivation, and modular target structure.

**Rationale**: Go projects conventionally use a Makefile with `build`, `test`, `clean`, and packaging targets. Version is derived once via `$(shell git describe ...)` and reused across targets. `.PHONY` declarations prevent conflicts with file names. The structure naturally extends to additional targets (`deb`, `flatpak`) per FR-007.

**Key patterns**:
- `VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")` — single source of truth
- `LDFLAGS := -X github.com/igorzel/mytets/internal/version.Version=$(VERSION)` — reusable flag
- `BUILD_DIR := ./bin` — consistent output location
- Target for `snap` invokes snapcraft from `packaging/snap/`
- `clean` target removes build artifacts

**Alternatives considered**:
- **Task runners (Just, Mage)**: More expressive but add a dependency. Rejected per constitution (minimize external dependencies).
- **Shell scripts**: No dependency tracking, harder to extend. Rejected.
- **Go-based build (Mage)**: Go-native but heavy for packaging orchestration. Rejected.

## Research Task 6: Go Snap Build Environment

**Decision**: Use `build-snaps: [go/latest/stable]` to install Go in the snap build environment.

**Rationale**: The `go` snap provides the Go toolchain inside the snap build container. Using `latest/stable` ensures compatibility with the project's Go 1.26.2 requirement. The Go snap is the standard way to get Go in snapcraft builds (recommended by Canonical). The binary is available on PATH automatically.

**Alternatives considered**:
- **`build-packages: [golang-go]`**: Uses the Ubuntu package. Problem: Ubuntu 24.04 ships Go 1.22, which is older than the project's Go 1.26.2 requirement. Rejected.
- **`build-snaps: [go/1.26/stable]`**: Pin to specific version. May not be available as a snap channel. Using `latest/stable` is safer and will include Go 1.26+ when available.
