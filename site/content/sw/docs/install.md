---
title: Sakinisha ccam
metaTitle: Sakinisha ccam kwenye Linux na macOS
description: Sakinisha ccam kwenye Linux au macOS kwa amri moja, kifurushi cha .deb au msimbo chanzo, na uisasishe kwa ccam update. Sakinisha Claude Code kwanza.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 7ecd1ed00e09
---

Endesha hii kwenye terminal ya Linux au macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Skripti hupakua toleo la mfumo wako, hukagua checksum yake ya SHA-256 na huweka `ccam` katika `~/.local/bin`, ambapo ccam pia huunda amri za njia fupi utakazochagua baadaye, kama `ccwork`. Ikiwa folda hiyo haiko kwenye `PATH` yako, skripti huchapisha mstari wa kuongeza kwenye `~/.bashrc` au `~/.zshrc`.

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

Kifurushi husakinisha `/usr/bin/ccam`.

**Kutoka msimbo chanzo.** Ukiwa na Go 1.26 au mpya zaidi:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
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

Ukamilishaji unashughulikia amri na chaguo zote za `ccam`, pamoja na majina ya akaunti.

## Kagua usakinishaji

```sh
ccam version
ccam doctor
```

`ccam doctor` hukagua kuwa Claude Code inapatikana na kwamba kila folda ya akaunti, kiungo na taarifa za kuingia viko sawa. Endelea na [mwongozo wa kuanza haraka](/docs/quick-start).

## Sasisha ccam

```sh
ccam update --check
ccam update
```

`ccam update --check` hukuambia tu kama kuna toleo jipya zaidi. `ccam update` hupakua toleo la mfumo wako, hulinganisha checksum yake ya SHA-256 na `checksums.txt` ya toleo hilo, kisha hubadilisha faili ya `ccam` kwa hatua moja. Ikiwa checksum hailingani, hakuna kinachobadilishwa. Njia fupi kama `ccwork` ni viungo vya ccam, kwa hiyo huendesha toleo jipya mara moja.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam kamwe haikagui masasisho yenyewe. Hutafuta toleo jipya tu unapoendesha amri hii. Ili kusakinisha toleo mahususi, hata la zamani, tumia `--version`. `ccam --dry-run update` huonyesha kile ambacho kingepakuliwa na kubadilishwa, bila kubadilisha chochote.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Kuendesha skripti ya usakinishaji tena pia husasisha nakala iliyo katika `~/.local/bin`. Ikiwa ulijenga ccam kwa `go install`, endesha amri ile ile tena.

**Kifurushi cha Debian na Ubuntu.** `ccam update` haibadilishi `/usr/bin/ccam` iliyosakinishwa na apt, na badala yake huchapisha hatua hizi. Pakua kifurushi kipya `ccam_<version>_linux_amd64.deb` kutoka [ukurasa wa matoleo](https://github.com/jhyunwoo/ccam/releases) na ukisakinishe juu ya kile cha zamani:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Ondoa ccam

Kwanza ondoa njia fupi za akaunti, wakati ccam bado imesakinishwa. `ccam list` huzionyesha katika safu ya Shortcut; ondoa kila moja kwa `ccam shortcut <account> --remove`, au futa viungo hivyo kutoka `~/.local/bin` mwenyewe. Kisha ondoa amri, na folda za data ikiwa huhitaji tena akaunti hizo:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Kufuta `~/.local/share/ccam` pia hufuta taarifa za kuingia na historia ya mazungumzo zilizohifadhiwa humo. Folda ya `~/.claude` ambayo [uliiingiza](/docs/commands) hubaki pale ilipo.
