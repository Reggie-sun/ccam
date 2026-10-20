---
title: Kukagua viendelezi
description: ccam ext huonyesha skills, commands, plugins, seva za MCP na hooks ambazo kila akaunti hupakia, kila kimoja kinatoka wapi, na nini kimeharibika.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` hujibu swali "kwa nini skill hii (au seva hii ya MCP) inafanya kazi katika akaunti moja na si katika nyingine?". Huorodhesha, kwa kila akaunti, viendelezi ambavyo Claude Code hupakia na kila kimoja kinatoka wapi: kundi la kushiriki, folda ya akaunti yenyewe, au mradi ulio ndani yake. Chini ya jedwali huorodhesha kinachozuia kiendelezi kupakiwa kama unavyotarajia. Husoma tu.

## Matumizi

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

Kila kisanduku hutaja chanzo cha kipengele hicho katika akaunti hiyo: kundi la kushiriki (`default`), `account` kinapokaa katika folda ya akaunti yenyewe, au `-` akaunti isipokuwa nacho. `!` huashiria kipengele chenye tatizo, kinachoelezwa chini ya jedwali.

## Kile inachosoma

| Aina | Wapi | Chanzo kinachoonyeshwa |
| --- | --- | --- |
| Skills | kila folda katika `skills/` yenye `SKILL.md` | kundi ambalo `skills/` huunganisha nalo, au `account` |
| Commands, agents | faili za `.md` chini ya `commands/` na `agents/`, folda ndogo kama `folder/name` | kundi, au `account` |
| Plugins | `plugins/installed_plugins.json` | kundi, au `account` |
| Seva za MCP | `mcp.json` ya kundi (MCP inaposhirikiwa) na `.claude.json` ya akaunti | kundi, au `account` |
| Hooks | `hooks` katika `settings.json`, zikihesabiwa kwa kila tukio | kundi, au `account` |
| Mradi huu | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` na `.mcp.json` katika folda ya sasa | huorodheshwa mara moja, kwa kila akaunti |

## Matatizo inayoripoti

- **Viungo vilivyovunjika**: kiungo katika `skills/`, `commands/` au `agents/` ambacho lengo lake halipo tena, au folda inayoshirikiwa (`skills`, `commands`, `agents`, `plugins`) inayounganisha na mahali pasipokuwepo (irekebishe kwa `ccam doctor --fix`).
- **Folda ya skill bila `SKILL.md`**: Claude Code huiruka.
- **Seva ya MCP iliyofafanuliwa tofauti katika sehemu mbili**: `mcp.json` ya kundi na `.claude.json` ya akaunti, au mojawapo yao na `.mcp.json` ya mradi. Mipangilio ile ile katika sehemu zote mbili (ambayo `mcp_strategy = "sync"` huandika) si tatizo. Hamisha seva ya akaunti yenyewe kwenda kwenye kundi kwa `ccam mcp import <account> <server>`.
- **`installed_plugins.json` ambayo si JSON halali.**

## Usalama na hali maalum

- `ccam ext` haibadilishi chochote na haifanyi maombi ya mtandao. Husoma majina ya seva na idadi ya hooks, kamwe si thamani za `env` au headers katika ufafanuzi wa MCP.
- Huonyesha kilichopo kwenye diski. Plugin iliyosakinishwa lakini imezimwa katika mipangilio, au seva iliyozuiwa na managed settings, bado huonekana.
- Iendeshe katika folda ya mradi ili kujumuisha folda ya `.claude` na `.mcp.json` ya mradi huo.

## Zinazohusiana

- [Kushiriki kati ya akaunti](/features/share)
- [Seva za MCP zinazoshirikiwa](/features/mcp)
- [Kukagua na kurekebisha](/features/doctor)
