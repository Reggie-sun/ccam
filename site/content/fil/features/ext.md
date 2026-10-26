---
title: Pagsusuri ng mga extension
description: Ipinapakita ng ccam ext ang mga skill, command, plugin, MCP server at hook na nilo-load ng bawat account, saan galing ang bawat isa, at ano ang sira.
updated: 2026-10-26
order: 14
translation: machine
sourceHash: "7ca357a7f38a"
---

Sinasagot ng `ccam ext` ang tanong na "Bakit gumagana ang skill (o MCP server) na ito sa isang account pero hindi sa isa?" Inililista nito, para sa bawat account, ang mga extension na nilo-load ng Claude Code at kung saan galing ang bawat isa: isang share group, ang sariling folder ng account, o ang project na kinaroroonan mo. Sa ilalim ng table, inililista nito ang mga pumipigil sa isang extension na mag-load gaya ng inaasahan mo. Wala itong binabago.

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

Pinapangalanan ng bawat cell ang pinagmulan ng item na iyon sa account na iyon: ang share group (`default`), `account` kapag nasa sariling folder ito ng account, `local` para sa MCP server na mayroon ang account para sa folder na ito lang, o `-` kapag wala nito ang account. Nagpapakita ng `(off)` pagkatapos ng pinagmulan nito ang plugin na naka-install pero hindi naka-on. Minamarkahan ng `!` ang item na may problema, na ipinapaliwanag sa ilalim ng table; huling lumalabas, nang walang pangalan, ang mga problemang hindi pag-aari ng anumang account, gaya ng mga server na ang karaniwang `claude` lang ang mayroon. Ipi-print ng `ccam ext --json` ang parehong ulat (`accounts`, ang project, `problems`), na may `remote` sa mga MCP server na inaabot ng Claude Code sa network; kapag walang account, ipi-print nito ang ulat na may `accounts: []`.

## Ano ang binabasa nito

| Uri | Saan | Ipinapakitang pinagmulan |
| --- | --- | --- |
| Mga skill | bawat folder sa `skills/` na may `SKILL.md` | ang group na nili-link-an ng `skills/`, o `account` |
| Mga command, agent | mga `.md` file sa ilalim ng `commands/` at `agents/`, mga subfolder bilang `folder/name` | ang group, o `account` |
| Mga plugin | `plugins/installed_plugins.json`, na ino-on o ino-off ng `enabledPlugins` sa settings ng account at ng project; lumalabas lang sa project na iyon ang plugin na naka-install para sa iisang project lang | ang group, o `account`; `(off)` kapag hindi naka-on |
| Mga MCP server | ang `mcp.json` ng group (kapag shared ang MCP), ang `.claude.json` ng account (ang sarili nitong mga server, at ang mga mayroon ito para sa folder na ito lang), at ang mga server ng mga plugin nito, bilang `plugin:<plugin>:<server>` | ang group, `account` o `local` |
| Mga hook | `hooks` sa `settings.json`, binibilang bawat event | ang group, o `account` |
| Ang project na ito | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` at `.mcp.json` sa kasalukuyang folder | nakalista nang isang beses; nagkakaroon ng problema para sa mga `.mcp.json` server nito ang account na hindi pa nagtitiwala sa folder |

## Mga problemang iniuulat nito

- **Mga sirang link**: isang link sa `skills/`, `commands/` o `agents/` na wala na ang target, o isang shared na folder (`skills`, `commands`, `agents`, `plugins`) na naka-link sa lugar na wala na (ayusin ito gamit ang `ccam doctor --fix`, na muling gumagawa ng nawawalang group folder).
- **Skill folder na walang `SKILL.md`**: nilalaktawan ito ng Claude Code.
- **MCP server na magkaiba ang depinisyon sa dalawang lugar**: ang `mcp.json` ng group at ang `.claude.json` ng account, o isa sa mga iyon at ang `.mcp.json` ng project. Hindi problema ang magkaparehong depinisyon sa dalawang lugar (na isinusulat ng `mcp_strategy = "sync"`). Ilipat sa group ang sariling server ng isang account gamit ang `ccam mcp import <account> <server>`.
- **`installed_plugins.json` na hindi valid na JSON**, plugin na hindi pinapangalanan ng `enabledPlugins`, at plugin na nawawala ang folder.
- **Mga server na iisang account lang ang naglo-load**: server na nasa `.claude.json` lang ng account na ito, na hindi nilo-load ng ibang account ng group nito, at server na nasa local scope ng account para sa folder na ito (ang ginagawa ng `claude mcp add` maliban kung bibigyan ng `-s user`). May kasamang `ccam mcp import` command na nagshe-share nito ang bawat isa. May kasamang `ccam migrate --mcp` command ang mga server sa `~/.claude.json` ng karaniwang `claude` na hindi nilo-load ng anumang account.
- **Project server na naghihintay ng pag-apruba**: hindi pa nagtitiwala ang account sa folder, kaya magtatanong nang isang beses ang Claude Code, at hanggang doon ay naghihintay ang mga server sa `.mcp.json` ng project. Sinasabi ng problema kung saan sila aaprubahan: `"enabledMcpjsonServers"` sa `settings.json` ng mga account ng settings group, o ng account na ito. Ang tiwalang ibinigay sa main worktree ng isang git repository ay nalalapat din sa iba pa nitong worktree.
- **Mga sign-in at connector**: MCP server na nangangailangan ng sign-in sa account na ito habang naka-sign in dito ang ibang account (para sa iisang account lang ang isang sign-in; mag-sign in gamit ang `/mcp`), at mga claude.ai connector na nagamit na ng ibang account pero hindi pa ng account na ito.

## Kaligtasan at mga espesyal na kaso

- Walang binabago ang `ccam ext` at hindi ito gumagawa ng mga network request. Binabasa nito ang mga pangalan ng server at bilang ng hook, hindi kailanman ang mga value ng `env` o mga header sa mga MCP definition. Para sa pagsusuri ng sign-in, binabasa nito ang mga pangalan ng mga MCP sign-in sa login ng bawat account (sa macOS, mula sa Keychain), hindi kailanman ang mga token, at para lang sa mga account na naglo-load ng server na nilo-load din ng ibang account.
- Ipinapakita nito kung ano ang nasa disk. Lumalabas pa rin ang server na hinarang ng managed settings; nagpapakita ng `(off)` ang plugin na naka-off sa settings.
- Patakbuhin ito sa folder ng project para isama ang `.claude` folder at `.mcp.json` ng project na iyon.

## Kaugnay

- [Pag-share gamit ang mga group](/features/share)
- [Mga shared na MCP server](/features/mcp)
- [Pagsusuri at pagkukumpuni](/features/doctor)
