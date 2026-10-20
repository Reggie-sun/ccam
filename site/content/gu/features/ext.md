---
title: extensions તપાસવાં
description: ccam ext દરેક એકાઉન્ટ લોડ કરે છે તે skills, commands, plugins, MCP સર્વર અને hooks, દરેક ક્યાંથી આવે છે, અને શું તૂટેલું છે તે બતાવે છે.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` આ પ્રશ્નનો જવાબ આપે છે: "આ skill (કે MCP સર્વર) એક એકાઉન્ટમાં કેમ ચાલે છે અને બીજામાં કેમ નહીં?". તે દરેક એકાઉન્ટ માટે Claude Code લોડ કરે છે તે extensions અને દરેક ક્યાંથી આવે છે તેની યાદી આપે છે: કોઈ share group, એકાઉન્ટનું પોતાનું ફોલ્ડર, કે તમે જેમાં છો તે પ્રોજેક્ટ. table ની નીચે તે એ બધું બતાવે છે જે કોઈ extension ને તમારી અપેક્ષા મુજબ લોડ થતાં રોકે છે. તે ફક્ત વાંચે છે.

## ઉપયોગ

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

દરેક ખાનું એ એકાઉન્ટમાં એ આઇટમનો source બતાવે છે: share group (`default`), આઇટમ એકાઉન્ટના પોતાના ફોલ્ડરમાં હોય ત્યારે `account`, અથવા એકાઉન્ટ પાસે તે ન હોય ત્યારે `-`. `!` સમસ્યાવાળી આઇટમ દર્શાવે છે, જે table ની નીચે સમજાવેલી હોય છે.

## તે શું વાંચે છે

| પ્રકાર | ક્યાં | બતાવેલો source |
| --- | --- | --- |
| Skills | `skills/` માંનું `SKILL.md` વાળું દરેક ફોલ્ડર | `skills/` જે group સાથે લિંક છે તે, અથવા `account` |
| Commands, agents | `commands/` અને `agents/` હેઠળની `.md` ફાઇલો, પેટા-ફોલ્ડર `folder/name` તરીકે | group, અથવા `account` |
| Plugins | `plugins/installed_plugins.json` | group, અથવા `account` |
| MCP સર્વર | group ની `mcp.json` (MCP શેર થતું હોય ત્યારે) અને એકાઉન્ટની `.claude.json` | group, અથવા `account` |
| Hooks | `settings.json` માં `hooks`, event પ્રમાણે ગણેલા | group, અથવા `account` |
| આ પ્રોજેક્ટ | હાલના ફોલ્ડરમાં `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` અને `.mcp.json` | એક જ વાર, બધાં એકાઉન્ટ માટે |

## તે કઈ સમસ્યાઓ બતાવે છે

- **તૂટેલી લિંક**: `skills/`, `commands/` કે `agents/` માંની એવી લિંક જેનું લક્ષ્ય નથી રહ્યું, અથવા એવું શેર કરેલું ફોલ્ડર (`skills`, `commands`, `agents`, `plugins`) જે ન હોય તેવી જગ્યાએ લિંક થાય છે (`ccam doctor --fix` વડે સુધારો).
- **`SKILL.md` વગરનું skill ફોલ્ડર**: Claude Code તેને છોડી દે છે.
- **બે જગ્યાએ જુદી રીતે વ્યાખ્યાયિત MCP સર્વર**: group ની `mcp.json` અને એકાઉન્ટની `.claude.json`, અથવા તેમાંથી એક અને પ્રોજેક્ટની `.mcp.json`. બંને જગ્યાએ એ જ સેટિંગ્સ હોય (જે `mcp_strategy = "sync"` લખે છે) તો તે સમસ્યા નથી. એકાઉન્ટના પોતાના સર્વરને `ccam mcp import <account> <server>` વડે group માં ખસેડો.
- **માન્ય JSON ન હોય તેવી `installed_plugins.json`.**

## સલામતી અને ખાસ કિસ્સા

- `ccam ext` કંઈ બદલતું નથી અને કોઈ network request કરતું નથી. તે સર્વરનાં નામ અને hooks ની સંખ્યા વાંચે છે, MCP વ્યાખ્યાઓમાંનાં `env` કે headers નાં મૂલ્યો ક્યારેય નહીં.
- તે disk પર જે છે તે બતાવે છે. ઇન્સ્ટૉલ થયેલું પણ સેટિંગ્સમાં બંધ કરેલું plugin, કે managed settings વડે અવરોધાયેલું સર્વર, હજી પણ દેખાય છે.
- એ પ્રોજેક્ટનું `.claude` ફોલ્ડર અને `.mcp.json` સામેલ કરવા માટે તેને પ્રોજેક્ટ ફોલ્ડરમાં ચલાવો.

## સંબંધિત

- [એકાઉન્ટ વચ્ચે શેરિંગ](/features/share)
- [શેર કરેલા MCP સર્વર](/features/mcp)
- [તપાસ અને સમારકામ](/features/doctor)
