---
title: extensions ను పరిశీలించడం
description: ccam ext ప్రతి ఖాతా లోడ్ చేసే skills, commands, plugins, MCP సర్వర్లు, hooks ను, ప్రతిదీ ఎక్కడి నుంచి వస్తుందో, ఏది విరిగిందో చూపిస్తుంది.
updated: 2026-10-23
order: 14
translation: machine
sourceHash: 45740e4efbef
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
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

ప్రతి గడి ఆ ఖాతాలో ఆ అంశం సోర్స్‌ను సూచిస్తుంది: షేర్ గ్రూప్ (`default`), అది ఖాతా సొంత ఫోల్డర్‌లో ఉంటే `account`, ఖాతాకు ఈ ఫోల్డర్ కోసం మాత్రమే ఉన్న MCP సర్వర్‌కు `local`, లేదా ఖాతాకు అది లేకపోతే `-`. ఇన్‌స్టాల్ అయి ఉండి ఆన్ చేయని plugin దాని సోర్స్ తర్వాత `(off)` ను చూపిస్తుంది. `!` సమస్య ఉన్న అంశాన్ని సూచిస్తుంది, అది టేబుల్ కింద వివరించబడుతుంది; సాధారణ `claude` కు మాత్రమే ఉన్న సర్వర్ల వంటి, ఏ ఖాతాకూ చెందని సమస్యలు పేరు లేకుండా చివరలో వస్తాయి. `ccam ext --json` అదే నివేదికను ప్రింట్ చేస్తుంది (`accounts`, ప్రాజెక్ట్, `problems`), Claude Code నెట్‌వర్క్ ద్వారా చేరే MCP సర్వర్లపై `remote` తో; ఖాతాలు లేకపోతే అది `accounts: []` తో నివేదికను ప్రింట్ చేస్తుంది.

## ఇది ఏమి చదువుతుంది

| రకం | ఎక్కడ | చూపించే సోర్స్ |
| --- | --- | --- |
| Skills | `skills/` లో `SKILL.md` ఉన్న ప్రతి ఫోల్డర్ | `skills/` లింక్ అయిన గ్రూప్, లేదా `account` |
| Commands, agents | `commands/`, `agents/` కింద ఉన్న `.md` ఫైల్‌లు, ఉప ఫోల్డర్‌లు `folder/name` గా | గ్రూప్, లేదా `account` |
| Plugins | `plugins/installed_plugins.json`, ఖాతా, ప్రాజెక్ట్ సెట్టింగ్‌లలోని `enabledPlugins` తో ఆన్ లేదా ఆఫ్; ఒకే ప్రాజెక్ట్ కోసం ఇన్‌స్టాల్ చేసిన plugin ఆ ప్రాజెక్ట్‌లో మాత్రమే కనిపిస్తుంది | గ్రూప్, లేదా `account`; ఆన్ చేయనప్పుడు `(off)` |
| MCP సర్వర్లు | గ్రూప్ యొక్క `mcp.json` (MCP షేర్ అయినప్పుడు), ఖాతా యొక్క `.claude.json` (దాని సొంత సర్వర్లు, ఈ ఫోల్డర్ కోసం మాత్రమే దానికి ఉన్నవి), దాని plugins సర్వర్లు, `plugin:<plugin>:<server>` గా | గ్రూప్, `account` లేదా `local` |
| Hooks | `settings.json` లోని `hooks`, ఈవెంట్ వారీగా లెక్కించి | గ్రూప్, లేదా `account` |
| ఈ ప్రాజెక్ట్ | ప్రస్తుత ఫోల్డర్‌లోని `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json`, `.mcp.json` | ఒక్కసారి జాబితా అవుతుంది; ఫోల్డర్‌ను ఇంకా trust చేయని ఖాతాకు ఆ ఫోల్డర్ `.mcp.json` సర్వర్ల కోసం ఒక సమస్య చూపబడుతుంది |

## ఇది నివేదించే సమస్యలు

