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

# ── Setup ─────────────────────────────────────────────────────────────────────
setup:
	@echo "Setting up the build environment..."
	@if ! sudo snap list snapcraft >/dev/null 2>&1; then sudo snap install snapcraft --classic; fi
	@if ! sudo snap list lxd >/dev/null 2>&1; then sudo snap install lxd; fi
	sudo usermod -aG lxd "$$USER"
	sudo apt-get --yes install util-linux-extra
	@if ! sg lxd -c "lxc storage show default" >/dev/null 2>&1 || ! sg lxd -c "lxc network show lxdbr0" >/dev/null 2>&1; then \
		echo "Initializing LXD storage and networking..."; \
		sg lxd -c "lxd init --auto"; \
	fi
	@if ! sg lxd -c "lxc project show snapcraft" >/dev/null 2>&1; then sg lxd -c "lxc project create snapcraft"; fi
	sg lxd -c "lxc project set snapcraft features.images=true features.profiles=true features.storage.buckets=true features.storage.volumes=true"
	@if ! sg lxd -c "lxc --project snapcraft profile device get default root path" >/dev/null 2>&1; then \
		sg lxd -c "lxc --project snapcraft profile device add default root disk path=/ pool=default"; \
	fi
	sg lxd -c "lxc --project snapcraft profile device set default root path=/ pool=default"
	@if ! sg lxd -c "lxc --project snapcraft profile device get default eth0 type" >/dev/null 2>&1; then \
		sg lxd -c "lxc --project snapcraft profile device add default eth0 nic name=eth0 network=lxdbr0"; \
	fi
	sg lxd -c "lxc --project snapcraft profile device set default eth0 name=eth0 network=lxdbr0"
	@sudo iptables -C FORWARD -i lxdbr0 -j ACCEPT 2>/dev/null || sudo iptables -I FORWARD -i lxdbr0 -j ACCEPT
	@sudo iptables -C FORWARD -o lxdbr0 -m state --state RELATED,ESTABLISHED -j ACCEPT 2>/dev/null || sudo iptables -I FORWARD -o lxdbr0 -m state --state RELATED,ESTABLISHED -j ACCEPT
	@if ! id -nG "$$USER" | tr ' ' '\n' | grep -qx lxd; then echo "Log out and back in for lxd group membership to take effect."; fi
	@echo "Build environment setup complete."


# ── Build ─────────────────────────────────────────────────────────────────────

build:
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/mytets

# ── Packaging ─────────────────────────────────────────────────────────────────

snap: snap-clean
	sg lxd -c "snapcraft pack"

snap-register:
	sg lxd -c "snapcraft register $(SNAP_NAME)"

snap-login:
	sg lxd -c "snapcraft login"

snap-login-file:
	@if [ -z "$(STORE_CREDS_FILE)" ]; then \
		echo "Usage: make snap-login-file STORE_CREDS_FILE=./snapcraft.login"; \
		exit 1; \
	fi
	sg lxd -c "snapcraft login --with '$(STORE_CREDS_FILE)'"

snap-upload: snap
	@SNAP_FILE="$$(ls -1t $(SNAP_NAME)_*.snap 2>/dev/null | head -n1)"; \
	if [ -z "$$SNAP_FILE" ]; then \
		echo "No snap artifact found. Build first with 'make snap'."; \
		exit 1; \
	fi; \
	echo "Uploading $$SNAP_FILE"; \
	sg lxd -c "snapcraft upload '$$SNAP_FILE'"

snap-publish: snap
	@SNAP_FILE="$$(ls -1t $(SNAP_NAME)_*.snap 2>/dev/null | head -n1)"; \
	if [ -z "$$SNAP_FILE" ]; then \
		echo "No snap artifact found. Build first with 'make snap'."; \
		exit 1; \
	fi; \
	echo "Uploading and releasing $$SNAP_FILE to channel $(SNAP_CHANNEL)"; \
	sg lxd -c "snapcraft upload '$$SNAP_FILE' --release='$(SNAP_CHANNEL)'"

snap-release:
	@if [ -z "$(REVISION)" ]; then \
		echo "Usage: make snap-release REVISION=<store-revision> SNAP_CHANNEL=<edge|beta|candidate|stable>"; \
		exit 1; \
	fi
	sg lxd -c "snapcraft release $(SNAP_NAME) '$(REVISION)' '$(SNAP_CHANNEL)'"

snap-status:
	sg lxd -c "snapcraft status $(SNAP_NAME)"

snap-clean:
	sg lxd -c "snapcraft clean mytets"

# ── Maintenance ───────────────────────────────────────────────────────────────

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(SNAP_NAME)_*.snap

.PHONY: setup build snap snap-register snap-login snap-login-file snap-upload snap-publish snap-release snap-status snap-clean clean
