# Feature Specification: Snap Packaging

**Feature Branch**: `009-snap-packaging`  
**Created**: 2026-04-22  
**Status**: Draft  
**Input**: User description: "Implement snap packaging so the mytets application can be uploaded to the Linux Snap Store"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Build a Snap Package Locally (Priority: P1)

As a developer, I want to build a snap package of the mytets application on my local machine by running a single command, so that I can produce a distributable `.snap` file without remembering complex build steps.

**Why this priority**: This is the core capability — without a working local snap build, nothing else in this feature is possible. It delivers the fundamental value of packaging the application for distribution.

**Independent Test**: Can be fully tested by running `make snap` on a machine with snapcraft installed and verifying a `.snap` file is produced in the expected location.

**Acceptance Scenarios**:

1. **Given** the developer has snapcraft installed and is in the project root, **When** they run `make snap`, **Then** a `.snap` file for the mytets application is produced in the project directory.
2. **Given** the developer runs `make snap`, **When** the build completes successfully, **Then** the snap file name includes the application name and version (e.g., `mytets_1.0.1_amd64.snap`).
3. **Given** the developer runs `make snap` without any prior setup beyond having snapcraft, **When** the build process runs, **Then** it compiles the Go application from source as part of the snap build and embeds the correct version.

---

### User Story 2 - Install and Run the Snap Locally (Priority: P2)

As a developer, I want to install the locally built snap and verify that all mytets commands work correctly, so I can validate the package before uploading it to the Snap Store.

**Why this priority**: Verifying the snap works after building is essential for quality assurance before public distribution.

**Independent Test**: Can be tested by installing the `.snap` file with `snap install --dangerous` and running `mytets version`, `mytets one`, and `mytets list` to verify all commands function correctly.

**Acceptance Scenarios**:

1. **Given** a `.snap` file has been built, **When** the developer installs it locally using `snap install --dangerous <snap-file>`, **Then** the `mytets` command becomes available system-wide.
2. **Given** the snap is installed locally, **When** the developer runs `mytets version`, **Then** the correct version string is returned (not "dev").
3. **Given** the snap is installed locally, **When** the developer runs `mytets one` and `mytets list`, **Then** the commands produce the expected output identical to running the Go binary directly.

---

### User Story 3 - Upload Snap to the Snap Store (Priority: P3)

As a developer, I want the snap metadata to be configured for public distribution on the Snap Store, so that anyone can discover and install mytets via `snap install mytets`.

**Why this priority**: Public availability is the end goal of this feature, but it depends on a correctly built and tested snap first.

**Independent Test**: Can be tested by reviewing the snap metadata (name, summary, description, license, confinement level) and confirming it meets Snap Store publishing requirements.

**Acceptance Scenarios**:

1. **Given** the snap is built, **When** the snap metadata is inspected, **Then** it contains a meaningful summary, description, and license field suitable for public listing.
2. **Given** the snap configuration, **When** reviewed against Snap Store requirements, **Then** the confinement level and all mandatory metadata fields are correctly set for public distribution.
3. **Given** the developer has a Snap Store account and has run `snapcraft login`, **When** they run `snapcraft upload <snap-file> --release=stable`, **Then** the snap is uploaded and available for public installation.

---

### User Story 4 - Build the Go Binary via Makefile (Priority: P1)

As a developer, I want a Makefile that standardizes building the Go binary, so that both manual builds and the snap build process use the same consistent build steps.

**Why this priority**: The Makefile is a prerequisite for `make snap` and establishes the project's build system. It is equally critical to User Story 1.

**Independent Test**: Can be tested by running `make build` and verifying the Go binary is produced with the correct version injected.

**Acceptance Scenarios**:

1. **Given** the developer is in the project root, **When** they run `make build`, **Then** the Go binary is compiled to a well-known output path.
2. **Given** the developer runs `make build`, **When** the build completes, **Then** the binary has the application version injected via ldflags (not "dev").

---

### Edge Cases

