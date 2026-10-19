---
title: ccam im Vergleich mit Shell-Aliasen
description: Shell-Aliase mit CLAUDE_CONFIG_DIR und ccam im direkten Vergleich, um mehrere Claude Code-Konten zu betreiben.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: f751716e46c9
---

Ein Shell-Alias pro Konto ist der einfachste Weg, mehrere Claude Code-Konten zu betreiben. ccam macht darunter dasselbe und ergänzt Teilen, Wechseln und das Verfolgen der Nutzung.

| | Shell-Alias | ccam |
| --- | --- | --- |
| Getrennte Anmeldungen | Ja | Ja |
| Einrichtung pro Konto | `~/.bashrc` bearbeiten, anmelden | `ccam add` |
| Funktioniert in Skripten und anderen Programmen | Nein, Aliase gibt es nur in der Shell | Ja, `ccam work` und Kurzbefehle wie `ccwork` sind echte Befehle |
| Geteilte Einstellungen, `CLAUDE.md`, Erweiterungen | Von Hand kopieren oder verlinken | Pro Element, mit Sicherungen |
| Geteilte MCP-Server | Von Hand, pro Konto | Einmal `ccam mcp add` |
| Gespräch in einem anderen Konto fortsetzen | Nein | `ccam next --continue` |
| Nutzung aller Konten | Nein | `ccam usage` |
| Konto mit dem meisten Spielraum wählen | Nein | `ccam next --auto` |
| Kontoname in der Statuszeile | Nein | `ccam statusline install` |
| Defekte Links reparieren | Von Hand | `ccam doctor --fix` |
| Zusätzliche Software | Keine | Eine statische Binärdatei |

Nimm Aliase, wenn du zwei Konten hast und nichts zwischen ihnen teilen willst. Nimm ccam, wenn sich die Konten wie eine einzige Einrichtung anfühlen sollen. Du kannst mit Aliasen anfangen und später zu ccam wechseln: `ccam import` registriert vorhandene `~/.claude-*`-Ordner dort, wo sie liegen.

[Zwei Konten auf einem Computer nutzen](/guides/multiple-claude-code-accounts) erklärt beide Ansätze Schritt für Schritt.
