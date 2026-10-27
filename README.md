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

These steps are for project maintainers. First prepare the LXD environment on
each machine that will build Snap packages:

```bash
make setup
```

This setup requires `sudo`. Authenticate once per machine, then register the
Snap name:

```bash
make snap-export-login                 # one-time login; writes ./snapcraft.login (keep it secret)
make snap-register
```

`make snap-export-login` prompts for your Ubuntu One email/password/2FA and
stores the resulting credentials in `./snapcraft.login`. Every other
`snap-*` target below picks these credentials up automatically — no manual
`export SNAPCRAFT_STORE_CREDENTIALS=...` is needed.

Alternatively, `make snap-login` logs in via the OS keyring (Secret Service /
gnome-keyring / KWallet) instead of a file. This works fine on a desktop
Linux machine, but fails with `No keyring found` on machines without one —
e.g. WSL2 or a headless server — so `make snap-export-login` is the more
portable option and is what the rest of this guide assumes.

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
