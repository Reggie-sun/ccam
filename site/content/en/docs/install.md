---
title: Install ccam
metaTitle: Install ccam on Linux, macOS and Windows
description: Install ccam on Linux, macOS or Windows with one command, a .deb package or from source, and update it with ccam update. Claude Code must be installed first.
updated: 2026-10-20
published: 2026-10-19
order: 1
translation: source
---

Run this in a terminal on Linux or macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

The script downloads the release for your system, checks its SHA-256 checksum and puts `ccam` in `~/.local/bin`, where ccam also creates the shortcut commands you choose later, such as `ccwork`. If that folder is not on your `PATH`, the script prints the line to add to `~/.bashrc` or `~/.zshrc`.

On Windows, run this in PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

It downloads `ccam_windows_<arch>.zip`, checks its SHA-256 checksum, puts `ccam.exe` in `%USERPROFILE%\.local\bin`, the folder where Claude Code's installer puts `claude.exe`, and adds that folder to your user `PATH` if it is not there yet. Open a new terminal afterwards.

ccam starts Claude Code for you, so install [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) first and make sure `claude --version` works.

## Supported systems

- Linux on x86-64 and ARM64: Ubuntu, Debian, Fedora, Arch, Alpine and other distributions. The binary is static and has no dependencies.
- macOS on Apple silicon and Intel.
- Windows 10 and 11 on x86-64 and ARM64, in PowerShell, Windows Terminal or Git Bash. Shared folders work for every user. Shared files (`settings.json`, `CLAUDE.md`, `history.jsonl`) need Developer Mode (Settings > System > For developers), because Windows lets only that mode, or an administrator, create links to files; see [sharing on Windows](/features/share#on-windows).
- Inside WSL, ccam works like on any Linux system.

## Other ways to install

**Debian and Ubuntu package.** Download `ccam_<version>_linux_amd64.deb` from the [releases page](https://github.com/jhyunwoo/ccam/releases) and install it:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

The package installs `/usr/bin/ccam`.

**From source.** With Go 1.26 or newer:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**A specific version.** Set `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Shell completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Completion covers every `ccam` command and option, including account names.

## Check the installation

```sh
ccam version
ccam doctor
```

`ccam doctor` checks that Claude Code is found and that every account folder, link and login is in order. Continue with the [quick start](/docs/quick-start).

## Update

```sh
ccam update --check
ccam update
```

`ccam update --check` only tells you whether a newer release exists. `ccam update` downloads the release for your system, checks its SHA-256 checksum against the release's `checksums.txt` and then replaces the `ccam` file in one step. If the checksum does not match, nothing is changed. Shortcuts such as `ccwork` are links to ccam, so they run the new version right away.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam never checks for updates on its own; it looks for a new release only when you run this command. To install a particular release, including an older one, pass `--version`. `ccam --dry-run update` shows what would be downloaded and replaced without changing anything.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Running the install script again also updates the copy in `~/.local/bin`. If you built ccam with `go install`, run the same command again.

**Windows.** `ccam update` downloads the `.zip` for Windows. A running program cannot be replaced on Windows, so the old binary is renamed to `ccam.exe.old`, which the next update deletes.

**Debian and Ubuntu package.** `ccam update` does not replace a `/usr/bin/ccam` that apt installed, and prints these steps instead. Download the new `ccam_<version>_linux_amd64.deb` from the [releases page](https://github.com/jhyunwoo/ccam/releases) and install it over the old one:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Uninstall

Remove the account shortcuts first, while ccam is still installed. `ccam list` shows them in the Shortcut column; remove each one with `ccam shortcut <account> --remove`, or delete those links from `~/.local/bin` yourself. Then remove the command, and the data folders if you no longer need the accounts:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

On Windows, in PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Deleting `~/.local/share/ccam` also deletes the logins and conversation history stored there. A `~/.claude` folder that you [imported](/docs/commands) stays where it is.
