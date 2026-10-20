---
title: Girka ccam
metaTitle: Girka ccam a Linux da macOS
description: Girka ccam a Linux ko macOS da umarni ɗaya, kunshin .deb ko lambar tushe, ka kuma sabunta shi da ccam update. Fara girka Claude Code tukuna.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 7ecd1ed00e09
---

Gudanar da wannan a terminal a Linux ko macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Rubutun yana sauke sakin da ya dace da tsarinka, yana duba checksum na SHA-256 nasa, kuma yana sa `ccam` a `~/.local/bin`, inda ccam kuma yake ƙirƙirar gajerun umarnin da ka zaɓa daga baya, kamar `ccwork`. Idan wannan foldar ba ta cikin `PATH` naka, rubutun yana buga layin da za ka ƙara a `~/.bashrc` ko `~/.zshrc`.

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

## Sabuntawa

```sh
ccam update --check
ccam update
```

`ccam update --check` yana gaya maka kawai ko akwai sabon saki. `ccam update` yana sauke sakin da ya dace da tsarinka, yana duba checksum na SHA-256 nasa da `checksums.txt` na sakin, sannan yana maye gurbin fayil ɗin `ccam` a mataki ɗaya. Idan checksum bai yi daidai ba, ba a canza komai ba. Gajerun umarni kamar `ccwork` mahaɗai ne zuwa ccam, don haka nan take suke gudanar da sabuwar siga.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam ba ya taɓa duba sabuntawa da kansa; yana neman sabon saki ne kawai idan ka gudanar da wannan umarnin. Don girka takamaiman saki, har da tsoho, ka ba da `--version`. `ccam --dry-run update` yana nuna abin da za a sauke da abin da za a maye gurbinsa, ba tare da canza komai ba.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Sake gudanar da rubutun girkawa shi ma yana sabunta kwafin da ke `~/.local/bin`. Idan ka gina ccam da `go install`, sake gudanar da wannan umarnin.

**Kunshin Debian da Ubuntu.** `ccam update` ba ya maye gurbin `/usr/bin/ccam` da apt ya girka, a maimakon haka yana buga waɗannan matakan. Sauke sabon `ccam_<version>_linux_amd64.deb` daga [shafin saki](https://github.com/jhyunwoo/ccam/releases) sannan ka girka shi a kan tsohon:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Cire ccam

Da farko ka cire gajerun umarnin asusai, tun ccam yana nan a kwamfutarka. `ccam list` yana nuna su a ginshiƙin Shortcut; cire kowanne da `ccam shortcut <account> --remove`, ko ka goge waɗannan mahaɗai daga `~/.local/bin` da kanka. Sannan ka cire umarnin, da foldojin bayanai idan ba ka buƙatar asusan kuma:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Share `~/.local/share/ccam` yana share bayanan shiga da tarihin tattaunawa da aka ajiye a can ma. Foldar `~/.claude` da ka [shigo da ita](/docs/commands) tana nan a inda take.
