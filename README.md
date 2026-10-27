# mytets

[Українською](README_uk.md)

`mytets` is a small command-line utility that prints random phrases from the
works of Ukrainian writer and playwright [Les Podervianskyi](https://uk.wikipedia.org/wiki/%D0%9B%D0%B5%D1%81%D1%8C_%D0%9F%D0%BE%D0%B4%D0%B5%D1%80%D0%B2%27%D1%8F%D0%BD%D1%81%D1%8C%D0%BA%D0%B8%D0%B9).
Get a single phrase, print a batch of unique phrases, or choose JSON output for
use in scripts. The phrases are included with the program, so you do not need
an account or an internet connection to use it.

## Install

The recommended option is the Snap Store:

```bash
sudo snap install mytets
```

To try the latest development build instead of the stable release:

```bash
sudo snap install mytets --channel=edge
```

You can also install with Go:

```bash
go install github.com/igorzel/mytets/cmd/mytets@latest
```

Make sure Go's bin directory is on your `PATH` if the `mytets` command is not
found after installation.

## Quick Start

Print one random phrase:

```bash
mytets one
```

Print five unique phrases (the default):

```bash
mytets list
```

Choose how many phrases to print:

```bash
mytets list --count 3
```

Each command prints plain text by default. Add `--output json` when you need
machine-readable output:

```bash
mytets one --output json
mytets list --count 3 --output json
mytets version --output json
```

For example, `mytets one --output json` returns:

```json
{"message":"Ето кайф, марожено!"}
```

## Commands

| Command | What it does |
|---------|--------------|
| `mytets one` | Print one random phrase. |
| `mytets list` | Print up to five unique phrases. |
| `mytets list --count N` | Print up to `N` unique phrases. |
| `mytets version` | Print the installed version. |
| `mytets --help` | Show available commands and options. |

Use `--output json` with `one`, `list`, or `version` to get JSON. Use
`mytets <command> --help` for help with a specific command.

## Build From Source

To build the executable yourself, install Go and run:

```bash
git clone https://github.com/igorzel/mytets.git
cd mytets
make build
./bin/mytets one
```

To build and install a local Snap package:

```bash
make snap
sudo snap install ./mytets_*.snap --dangerous
```

## Snap Store Publishing

These steps are for project maintainers who publish new releases of `mytets`
to the Snap Store. Run them in order on a new machine; steps 1–3 only need
to be done once, step 4 is the routine release workflow.

### 1. Prepare the build environment (once per machine, needs `sudo`)

```bash
make setup
```

This installs `snapcraft`/`lxd` and configures LXD. On first run it may
print "Log out and back in for lxd group membership to take effect" — you
can ignore this and continue right away: every command below runs through
`sg lxd`, which picks up the new group membership immediately without a
fresh login session.

### 2. Log in to the Snap Store (once per machine)

You need a Snap Store (Ubuntu One) account that's already an owner or
collaborator of the `mytets` snap name — if you're a new maintainer, ask an
existing one to add you via https://dashboard.snapcraft.io/ first,
otherwise the commands in step 4 will fail with a permission error.

```bash
make snap-login
```

`make snap-login` first tries the normal OS-keyring login (Secret Service /
gnome-keyring / KWallet), prompting for your Ubuntu One email/password/2FA.
On machines without a working keyring — e.g. WSL2 or a headless server —
snapcraft fails instantly with `No keyring found`, and `make snap-login`
transparently falls back to a keyring-free login instead, writing the
credentials to `./snapcraft.login` (keep this file secret; it's already
git-ignored). Either way, every other `snap-*` target below picks up the
resulting credentials automatically — no manual
`export SNAPCRAFT_STORE_CREDENTIALS=...` is ever needed. If you already
have a `./snapcraft.login` file, re-running `make snap-login` on a
keyring-less machine is a no-op; remove the file first to re-authenticate.

### 3. Register the snap name (already done, one-time-ever)

The `mytets` name is already registered on the Snap Store, so you normally
never need this. It's only required once, ever, by whoever first claimed
the name — running it again will fail with "name already taken", so skip
this step unless you're publishing under a brand-new snap name:

```bash
make snap-register
```

### 4. Release a new version

Publish a build to the edge channel:

```bash
make snap-publish SNAP_CHANNEL=edge
```

Check the currently released revision per channel:

```bash
make snap-status
```

To promote an existing revision later, for example to stable:

```bash
make snap-release REVISION=<store-revision> SNAP_CHANNEL=stable
```

`make snap-status` only reports the revision currently live per channel. To
see every uploaded revision, including ones not released to any channel:

```bash
make snap-revisions
```
