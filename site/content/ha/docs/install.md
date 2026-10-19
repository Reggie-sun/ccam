---
title: Girka ccam
description: Girka ccam a Linux ko macOS da umarni ɗaya, daga kunshin .deb, ko daga lambar tushe. Dole ne ka fara girka Claude Code.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: fc63978e8f9e
---

Gudanar da wannan a terminal a Linux ko macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Rubutun (script) yana sauke sakin da ya dace da tsarinka, yana duba checksum ɗinsa na SHA-256, sannan yana saka `ccam` a `~/.local/bin`, inda daga baya `ccam add` ke ƙirƙirar gajerun umarnin `claude1`, `claude2` da sauransu. Idan wannan foldar ba ta cikin `PATH` ɗinka, rubutun yana buga layin da za ka ƙara a `~/.bashrc` ko `~/.zshrc`.

ccam ne ke fara maka Claude Code, don haka ka fara girka [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) kuma ka tabbata `claude --version` yana aiki.

## Tsarukan da ake tallafawa

- Linux a x86-64 da ARM64: Ubuntu, Debian, Fedora, Arch, Alpine da sauran distributions. Binary ɗin static ne kuma ba shi da dependencies.
- macOS a Apple silicon da Intel.
- Ba a tallafa wa Windows kai tsaye ba. A cikin WSL, ccam yana aiki kamar a kowane tsarin Linux.

## Wasu hanyoyin girkawa

**Kunshin Debian da Ubuntu.** Sauke `ccam_<version>_linux_amd64.deb` daga [shafin saki](https://github.com/jhyunwoo/ccam/releases) sannan ka girka shi:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Kunshin yana girka `/usr/bin/ccam`.

**Daga lambar tushe.** Da Go 1.26 ko sabo:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**Wata siga ta musamman.** Saita `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Cike umarni a shell (shell completion)

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Completion yana aiki ga dukkan umarni da zaɓuɓɓukan `ccam`, har da sunayen asusai.

## Duba girkawar

```sh
ccam version
ccam doctor
```

`ccam doctor` yana duba cewa an sami Claude Code, kuma cewa kowace foldar asusu, kowane mahaɗi da kowane shiga suna cikin tsari. Ci gaba da [farawa da sauri](/docs/quick-start).

## Cire ccam

Cire umarnin da gajerun umarnin asusai, sannan foldojin bayanai idan ba ka buƙatar asusan kuma:

```sh
rm ~/.local/bin/ccam ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Share `~/.local/share/ccam` yana share bayanan shiga da tarihin tattaunawa da aka ajiye a can ma. Foldar `~/.claude` da ka [shigo da ita](/docs/commands) tana nan a inda take.
