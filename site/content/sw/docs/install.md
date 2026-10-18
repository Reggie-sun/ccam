---
title: Sakinisha ccam
description: Sakinisha ccam kwenye Linux au macOS kwa amri moja, kutoka kifurushi cha .deb, au kutoka msimbo chanzo. Claude Code lazima isakinishwe kwanza.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Endesha hii kwenye terminal ya Linux au macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Skripti hupakua toleo la mfumo wako, hukagua checksum yake ya SHA-256 na huweka amri mbili katika `~/.local/bin`: `ccam` na jina lake fupi mbadala `ccs`. Ikiwa folda hiyo haipo kwenye `PATH` yako, skripti huchapisha mstari wa kuongeza kwenye `~/.bashrc` au `~/.zshrc`.

ccam hukuanzishia Claude Code, kwa hiyo sakinisha [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) kwanza na uhakikishe kuwa `claude --version` inafanya kazi.

## Mifumo inayotumika

- Linux kwenye x86-64 na ARM64: Ubuntu, Debian, Fedora, Arch, Alpine na distribution nyingine. Faili ya binary ni static na haina utegemezi wowote.
- macOS kwenye Apple silicon na Intel.
- Windows haitumiki moja kwa moja. Ndani ya WSL, ccam hufanya kazi kama kwenye mfumo wowote wa Linux.

## Njia nyingine za kusakinisha

**Kifurushi cha Debian na Ubuntu.** Pakua `ccam_<version>_linux_amd64.deb` kutoka [ukurasa wa matoleo](https://github.com/jhyunwoo/ccam/releases) na ukisakinishe:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Kifurushi husakinisha `/usr/bin/ccam` na kiungo `/usr/bin/ccs`.

**Kutoka msimbo chanzo.** Ukiwa na Go 1.26 au mpya zaidi:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
```

**Toleo mahususi.** Weka `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Ukamilishaji wa shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Ukamilishaji unashughulikia `ccam` na `ccs`, pamoja na majina ya akaunti.

## Kagua usakinishaji

```sh
ccam version
ccam doctor
```

`ccam doctor` hukagua kuwa Claude Code inapatikana na kwamba kila folda ya akaunti, kiungo na taarifa za kuingia viko sawa. Endelea na [mwongozo wa kuanza haraka](/docs/quick-start).

## Ondoa ccam

Ondoa amri hizo mbili na njia fupi za akaunti, kisha folda za data ikiwa huhitaji tena akaunti hizo:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Kufuta `~/.local/share/ccam` pia hufuta taarifa za kuingia na historia ya mazungumzo zilizohifadhiwa humo. Folda ya `~/.claude` ambayo [uliiingiza](/docs/commands) hubaki pale ilipo.
