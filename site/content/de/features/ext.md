---
title: Erweiterungen untersuchen
description: "ccam ext zeigt die Skills, Befehle, Plugins, MCP-Server und Hooks, die jedes Konto lädt, woher sie jeweils kommen und was defekt ist."
updated: 2026-10-26
order: 14
translation: machine
sourceHash: 7ca357a7f38a
---

`ccam ext` beantwortet die Frage „Warum funktioniert dieser Skill (oder MCP-Server) in einem Konto, aber im anderen nicht?“ Es listet für jedes Konto die Erweiterungen auf, die Claude Code lädt, und woher jede kommt: eine geteilte Gruppe, der eigene Ordner des Kontos oder das Projekt, in dem du bist. Unter der Tabelle listet es auf, was verhindert, dass eine Erweiterung wie erwartet geladen wird. Es ändert nichts.

## Verwendung

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

Eine Zelle nennt die Quelle dieses Elements in diesem Konto: die geteilte Gruppe (`default`), `account`, wenn es im eigenen Ordner des Kontos liegt, `local` für einen MCP-Server, den das Konto nur für diesen Ordner hat, oder `-`, wenn das Konto es nicht hat. Ein Plugin, das installiert, aber nicht eingeschaltet ist, zeigt `(off)` nach seiner Quelle. `!` markiert ein Element mit einem Problem, das unter der Tabelle erklärt wird; Probleme, die zu keinem Konto gehören, etwa Server, die nur ein einfaches `claude` hat, stehen ohne Namen am Ende. `ccam ext --json` gibt denselben Bericht aus (`accounts`, das Projekt, `problems`), mit `remote` bei MCP-Servern, die Claude Code über das Netzwerk erreicht; ohne Konten gibt es den Bericht mit `accounts: []` aus.

## Was es liest

| Art | Ort | Angezeigte Quelle |
| --- | --- | --- |
| Skills | jeder Ordner in `skills/` mit einer `SKILL.md` | die Gruppe, auf die `skills/` verweist, oder `account` |
| Befehle, Agents | `.md`-Dateien unter `commands/` und `agents/`, Unterordner als `folder/name` | die Gruppe oder `account` |
| Plugins | `plugins/installed_plugins.json`, ein- oder ausgeschaltet durch `enabledPlugins` in den Einstellungen des Kontos und des Projekts; ein Plugin, das nur für ein Projekt installiert ist, erscheint nur in diesem Projekt | die Gruppe oder `account`; `(off)`, wenn nicht eingeschaltet |
| MCP-Server | das `mcp.json` der Gruppe (wenn MCP geteilt wird), die `.claude.json` des Kontos (seine eigenen Server und die, die es nur für diesen Ordner hat) und die Server seiner Plugins, als `plugin:<plugin>:<server>` | die Gruppe, `account` oder `local` |
| Hooks | `hooks` in `settings.json`, pro Ereignis gezählt | die Gruppe oder `account` |
| Dieses Projekt | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` und `.mcp.json` im aktuellen Ordner | einmal aufgeführt; ein Konto, das dem Ordner nicht vertraut hat, bekommt ein Problem für die Server seiner `.mcp.json` |

## Gemeldete Probleme

- **Defekte Links**: ein Link in `skills/`, `commands/` oder `agents/`, dessen Ziel fehlt, oder ein geteilter Ordner (`skills`, `commands`, `agents`, `plugins`), der auf einen fehlenden Ort verweist (repariere ihn mit `ccam doctor --fix`, das einen fehlenden Gruppenordner neu anlegt).
- **Ein Skill-Ordner ohne `SKILL.md`**: Claude Code überspringt ihn.
- **Ein MCP-Server, der an zwei Stellen unterschiedlich definiert ist**: im `mcp.json` der Gruppe und in der `.claude.json` des Kontos oder in einem davon und in der `.mcp.json` des Projekts. Identische Definitionen an beiden Stellen (was `mcp_strategy = "sync"` schreibt) sind kein Problem. Verschiebe einen eigenen Server eines Kontos mit `ccam mcp import <account> <server>` in die Gruppe.
- **Eine `installed_plugins.json`, die kein gültiges JSON ist**, ein Plugin, das `enabledPlugins` nicht nennt, und ein Plugin, dessen Ordner fehlt.
- **Server, die nur ein Konto lädt**: ein Server nur in der `.claude.json` dieses Kontos, den die anderen Konten seiner Gruppe nicht laden, und einer im lokalen Bereich des Kontos für diesen Ordner (was `claude mcp add` ohne `-s user` tut). Zu jedem gibt es den Befehl `ccam mcp import`, der ihn teilt. Zu Servern in der `~/.claude.json` des einfachen `claude`, die kein Konto lädt, gibt es den Befehl `ccam migrate --mcp`.
- **Ein Projektserver, der auf Freigabe wartet**: Das Konto hat dem Ordner noch nicht vertraut, also fragt Claude Code einmal nach, und bis dahin warten die Server der `.mcp.json` des Projekts. Das Problem sagt, wo du sie freigibst: `"enabledMcpjsonServers"` in der `settings.json` der Konten der Einstellungsgruppe oder dieses Kontos. Vertrauen, das dem Haupt-Worktree eines Git-Repositorys gegeben wurde, gilt auch in seinen anderen Worktrees.
- **Anmeldungen und Connectors**: ein MCP-Server, der in diesem Konto eine Anmeldung braucht, während ein anderes Konto bei ihm angemeldet ist (eine Anmeldung gilt nur für ein Konto; melde dich mit `/mcp` an), und claude.ai-Connectors, die ein anderes Konto verwendet hat und dieses nicht.

## Sicherheit und Sonderfälle

- `ccam ext` ändert nichts und stellt keine Netzwerkanfragen. Es liest Servernamen und die Anzahl der Hooks, nie die Werte von `env` oder Headern in MCP-Definitionen. Für die Prüfung der Anmeldungen liest es die Namen der MCP-Anmeldungen in der Anmeldung jedes Kontos (unter macOS aus dem Schlüsselbund), nie die Tokens, und nur für Konten, die einen Server laden, den auch ein anderes Konto lädt.
- Es zeigt, was auf der Festplatte liegt. Ein durch verwaltete Einstellungen blockierter Server erscheint trotzdem; ein in den Einstellungen ausgeschaltetes Plugin zeigt `(off)`.
- Führe es im Projektordner aus, um den `.claude`-Ordner und die `.mcp.json` dieses Projekts einzubeziehen.

## Verwandte Themen

- [Teilen mit Gruppen](/features/share)
- [Geteilte MCP-Server](/features/mcp)
- [Prüfen und reparieren](/features/doctor)
