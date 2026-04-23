# ──────────────────────────────────────────────────────────────────────────────
# mytets — Build & Packaging
# ──────────────────────────────────────────────────────────────────────────────

MODULE      := github.com/igorzel/mytets
BINARY_NAME := mytets
BUILD_DIR   := ./bin
VERSION     := $(shell git describe --tags --always 2>/dev/null || echo "dev")
LDFLAGS     := -X $(MODULE)/internal/version.Version=$(VERSION)
SNAP_NAME   := mytets
SNAP_CHANNEL ?= edge
STORE_CREDS_FILE ?=

# ── Build ─────────────────────────────────────────────────────────────────────

build:
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/mytets

# ── Packaging ─────────────────────────────────────────────────────────────────

snap: snap-clean
	snapcraft pack

snap-register:
	snapcraft register $(SNAP_NAME)

snap-login:
	snapcraft login

snap-login-file:
	@if [ -z "$(STORE_CREDS_FILE)" ]; then \
		echo "Usage: make snap-login-file STORE_CREDS_FILE=./snapcraft.login"; \
		exit 1; \
	fi
	snapcraft login --with "$(STORE_CREDS_FILE)"

snap-upload: snap
	@SNAP_FILE="$$(ls -1t $(SNAP_NAME)_*.snap 2>/dev/null | head -n1)"; \
	if [ -z "$$SNAP_FILE" ]; then \
		echo "No snap artifact found. Build first with 'make snap'."; \
		exit 1; \
	fi; \
	echo "Uploading $$SNAP_FILE"; \
	snapcraft upload "$$SNAP_FILE"

snap-publish: snap
	@SNAP_FILE="$$(ls -1t $(SNAP_NAME)_*.snap 2>/dev/null | head -n1)"; \
	if [ -z "$$SNAP_FILE" ]; then \
		echo "No snap artifact found. Build first with 'make snap'."; \
		exit 1; \
	fi; \
	echo "Uploading and releasing $$SNAP_FILE to channel $(SNAP_CHANNEL)"; \
	snapcraft upload "$$SNAP_FILE" --release="$(SNAP_CHANNEL)"

snap-release:
	@if [ -z "$(REVISION)" ]; then \
		echo "Usage: make snap-release REVISION=<store-revision> SNAP_CHANNEL=<edge|beta|candidate|stable>"; \
		exit 1; \
	fi
	snapcraft release $(SNAP_NAME) "$(REVISION)" "$(SNAP_CHANNEL)"

snap-status:
	snapcraft status $(SNAP_NAME)

snap-clean:
	snapcraft clean mytets

# ── Maintenance ───────────────────────────────────────────────────────────────

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(SNAP_NAME)_*.snap

.PHONY: build snap snap-register snap-login snap-login-file snap-upload snap-publish snap-release snap-status snap-clean clean
