---
title: ccam ను ఇన్‌స్టాల్ చేయండి
description: Linux లేదా macOS లో ccam ను ఒకే కమాండ్‌తో, .deb ప్యాకేజీ నుంచి లేదా సోర్స్ నుంచి ఇన్‌స్టాల్ చేయండి. ముందుగా Claude Code ఇన్‌స్టాల్ అయి ఉండాలి.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 79dc1b55e774
---

Linux లేదా macOS లో టెర్మినల్‌లో ఇది రన్ చేయండి:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

ఈ స్క్రిప్ట్ మీ సిస్టమ్‌కు సరిపడే రిలీజ్‌ను డౌన్‌లోడ్ చేసి, దాని SHA-256 చెక్‌సమ్‌ను తనిఖీ చేసి, `ccam` ను `~/.local/bin` లో ఉంచుతుంది; మీరు తర్వాత ఎంచుకునే `ccwork` వంటి షార్ట్‌కట్ కమాండ్‌లను కూడా ccam అదే ఫోల్డర్‌లో సృష్టిస్తుంది. ఆ ఫోల్డర్ మీ `PATH` లో లేకపోతే, `~/.bashrc` లేదా `~/.zshrc` కు జోడించాల్సిన లైన్‌ను స్క్రిప్ట్ ప్రింట్ చేస్తుంది.

ccam మీ కోసం Claude Code ను ప్రారంభిస్తుంది, కాబట్టి ముందుగా [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) ను ఇన్‌స్టాల్ చేసి, `claude --version` పనిచేస్తోందని నిర్ధారించుకోండి.

## మద్దతు ఉన్న సిస్టమ్‌లు

- x86-64, ARM64 పై Linux: Ubuntu, Debian, Fedora, Arch, Alpine, ఇతర డిస్ట్రిబ్యూషన్లు. బైనరీ స్టాటిక్, దానికి ఎలాంటి డిపెండెన్సీలు లేవు.
- Apple silicon, Intel పై macOS.
- Windows కు నేరుగా మద్దతు లేదు. WSL లోపల, ccam ఏ Linux సిస్టమ్‌లోనైనా పనిచేసినట్లే పనిచేస్తుంది.

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

## షెల్ కంప్లీషన్

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

కంప్లీషన్ `ccam` యొక్క అన్ని కమాండ్‌లు, ఆప్షన్‌లకు పనిచేస్తుంది, ఖాతా పేర్లతో సహా.

## ఇన్‌స్టాలేషన్‌ను తనిఖీ చేయండి

```sh
ccam version
ccam doctor
```

`ccam doctor` Claude Code కనబడుతోందా, ప్రతి ఖాతా ఫోల్డర్, లింక్, లాగిన్ సరిగ్గా ఉన్నాయా అని తనిఖీ చేస్తుంది. తర్వాత [క్విక్ స్టార్ట్](/docs/quick-start) తో కొనసాగండి.

## అన్‌ఇన్‌స్టాల్ చేయడం

ccam ఇంకా ఇన్‌స్టాల్ అయి ఉండగానే, ముందుగా ఖాతా షార్ట్‌కట్‌లను తీసివేయండి. `ccam list` వాటిని Shortcut కాలమ్‌లో చూపిస్తుంది; ప్రతిదాన్ని `ccam shortcut <account> --remove` తో తీసివేయండి, లేదా ఆ లింక్‌లను `~/.local/bin` నుంచి మీరే తొలగించండి. తర్వాత కమాండ్‌ను, ఖాతాలు ఇక అవసరం లేకపోతే డేటా ఫోల్డర్‌లను తొలగించండి:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam` ను తొలగిస్తే, అక్కడ నిల్వ ఉన్న లాగిన్‌లు, సంభాషణ చరిత్ర కూడా తొలగిపోతాయి. మీరు [ఇంపోర్ట్ చేసిన](/docs/commands) `~/.claude` ఫోల్డర్ ఉన్న చోటే ఉంటుంది.
