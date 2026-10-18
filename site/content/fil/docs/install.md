---
title: I-install ang ccam
description: I-install ang ccam sa Linux o macOS gamit ang isang command, mula sa .deb package, o mula sa source. Kailangang naka-install muna ang Claude Code.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Patakbuhin ito sa terminal sa Linux o macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Dina-download ng script ang release para sa iyong system, sinusuri ang SHA-256 checksum nito at naglalagay ng dalawang command sa `~/.local/bin`: ang `ccam` at ang maikling alias nitong `ccs`. Kung wala sa iyong `PATH` ang folder na iyon, ipi-print ng script ang linyang idadagdag sa `~/.bashrc` o `~/.zshrc`.

Ang ccam ang nagsisimula ng Claude Code para sa iyo, kaya i-install muna ang [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) at tiyaking gumagana ang `claude --version`.

## Mga suportadong system

- Linux sa x86-64 at ARM64: Ubuntu, Debian, Fedora, Arch, Alpine at iba pang distribution. Static ang binary at wala itong dependency.
- macOS sa Apple silicon at Intel.
- Hindi direktang suportado ang Windows. Sa loob ng WSL, gumagana ang ccam gaya ng sa anumang Linux system.

## Iba pang paraan ng pag-install

**Package para sa Debian at Ubuntu.** I-download ang `ccam_<version>_linux_amd64.deb` mula sa [pahina ng mga release](https://github.com/jhyunwoo/ccam/releases) at i-install ito:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Ini-install ng package ang `/usr/bin/ccam` at ang link na `/usr/bin/ccs`.

**Mula sa source.** Gamit ang Go 1.26 o mas bago:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
```

**Isang partikular na bersyon.** Itakda ang `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Shell completion

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Saklaw ng completion ang `ccam` at `ccs`, pati ang mga pangalan ng account.

## Suriin ang installation

```sh
ccam version
ccam doctor
```

Sinusuri ng `ccam doctor` na nahahanap ang Claude Code at na maayos ang bawat account folder, link at login. Magpatuloy sa [quick start](/docs/quick-start).

## Pag-uninstall

Alisin ang dalawang command at ang mga shortcut ng account, saka ang mga data folder kung hindi mo na kailangan ang mga account:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Kapag binura mo ang `~/.local/share/ccam`, mabubura rin ang mga login at kasaysayan ng usapan na naka-store doon. Mananatili kung nasaan ito ang `~/.claude` folder na iyong [in-import](/docs/commands).
