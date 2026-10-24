---
title: extensions તપાસવાં
description: ccam ext દરેક એકાઉન્ટ લોડ કરે છે તે skills, commands, plugins, MCP સર્વર અને hooks, દરેક ક્યાંથી આવે છે, અને શું તૂટેલું છે તે બતાવે છે.
updated: 2026-10-23
order: 14
translation: machine
sourceHash: 45740e4efbef
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

દરેક ખાનું એ એકાઉન્ટમાં એ આઇટમનો source બતાવે છે: share group (`default`), આઇટમ એકાઉન્ટના પોતાના ફોલ્ડરમાં હોય ત્યારે `account`, એકાઉન્ટ પાસે ફક્ત આ ફોલ્ડર માટે હોય તેવા MCP સર્વર માટે `local`, અથવા એકાઉન્ટ પાસે તે ન હોય ત્યારે `-`. ઇન્સ્ટૉલ થયેલું પણ ચાલુ ન કરેલું plugin તેના source પછી `(off)` બતાવે છે. `!` સમસ્યાવાળી આઇટમ દર્શાવે છે, જે table ની નીચે સમજાવેલી હોય છે; કોઈ એકાઉન્ટની ન હોય તેવી સમસ્યાઓ, જેમ કે ફક્ત સાદા `claude` પાસે હોય તેવા સર્વર, નામ વગર છેલ્લે આવે છે. `ccam ext --json` એ જ report છાપે છે (`accounts`, પ્રોજેક્ટ, `problems`), Claude Code જે MCP સર્વર સુધી network દ્વારા પહોંચે છે તેમના પર `remote` સાથે; કોઈ એકાઉન્ટ ન હોય ત્યારે તે `accounts: []` સાથેનો report છાપે છે.

## તે શું વાંચે છે

| પ્રકાર | ક્યાં | બતાવેલો source |
| --- | --- | --- |
| Skills | `skills/` માંનું `SKILL.md` વાળું દરેક ફોલ્ડર | `skills/` જે group સાથે લિંક છે તે, અથવા `account` |
| Commands, agents | `commands/` અને `agents/` હેઠળની `.md` ફાઇલો, પેટા-ફોલ્ડર `folder/name` તરીકે | group, અથવા `account` |
| Plugins | `plugins/installed_plugins.json`, જેને એકાઉન્ટનાં અને પ્રોજેક્ટનાં સેટિંગ્સમાંનું `enabledPlugins` ચાલુ કે બંધ કરે છે; ફક્ત એક પ્રોજેક્ટ માટે ઇન્સ્ટૉલ કરેલું plugin ફક્ત એ પ્રોજેક્ટમાં દેખાય છે | group, અથવા `account`; ચાલુ ન હોય ત્યારે `(off)` |
| MCP સર્વર | group ની `mcp.json` (MCP શેર થતું હોય ત્યારે), એકાઉન્ટની `.claude.json` (તેના પોતાના સર્વર, અને ફક્ત આ ફોલ્ડર માટેના તેના સર્વર), અને તેનાં plugins ના સર્વર, `plugin:<plugin>:<server>` તરીકે | group, `account` કે `local` |
| Hooks | `settings.json` માં `hooks`, event પ્રમાણે ગણેલા | group, અથવા `account` |
| આ પ્રોજેક્ટ | હાલના ફોલ્ડરમાં `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` અને `.mcp.json` | એક જ વાર; જે એકાઉન્ટે હજી ફોલ્ડરને trust નથી કર્યું તેને ફોલ્ડરના `.mcp.json` સર્વર માટે સમસ્યા બતાવાય છે |

## તે કઈ સમસ્યાઓ બતાવે છે

