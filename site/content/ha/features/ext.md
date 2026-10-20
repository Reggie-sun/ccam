---
title: Duba extensions
description: ccam ext yana nuna skills, commands, plugins, sabobin MCP da hooks da kowane asusu ke lodawa, inda kowanne ya fito, da abin da ya lalace.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` yana amsa tambayar "me ya sa wannan skill (ko saban MCP) yake aiki a wani asusu amma ba a ɗayan ba?". Yana lissafa, ga kowane asusu, extensions da Claude Code ke lodawa da inda kowanne ya fito: rukunin rabawa, foldar asusun kansa, ko project ɗin da kake ciki. A ƙarƙashin teburin yana lissafa abin da ke hana extension lodawa kamar yadda kake tsammani. Yana karantawa ne kawai.

## Amfani

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

Kowane ɗaki yana faɗin tushen wannan abu a wannan asusu: rukunin rabawa (`default`), `account` idan yana zaune a foldar asusun kansa, ko `-` idan asusun ba shi da shi. `!` yana alamta abu mai matsala, wanda aka bayyana a ƙarƙashin teburin.

## Abin da yake karantawa

| Iri | Inda | Tushen da ake nunawa |
| --- | --- | --- |
| Skills | kowace folda a `skills/` mai `SKILL.md` | rukunin da `skills/` ke haɗawa da shi, ko `account` |
| Commands, agents | fayilolin `.md` a ƙarƙashin `commands/` da `agents/`, ƙananan foldoji a matsayin `folder/name` | rukunin, ko `account` |
| Plugins | `plugins/installed_plugins.json` | rukunin, ko `account` |
| Sabobin MCP | `mcp.json` na rukunin (idan ana raba MCP) da `.claude.json` na asusun | rukunin, ko `account` |
| Hooks | `hooks` a `settings.json`, ana ƙirga su ga kowane event | rukunin, ko `account` |
| Wannan project | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` da `.mcp.json` a foldar yanzu | ana lissafa su sau ɗaya, ga kowane asusu |

## Matsalolin da yake bayarwa

- **Mahaɗai da suka lalace**: mahaɗi a `skills/`, `commands/` ko `agents/` wanda abin da yake nunawa ya ɓace, ko folda da aka raba (`skills`, `commands`, `agents`, `plugins`) da ke haɗawa da wurin da babu shi (gyara shi da `ccam doctor --fix`).
- **Foldar skill ba tare da `SKILL.md` ba**: Claude Code yana tsallake ta.
- **Saban MCP da aka ayyana daban a wurare biyu**: `mcp.json` na rukunin da `.claude.json` na asusun, ko ɗaya daga cikinsu da `.mcp.json` na project ɗin. Saituna iri ɗaya a wuraren biyu (waɗanda `mcp_strategy = "sync"` ke rubutawa) ba matsala ba ne. Matsar da saban asusu na kansa zuwa rukunin da `ccam mcp import <account> <server>`.
- **`installed_plugins.json` da ba JSON mai inganci ba.**

## Tsaro da yanayi na musamman

- `ccam ext` ba ya canza komai kuma ba ya yin buƙatun network. Yana karanta sunayen sabobi da adadin hooks, ba ya taɓa karanta ƙimomin `env` ko headers a cikin ma'anar MCP.
- Yana nuna abin da ke kan disk. Plugin da aka shigar amma aka kashe a saituna, ko saban da managed settings suka toshe, har yanzu yana bayyana.
- Gudanar da shi a foldar project ɗin don haɗa foldar `.claude` da `.mcp.json` na wannan project ɗin.

## Masu alaƙa

- [Rabawa tsakanin asusai](/features/share)
- [Sabobin MCP da aka raba](/features/mcp)
- [Dubawa da gyarawa](/features/doctor)