- What happens when snapcraft is not installed and the developer runs `make snap`? The build should fail with a clear error indicating snapcraft is required.
- What happens when the Go toolchain is not available during snap build? The snapcraft configuration declares `build-snaps: [go/latest/stable]` to install the Go toolchain inside the snap build environment.
- What happens when building on an architecture other than amd64? The snap should build for the host architecture by default.
- What happens when the version string is not set? The build process derives the version from git tags via `git describe --tags`; if no tags exist, it should fall back to a sensible default (e.g., `dev`).

## Clarifications

### Session 2026-04-22

- Q: How should make build and make snap determine the version string? → A: Derive from git tags via `git describe --tags`
- Q: Where should the snapcraft.yaml file live? → A: `packaging/snap/snapcraft.yaml` (custom packaging directory)
- Q: Which snap base image should be used? → A: `core24` (Ubuntu 24.04 LTS)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The project MUST include a `Makefile` at the repository root with at least a `snap` target that builds a snap package and a `build` target that compiles the Go binary.
- **FR-002**: Running `make snap` MUST produce a `.snap` file for the mytets application using snapcraft.
- **FR-003**: The snap MUST include the compiled `mytets` binary as its only application entry point, exposing all existing CLI commands (`version`, `one`, `list`).
- **FR-004**: The snap configuration MUST define the application name as `mytets`, with a user-facing summary (≤79 chars), multi-line description, and SPDX license identifier, as required by the Snap Store for public listing.
- **FR-005**: The snap MUST use `strict` confinement and `core24` (Ubuntu 24.04 LTS) as the base to meet Snap Store requirements for public distribution.
- **FR-006**: The snap build process MUST compile the Go source code from the repository (not rely on a pre-built binary) and inject the version string via ldflags, deriving the version from git tags using `git describe --tags`.
- **FR-007**: The Makefile MUST be designed to allow additional packaging targets in the future (e.g., `make deb`, `make flatpak`) without restructuring. Packaging configurations MUST reside under a `packaging/` directory (e.g., `packaging/snap/`, `packaging/debian/`).
- **FR-008**: The snap metadata MUST declare the `mytets` command so that after installation, users can run `mytets` directly from the command line.
- **FR-009**: The snap MUST be buildable on a developer's local machine using only `make snap` (assuming snapcraft is installed).

### Key Entities

- **Snap Package**: The distributable `.snap` file containing the compiled mytets binary and all metadata required by the Snap Store.
- **Snapcraft Configuration**: The `snapcraft.yaml` file located at `packaging/snap/snapcraft.yaml` that defines how the snap is built, its metadata, parts, and application entry points.
- **Makefile**: The project build file that provides standardized targets for building the binary and creating the snap package.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A developer can produce a working `.snap` file by running a single command (`make snap`) from the project root.
- **SC-002**: The snap-packaged application passes all functional tests — `mytets version`, `mytets one`, and `mytets list` produce correct output when run from the installed snap.
- **SC-003**: The snap package contains all required metadata for public Snap Store listing (name, version, summary, description, license, confinement).
- **SC-004**: The snap can be uploaded to the Snap Store and installed by any user via `snap install mytets`.
- **SC-005**: Adding a new packaging target (e.g., `make deb`) requires only adding a new Makefile target and packaging config, not restructuring existing files.

## Assumptions

- The developer has `snapcraft` installed on their local Linux machine (snap packaging is Linux-only).
- The developer has a registered Snap Store account and has authenticated via `snapcraft login` for uploading (uploading is a manual step, not automated by this feature).
- The application targets `amd64` architecture as the primary build target; multi-architecture support is out of scope for this feature.
- GitHub Actions automation for snap builds is out of scope for this feature and will be addressed in a future iteration.
- Other packaging formats (Debian, Flatpak, Docker, etc.) are out of scope for this feature; the Makefile is designed to accommodate them later.
- The application version is derived from git tags (`git describe --tags`) and injected at build time via Go ldflags, consistent with the existing build process.
- The `mytets` snap name is available on the Snap Store (or will be registered by the developer).
