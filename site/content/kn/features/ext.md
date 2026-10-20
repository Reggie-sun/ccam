---
title: ಎಕ್ಸ್‌ಟೆನ್ಶನ್‌ಗಳ ಪರಿಶೀಲನೆ
description: ccam ext ಪ್ರತಿ ಖಾತೆ ಲೋಡ್ ಮಾಡುವ ಸ್ಕಿಲ್, ಕಮಾಂಡ್, ಪ್ಲಗಿನ್, MCP ಸರ್ವರ್ ಮತ್ತು ಹುಕ್‌ಗಳನ್ನು, ಅವುಗಳ ಮೂಲವನ್ನು ಮತ್ತು ಮುರಿದಿರುವುದನ್ನು ತೋರಿಸುತ್ತದೆ.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

"ಈ ಸ್ಕಿಲ್ (ಅಥವಾ MCP ಸರ್ವರ್) ಒಂದು ಖಾತೆಯಲ್ಲಿ ಕೆಲಸ ಮಾಡುತ್ತದೆ, ಇನ್ನೊಂದರಲ್ಲಿ ಏಕೆ ಇಲ್ಲ?" ಎಂಬ ಪ್ರಶ್ನೆಗೆ `ccam ext` ಉತ್ತರಿಸುತ್ತದೆ. ಅದು ಪ್ರತಿ ಖಾತೆಗೂ Claude Code ಲೋಡ್ ಮಾಡುವ ಎಕ್ಸ್‌ಟೆನ್ಶನ್‌ಗಳನ್ನು ಮತ್ತು ಪ್ರತಿಯೊಂದೂ ಎಲ್ಲಿಂದ ಬಂದಿದೆ ಎಂಬುದನ್ನು ಪಟ್ಟಿ ಮಾಡುತ್ತದೆ: ಹಂಚಿಕೆ ಗುಂಪು, ಖಾತೆಯ ಸ್ವಂತ ಫೋಲ್ಡರ್, ಅಥವಾ ನೀವು ಇರುವ ಪ್ರಾಜೆಕ್ಟ್. ಟೇಬಲ್‌ನ ಕೆಳಗೆ, ಎಕ್ಸ್‌ಟೆನ್ಶನ್ ನೀವು ನಿರೀಕ್ಷಿಸಿದಂತೆ ಲೋಡ್ ಆಗುವುದನ್ನು ತಡೆಯುವ ಸಂಗತಿಗಳನ್ನು ಪಟ್ಟಿ ಮಾಡುತ್ತದೆ. ಅದು ಓದುತ್ತದೆ ಮಾತ್ರ.

## ಬಳಕೆ

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

ಪ್ರತಿ ಸೆಲ್ ಆ ಖಾತೆಯಲ್ಲಿ ಆ ಐಟಂನ ಮೂಲವನ್ನು ಹೆಸರಿಸುತ್ತದೆ: ಹಂಚಿಕೆ ಗುಂಪು (`default`), ಖಾತೆಯ ಸ್ವಂತ ಫೋಲ್ಡರ್‌ನಲ್ಲಿದ್ದರೆ `account`, ಅಥವಾ ಖಾತೆಯಲ್ಲಿ ಅದು ಇಲ್ಲದಿದ್ದರೆ `-`. `!` ಸಮಸ್ಯೆಯಿರುವ ಐಟಂ ಅನ್ನು ಗುರುತಿಸುತ್ತದೆ, ಅದನ್ನು ಟೇಬಲ್‌ನ ಕೆಳಗೆ ವಿವರಿಸಲಾಗುತ್ತದೆ.

## ಅದು ಏನನ್ನು ಓದುತ್ತದೆ

| ಪ್ರಕಾರ | ಎಲ್ಲಿ | ತೋರಿಸುವ ಮೂಲ |
| --- | --- | --- |
| ಸ್ಕಿಲ್‌ಗಳು | `SKILL.md` ಇರುವ `skills/` ನಲ್ಲಿನ ಪ್ರತಿ ಫೋಲ್ಡರ್ | `skills/` ಲಿಂಕ್ ಆಗಿರುವ ಗುಂಪು, ಅಥವಾ `account` |
| ಕಮಾಂಡ್‌ಗಳು, ಏಜೆಂಟ್‌ಗಳು | `commands/` ಮತ್ತು `agents/` ಅಡಿಯ `.md` ಫೈಲ್‌ಗಳು, ಉಪಫೋಲ್ಡರ್‌ಗಳು `folder/name` ಆಗಿ | ಗುಂಪು, ಅಥವಾ `account` |
| ಪ್ಲಗಿನ್‌ಗಳು | `plugins/installed_plugins.json` | ಗುಂಪು, ಅಥವಾ `account` |
| MCP ಸರ್ವರ್‌ಗಳು | ಗುಂಪಿನ `mcp.json` (MCP ಹಂಚಿಕೆಯಾಗಿದ್ದಾಗ) ಮತ್ತು ಖಾತೆಯ `.claude.json` | ಗುಂಪು, ಅಥವಾ `account` |
| ಹುಕ್‌ಗಳು | `settings.json` ನಲ್ಲಿನ `hooks`, ಪ್ರತಿ ಈವೆಂಟ್‌ಗೆ ಎಣಿಸಲಾಗಿದೆ | ಗುಂಪು, ಅಥವಾ `account` |
| ಈ ಪ್ರಾಜೆಕ್ಟ್ | ಪ್ರಸ್ತುತ ಫೋಲ್ಡರ್‌ನಲ್ಲಿನ `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` ಮತ್ತು `.mcp.json` | ಒಮ್ಮೆ ಪಟ್ಟಿ ಮಾಡಲಾಗುತ್ತದೆ, ಪ್ರತಿ ಖಾತೆಗೂ |

