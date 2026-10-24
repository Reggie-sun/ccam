---
title: Duba extensions
description: ccam ext yana nuna skills, commands, plugins, sabobin MCP da hooks da kowane asusu ke lodawa, inda kowanne ya fito, da abin da ya lalace.
updated: 2026-10-23
order: 14
translation: machine
sourceHash: 45740e4efbef
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

Kowane ɗaki yana faɗin tushen wannan abu a wannan asusu: rukunin rabawa (`default`), `account` idan yana zaune a foldar asusun kansa, `local` ga saban MCP da asusun ke da shi don wannan foldar kawai, ko `-` idan asusun ba shi da shi. Plugin da aka shigar amma ba a kunna ba yana nuna `(off)` bayan tushensa. `!` yana alamta abu mai matsala, wanda aka bayyana a ƙarƙashin teburin; matsalolin da ba na wani asusu ba, kamar sabobin da `claude` na yau da kullum kaɗai ke da su, suna zuwa a ƙarshe ba tare da suna ba. `ccam ext --json` yana buga rahoto iri ɗaya (`accounts`, project ɗin, `problems`), tare da `remote` a kan sabobin MCP da Claude Code ke kaiwa ta network; idan babu asusai yana buga rahoton da `accounts: []`.

## Abin da yake karantawa

| Iri | Inda | Tushen da ake nunawa |
| --- | --- | --- |
| Skills | kowace folda a `skills/` mai `SKILL.md` | rukunin da `skills/` ke haɗawa da shi, ko `account` |
| Commands, agents | fayilolin `.md` a ƙarƙashin `commands/` da `agents/`, ƙananan foldoji a matsayin `folder/name` | rukunin, ko `account` |
| Plugins | `plugins/installed_plugins.json`, ana kunnawa ko kashewa ta `enabledPlugins` a saitunan asusun da na project ɗin; plugin da aka shigar don project ɗaya kawai yana bayyana ne kawai a wannan project ɗin | rukunin, ko `account`; `(off)` idan ba a kunna ba |
| Sabobin MCP | `mcp.json` na rukunin (idan ana raba MCP), `.claude.json` na asusun (sabobinsa na kansa, da waɗanda yake da su don wannan foldar kawai), da sabobin plugins ɗinsa, a matsayin `plugin:<plugin>:<server>` | rukunin, `account` ko `local` |
| Hooks | `hooks` a `settings.json`, ana ƙirga su ga kowane event | rukunin, ko `account` |
| Wannan project | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` da `.mcp.json` a foldar yanzu | ana lissafa su sau ɗaya; asusun da bai amince da foldar ba yana samun matsala game da sabobin `.mcp.json` nata |

## Matsalolin da yake bayarwa

- **Mahaɗai da suka lalace**: mahaɗi a `skills/`, `commands/` ko `agents/` wanda abin da yake nunawa ya ɓace, ko folda da aka raba (`skills`, `commands`, `agents`, `plugins`) da ke haɗawa da wurin da babu shi (gyara shi da `ccam doctor --fix`, wanda ke sake ƙirƙirar foldar rukuni da ta ɓace).
- **Foldar skill ba tare da `SKILL.md` ba**: Claude Code yana tsallake ta.
- **Saban MCP da aka ayyana daban a wurare biyu**: `mcp.json` na rukunin da `.claude.json` na asusun, ko ɗaya daga cikinsu da `.mcp.json` na project ɗin. Saituna iri ɗaya a wuraren biyu (waɗanda `mcp_strategy = "sync"` ke rubutawa) ba matsala ba ne. Matsar da saban asusu na kansa zuwa rukunin da `ccam mcp import <account> <server>`.
- **`installed_plugins.json` da ba JSON mai inganci ba**, plugin da `enabledPlugins` bai ambata ba, da plugin wanda foldarsa ta ɓace.
- **Sabobin da asusu ɗaya kawai ke lodawa**: saban da ke cikin `.claude.json` na wannan asusun kawai, wanda sauran asusan rukuninsa ba sa lodawa, da wanda ke cikin local scope na asusun don wannan foldar (abin da `claude mcp add` ke yi sai dai idan an ba shi `-s user`). Kowanne yana zuwa da umarnin `ccam mcp import` da ke raba shi. Sabobin da ke cikin `~/.claude.json` na `claude` na yau da kullum waɗanda babu asusun da ke lodawa suna zuwa da umarnin `ccam migrate --mcp`.
- **Saban project da ke jiran amincewa**: asusun bai amince da foldar ba tukuna, don haka Claude Code yana tambaya sau ɗaya, kuma har zuwa lokacin sabobin `.mcp.json` na project ɗin suna jira. Matsalar tana faɗin inda za a amince da su: `"enabledMcpjsonServers"` a `settings.json` na asusan rukunin saituna, ko na wannan asusun. Amincewar da aka ba babban worktree na git repository tana aiki a sauran worktrees ɗinsa.
- **Shiga da connectors**: saban MCP da ke buƙatar shiga a wannan asusun alhali wani asusu ya shiga cikinsa (shiga tana aiki ga asusu ɗaya kawai; shiga da `/mcp`), da connectors na claude.ai da wani asusu ya yi amfani da su amma wannan bai yi ba.

## Tsaro da yanayi na musamman

- `ccam ext` ba ya canza komai kuma ba ya yin buƙatun network. Yana karanta sunayen sabobi da adadin hooks, ba ya taɓa karanta ƙimomin `env` ko headers a cikin ma'anar MCP. Don duba shiga yana karanta sunayen shigar MCP a cikin bayanan shigar kowane asusu (a macOS daga Keychain), ba token ba, kuma ga asusan da ke lodawa saban da wani asusu ke lodawa ma kawai.
- Yana nuna abin da ke kan disk. Saban da managed settings suka toshe har yanzu yana bayyana; plugin da aka kashe a saituna yana nuna `(off)`.
- Gudanar da shi a foldar project ɗin don haɗa foldar `.claude` da `.mcp.json` na wannan project ɗin.

## Masu alaƙa

- [Rabawa tsakanin asusai](/features/share)
- [Sabobin MCP da aka raba](/features/mcp)
- [Dubawa da gyarawa](/features/doctor)
