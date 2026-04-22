# Implementation Plan: Snap Packaging

**Branch**: `009-snap-packaging` | **Date**: 2026-04-22 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/009-snap-packaging/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Add snap packaging support so the mytets CLI application can be distributed via the Linux Snap Store. Introduces a `Makefile` at the repository root with `build` and `snap` targets, and a `snapcraft.yaml` configuration at `packaging/snap/snapcraft.yaml`. The snap uses `core24` base, `strict` confinement, and the `nil` plugin with `override-build` to compile from source with version injection from git tags.

## Technical Context

**Language/Version**: Go 1.26.2 (with Go 1.25+ compatibility)
**Primary Dependencies**: `snapcraft` (build-time tool, not a Go dependency), `make` (build orchestration)
**Storage**: N/A (no data storage changes)
**Testing**: Manual verification — `make build`, `make snap`, `snap install --dangerous`, command invocation
**Target Platform**: Linux amd64 (snap packages)
**Project Type**: CLI tool — packaging/distribution concern
**Performance Goals**: N/A (packaging feature, no runtime performance impact)
**Constraints**: Snap must build on developer's local Linux machine; snapcraft must be installed
**Scale/Scope**: 2 new files (Makefile, snapcraft.yaml), 0 Go code changes

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Clean, Self-Explained Code | PASS | Makefile targets and snapcraft.yaml are self-documenting with clear naming |
| II. Simplicity | PASS | Single `make snap` command; no unnecessary flags or options |
| III. Reliability | PASS | Makefile checks for snapcraft availability; clear errors on failure |
| IV. Performance | PASS | No runtime impact; packaging is build-time only |
| V. Extensibility | PASS | `packaging/` directory structure allows adding `debian/`, `flatpak/`, etc. later |
| VI. Documentation | PASS | README update with snap installation instructions |
| VII. Distribution | PASS | Directly supports distribution principle — adds snap as distribution channel. Single static binary preserved inside the snap |
| VIII. Go Best Practices | PASS | No Go code changes; Makefile uses same ldflags build approach |
| Constraints: Go 1.25+ | PASS | No Go source changes |
| Constraints: No runtime file I/O | PASS | Snap packaging is build-time only |
| Constraints: Version via ldflags | PASS | `git describe --tags` feeds into `-ldflags` |
| Testing: 80% coverage | PASS | No new Go code to cover |
| Security: Input validation | PASS | No new user inputs |

**Gate result**: ALL PASS — proceed to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/009-snap-packaging/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── cli-snap-contract.md
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
Makefile                       # NEW — build orchestration (make build, make snap)
packaging/
└── snap/
    └── snapcraft.yaml         # NEW — snap build configuration
```

**Structure Decision**: The `packaging/` directory follows the user's chosen convention from clarification (Option C). Each packaging format gets its own subdirectory (`packaging/snap/`, future `packaging/debian/`, etc.). The Makefile lives at the repository root as is conventional for Make-based projects.

## Complexity Tracking

No constitution violations — this table is intentionally empty.
