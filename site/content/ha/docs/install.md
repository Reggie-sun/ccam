---
title: Girka ccam
metaTitle: Girka ccam a Linux, macOS da Windows
description: Girka ccam a Linux, macOS ko Windows da umarni ɗaya, kunshin .deb ko lambar tushe, ka kuma sabunta shi da ccam update. Fara girka Claude Code tukuna.
updated: 2026-10-21
order: 1
translation: machine
sourceHash: 1e3a53623de0
---

Gudanar da wannan a terminal a Linux ko macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Rubutun yana sauke sakin da ya dace da tsarinka, yana duba checksum na SHA-256 nasa, kuma yana sa `ccam` a `~/.local/bin`, inda ccam kuma yake ƙirƙirar gajerun umarnin da ka zaɓa daga baya, kamar `ccwork`. Idan wannan foldar ba ta cikin `PATH` naka, rubutun yana buga layin da za ka ƙara a `~/.bashrc` ko `~/.zshrc`.

A Windows, gudanar da wannan a PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Yana sauke `ccam_windows_<arch>.zip`, yana duba checksum na SHA-256 nasa, yana sa `ccam.exe` a `%USERPROFILE%\.local\bin`, foldar da mai girka Claude Code yake sa `claude.exe`, kuma yana ƙara wannan foldar a `PATH` na mai amfani idan ba ta can tukuna. Buɗe sabon terminal bayan haka.

ccam ne ke fara maka Claude Code, don haka ka fara girka [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) kuma ka tabbata `claude --version` yana aiki.

## Tsarukan da ake tallafawa

- Linux a x86-64 da ARM64: Ubuntu, Debian, Fedora, Arch, Alpine da sauran distributions. Binary ɗin static ne kuma ba shi da dependencies.
- macOS a Apple silicon da Intel.
- Windows 10 da 11 a x86-64 da ARM64, a PowerShell, Windows Terminal ko Git Bash. Foldojin da aka raba suna aiki ga kowane mai amfani. Fayilolin da aka raba (`settings.json`, `CLAUDE.md`, `history.jsonl`) suna buƙatar Developer Mode (Settings > System > For developers), saboda Windows yana bari wannan yanayin kaɗai, ko administrator, ya ƙirƙiri mahaɗai zuwa fayiloli; duba [rabawa a Windows](/features/share#on-windows).
- A cikin WSL, ccam yana aiki kamar a kowane tsarin Linux.

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

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Cike umarni a shell (shell completion)

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Completion yana aiki ga dukkan umarni da zaɓuɓɓukan `ccam`, har da sunayen asusai.

## Plugin na Claude Code

Plugin na ccam yana koya wa Claude Code kansa yin amfani da ccam: don duba wane asusu ke da sarari, don fara workers na `claude -p` a lokaci guda a asusai daban-daban da `ccam pick`, da kuma bayyana dalilin da wani skill ko saban MCP ya ɓace da `ccam ext`. Shigar da shi a cikin Claude Code:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

Yana ƙara skill, umarnan `/ccam:usage`, `/ccam:ext` da `/ccam:switch`, da hooks biyu: idan zama ya tsaya a iyakar amfani kana samun sanarwa a terminal, kuma `/ccam:switch [account]` yana miƙa tattaunawar ga wani asusu ([cikakken bayani](/features/next#inside-claude-code)). Hooks ɗin suna gudanar da `ccam`, don haka dole ne ya kasance a `PATH` da Claude Code ke gani. Skill ɗin ba ya taɓa gudanar da umarnan da ke shiga ko motsa fayiloli; yana ba da shawarar su maimakon haka.

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

Dashboard (`ccam` a terminal) ma yana neman sabon saki duk lokacin da ya buɗe kuma, idan akwai, yana ba da shi: danna `U` ka tabbatar. Babu wani umarni dabam da ke nema, kuma ccam ba ya taɓa sabunta kansa ba tare da tambaya ba. Don girka takamaiman saki, har da tsoho, ka ba da `--version`. `ccam --dry-run update` yana nuna abin da za a sauke da abin da za a maye gurbinsa, ba tare da canza komai ba.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Sake gudanar da rubutun girkawa shi ma yana sabunta kwafin da ke `~/.local/bin`. Idan ka gina ccam da `go install`, sake gudanar da wannan umarnin.

**Windows.** `ccam update` yana sauke `.zip` na Windows. Ba za a iya maye gurbin shirin da ke gudana a Windows ba, don haka ana sauya sunan tsohon binary zuwa `ccam.exe.old`, wanda sabuntawa ta gaba ke gogewa.

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

A Windows, a PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Share `~/.local/share/ccam` yana share bayanan shiga da tarihin tattaunawa da aka ajiye a can ma. Foldar `~/.claude` da ka [shigo da ita](/docs/commands) tana nan a inda take.
