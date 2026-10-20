---
title: Pagsusuri ng mga extension
description: Ipinapakita ng ccam ext ang mga skill, command, plugin, MCP server at hook na nilo-load ng bawat account, saan galing ang bawat isa, at ano ang sira.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

Sinasagot ng `ccam ext` ang tanong na "bakit gumagana ang skill (o MCP server) na ito sa isang account pero hindi sa isa?". Inililista nito, para sa bawat account, ang mga extension na nilo-load ng Claude Code at kung saan galing ang bawat isa: isang share group, ang sariling folder ng account, o ang project na kinaroroonan mo. Sa ilalim ng table, inililista nito ang mga pumipigil sa isang extension na mag-load gaya ng inaasahan mo. Nagbabasa lang ito.

## Paggamit

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

Pinapangalanan ng bawat cell ang pinagmulan ng item na iyon sa account na iyon: ang share group (`default`), `account` kapag nasa sariling folder ito ng account, o `-` kapag wala nito ang account. Minamarkahan ng `!` ang item na may problema, na ipinapaliwanag sa ilalim ng table.

## Ano ang binabasa nito

| Uri | Saan | Ipinapakitang pinagmulan |
| --- | --- | --- |
| Mga skill | bawat folder sa `skills/` na may `SKILL.md` | ang group na nili-link-an ng `skills/`, o `account` |
| Mga command, agent | mga `.md` file sa ilalim ng `commands/` at `agents/`, mga subfolder bilang `folder/name` | ang group, o `account` |
| Mga plugin | `plugins/installed_plugins.json` | ang group, o `account` |
| Mga MCP server | ang `mcp.json` ng group (kapag shared ang MCP) at ang `.claude.json` ng account | ang group, o `account` |
| Mga hook | `hooks` sa `settings.json`, binibilang bawat event | ang group, o `account` |
| Ang project na ito | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` at `.mcp.json` sa kasalukuyang folder | nakalista nang isang beses, para sa bawat account |

## Mga problemang iniuulat nito

- **Mga sirang link**: isang link sa `skills/`, `commands/` o `agents/` na wala na ang target, o isang shared na folder (`skills`, `commands`, `agents`, `plugins`) na naka-link sa lugar na wala na (ayusin ito gamit ang `ccam doctor --fix`).
- **Skill folder na walang `SKILL.md`**: nilalaktawan ito ng Claude Code.
- **MCP server na magkaiba ang depinisyon sa dalawang lugar**: ang `mcp.json` ng group at ang `.claude.json` ng account, o isa sa mga iyon at ang `.mcp.json` ng project. Hindi problema ang parehong setting sa dalawang lugar (na isinusulat ng `mcp_strategy = "sync"`). Ilipat sa group ang sariling server ng isang account gamit ang `ccam mcp import <account> <server>`.
- **`installed_plugins.json` na hindi valid na JSON.**

## Kaligtasan at mga espesyal na kaso

- Walang binabago ang `ccam ext` at hindi ito gumagawa ng mga network request. Binabasa nito ang mga pangalan ng server at bilang ng hook, hindi kailanman ang mga value ng `env` o mga header sa mga MCP definition.
- Ipinapakita nito kung ano ang nasa disk. Lumalabas pa rin ang plugin na naka-install pero naka-off sa settings, o ang server na hinarang ng managed settings.
- Patakbuhin ito sa folder ng project para isama ang `.claude` folder at `.mcp.json` ng project na iyon.

## Kaugnay

- [Sharing sa pagitan ng mga account](/features/share)
- [Mga shared na MCP server](/features/mcp)
- [Pagsusuri at pagkukumpuni](/features/doctor)