- **તૂટેલી લિંક**: `skills/`, `commands/` કે `agents/` માંની એવી લિંક જેનું લક્ષ્ય નથી રહ્યું, અથવા એવું શેર કરેલું ફોલ્ડર (`skills`, `commands`, `agents`, `plugins`) જે ન હોય તેવી જગ્યાએ લિંક થાય છે (`ccam doctor --fix` વડે સુધારો, જે ખૂટતું group ફોલ્ડર ફરી બનાવે છે).
- **`SKILL.md` વગરનું skill ફોલ્ડર**: Claude Code તેને છોડી દે છે.
- **બે જગ્યાએ જુદી રીતે વ્યાખ્યાયિત MCP સર્વર**: group ની `mcp.json` અને એકાઉન્ટની `.claude.json`, અથવા તેમાંથી એક અને પ્રોજેક્ટની `.mcp.json`. બંને જગ્યાએ એ જ સેટિંગ્સ હોય (જે `mcp_strategy = "sync"` લખે છે) તો તે સમસ્યા નથી. એકાઉન્ટના પોતાના સર્વરને `ccam mcp import <account> <server>` વડે group માં ખસેડો.
- **માન્ય JSON ન હોય તેવી `installed_plugins.json`**, `enabledPlugins` જેનું નામ ન આપતું હોય તેવું plugin, અને જેનું ફોલ્ડર ખૂટતું હોય તેવું plugin.
- **ફક્ત એક એકાઉન્ટ લોડ કરતું હોય તેવા સર્વર**: ફક્ત આ એકાઉન્ટની `.claude.json` માં હોય તેવો સર્વર, જે તેના group નાં બીજાં એકાઉન્ટ લોડ કરતાં નથી, અને આ ફોલ્ડર માટે એકાઉન્ટના local scope માં હોય તેવો સર્વર (`-s user` ન આપ્યું હોય તો `claude mcp add` આવું જ કરે છે). દરેક સાથે તેને શેર કરતો `ccam mcp import` કમાન્ડ આવે છે. સાદા `claude` ની `~/.claude.json` માંના જે સર્વર કોઈ એકાઉન્ટ લોડ કરતું નથી તેમની સાથે `ccam migrate --mcp` કમાન્ડ આવે છે.
- **મંજૂરીની રાહ જોતો પ્રોજેક્ટ સર્વર**: એકાઉન્ટે હજી ફોલ્ડરને trust નથી કર્યું, તેથી Claude Code એક વાર પૂછે છે, અને ત્યાં સુધી પ્રોજેક્ટની `.mcp.json` ના સર્વર રાહ જુએ છે. સમસ્યા જણાવે છે કે તેમને ક્યાં મંજૂર કરવા: settings group નાં એકાઉન્ટની, કે આ એકાઉન્ટની `settings.json` માં `"enabledMcpjsonServers"`. git repository ના મુખ્ય worktree ને આપેલો trust તેનાં બીજાં worktrees માં પણ ગણાય છે.
- **સાઇન-ઇન અને connectors**: જે MCP સર્વરમાં બીજું એકાઉન્ટ સાઇન ઇન છે પણ આ એકાઉન્ટમાં સાઇન-ઇન જોઈએ છે (સાઇન-ઇન ફક્ત એક એકાઉન્ટ માટે જ ચાલે છે; `/mcp` વડે સાઇન ઇન કરો), અને બીજા એકાઉન્ટે વાપરેલા પણ આ એકાઉન્ટે ન વાપરેલા claude.ai connectors.

## સલામતી અને ખાસ કિસ્સા

- `ccam ext` કંઈ બદલતું નથી અને કોઈ network request કરતું નથી. તે સર્વરનાં નામ અને hooks ની સંખ્યા વાંચે છે, MCP વ્યાખ્યાઓમાંનાં `env` કે headers નાં મૂલ્યો ક્યારેય નહીં. સાઇન-ઇનની તપાસ માટે તે દરેક એકાઉન્ટના લૉગિનમાંથી (macOS પર Keychain માંથી) MCP સાઇન-ઇનનાં નામ વાંચે છે, ટોકન ક્યારેય નહીં, અને ફક્ત એવાં એકાઉન્ટ માટે જે બીજું એકાઉન્ટ પણ લોડ કરતું હોય તેવો સર્વર લોડ કરે છે.
- તે disk પર જે છે તે બતાવે છે. managed settings વડે અવરોધાયેલું સર્વર હજી પણ દેખાય છે; સેટિંગ્સમાં બંધ કરેલું plugin `(off)` બતાવે છે.
- એ પ્રોજેક્ટનું `.claude` ફોલ્ડર અને `.mcp.json` સામેલ કરવા માટે તેને પ્રોજેક્ટ ફોલ્ડરમાં ચલાવો.

## સંબંધિત

- [એકાઉન્ટ વચ્ચે શેરિંગ](/features/share)
- [શેર કરેલા MCP સર્વર](/features/mcp)
- [તપાસ અને સમારકામ](/features/doctor)
