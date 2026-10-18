---
title: MCP-Server, Einstellungen und Verlauf teilen
description: Wie ccam MCP-Server, settings.json, CLAUDE.md, Agents, Skills und den Gesprächsverlauf zwischen Claude Code-Konten teilt und jede Anmeldung getrennt hält.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: e4b31e5098e3
---

Jedes Konto behält seine eigene Anmeldung. Alles andere kann Element für Element geteilt werden: ccam verlinkt die Datei oder den Ordner des Kontos mit einem geteilten Gruppenordner, sodass eine Änderung in einem Konto in allen Konten dieser Gruppe erscheint.

## Was geteilt werden kann

| Element | Option | Was geteilt wird |
| --- | --- | --- |
| MCP-Server | `--mcp` | Die `mcp.json` der Gruppe, die beim Start an Claude Code übergeben wird |
| Einstellungen | `--settings` | `settings.json`: Berechtigungen, Hooks, Modell, Statuszeile |
| Globales Gedächtnis | `--memory` | `CLAUDE.md` |
| Erweiterungen | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| Gesprächsverlauf | `--history` | `projects/` und `history.jsonl` |

Neue Konten teilen alles außer dem Gesprächsverlauf.

Die Anmeldedateien `.credentials.json` und `.claude.json` werden nie geteilt. ccam verlinkt nur Pfade aus einer festen Liste, und diese beiden stehen nicht darauf.

## Gruppen

`on` steht für die Gruppe `default`. Du kannst auch einen Namen angeben, sodass nur einige Konten ein Element teilen:

```sh
ccam share work --history team
ccam share side --history team
```

`work` und `side` teilen jetzt Gespräche miteinander, aber nicht mit deinen anderen Konten.

## Teilen einschalten

Wenn ein Konto anfängt, ein Element zu teilen, macht ccam Folgendes:

1. Es gibt einen Plan mit jeder Datei aus, die es verschieben, kopieren oder verlinken wird, und fragt, bevor es beginnt.
2. Es verschiebt die Dateien des Kontos in den Gruppenordner, wenn die Gruppe noch leer ist, oder fügt die Dateien hinzu, die die Gruppe noch nicht hat.
3. Wenn sich dieselbe Datei unterscheidet, fragt es, ob die Version der Gruppe behalten, deine verwendet oder das JSON zusammengeführt werden soll.
4. Es bewahrt dein Original in `.ccam-backup/<time>/` im Kontoordner auf.
5. Es ersetzt das Original durch einen Link auf die Gruppe.

Gesprächsdateien sind nach ihrer Sitzungs-ID benannt, daher lassen sich Verläufe aus mehreren Konten ohne Konflikte zusammenführen.

## Teilen ausschalten

```sh
ccam share work --history off
```

Der Link wird entfernt, und das Konto erhält eine eigene Kopie des aktuellen Inhalts der Gruppe. Füge `--fresh` hinzu, um stattdessen leer zu beginnen. Die Gruppe behält ihren Inhalt für die anderen Konten.

## MCP-Server

Claude Code speichert MCP-Server in `.claude.json`, zusammen mit der Anmeldung, daher kann diese Datei nicht geteilt werden. ccam hält geteilte Server in der `mcp.json` der Gruppe und startet Claude Code mit `--mcp-config=<group>/mcp.json`. Deine eigenen `--mcp-config`-Optionen gelten zusätzlich weiterhin.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Ein Server, den du innerhalb einer Sitzung mit `claude mcp add` hinzufügst, wird nur für dieses Konto gespeichert. Beim nächsten Start dieses Kontos nennt ccam den Server und den Befehl, der ihn in die Gruppe verschiebt:

```sh
ccam mcp import work gh
```

MCP-Server, die sich per OAuth anmelden, speichern ihre Tokens zusammen mit der Anmeldung des Kontos, autorisiere sie also in jedem Konto einmal.

Wenn die Server in `claude mcp list` und `/mcp` als gewöhnliche Benutzer-Server erscheinen sollen, setze `mcp_strategy = "sync"` in `~/.config/ccam/config.toml`. ccam schreibt sie dann vor dem Start in die `.claude.json` jedes Kontos, unter der eigenen Dateisperre von Claude Code und mit einer Sicherung.

## Zwei Konten im selben Ordner

Konten, die den Gesprächsverlauf teilen, schreiben in dieselben Dateien. Wenn du ein zweites Claude Code in einem Ordner startest, in dem bereits ein anderes Konto derselben Gruppe läuft, warnt ccam dich und fragt, ob du fortfahren willst.

## Wenn etwas kaputtgeht

```sh
ccam doctor
ccam doctor --fix
```

Wenn ein geteilter Pfad zu einer gewöhnlichen Datei geworden ist, zum Beispiel nachdem du eine Sicherung von Hand wiederhergestellt hast, hält ccam vor dem Start von Claude Code an und verweist auf `ccam doctor --fix`. Dieser Befehl führt die Datei mit der Gruppe zusammen und behält eine Sicherung.
