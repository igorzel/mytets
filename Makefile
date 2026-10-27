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
STORE_CREDS_FILE ?= ./snapcraft.login

# Loads Store credentials from STORE_CREDS_FILE (created by snap-login's
# keyring-free fallback) into SNAPCRAFT_STORE_CREDENTIALS for the rest of the
# recipe, so maintainers never need to export it manually. If the file is
# missing, falls back to whatever snapcraft already has in the OS keyring.
LOAD_CREDS = if [ -f "$(STORE_CREDS_FILE)" ]; then export SNAPCRAFT_STORE_CREDENTIALS="$$(cat $(STORE_CREDS_FILE))"; fi

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
	@$(LOAD_CREDS); \
	sg lxd -c "snapcraft register $(SNAP_NAME)"

# Logs in to the Snap Store, picking the right method automatically:
# tries the normal OS-keyring login first (interactive, same as before);
# if the environment has no working keyring (WSL2, headless servers, ...),
# snapcraft fails instantly with "No keyring found" before prompting for
# anything, so we transparently retry via the keyring-free export-login
# flow instead, writing credentials to STORE_CREDS_FILE. Every other
# snap-* target below picks these up automatically — no manual
# `export SNAPCRAFT_STORE_CREDENTIALS=...` is ever needed.
snap-login:
	@ERR_LOG="$$(mktemp)"; \
	sg lxd -c "snapcraft login" 2>"$$ERR_LOG"; \
	STATUS=$$?; \
	if [ $$STATUS -eq 0 ]; then \
		cat "$$ERR_LOG" >&2; rm -f "$$ERR_LOG"; \
	elif grep -qi "No keyring found" "$$ERR_LOG"; then \
		rm -f "$$ERR_LOG"; \
		if [ -f "$(STORE_CREDS_FILE)" ]; then \
			echo "No OS keyring available here; already logged in via $(STORE_CREDS_FILE)."; \
			echo "(Remove that file first if you need to re-authenticate.)"; \
		else \
			echo "No OS keyring available here; falling back to file-based login."; \
			sg lxd -c "snapcraft export-login '$(STORE_CREDS_FILE)'"; \
			echo ""; \
			echo "Credentials written to $(STORE_CREDS_FILE) — keep this file secret, do not commit it."; \
		fi; \
	else \
		cat "$$ERR_LOG" >&2; rm -f "$$ERR_LOG"; \
		exit $$STATUS; \
	fi

# Builds, uploads, and releases a new revision to SNAP_CHANNEL in one step.
snap-publish: snap
	@$(LOAD_CREDS); \
	SNAP_FILE="$$(ls -1t $(SNAP_NAME)_*.snap 2>/dev/null | head -n1)"; \
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
	fi; \
	$(LOAD_CREDS); \
	sg lxd -c "snapcraft release $(SNAP_NAME) '$(REVISION)' '$(SNAP_CHANNEL)'"

snap-status:
	@$(LOAD_CREDS); \
	sg lxd -c "snapcraft status $(SNAP_NAME)"

# Lists every uploaded revision, including ones not released to any channel
# yet (snap-status only shows the revision currently live per channel).
snap-revisions:
	@$(LOAD_CREDS); \
	sg lxd -c "snapcraft revisions $(SNAP_NAME)"

snap-clean:
	sg lxd -c "snapcraft clean mytets"

# ── Maintenance ───────────────────────────────────────────────────────────────

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(SNAP_NAME)_*.snap

.PHONY: setup build snap snap-register snap-login snap-publish snap-release snap-status snap-revisions snap-clean clean
