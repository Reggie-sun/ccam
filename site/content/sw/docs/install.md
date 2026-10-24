---
title: Sakinisha ccam
metaTitle: Sakinisha ccam kwenye Linux, macOS na Windows
description: Sakinisha ccam kwenye Linux, macOS au Windows kwa amri moja, kifurushi cha .deb au msimbo chanzo, na uisasishe kwa ccam update. Sakinisha Claude Code kwanza.
updated: 2026-10-23
order: 1
translation: machine
sourceHash: "af21dcc24169"
---

Endesha hii kwenye terminal ya Linux au macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Skripti hupakua toleo la mfumo wako, hukagua checksum yake ya SHA-256 na huweka `ccam` katika `~/.local/bin`, ambapo ccam pia huunda amri za njia fupi utakazochagua baadaye, kama `ccwork`. Ikiwa folda hiyo haiko kwenye `PATH` yako, skripti huchapisha mstari wa kuongeza kwenye `~/.bashrc` au `~/.zshrc`.

Kwenye Windows, endesha hii katika PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Skripti hupakua `ccam_windows_<arch>.zip`, hukagua checksum yake ya SHA-256, huweka `ccam.exe` katika `%USERPROFILE%\.local\bin`, folda ambayo kisakinishi cha Claude Code huweka `claude.exe`, na huongeza folda hiyo kwenye `PATH` ya mtumiaji wako ikiwa haipo tayari, huku ikiacha vipengele kama `%USERPROFILE%\bin` jinsi vilivyo. Fungua terminal mpya baadaye. Kuiendesha tena hubadilisha ccam, hata wakati kipindi cha Claude Code bado kinaendeshwa kutoka binary ya zamani (hiyo hupewa jina jipya `ccam.exe.old`, au `ccam.exe.old-<number>` wakati ya zamani zaidi bado inatumika), na husasisha njia fupi ambazo ni hard links: `OK 2 shortcuts now run the new ccam`.

ccam hukuanzishia Claude Code, kwa hiyo sakinisha [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) kwanza na uhakikishe kuwa `claude --version` inafanya kazi.

## Mifumo inayotumika

- Linux kwenye x86-64 na ARM64: Ubuntu, Debian, Fedora, Arch, Alpine na distribution nyingine. Faili ya binary ni static na haina utegemezi wowote.
- macOS kwenye Apple silicon na Intel.
- Windows 10 na 11 kwenye x86-64 na ARM64, katika PowerShell, Windows Terminal au Git Bash. Folda zinazoshirikiwa hufanya kazi kwa kila mtumiaji. Faili zinazoshirikiwa (`settings.json`, `CLAUDE.md`, `history.jsonl`) zinahitaji Developer Mode (Settings > System > For developers), kwa sababu Windows huruhusu hali hiyo pekee, au msimamizi, kuunda viungo vya faili; tazama [kushiriki kwenye Windows](/features/share#on-windows). Bila hiyo, `ccam add` ya kwanza huunda kikundi `default` chenye seva za MCP na viendelezi pekee, na husema jinsi ya kuongeza vilivyobaki Developer Mode ikishawashwa.
- Ndani ya WSL, ccam hufanya kazi kama kwenye mfumo wowote wa Linux.

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

Ujenzi wa tagi ya toleo (`@latest` au `@vX.Y.Z`) huripoti toleo hilo, kwa hiyo `ccam update` hufanya kazi kwake; ujenzi wa `@main` au wa nakala ya ndani (checkout) iliyo mbele ya tagi huripoti `dev`.

**Toleo mahususi.** Weka `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Ukamilishaji wa shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Ukamilishaji unashughulikia amri na chaguo zote za `ccam`, pamoja na majina ya akaunti.

## Plugin ya Claude Code

Plugin ya ccam humfundisha Claude Code yenyewe kutumia ccam: kuangalia akaunti ipi ina nafasi, kuanzisha workers wa `claude -p` sambamba kwenye akaunti tofauti kwa `ccam pick`, na kueleza kwa nini skill au seva ya MCP haipo kwa `ccam ext`. Isakinishe ndani ya Claude Code:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

Huongeza skill moja, amri `/ccam:usage`, `/ccam:ext` na `/ccam:switch`, na hooks mbili: kipindi kinaposimama kwenye kikomo cha matumizi unapata arifa kwenye terminal, na `/ccam:switch [account]` hukabidhi mazungumzo kwa akaunti nyingine ([maelezo](/features/next#inside-claude-code)). Hooks huendesha `ccam`, kwa hiyo lazima iwe kwenye `PATH` ambayo Claude Code huiona. Skill hiyo kamwe haiendeshi amri zinazoingia kwenye akaunti au kuhamisha faili; badala yake huzipendekeza.

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

Dashibodi (`ccam` kwenye terminal) pia hutafuta toleo jipya kila inapofunguliwa na, likiwepo, hulitoa: bonyeza `U` kisha uthibitishe. Hakuna amri nyingine inayotafuta, na ccam kamwe haijisasishi bila kuuliza. Ili kusakinisha toleo mahususi, hata la zamani, tumia `--version`. `ccam --dry-run update` huonyesha kile ambacho kingepakuliwa na kubadilishwa, bila kubadilisha chochote.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Kuendesha skripti ya usakinishaji tena pia husasisha nakala iliyo katika `~/.local/bin`. Ikiwa ulijenga ccam kwa `go install` kutoka tagi ya toleo, `ccam update` hufanya kazi pia, au endesha amri ile ile tena.

**Windows.** `ccam update` hupakua `.zip` ya Windows. Programu inayoendeshwa haiwezi kubadilishwa kwenye Windows, kwa hiyo binary ya zamani hupewa jina jipya `ccam.exe.old`, ambalo sasisho linalofuata hulifuta; wakati kipindi bado kinaendeshwa kutoka `ccam.exe.old` ya awali, jina huwa `ccam.exe.old-<time>`. Njia fupi ambazo ni hard links husasishwa mara baada ya hapo.

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

Kwenye Windows, katika PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Kufuta `~/.local/share/ccam` pia hufuta taarifa za kuingia na historia ya mazungumzo zilizohifadhiwa humo. Folda ya `~/.claude` ambayo [uliiingiza](/docs/commands) hubaki pale ilipo, lakini ikiwa katika kikundi, faili inazoshiriki huishi katika `~/.local/share/ccam/groups/`. Kwanza endesha `ccam remove <name>` kwa ajili yake (au `ccam group leave`): huacha vikundi vyake ikiwa na nakala ya kila kitu ilichoshiriki.
