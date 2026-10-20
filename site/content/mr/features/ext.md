---
title: extensions तपासणे
description: ccam ext प्रत्येक खाते कोणते skills, commands, plugins, MCP सर्व्हर आणि hooks लोड करते, प्रत्येक कुठून येते आणि काय बिघडले आहे ते दाखवते.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` या प्रश्नाचे उत्तर देते: "हे skill (किंवा MCP सर्व्हर) एका खात्यात चालते आणि दुसऱ्यात का नाही?". ते प्रत्येक खात्यासाठी Claude Code लोड करत असलेली extensions आणि प्रत्येक कुठून येते ते दाखवते: एखादा share group, खात्याचा स्वतःचा फोल्डर, किंवा तुम्ही ज्या प्रकल्पात आहात तो. table खाली ते अशा गोष्टींची यादी देते ज्या एखादे extension तुमच्या अपेक्षेप्रमाणे लोड होण्यापासून रोखतात. ते फक्त वाचते.

## वापर

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

प्रत्येक cell त्या खात्यातील त्या आयटमचा source सांगतो: share group (`default`), खात्याच्या स्वतःच्या फोल्डरमध्ये असल्यास `account`, किंवा खात्याकडे तो नसल्यास `-`. `!` समस्या असलेला आयटम दर्शवते, ज्याचे स्पष्टीकरण table खाली असते.

## ते काय वाचते

| प्रकार | कुठे | दाखवलेला source |
| --- | --- | --- |
| Skills | `skills/` मधील `SKILL.md` असलेला प्रत्येक फोल्डर | `skills/` ज्या group कडे link करतो तो, किंवा `account` |
| Commands, agents | `commands/` आणि `agents/` खालील `.md` फाइल्स, उप-फोल्डर `folder/name` म्हणून | group, किंवा `account` |
| Plugins | `plugins/installed_plugins.json` | group, किंवा `account` |
| MCP सर्व्हर | group चा `mcp.json` (MCP शेअर केलेले असताना) आणि खात्याचा `.claude.json` | group, किंवा `account` |
| Hooks | `settings.json` मधील `hooks`, प्रत्येक event नुसार मोजलेले | group, किंवा `account` |
| हा प्रकल्प | सध्याच्या फोल्डरमधील `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` आणि `.mcp.json` | एकदाच दाखवले जाते, प्रत्येक खात्यासाठी |

## ते कोणत्या समस्या सांगते

- **तुटलेल्या links**: `skills/`, `commands/` किंवा `agents/` मधील अशी link जिचे target आता नाही, किंवा अस्तित्वात नसलेल्या जागेकडे link करणारा शेअर केलेला फोल्डर (`skills`, `commands`, `agents`, `plugins`) (`ccam doctor --fix` ने दुरुस्त करा).
- **`SKILL.md` नसलेला skill फोल्डर**: Claude Code तो वगळते.
- **दोन ठिकाणी वेगवेगळ्या प्रकारे परिभाषित केलेला MCP सर्व्हर**: group चा `mcp.json` आणि खात्याचा `.claude.json`, किंवा त्यांपैकी एक आणि प्रकल्पाचा `.mcp.json`. दोन्ही ठिकाणी सारख्याच सेटिंग्ज (ज्या `mcp_strategy = "sync"` लिहिते) ही समस्या नाही. खात्याचा स्वतःचा सर्व्हर `ccam mcp import <account> <server>` ने group मध्ये हलवा.
- **वैध JSON नसलेली `installed_plugins.json`.**

## सुरक्षितता आणि अपवादात्मक स्थिती

- `ccam ext` काहीही बदलत नाही आणि कोणतीही network request करत नाही. ते सर्व्हरची नावे आणि hooks ची संख्या वाचते, MCP definitions मधील `env` किंवा headers ची मूल्ये कधीच नाही.
- disk वर जे आहे तेच ते दाखवते. इन्स्टॉल केलेले पण सेटिंग्जमध्ये बंद केलेले plugin, किंवा managed settings ने अडवलेला सर्व्हर, तरीही दिसतो.
- त्या प्रकल्पाचा `.claude` फोल्डर आणि `.mcp.json` समाविष्ट करण्यासाठी ते प्रकल्पाच्या फोल्डरमध्ये चालवा.

## संबंधित

- [खात्यांमध्ये शेअरिंग](/features/share)
- [शेअर केलेले MCP सर्व्हर](/features/mcp)
- [तपासणी आणि दुरुस्ती](/features/doctor)
