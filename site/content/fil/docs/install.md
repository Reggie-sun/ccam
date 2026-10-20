---
title: I-install ang ccam
metaTitle: I-install ang ccam sa Linux, macOS at Windows
description: I-install ang ccam sa Linux, macOS o Windows gamit ang isang command, .deb package o source, at i-update sa ccam update. Kailangan munang may Claude Code.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: 07eadc0769c8
---

Patakbuhin ito sa terminal sa Linux o macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Dina-download ng script ang release para sa iyong system, sinusuri ang SHA-256 checksum nito at inilalagay ang `ccam` sa `~/.local/bin`, kung saan gagawa rin ang ccam ng mga shortcut command na pipiliin mo mamaya, gaya ng `ccwork`. Kung wala sa iyong `PATH` ang folder na iyon, ipi-print ng script ang linyang idadagdag sa `~/.bashrc` o `~/.zshrc`.

Sa Windows, patakbuhin ito sa PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Dina-download nito ang `ccam_windows_<arch>.zip`, sinusuri ang SHA-256 checksum nito, inilalagay ang `ccam.exe` sa `%USERPROFILE%\.local\bin`, ang folder kung saan inilalagay ng installer ng Claude Code ang `claude.exe`, at idinadagdag ang folder na iyon sa user `PATH` mo kung wala pa ito roon. Magbukas ng bagong terminal pagkatapos.

Ang ccam ang nagsisimula ng Claude Code para sa iyo, kaya i-install muna ang [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) at tiyaking gumagana ang `claude --version`.

## Mga suportadong system

- Linux sa x86-64 at ARM64: Ubuntu, Debian, Fedora, Arch, Alpine at iba pang distribution. Static ang binary at wala itong dependency.
- macOS sa Apple silicon at Intel.
- Windows 10 at 11 sa x86-64 at ARM64, sa PowerShell, Windows Terminal o Git Bash. Gumagana ang mga shared na folder para sa bawat user. Kailangan ng Developer Mode (Settings > System > For developers) ang mga shared na file (`settings.json`, `CLAUDE.md`, `history.jsonl`), dahil ang mode na iyon lang, o isang administrator, ang pinapayagan ng Windows na gumawa ng mga link sa mga file; tingnan ang [sharing sa Windows](/features/share#on-windows).
- Sa loob ng WSL, gumagana ang ccam gaya ng sa anumang Linux system.

## Iba pang paraan ng pag-install

**Package para sa Debian at Ubuntu.** I-download ang `ccam_<version>_linux_amd64.deb` mula sa [pahina ng mga release](https://github.com/jhyunwoo/ccam/releases) at i-install ito:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Ini-install ng package ang `/usr/bin/ccam`.

**Mula sa source.** Gamit ang Go 1.26 o mas bago:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**Isang partikular na bersyon.** Itakda ang `CCAM_VERSION`:

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

Saklaw ng completion ang lahat ng command at option ng `ccam`, pati ang mga pangalan ng account.

## Suriin ang installation

```sh
ccam version
ccam doctor
```

Sinusuri ng `ccam doctor` na nahahanap ang Claude Code at na maayos ang bawat account folder, link at login. Magpatuloy sa [quick start](/docs/quick-start).

## Pag-update

```sh
ccam update --check
ccam update
```

Sinasabi lang ng `ccam update --check` kung may mas bagong release. Dina-download ng `ccam update` ang release para sa iyong system, sinusuri ang SHA-256 checksum nito laban sa `checksums.txt` ng release, saka pinapalitan ang `ccam` file sa iisang hakbang. Kung hindi tugma ang checksum, walang binabago. Mga link sa ccam ang mga shortcut gaya ng `ccwork`, kaya agad nilang pinapatakbo ang bagong bersyon.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

Hindi kailanman kusang naghahanap ng update ang ccam; naghahanap lang ito ng bagong release kapag pinatakbo mo ang command na ito. Para mag-install ng partikular na release, kahit mas luma, ipasa ang `--version`. Ipinapakita ng `ccam --dry-run update` kung ano ang ida-download at papalitan, nang walang binabago.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Kapag pinatakbo mong muli ang install script, maa-update din ang kopya sa `~/.local/bin`. Kung binuild mo ang ccam gamit ang `go install`, patakbuhin muli ang parehong command.

**Windows.** Dina-download ng `ccam update` ang `.zip` para sa Windows. Hindi mapapalitan sa Windows ang isang tumatakbong program, kaya pinapalitan ang pangalan ng lumang binary ng `ccam.exe.old`, na binubura ng susunod na update.

**Package para sa Debian at Ubuntu.** Hindi pinapalitan ng `ccam update` ang `/usr/bin/ccam` na in-install ng apt, at ipi-print na lang nito ang mga hakbang na ito. I-download ang bagong `ccam_<version>_linux_amd64.deb` mula sa [pahina ng mga release](https://github.com/jhyunwoo/ccam/releases) at i-install ito sa ibabaw ng luma:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Pag-uninstall

Alisin muna ang mga shortcut ng account habang naka-install pa ang ccam. Ipinapakita ng `ccam list` ang mga ito sa column na Shortcut; alisin ang bawat isa gamit ang `ccam shortcut <account> --remove`, o ikaw mismo ang magbura ng mga link na iyon mula sa `~/.local/bin`. Saka alisin ang command, at ang mga data folder kung hindi mo na kailangan ang mga account:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Sa Windows, sa PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Kapag binura mo ang `~/.local/share/ccam`, mabubura rin ang mga login at kasaysayan ng usapan na naka-store doon. Mananatili kung nasaan ito ang `~/.claude` folder na iyong [in-import](/docs/commands).