- **విరిగిన లింక్‌లు**: `skills/`, `commands/` లేదా `agents/` లో లక్ష్యం లేని లింక్, లేదా లేని చోటుకు లింక్ అయిన షేర్ చేసిన ఫోల్డర్ (`skills`, `commands`, `agents`, `plugins`) (`ccam doctor --fix` తో సరిచేయండి, అది లేని గ్రూప్ ఫోల్డర్‌ను మళ్లీ సృష్టిస్తుంది).
- **`SKILL.md` లేని skill ఫోల్డర్**: Claude Code దాన్ని వదిలేస్తుంది.
- **రెండు చోట్ల వేర్వేరుగా నిర్వచించిన MCP సర్వర్**: గ్రూప్ యొక్క `mcp.json`, ఖాతా యొక్క `.claude.json`, లేదా వాటిలో ఒకటి, ప్రాజెక్ట్ యొక్క `.mcp.json`. రెండు చోట్లా అవే సెట్టింగ్‌లు ఉంటే (`mcp_strategy = "sync"` రాసేది అదే) అది సమస్య కాదు. ఖాతా సొంత సర్వర్‌ను `ccam mcp import <account> <server>` తో గ్రూప్‌లోకి తరలించండి.
- **చెల్లుబాటు అయ్యే JSON కాని `installed_plugins.json`**, `enabledPlugins` పేర్కొనని plugin, ఫోల్డర్ లేని plugin.
- **ఒకే ఖాతా లోడ్ చేసే సర్వర్లు**: ఈ ఖాతా `.claude.json` లో మాత్రమే ఉండి, దాని గ్రూప్‌లోని ఇతర ఖాతాలు లోడ్ చేయని సర్వర్, ఈ ఫోల్డర్ కోసం ఖాతా local scope లో ఉన్న సర్వర్ (`-s user` ఇవ్వకపోతే `claude mcp add` చేసేది అదే). ప్రతిదానితో, దాన్ని షేర్ చేసే `ccam mcp import` కమాండ్ వస్తుంది. ఏ ఖాతా లోడ్ చేయని, సాధారణ `claude` యొక్క `~/.claude.json` లోని సర్వర్లతో `ccam migrate --mcp` కమాండ్ వస్తుంది.
- **ఆమోదం కోసం వేచి ఉన్న ప్రాజెక్ట్ సర్వర్**: ఖాతా ఇంకా ఫోల్డర్‌ను trust చేయలేదు, కాబట్టి Claude Code ఒక్కసారి అడుగుతుంది, అప్పటి వరకు ప్రాజెక్ట్ `.mcp.json` సర్వర్లు వేచి ఉంటాయి. వాటిని ఎక్కడ ఆమోదించాలో సమస్య చెబుతుంది: సెట్టింగ్‌ల గ్రూప్ ఖాతాల, లేదా ఈ ఖాతా యొక్క `settings.json` లో `"enabledMcpjsonServers"`. git repository ప్రధాన worktree కు ఇచ్చిన trust దాని ఇతర worktrees లోనూ వర్తిస్తుంది.
- **సైన్-ఇన్‌లు, connectors**: మరో ఖాతా సైన్ ఇన్ అయి ఉండగా ఈ ఖాతాలో సైన్-ఇన్ అవసరమైన MCP సర్వర్ (ఒక సైన్-ఇన్ ఒక ఖాతాకు మాత్రమే వర్తిస్తుంది; `/mcp` తో సైన్ ఇన్ చేయండి), మరో ఖాతా ఉపయోగించి ఈ ఖాతా ఉపయోగించని claude.ai connectors.

## భద్రత, ప్రత్యేక సందర్భాలు

- `ccam ext` ఏమీ మార్చదు, నెట్‌వర్క్ అభ్యర్థనలు చేయదు. అది సర్వర్ పేర్లను, hook సంఖ్యలను చదువుతుంది, MCP నిర్వచనాల్లోని `env` లేదా headers విలువలను ఎప్పుడూ చదవదు. సైన్-ఇన్ తనిఖీ కోసం అది ప్రతి ఖాతా లాగిన్‌లోని MCP సైన్-ఇన్‌ల పేర్లను చదువుతుంది (macOS లో Keychain నుంచి), టోకెన్‌లను ఎప్పుడూ చదవదు, అది కూడా మరో ఖాతా కూడా లోడ్ చేసే సర్వర్‌ను లోడ్ చేసే ఖాతాలకు మాత్రమే.
- అది డిస్క్‌పై ఉన్నదాన్ని చూపిస్తుంది. managed settings అడ్డుకున్న సర్వర్ అయినా కనిపిస్తుంది; సెట్టింగ్‌లలో ఆఫ్ చేసిన plugin `(off)` ను చూపిస్తుంది.
- ఆ ప్రాజెక్ట్ యొక్క `.claude` ఫోల్డర్‌ను, `.mcp.json` ను చేర్చడానికి దాన్ని ప్రాజెక్ట్ ఫోల్డర్‌లో రన్ చేయండి.

## సంబంధిత పేజీలు

- [ఖాతాల మధ్య షేరింగ్](/features/share)
- [షేర్ చేసిన MCP సర్వర్లు](/features/mcp)
- [తనిఖీ చేయడం, సరిచేయడం](/features/doctor)
