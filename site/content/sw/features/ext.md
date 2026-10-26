---
title: Kukagua viendelezi
description: ccam ext huonyesha skills, commands, plugins, seva za MCP na hooks ambazo kila akaunti hupakia, kila kimoja kinatoka wapi, na nini kimeharibika.
updated: 2026-10-26
order: 14
translation: machine
sourceHash: "7ca357a7f38a"
---

`ccam ext` hujibu swali "Kwa nini skill hii (au seva hii ya MCP) inafanya kazi katika akaunti moja lakini si katika nyingine?" Huorodhesha, kwa kila akaunti, viendelezi ambavyo Claude Code hupakia na kila kimoja kinatoka wapi: kundi la kushiriki, folda ya akaunti yenyewe, au mradi ulio ndani yake. Chini ya jedwali huorodhesha kinachozuia kiendelezi kupakiwa kama unavyotarajia. Haibadilishi chochote.

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

Kila kisanduku hutaja chanzo cha kipengele hicho katika akaunti hiyo: kundi la kushiriki (`default`), `account` kinapokaa katika folda ya akaunti yenyewe, `local` kwa seva ya MCP ambayo akaunti inayo kwa folda hii pekee, au `-` akaunti isipokuwa nacho. Plugin iliyosakinishwa lakini haijawashwa huonyesha `(off)` baada ya chanzo chake. `!` huashiria kipengele chenye tatizo, kinachoelezwa chini ya jedwali; matatizo yasiyo ya akaunti yoyote, kama seva ambazo `claude` ya kawaida pekee inazo, huja mwisho bila jina. `ccam ext --json` huchapisha ripoti ile ile (`accounts`, mradi, `problems`), pamoja na `remote` kwenye seva za MCP ambazo Claude Code huzifikia kupitia mtandao; bila akaunti yoyote huchapisha ripoti yenye `accounts: []`.

## Kile inachosoma

| Aina | Wapi | Chanzo kinachoonyeshwa |
| --- | --- | --- |
| Skills | kila folda katika `skills/` yenye `SKILL.md` | kundi ambalo `skills/` huunganisha nalo, au `account` |
| Commands, agents | faili za `.md` chini ya `commands/` na `agents/`, folda ndogo kama `folder/name` | kundi, au `account` |
| Plugins | `plugins/installed_plugins.json`, zikiwashwa au kuzimwa na `enabledPlugins` katika mipangilio ya akaunti na ya mradi; plugin iliyosakinishwa kwa mradi mmoja tu huonekana katika mradi huo | kundi, au `account`; `(off)` isipowashwa |
| Seva za MCP | `mcp.json` ya kundi (MCP inaposhirikiwa), `.claude.json` ya akaunti (seva zake yenyewe, na zile ilizo nazo kwa folda hii pekee), na seva za plugins zake, kama `plugin:<plugin>:<server>` | kundi, `account` au `local` |
| Hooks | `hooks` katika `settings.json`, zikihesabiwa kwa kila tukio | kundi, au `account` |
| Mradi huu | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` na `.mcp.json` katika folda ya sasa | huorodheshwa mara moja; akaunti ambayo haijaiamini folda hupata tatizo kwa seva zake za `.mcp.json` |

## Matatizo inayoripoti

- **Viungo vilivyovunjika**: kiungo katika `skills/`, `commands/` au `agents/` ambacho lengo lake halipo tena, au folda inayoshirikiwa (`skills`, `commands`, `agents`, `plugins`) inayounganisha na mahali pasipokuwepo (irekebishe kwa `ccam doctor --fix`, ambayo huunda upya folda ya kundi inayokosekana).
- **Folda ya skill bila `SKILL.md`**: Claude Code huiruka.
- **Seva ya MCP iliyofafanuliwa tofauti katika sehemu mbili**: `mcp.json` ya kundi na `.claude.json` ya akaunti, au mojawapo yao na `.mcp.json` ya mradi. Ufafanuzi unaofanana kabisa katika sehemu zote mbili (ambao `mcp_strategy = "sync"` huandika) si tatizo. Hamisha seva ya akaunti yenyewe kwenda kwenye kundi kwa `ccam mcp import <account> <server>`.
- **`installed_plugins.json` ambayo si JSON halali**, plugin ambayo `enabledPlugins` haiitaji, na plugin ambayo folda yake inakosekana.
- **Seva ambazo akaunti moja tu hupakia**: seva iliyo katika `.claude.json` ya akaunti hii pekee, ambayo akaunti nyingine za kundi lake hazipakii, na seva iliyo katika local scope ya akaunti kwa folda hii (kile ambacho `claude mcp add` hufanya isipopewa `-s user`). Kila moja huja pamoja na amri ya `ccam mcp import` inayoishiriki. Seva zilizo katika `~/.claude.json` ya `claude` ya kawaida ambazo hakuna akaunti inayozipakia huja pamoja na amri ya `ccam migrate --mcp`.
- **Seva ya mradi inayosubiri idhini**: akaunti bado haijaiamini folda, kwa hiyo Claude Code huuliza mara moja, na hadi wakati huo seva za `.mcp.json` za mradi husubiri. Tatizo husema mahali pa kuziidhinisha: `"enabledMcpjsonServers"` katika `settings.json` ya akaunti za kundi la mipangilio, au ya akaunti hii. Imani iliyotolewa kwa worktree kuu ya hazina ya git pia hutumika kwa worktrees zake nyingine.
- **Kuingia na connectors**: seva ya MCP inayohitaji kuingia katika akaunti hii wakati akaunti nyingine imeingia kwayo (kuingia huhusu akaunti moja tu; ingia kwa `/mcp`), na connectors za claude.ai ambazo akaunti nyingine imezitumia na hii haijazitumia.

## Usalama na hali maalum

- `ccam ext` haibadilishi chochote na haifanyi maombi ya mtandao. Husoma majina ya seva na idadi ya hooks, kamwe si thamani za `env` au headers katika ufafanuzi wa MCP. Kwa ukaguzi wa kuingia husoma majina ya kuingia kwa MCP katika taarifa za kuingia za kila akaunti (kwenye macOS kutoka Keychain), kamwe si tokeni, na tu kwa akaunti zinazopakia seva ambayo akaunti nyingine pia huipakia.
- Huonyesha kilichopo kwenye diski. Seva iliyozuiwa na managed settings bado huonekana; plugin iliyozimwa katika mipangilio huonyesha `(off)`.
- Iendeshe katika folda ya mradi ili kujumuisha folda ya `.claude` na `.mcp.json` ya mradi huo.

## Zinazohusiana

- [Kushiriki kwa vikundi](/features/share)
- [Seva za MCP zinazoshirikiwa](/features/mcp)
- [Kukagua na kurekebisha](/features/doctor)
