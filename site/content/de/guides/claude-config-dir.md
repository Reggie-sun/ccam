---
title: CLAUDE_CONFIG_DIR erklärt
description: Was CLAUDE_CONFIG_DIR in Claude Code bewirkt, was im Ordner liegt, wo macOS die Anmeldung speichert und welche Dateien Konten teilen können.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` sagt Claude Code, welchen Ordner es statt `~/.claude` verwenden soll. Alles Kontospezifische liegt dort: die Anmeldung, Einstellungen, MCP-Server, Plugins und der Gesprächsverlauf. Wenn du die Variable auf verschiedene Ordner zeigen lässt, bekommst du getrennte Konten.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Was im Ordner liegt

Geprüft mit Claude Code 2.1.302:

| Pfad | Inhalt | Teilbar |
| --- | --- | --- |
| `.credentials.json` | Anmelde-Tokens (Linux), MCP-OAuth-Tokens | Nein |
| `.claude.json` | Kontodaten, Benutzer-MCP-Server, Zustand pro Projekt | Nein (MCP-Server separat) |
| `settings.json` | Berechtigungen, Hooks, Modell, Statuszeile, Plugins | Ja |
| `CLAUDE.md` | Globale Anweisungen | Ja |
| `agents/`, `commands/`, `skills/`, `plugins/` | Erweiterungen | Ja |
| `projects/` | Gespräche als `.jsonl`-Dateien, Gedächtnis pro Projekt | Ja |
| `history.jsonl` | Prompt-Verlauf | Ja |
| `sessions/`, `shell-snapshots/`, `cache/` und weitere | Laufzeitzustand | Nein |

## Zwei Details, über die man leicht stolpert

**Ohne `CLAUDE_CONFIG_DIR` liegt `.claude.json` nicht in `~/.claude`.** Claude Code speichert die Datei dann unter `~/.claude.json`, neben dem Ordner. Wenn du anfängst, `CLAUDE_CONFIG_DIR=~/.claude` zu verwenden, sucht Claude Code stattdessen nach `~/.claude/.claude.json` und findet dein Konto nicht. `ccam import` löst das, indem es die beiden verlinkt.

**Unter macOS liegt die Anmeldung im Schlüsselbund, nicht im Ordner.** Claude Code speichert sie unter dem Namen `Claude Code-credentials`, gefolgt von einem kurzen Hash des Ordnerpfads, wenn `CLAUDE_CONFIG_DIR` gesetzt ist. Getrennte Ordner bekommen daher getrennte Schlüsselbundeinträge. Wenn der Schlüsselbund das Schreiben verweigert, weicht Claude Code auf `.credentials.json` im Ordner aus.

## Zwischen Ordnern teilen

Symbolische Links funktionieren: Claude Code liest und schreibt über eine verlinkte `settings.json` oder einen verlinkten `projects/`-Ordner und behält den Link bei. MCP-Server sind die Ausnahme, weil sie sich `.claude.json` mit der Anmeldung teilen; übergib sie stattdessen mit `--mcp-config`.

ccam baut genau auf diesen Regeln auf. Unter [Teilen](/docs/sharing) steht, wie es Ordner verlinkt, und [die Anleitung zu mehreren Konten](/guides/multiple-claude-code-accounts) zeigt das Gesamtbild.
