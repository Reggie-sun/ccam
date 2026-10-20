---
title: extensions ను పరిశీలించడం
description: ccam ext ప్రతి ఖాతా లోడ్ చేసే skills, commands, plugins, MCP సర్వర్లు, hooks ను, ప్రతిదీ ఎక్కడి నుంచి వస్తుందో, ఏది విరిగిందో చూపిస్తుంది.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

"ఈ skill (లేదా MCP సర్వర్) ఒక ఖాతాలో ఎందుకు పనిచేస్తుంది, ఇంకో దానిలో ఎందుకు పనిచేయదు?" అనే ప్రశ్నకు `ccam ext` సమాధానం ఇస్తుంది. ప్రతి ఖాతాకు, Claude Code లోడ్ చేసే extensions ను, ప్రతిదీ ఎక్కడి నుంచి వస్తుందో అది జాబితా చేస్తుంది: ఒక షేర్ గ్రూప్, ఖాతా సొంత ఫోల్డర్, లేదా మీరు ఉన్న ప్రాజెక్ట్. టేబుల్ కింద, ఒక extension మీరు ఆశించినట్లు లోడ్ కాకుండా ఆపే వాటిని అది జాబితా చేస్తుంది. అది చదువుతుంది మాత్రమే.

## వినియోగం

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

ప్రతి గడి ఆ ఖాతాలో ఆ అంశం సోర్స్‌ను సూచిస్తుంది: షేర్ గ్రూప్ (`default`), అది ఖాతా సొంత ఫోల్డర్‌లో ఉంటే `account`, లేదా ఖాతాకు అది లేకపోతే `-`. `!` సమస్య ఉన్న అంశాన్ని సూచిస్తుంది, అది టేబుల్ కింద వివరించబడుతుంది.

## ఇది ఏమి చదువుతుంది

| రకం | ఎక్కడ | చూపించే సోర్స్ |
| --- | --- | --- |
| Skills | `skills/` లో `SKILL.md` ఉన్న ప్రతి ఫోల్డర్ | `skills/` లింక్ అయిన గ్రూప్, లేదా `account` |
| Commands, agents | `commands/`, `agents/` కింద ఉన్న `.md` ఫైల్‌లు, ఉప ఫోల్డర్‌లు `folder/name` గా | గ్రూప్, లేదా `account` |
| Plugins | `plugins/installed_plugins.json` | గ్రూప్, లేదా `account` |
| MCP సర్వర్లు | గ్రూప్ యొక్క `mcp.json` (MCP షేర్ అయినప్పుడు), ఖాతా యొక్క `.claude.json` | గ్రూప్, లేదా `account` |
| Hooks | `settings.json` లోని `hooks`, ఈవెంట్ వారీగా లెక్కించి | గ్రూప్, లేదా `account` |
| ఈ ప్రాజెక్ట్ | ప్రస్తుత ఫోల్డర్‌లోని `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json`, `.mcp.json` | ఒక్కసారి, అన్ని ఖాతాలకు |

## ఇది నివేదించే సమస్యలు

- **విరిగిన లింక్‌లు**: `skills/`, `commands/` లేదా `agents/` లో లక్ష్యం లేని లింక్, లేదా లేని చోటుకు లింక్ అయిన షేర్ చేసిన ఫోల్డర్ (`skills`, `commands`, `agents`, `plugins`) (`ccam doctor --fix` తో సరిచేయండి).
- **`SKILL.md` లేని skill ఫోల్డర్**: Claude Code దాన్ని వదిలేస్తుంది.
- **రెండు చోట్ల వేర్వేరుగా నిర్వచించిన MCP సర్వర్**: గ్రూప్ యొక్క `mcp.json`, ఖాతా యొక్క `.claude.json`, లేదా వాటిలో ఒకటి, ప్రాజెక్ట్ యొక్క `.mcp.json`. రెండు చోట్లా అవే సెట్టింగ్‌లు ఉంటే (`mcp_strategy = "sync"` రాసేది అదే) అది సమస్య కాదు. ఖాతా సొంత సర్వర్‌ను `ccam mcp import <account> <server>` తో గ్రూప్‌లోకి తరలించండి.
- **చెల్లుబాటు అయ్యే JSON కాని `installed_plugins.json`.**

## భద్రత, ప్రత్యేక సందర్భాలు

- `ccam ext` ఏమీ మార్చదు, నెట్‌వర్క్ అభ్యర్థనలు చేయదు. అది సర్వర్ పేర్లను, hook సంఖ్యలను చదువుతుంది, MCP నిర్వచనాల్లోని `env` లేదా headers విలువలను ఎప్పుడూ చదవదు.
- అది డిస్క్‌పై ఉన్నదాన్ని చూపిస్తుంది. ఇన్‌స్టాల్ అయి కానీ సెట్టింగ్‌లలో ఆఫ్ చేసిన plugin, లేదా managed settings అడ్డుకున్న సర్వర్ అయినా కనిపిస్తుంది.
- ఆ ప్రాజెక్ట్ యొక్క `.claude` ఫోల్డర్‌ను, `.mcp.json` ను చేర్చడానికి దాన్ని ప్రాజెక్ట్ ఫోల్డర్‌లో రన్ చేయండి.

## సంబంధిత పేజీలు

- [ఖాతాల మధ్య షేరింగ్](/features/share)
- [షేర్ చేసిన MCP సర్వర్లు](/features/mcp)
- [తనిఖీ చేయడం, సరిచేయడం](/features/doctor)
