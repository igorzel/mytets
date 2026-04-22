# ──────────────────────────────────────────────────────────────────────────────
# mytets — Build & Packaging
# ──────────────────────────────────────────────────────────────────────────────

MODULE      := github.com/igorzel/mytets
BINARY_NAME := mytets
BUILD_DIR   := ./bin
VERSION     := $(shell git describe --tags --always 2>/dev/null || echo "dev")
LDFLAGS     := -X $(MODULE)/internal/version.Version=$(VERSION)

# ── Build ─────────────────────────────────────────────────────────────────────

build:
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/mytets

# ── Packaging ─────────────────────────────────────────────────────────────────

snap:
	cd packaging/snap && snapcraft pack

# ── Maintenance ───────────────────────────────────────────────────────────────

clean:
	rm -rf $(BUILD_DIR)

.PHONY: build snap clean
