---
title: ccam ను ఇన్‌స్టాల్ చేయండి
metaTitle: Linux, macOS, Windows లో ccam ఇన్‌స్టాల్ చేయండి
description: Linux, macOS లేదా Windows లో ccam ను ఒకే కమాండ్, .deb ప్యాకేజీ లేదా సోర్స్ నుంచి ఇన్‌స్టాల్ చేసి, ccam update తో అప్‌డేట్ చేయండి. ముందుగా Claude Code ఉండాలి.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: aea4984c1576
---

Linux లేదా macOS లో టెర్మినల్‌లో ఇది రన్ చేయండి:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

ఈ స్క్రిప్ట్ మీ సిస్టమ్‌కు సరిపడే రిలీజ్‌ను డౌన్‌లోడ్ చేసి, దాని SHA-256 చెక్‌సమ్‌ను తనిఖీ చేసి, `ccam` ను `~/.local/bin` లో ఉంచుతుంది; మీరు తర్వాత ఎంచుకునే `ccwork` వంటి షార్ట్‌కట్ కమాండ్‌లను కూడా ccam అదే ఫోల్డర్‌లో సృష్టిస్తుంది. ఆ ఫోల్డర్ మీ `PATH` లో లేకపోతే, `~/.bashrc` లేదా `~/.zshrc` కు జోడించాల్సిన లైన్‌ను స్క్రిప్ట్ ప్రింట్ చేస్తుంది.

Windows లో, PowerShell లో ఇది రన్ చేయండి:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

ఇది `ccam_windows_<arch>.zip` ను డౌన్‌లోడ్ చేసి, దాని SHA-256 చెక్‌సమ్‌ను తనిఖీ చేసి, Claude Code ఇన్‌స్టాలర్ `claude.exe` ను ఉంచే అదే ఫోల్డర్ అయిన `%USERPROFILE%\.local\bin` లో `ccam.exe` ను ఉంచుతుంది, ఆ ఫోల్డర్ మీ యూజర్ `PATH` లో ఇంకా లేకపోతే దాన్ని జోడిస్తుంది. తర్వాత కొత్త టెర్మినల్ తెరవండి.

ccam మీ కోసం Claude Code ను ప్రారంభిస్తుంది, కాబట్టి ముందుగా [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) ను ఇన్‌స్టాల్ చేసి, `claude --version` పనిచేస్తోందని నిర్ధారించుకోండి.

## మద్దతు ఉన్న సిస్టమ్‌లు

