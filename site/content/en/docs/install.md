---
title: Install ccam
description: Install ccam on Linux or macOS with one command, from a .deb package, or from source. Claude Code must be installed first.
updated: 2026-10-19
order: 1
translation: source
---

Run this in a terminal on Linux or macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

The script downloads the release for your system, checks its SHA-256 checksum and puts two commands in `~/.local/bin`: `ccam` and its short alias `ccs`. If that folder is not on your `PATH`, the script prints the line to add to `~/.bashrc` or `~/.zshrc`.

ccam starts Claude Code for you, so install [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) first and make sure `claude --version` works.

## Supported systems

- Linux on x86-64 and ARM64: Ubuntu, Debian, Fedora, Arch, Alpine and other distributions. The binary is static and has no dependencies.
- macOS on Apple silicon and Intel.
- Windows is not supported directly. Inside WSL, ccam works like on any Linux system.

## Other ways to install

**Debian and Ubuntu package.** Download `ccam_<version>_linux_amd64.deb` from the [releases page](https://github.com/jhyunwoo/ccam/releases) and install it:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

The package installs `/usr/bin/ccam` and the `/usr/bin/ccs` link.

**From source.** With Go 1.26 or newer:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
```

**A specific version.** Set `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Shell completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Completion covers both `ccam` and `ccs`, including account names.

## Check the installation

```sh
ccam version
ccam doctor
```

`ccam doctor` checks that Claude Code is found and that every account folder, link and login is in order. Continue with the [quick start](/docs/quick-start).

## Uninstall

Remove the two commands and the account shortcuts, then the data folders if you no longer need the accounts:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Deleting `~/.local/share/ccam` also deletes the logins and conversation history stored there. A `~/.claude` folder that you [imported](/docs/commands) stays where it is.
