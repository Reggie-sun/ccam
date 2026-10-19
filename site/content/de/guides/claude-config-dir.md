---
title: CLAUDE_CONFIG_DIR erklärt
description: Was CLAUDE_CONFIG_DIR in Claude Code bewirkt, was im Ordner liegt, wo macOS die Anmeldung speichert und welche Dateien Konten teilen können.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: ef9317030e42
---

`CLAUDE_CONFIG_DIR` sagt Claude Code, welchen Ordner es statt `~/.claude` verwenden soll. Alles Kontospezifische liegt dort: die Anmeldung, Einstellungen, MCP-Server, Plugins und der Gesprächsverlauf. Wenn du die Variable auf verschiedene Ordner zeigen lässt, bekommst du getrennte Konten.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Wird CLAUDE_CONFIG_DIR offiziell unterstützt?

Ja. Claude Code führt die Variable unter seinen [Umgebungsvariablen](https://code.claude.com/docs/en/env-vars) als Weg auf, den Konfigurationsordner zu überschreiben, „nützlich, um mehrere Konten nebeneinander zu betreiben“, und die [Anleitung zur Authentifizierung](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) nutzt sie für die Anmeldung mit mehreren Konten. Du kannst sie in deiner Shell oder in den Benutzereinstellungen setzen, aber nicht in den Einstellungen eines Projekts.

## Was liegt im Ordner?

Geprüft mit Claude Code 2.1.302:

| Pfad | Inhalt | Teilbar |
| --- | --- | --- |
| `.credentials.json` | Anmelde-Tokens (Linux, Windows), MCP-OAuth-Tokens | Nein |
| `.claude.json` | Kontodaten, Benutzer-MCP-Server, Zustand pro Projekt | Nein (MCP-Server separat) |
| `settings.json` | Berechtigungen, Hooks, Modell, Statuszeile, Plugins | Ja |
| `CLAUDE.md` | Globale Anweisungen | Ja |
| `agents/`, `commands/`, `skills/`, `plugins/` | Erweiterungen | Ja |
| `projects/` | Gespräche als `.jsonl`-Dateien, Gedächtnis pro Projekt | Ja |
| `history.jsonl` | Prompt-Verlauf | Ja |
| `sessions/`, `shell-snapshots/`, `cache/` und weitere | Laufzeitzustand | Nein |

## Worüber stolpert man leicht?

**Ohne `CLAUDE_CONFIG_DIR` liegt `.claude.json` nicht in `~/.claude`.** Claude Code speichert die Datei dann unter `~/.claude.json`, neben dem Ordner. Wenn du anfängst, `CLAUDE_CONFIG_DIR=~/.claude` zu verwenden, sucht Claude Code stattdessen nach `~/.claude/.claude.json` und findet dein Konto nicht. `ccam import` löst das, indem es die beiden verlinkt.

**Unter macOS liegt die Anmeldung im Schlüsselbund, nicht im Ordner.** Claude Code speichert sie unter dem Namen `Claude Code-credentials`, gefolgt von einem kurzen Hash des Ordnerpfads, wenn `CLAUDE_CONFIG_DIR` gesetzt ist. Getrennte Ordner bekommen daher getrennte Schlüsselbundeinträge. Wenn der Schlüsselbund das Schreiben verweigert, weicht Claude Code auf `.credentials.json` im Ordner aus, wie [seine Dokumentation](https://code.claude.com/docs/en/authentication#credential-management) beschreibt.

**Nicht jede Anmeldung liegt im Ordner.** Eine Claude Console-Anmeldung ohne API-Schlüssel wird außerhalb gespeichert, getrennte Ordner halten zwei solche Anmeldungen also nicht auseinander. claude.ai-Anmeldungen und API-Schlüssel werden pro Ordner gespeichert.

## Wo liegt der Ordner unter Windows?

Standardmäßig unter `%USERPROFILE%\.claude`, mit der Anmeldung in `.credentials.json` darin. Setze `CLAUDE_CONFIG_DIR` für jedes zusätzliche Konto auf einen anderen Ordner, genau wie unter Linux und macOS.

## Können Ordner Dateien teilen?

Symbolische Links funktionieren: Claude Code liest und schreibt über eine verlinkte `settings.json` oder einen verlinkten `projects/`-Ordner und behält den Link bei. MCP-Server sind die Ausnahme, weil sie sich `.claude.json` mit der Anmeldung teilen; übergib sie stattdessen mit `--mcp-config`.

ccam baut genau auf diesen Regeln auf. Unter [Teilen](/docs/sharing) steht, wie es Ordner verlinkt, und [die Anleitung zu mehreren Konten](/guides/multiple-claude-code-accounts) zeigt das Gesamtbild.
