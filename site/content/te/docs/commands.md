---
title: కమాండ్ రిఫరెన్స్
description: ccam add, ccam <name> నుంచి ccam usage, ccam next, ccam doctor వరకు ccam యొక్క ప్రతి కమాండ్, ఆప్షన్.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 0c869c92fa63
---

ఏ కమాండ్ అయినా అసలు చేయకుండా అది ఏమి చేస్తుందో చూడటానికి, దాని ముందు `--dry-run` పెట్టండి, ఉదాహరణకు `ccam --dry-run work`.

<!-- commands -->

## Claude Code ను ప్రారంభించడం

`ccam <name|number> [arguments]` ఆ ఖాతాగా Claude Code ను ప్రారంభిస్తుంది. పేరు తర్వాత ఉన్న ఆర్గ్యుమెంట్‌లను ccam చదవదు; వాటిని అదే క్రమంలో `claude` కు పంపుతుంది. స్టాండర్డ్ ఇన్‌పుట్, ఎగ్జిట్ కోడ్, సిగ్నల్స్ అన్నీ `claude` తో నేరుగా ఉన్నట్లే పనిచేస్తాయి. పేరు తర్వాత `--` ఐచ్ఛికం.

ఖాతా పేరు లేకుండా ఆర్గ్యుమెంట్‌లతో `ccam`, ఉదాహరణకు `ccam -p "summarise"`, డిఫాల్ట్ ఖాతాను ఉపయోగిస్తుంది. దాన్ని `ccam use <name>` తో సెట్ చేయండి.

`ccwork` వంటి షార్ట్‌కట్ కూడా తన ఖాతాను ఇదే విధంగా ప్రారంభిస్తుంది: `ccam work --resume` చేసేదే `ccwork --resume` చేస్తుంది. [షార్ట్‌కట్‌లు](/docs/commands#shortcuts) చూడండి.

## ఖాతాలను జోడించడం, తీసివేయడం

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

ఫోల్డర్ ఇవ్వకుండా `ccam import` రన్ చేస్తే, అది `~/.claude`, `~/.claude-*` కోసం వెతికి ప్రతిదాని గురించి అడుగుతుంది. ఇంపోర్ట్ చేసిన ఫోల్డర్‌లు ఉన్న చోటే నమోదు అవుతాయి; ccam ఎప్పుడూ లాగిన్‌ను తరలించదు, కాపీ చేయదు.

## షార్ట్‌కట్‌లు

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

షార్ట్‌కట్ అంటే ఒక ఖాతాగా Claude Code ను ప్రారంభించే, మీరే పేరు పెట్టిన కమాండ్. `ccam work` లాగే, ఇది ప్రతి ఆర్గ్యుమెంట్‌ను మార్పు లేకుండా Claude Code కు పంపుతుంది, కాబట్టి `ccwork --resume`, `ccwork -p "hi"` మీరు ఆశించినట్లే పనిచేస్తాయి. ఇది `~/.local/bin` లో ccam కు ఉన్న లింక్.

షార్ట్‌కట్‌లు ఐచ్ఛికం. `ccam add` విజార్డ్ ఒకటి సృష్టించాలా అని అడుగుతుంది, `ccam add <name> --shortcut <command>` అడగకుండానే ఒకటి సృష్టిస్తుంది, `ccam list` ప్రతి ఖాతా షార్ట్‌కట్‌ను చూపిస్తుంది, లేకపోతే `-` చూపిస్తుంది.

పేర్లలో a-z, 0-9, `.`, `-`, `_` అక్షరాలు ఉపయోగించవచ్చు. పేరు ccam కమాండ్ కాకూడదు, మరో ఖాతా షార్ట్‌కట్ కాకూడదు, మీ `PATH` లోని మరో ప్రోగ్రామ్ పేరు కాకూడదు.

`ccam rename` షార్ట్‌కట్ పేరును అలాగే ఉంచుతుంది, `ccam remove` ఖాతాతో పాటు షార్ట్‌కట్‌ను తొలగిస్తుంది. `ccam shims install` అన్ని ఖాతాల లింక్‌లను మళ్లీ సృష్టిస్తుంది, ఉదాహరణకు మీరు ccam ను వేరే చోటికి తరలించిన తర్వాత.

## షేరింగ్

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

ప్రతి అంశం `on` (`default` గ్రూప్), `off`, లేదా ఒక గ్రూప్ పేరును తీసుకుంటుంది. [షేరింగ్](/docs/sharing) చూడండి.

## ఇప్పటికే ఉన్న సెటప్‌ను కాపీ చేయడం

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` ఒక Claude Code ఫోల్డర్ నుంచి `settings.json`, `CLAUDE.md`, ఎక్స్‌టెన్షన్లు, యూజర్ MCP సర్వర్లను ఒక గ్రూప్‌లోకి కాపీ చేస్తుంది, కాబట్టి ఆ అంశాలను షేర్ చేసే ప్రతి ఖాతాకు అవి అందుతాయి. ఫోల్డర్ చదవబడుతుంది మాత్రమే. `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` ఆప్షన్లు మీరు పేర్కొన్న అంశాలను మాత్రమే కాపీ చేస్తాయి; అవి లేకపోతే, సంభాషణ చరిత్ర తప్ప అన్నీ కాపీ అవుతాయి. `-g` లేదా `--group` గ్రూప్‌ను ఎంచుకుంటుంది, `-y` లేదా `--yes` అడగకుండా కాపీ చేసి, ఫైల్ వేరుగా ఉన్నప్పుడు గ్రూప్ కాపీని ఉంచుతుంది. [ఇప్పటికే ఉన్న సెటప్‌ను కాపీ చేయడం](/docs/sharing#copy-an-existing-setup) చూడండి.

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
ccam shims install   # recreate the shortcut links
```

## ఫైళ్లు

| పాత్ | విషయాలు |
| --- | --- |
| `~/.config/ccam/config.toml` | ఖాతాలు, అవి ఏమి షేర్ చేస్తాయి |
| `~/.local/share/ccam/profiles/<name>/` | ప్రతి ఖాతా యొక్క `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | షేర్ చేసిన సెట్టింగ్‌లు, MCP సర్వర్లు, ఎక్స్‌టెన్షన్లు, చరిత్ర |
| `~/.local/state/ccam/` | వినియోగ కాష్, చివరిగా ఉపయోగించిన సమయాలు |
| `~/.local/bin/<shortcut>` | `ccwork` వంటి షార్ట్‌కట్ కమాండ్‌లు, ccam కు లింక్‌లు |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_STATE_HOME` వేరియబుల్స్ ఈ ఫోల్డర్‌ల స్థానాన్ని మారుస్తాయి.