## ಅದು ವರದಿ ಮಾಡುವ ಸಮಸ್ಯೆಗಳು

- **ಮುರಿದ ಲಿಂಕ್‌ಗಳು**: ಗುರಿ ಇಲ್ಲವಾಗಿರುವ `skills/`, `commands/` ಅಥವಾ `agents/` ನಲ್ಲಿನ ಲಿಂಕ್, ಅಥವಾ ಕಾಣೆಯಾದ ಸ್ಥಳಕ್ಕೆ ಲಿಂಕ್ ಆಗಿರುವ ಹಂಚಿಕೊಂಡ ಫೋಲ್ಡರ್ (`skills`, `commands`, `agents`, `plugins`) (`ccam doctor --fix` ಮೂಲಕ ಸರಿಪಡಿಸಿ).
- **`SKILL.md` ಇಲ್ಲದ ಸ್ಕಿಲ್ ಫೋಲ್ಡರ್**: Claude Code ಅದನ್ನು ಬಿಟ್ಟುಬಿಡುತ್ತದೆ.
- **ಎರಡು ಸ್ಥಳಗಳಲ್ಲಿ ಬೇರೆ ಬೇರೆಯಾಗಿ ವ್ಯಾಖ್ಯಾನಿಸಿದ MCP ಸರ್ವರ್**: ಗುಂಪಿನ `mcp.json` ಮತ್ತು ಖಾತೆಯ `.claude.json`, ಅಥವಾ ಇವುಗಳಲ್ಲಿ ಒಂದು ಮತ್ತು ಪ್ರಾಜೆಕ್ಟ್‌ನ `.mcp.json`. ಎರಡೂ ಸ್ಥಳಗಳಲ್ಲಿ ಒಂದೇ ಸೆಟ್ಟಿಂಗ್‌ಗಳು ಇರುವುದು (`mcp_strategy = "sync"` ಬರೆಯುವಂತೆ) ಸಮಸ್ಯೆಯಲ್ಲ. ಖಾತೆಯ ಸ್ವಂತ ಸರ್ವರ್ ಅನ್ನು `ccam mcp import <account> <server>` ಮೂಲಕ ಗುಂಪಿಗೆ ಸ್ಥಳಾಂತರಿಸಿ.
- **ಮಾನ್ಯ JSON ಅಲ್ಲದ `installed_plugins.json`.**

## ಸುರಕ್ಷತೆ ಮತ್ತು ವಿಶೇಷ ಸಂದರ್ಭಗಳು

- `ccam ext` ಏನನ್ನೂ ಬದಲಾಯಿಸುವುದಿಲ್ಲ ಮತ್ತು ಯಾವುದೇ ನೆಟ್‌ವರ್ಕ್ ವಿನಂತಿ ಮಾಡುವುದಿಲ್ಲ. ಅದು ಸರ್ವರ್ ಹೆಸರುಗಳು ಮತ್ತು ಹುಕ್ ಎಣಿಕೆಗಳನ್ನು ಓದುತ್ತದೆ, MCP ವ್ಯಾಖ್ಯಾನಗಳಲ್ಲಿನ `env` ಅಥವಾ ಹೆಡರ್‌ಗಳ ಮೌಲ್ಯಗಳನ್ನು ಎಂದಿಗೂ ಓದುವುದಿಲ್ಲ.
- ಅದು ಡಿಸ್ಕ್‌ನಲ್ಲಿರುವುದನ್ನು ತೋರಿಸುತ್ತದೆ. ಇನ್‌ಸ್ಟಾಲ್ ಆಗಿದ್ದರೂ ಸೆಟ್ಟಿಂಗ್‌ಗಳಲ್ಲಿ ಆಫ್ ಮಾಡಿದ ಪ್ಲಗಿನ್, ಅಥವಾ ನಿರ್ವಹಿತ ಸೆಟ್ಟಿಂಗ್‌ಗಳಿಂದ ತಡೆಹಿಡಿದ ಸರ್ವರ್ ಕೂಡ ಕಾಣಿಸುತ್ತದೆ.
- ಆ ಪ್ರಾಜೆಕ್ಟ್‌ನ `.claude` ಫೋಲ್ಡರ್ ಮತ್ತು `.mcp.json` ಅನ್ನು ಸೇರಿಸಲು ಪ್ರಾಜೆಕ್ಟ್ ಫೋಲ್ಡರ್‌ನಲ್ಲಿ ಇದನ್ನು ರನ್ ಮಾಡಿ.

## ಸಂಬಂಧಿತ

- [ಖಾತೆಗಳ ನಡುವೆ ಹಂಚಿಕೆ](/features/share)
- [ಹಂಚಿಕೊಂಡ MCP ಸರ್ವರ್‌ಗಳು](/features/mcp)
- [ಪರಿಶೀಲನೆ ಮತ್ತು ದುರಸ್ತಿ](/features/doctor)
