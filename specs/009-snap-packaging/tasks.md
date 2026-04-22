# Tasks: Snap Packaging

**Input**: Design documents from `/specs/009-snap-packaging/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Not requested — no test tasks included.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create the packaging directory structure

- [X] T001 Create directory structure `packaging/snap/` at repository root

---

## Phase 2: Foundational — Makefile (US4, Blocking Prerequisite)

**Purpose**: Establish the Makefile build system that both `make build` and `make snap` depend on. Maps to **User Story 4** (P1).

**Goal**: A developer can run `make build` to compile the Go binary with version injection from git tags.

**Independent Test**: Run `make build && ./bin/mytets version` — should print the git-tag-derived version (or `dev`).

**⚠️ CRITICAL**: `make snap` (Phase 3) depends on the Makefile being complete.

- [X] T002 [US4] Create `Makefile` at repository root with variables: `VERSION` (from `git describe --tags --always`, fallback `dev`), `LDFLAGS` (linker flags for `github.com/igorzel/mytets/internal/version.Version`), `BINARY_NAME` (`mytets`), `BUILD_DIR` (`./bin`), `MODULE` (`github.com/igorzel/mytets`)
- [X] T003 [US4] Add `build` target to `Makefile` that compiles `./cmd/mytets` to `$(BUILD_DIR)/$(BINARY_NAME)` with `$(LDFLAGS)` via `go build`
- [X] T004 [US4] Add `clean` target to `Makefile` that removes `$(BUILD_DIR)`
- [X] T005 [US4] Add `.PHONY` declarations for all targets (`build`, `snap`, `clean`) in `Makefile`
- [X] T006 [US4] Add `bin/` to `.gitignore` at repository root (if not already present)

**Checkpoint**: `make build && ./bin/mytets version` produces correct version output

---

## Phase 3: User Story 1 — Build Snap Package Locally (Priority: P1) 🎯 MVP

**Goal**: Running `make snap` produces a `.snap` file for the mytets application.

**Independent Test**: Run `make snap` — a `mytets_<version>_amd64.snap` file is produced under `packaging/snap/`.

### Implementation for User Story 1

- [X] T007 [US1] Create `packaging/snap/snapcraft.yaml` with metadata: `name: mytets`, `adopt-info: mytets`, `base: core24`, `summary` (≤79 chars), multi-line `description`, `confinement: strict`, `license: GPL-3.0`
- [X] T008 [US1] Add `apps` section to `packaging/snap/snapcraft.yaml` declaring `mytets` app with `command: bin/mytets`
- [X] T009 [US1] Add `parts` section to `packaging/snap/snapcraft.yaml` with `mytets` part: `plugin: nil`, `source: ../..`, `source-type: git`, `build-snaps: [go/latest/stable]`, `override-build` scriptlet that derives version via `git describe --tags --always`, calls `craftctl set version=<version>`, runs `go build -ldflags "-X github.com/igorzel/mytets/internal/version.Version=<version>" -o $CRAFT_PART_INSTALL/bin/mytets ./cmd/mytets`
- [X] T010 [US1] Add `snap` target to `Makefile` that runs `cd packaging/snap && snapcraft`

**Checkpoint**: `make snap` completes successfully and produces a `.snap` file in `packaging/snap/`

---

## Phase 4: User Story 2 — Install and Run Snap Locally (Priority: P2)

**Goal**: The locally built snap installs and all commands (`version`, `one`, `list`) work correctly.

**Independent Test**: `sudo snap install packaging/snap/mytets_*.snap --dangerous && mytets version && mytets one && mytets list`

This story is verified manually — no code tasks. The snap configuration from Phase 3 must produce a working snap. Verification steps are documented in `specs/009-snap-packaging/quickstart.md`.

**Checkpoint**: All three commands produce correct output when run from the installed snap

---

## Phase 5: User Story 3 — Snap Store Metadata (Priority: P3)

**Goal**: The snap metadata meets all Snap Store requirements for public distribution.

**Independent Test**: Inspect `snapcraft.yaml` metadata fields and confirm all mandatory Snap Store fields are present and valid.

This story is satisfied by the snapcraft.yaml metadata created in T007. No additional code tasks needed — the metadata was designed to meet Snap Store requirements per contracts/cli-snap-contract.md (Contract 2).

**Checkpoint**: `snapcraft.yaml` contains name, version (via adopt-info), summary, description, license, confinement, and base — all required for public Snap Store listing

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and cleanup

- [X] T011 [P] Update `README.md` with snap installation section: install from Snap Store (`snap install mytets`), install from local build (`snap install --dangerous`), and Makefile usage (`make build`, `make snap`, `make clean`)
- [X] T012 Run `specs/009-snap-packaging/quickstart.md` validation: execute `make build`, `make snap`, local snap install, and verify all commands

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational / US4 (Phase 2)**: Depends on Phase 1 — BLOCKS Phase 3
- **US1 (Phase 3)**: Depends on Phase 2 (Makefile must exist for `snap` target)
- **US2 (Phase 4)**: Depends on Phase 3 (snap must be built to install)
- **US3 (Phase 5)**: Depends on Phase 3 (metadata in snapcraft.yaml)
- **Polish (Phase 6)**: Depends on Phase 3 completion

### Within Each Phase

- T002–T006 (Phase 2): T002 creates the file, T003–T005 add to it sequentially, T006 is independent [P]
- T007–T009 (Phase 3): T007 creates the file, T008–T009 add to it sequentially, T010 adds to Makefile

### Parallel Opportunities

- T006 (`.gitignore`) can run in parallel with T003–T005 (Makefile targets)
- T011 (README) can run in parallel with T012 (quickstart validation)
- Phase 4 and Phase 5 are both verification of Phase 3 output and can be validated in parallel

---

## Parallel Example: Phase 2

```bash
# Sequential (same file):
T002: Create Makefile with variables
T003: Add build target to Makefile
T004: Add clean target to Makefile
T005: Add .PHONY declarations to Makefile

# Parallel (different file):
T006: Add bin/ to .gitignore  (can run alongside T003–T005)
```

---

## Implementation Strategy

### MVP First (User Story 4 + User Story 1)

1. Complete Phase 1: Setup (directory structure)
2. Complete Phase 2: Makefile (US4) — `make build` works
3. Complete Phase 3: Snap config (US1) — `make snap` works
4. **STOP and VALIDATE**: Build snap, install locally, verify commands
5. Complete Phase 6: Polish (README, quickstart validation)

### Incremental Delivery

1. Setup + Makefile → `make build` works (US4 complete)
2. Add snapcraft.yaml + snap target → `make snap` works (US1 complete, MVP!)
3. Install and verify locally (US2 complete)
4. Review metadata for Snap Store (US3 complete)
5. Update docs → Feature complete

---

## Notes

- This feature introduces **0 Go code changes** — only build/packaging configuration files
- US2 and US3 are manual verification stories with no code tasks; they are satisfied by the quality of US1/US4 output
- The `nil` plugin with `override-build` is used instead of the Go plugin because the Go plugin has no ldflags support (see research.md, Task 1)
- `source-type: git` is critical for preserving git tags inside the snap build environment (see research.md, Task 2)
