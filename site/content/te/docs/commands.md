---
title: కమాండ్ రిఫరెన్స్
description: ccam add, ccam <name> నుంచి ccam usage, ccam next, ccam doctor వరకు ccam యొక్క ప్రతి కమాండ్, ఆప్షన్.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

ఏ కమాండ్ అయినా అసలు చేయకుండా అది ఏమి చేస్తుందో చూడటానికి, దాని ముందు `--dry-run` పెట్టండి, ఉదాహరణకు `ccam --dry-run work`.

<!-- commands -->

## Claude Code ను ప్రారంభించడం

`ccam <name|number> [arguments]` ఆ ఖాతాగా Claude Code ను ప్రారంభిస్తుంది. పేరు తర్వాత ఉన్న ఆర్గ్యుమెంట్‌లను ccam చదవదు; వాటిని అదే క్రమంలో `claude` కు పంపుతుంది. స్టాండర్డ్ ఇన్‌పుట్, ఎగ్జిట్ కోడ్, సిగ్నల్స్ అన్నీ `claude` తో నేరుగా ఉన్నట్లే పనిచేస్తాయి. పేరు తర్వాత `--` ఐచ్ఛికం.

ఖాతా పేరు లేకుండా ఆర్గ్యుమెంట్‌లతో `ccam`, ఉదాహరణకు `ccam -p "summarise"`, డిఫాల్ట్ ఖాతాను ఉపయోగిస్తుంది. దాన్ని `ccam use <name>` తో సెట్ చేయండి.

`claude1`, `claude2` మొదలైన షార్ట్‌కట్‌లు `~/.local/bin` లో ccam కు ఉన్న లింక్‌లు. `ccam shims install --names` ప్రతి ఖాతాకు `claude-<name>` ను కూడా సృష్టిస్తుంది.

## ఖాతాలను జోడించడం, తీసివేయడం

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

ఫోల్డర్ ఇవ్వకుండా `ccam import` రన్ చేస్తే, అది `~/.claude`, `~/.claude-*` కోసం వెతికి ప్రతిదాని గురించి అడుగుతుంది. ఇంపోర్ట్ చేసిన ఫోల్డర్‌లు ఉన్న చోటే నమోదు అవుతాయి; ccam ఎప్పుడూ లాగిన్‌ను తరలించదు, కాపీ చేయదు.

## షేరింగ్

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

ప్రతి అంశం `on` (`default` గ్రూప్), `off`, లేదా ఒక గ్రూప్ పేరును తీసుకుంటుంది. [షేరింగ్](/docs/sharing) చూడండి.

## MCP సర్వర్లు

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## వినియోగం, ఖాతా మార్చడం

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## నిర్వహణ

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## ఫైళ్లు

| పాత్ | విషయాలు |
| --- | --- |
| `~/.config/ccam/config.toml` | ఖాతాలు, అవి ఏమి షేర్ చేస్తాయి |
| `~/.local/share/ccam/profiles/<name>/` | ప్రతి ఖాతా యొక్క `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | షేర్ చేసిన సెట్టింగ్‌లు, MCP సర్వర్లు, ఎక్స్‌టెన్షన్లు, చరిత్ర |
| `~/.local/state/ccam/` | వినియోగ కాష్, చివరిగా ఉపయోగించిన సమయాలు |
| `~/.local/bin/claude1`, `claude2`… | ccam కు లింక్‌లు |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_STATE_HOME` వేరియబుల్స్ ఈ ఫోల్డర్‌ల స్థానాన్ని మారుస్తాయి.