- x86-64, ARM64 పై Linux: Ubuntu, Debian, Fedora, Arch, Alpine, ఇతర డిస్ట్రిబ్యూషన్లు. బైనరీ స్టాటిక్, దానికి ఎలాంటి డిపెండెన్సీలు లేవు.
- Apple silicon, Intel పై macOS.
- x86-64, ARM64 పై Windows 10, 11, PowerShell, Windows Terminal లేదా Git Bash లో. షేర్డ్ ఫోల్డర్‌లు ప్రతి యూజర్‌కూ పనిచేస్తాయి. షేర్డ్ ఫైల్‌లకు (`settings.json`, `CLAUDE.md`, `history.jsonl`) Developer Mode (Settings > System > For developers) అవసరం, ఎందుకంటే Windows ఆ మోడ్‌లో లేదా అడ్మినిస్ట్రేటర్‌కు మాత్రమే ఫైల్‌లకు లింక్‌లు సృష్టించనిస్తుంది; చూడండి [Windows లో షేరింగ్](/features/share#on-windows).
- WSL లోపల, ccam ఏ Linux సిస్టమ్‌లోనైనా పనిచేసినట్లే పనిచేస్తుంది.

## ఇన్‌స్టాల్ చేయడానికి ఇతర మార్గాలు

**Debian, Ubuntu ప్యాకేజీ.** [రిలీజ్‌ల పేజీ](https://github.com/jhyunwoo/ccam/releases) నుంచి `ccam_<version>_linux_amd64.deb` ను డౌన్‌లోడ్ చేసి ఇన్‌స్టాల్ చేయండి:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

ఈ ప్యాకేజీ `/usr/bin/ccam` ను ఇన్‌స్టాల్ చేస్తుంది.

**సోర్స్ నుంచి.** Go 1.26 లేదా కొత్త వెర్షన్‌తో:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**ఒక నిర్దిష్ట వెర్షన్.** `CCAM_VERSION` ను సెట్ చేయండి:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## షెల్ కంప్లీషన్

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

కంప్లీషన్ `ccam` యొక్క అన్ని కమాండ్‌లు, ఆప్షన్‌లకు పనిచేస్తుంది, ఖాతా పేర్లతో సహా.

## Claude Code plugin

ccam plugin Claude Code కే ccam ను ఉపయోగించడం నేర్పుతుంది: ఏ ఖాతాలో అవకాశం ఉందో తనిఖీ చేయడం, `ccam pick` తో వేర్వేరు ఖాతాలపై సమాంతర `claude -p` workers ను ప్రారంభించడం, కనిపించని skill లేదా MCP సర్వర్‌ను `ccam ext` తో వివరించడం. దాన్ని Claude Code లోపల ఇన్‌స్టాల్ చేయండి:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

ఇది ఒక skill ను, రెండు commands ను జోడిస్తుంది, `/ccam:usage`, `/ccam:ext`. ఆ skill లాగిన్ చేసే లేదా ఫైల్‌లను తరలించే కమాండ్‌లను ఎప్పుడూ రన్ చేయదు; బదులుగా వాటిని సూచిస్తుంది.

## ఇన్‌స్టాలేషన్‌ను తనిఖీ చేయండి

```sh
ccam version
ccam doctor
```

`ccam doctor` Claude Code కనబడుతోందా, ప్రతి ఖాతా ఫోల్డర్, లింక్, లాగిన్ సరిగ్గా ఉన్నాయా అని తనిఖీ చేస్తుంది. తర్వాత [క్విక్ స్టార్ట్](/docs/quick-start) తో కొనసాగండి.

## అప్‌డేట్ చేయడం

```sh
ccam update --check
ccam update
```

`ccam update --check` కొత్త రిలీజ్ ఉందో లేదో మాత్రమే చెబుతుంది. `ccam update` మీ సిస్టమ్‌కు సరిపడే రిలీజ్‌ను డౌన్‌లోడ్ చేసి, దాని SHA-256 చెక్‌సమ్‌ను రిలీజ్‌లోని `checksums.txt` తో సరిపోల్చి, తర్వాత `ccam` ఫైల్‌ను ఒకే దశలో మారుస్తుంది. చెక్‌సమ్ సరిపోలకపోతే, ఏదీ మార్చబడదు. `ccwork` వంటి షార్ట్‌కట్‌లు ccam కు లింక్‌లు కాబట్టి, అవి వెంటనే కొత్త వెర్షన్‌ను రన్ చేస్తాయి.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam తనంతట తానుగా ఎప్పుడూ అప్‌డేట్‌ల కోసం తనిఖీ చేయదు; మీరు ఈ కమాండ్ రన్ చేసినప్పుడు మాత్రమే కొత్త రిలీజ్ కోసం చూస్తుంది. పాత రిలీజ్‌తో సహా, ఒక నిర్దిష్ట రిలీజ్‌ను ఇన్‌స్టాల్ చేయడానికి `--version` ఇవ్వండి. `ccam --dry-run update` ఏదీ మార్చకుండా, ఏమి డౌన్‌లోడ్ అవుతుందో, ఏమి మారుతుందో చూపిస్తుంది.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

ఇన్‌స్టాల్ స్క్రిప్ట్‌ను మళ్లీ రన్ చేసినా `~/.local/bin` లోని కాపీ అప్‌డేట్ అవుతుంది. మీరు ccam ను `go install` తో బిల్డ్ చేసి ఉంటే, అదే కమాండ్‌ను మళ్లీ రన్ చేయండి.

**Windows.** `ccam update` Windows కోసం `.zip` ను డౌన్‌లోడ్ చేస్తుంది. Windows లో రన్ అవుతున్న ప్రోగ్రామ్‌ను మార్చలేము, కాబట్టి పాత బైనరీ పేరు `ccam.exe.old` గా మార్చబడుతుంది, తదుపరి అప్‌డేట్ దాన్ని తొలగిస్తుంది.

**Debian, Ubuntu ప్యాకేజీ.** apt ఇన్‌స్టాల్ చేసిన `/usr/bin/ccam` ను `ccam update` మార్చదు, దానికి బదులుగా ఈ దశలను ప్రింట్ చేస్తుంది. [రిలీజ్‌ల పేజీ](https://github.com/jhyunwoo/ccam/releases) నుంచి కొత్త `ccam_<version>_linux_amd64.deb` ను డౌన్‌లోడ్ చేసి, పాత దాని మీదే ఇన్‌స్టాల్ చేయండి:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## అన్‌ఇన్‌స్టాల్ చేయడం

ccam ఇంకా ఇన్‌స్టాల్ అయి ఉండగానే, ముందుగా ఖాతా షార్ట్‌కట్‌లను తీసివేయండి. `ccam list` వాటిని Shortcut కాలమ్‌లో చూపిస్తుంది; ప్రతిదాన్ని `ccam shortcut <account> --remove` తో తీసివేయండి, లేదా ఆ లింక్‌లను `~/.local/bin` నుంచి మీరే తొలగించండి. తర్వాత కమాండ్‌ను, ఖాతాలు ఇక అవసరం లేకపోతే డేటా ఫోల్డర్‌లను తొలగించండి:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Windows లో, PowerShell లో:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

`~/.local/share/ccam` ను తొలగిస్తే, అక్కడ నిల్వ ఉన్న లాగిన్‌లు, సంభాషణ చరిత్ర కూడా తొలగిపోతాయి. మీరు [ఇంపోర్ట్ చేసిన](/docs/commands) `~/.claude` ఫోల్డర్ ఉన్న చోటే ఉంటుంది.
