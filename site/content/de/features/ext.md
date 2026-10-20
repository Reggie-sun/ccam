---
title: Erweiterungen untersuchen
description: "ccam ext zeigt die Skills, Befehle, Plugins, MCP-Server und Hooks, die jedes Konto lädt, woher sie jeweils kommen und was defekt ist."
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` beantwortet die Frage „Warum funktioniert dieser Skill (oder MCP-Server) in einem Konto und im anderen nicht?“. Es listet für jedes Konto die Erweiterungen auf, die Claude Code lädt, und woher jede kommt: eine geteilte Gruppe, der eigene Ordner des Kontos oder das Projekt, in dem du bist. Unter der Tabelle listet es auf, was verhindert, dass eine Erweiterung wie erwartet geladen wird. Es liest nur.

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

Eine Zelle nennt die Quelle dieses Elements in diesem Konto: die geteilte Gruppe (`default`), `account`, wenn es im eigenen Ordner des Kontos liegt, oder `-`, wenn das Konto es nicht hat. `!` markiert ein Element mit einem Problem, das unter der Tabelle erklärt wird.

## Was es liest

| Art | Ort | Angezeigte Quelle |
| --- | --- | --- |
| Skills | jeder Ordner in `skills/` mit einer `SKILL.md` | die Gruppe, auf die `skills/` verweist, oder `account` |
| Befehle, Agents | `.md`-Dateien unter `commands/` und `agents/`, Unterordner als `folder/name` | die Gruppe oder `account` |
| Plugins | `plugins/installed_plugins.json` | die Gruppe oder `account` |
| MCP-Server | das `mcp.json` der Gruppe (wenn MCP geteilt wird) und die `.claude.json` des Kontos | die Gruppe oder `account` |
| Hooks | `hooks` in `settings.json`, pro Ereignis gezählt | die Gruppe oder `account` |
| Dieses Projekt | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` und `.mcp.json` im aktuellen Ordner | einmal aufgeführt, für jedes Konto |

## Gemeldete Probleme

- **Defekte Links**: ein Link in `skills/`, `commands/` oder `agents/`, dessen Ziel fehlt, oder ein geteilter Ordner (`skills`, `commands`, `agents`, `plugins`), der auf einen fehlenden Ort verweist (repariere ihn mit `ccam doctor --fix`).
- **Ein Skill-Ordner ohne `SKILL.md`**: Claude Code überspringt ihn.
- **Ein MCP-Server, der an zwei Stellen unterschiedlich definiert ist**: im `mcp.json` der Gruppe und in der `.claude.json` des Kontos oder in einem davon und in der `.mcp.json` des Projekts. Dieselben Einstellungen an beiden Stellen (was `mcp_strategy = "sync"` schreibt) sind kein Problem. Verschiebe einen eigenen Server eines Kontos mit `ccam mcp import <account> <server>` in die Gruppe.
- **Eine `installed_plugins.json`, die kein gültiges JSON ist.**

## Sicherheit und Sonderfälle

- `ccam ext` ändert nichts und stellt keine Netzwerkanfragen. Es liest Servernamen und die Anzahl der Hooks, nie die Werte von `env` oder Headern in MCP-Definitionen.
- Es zeigt, was auf der Festplatte liegt. Ein installiertes, aber in den Einstellungen ausgeschaltetes Plugin oder ein durch verwaltete Einstellungen blockierter Server erscheint trotzdem.
- Führe es im Projektordner aus, um den `.claude`-Ordner und die `.mcp.json` dieses Projekts einzubeziehen.

## Verwandte Themen

- [Teilen zwischen Konten](/features/share)
- [Geteilte MCP-Server](/features/mcp)
- [Prüfen und reparieren](/features/doctor)
