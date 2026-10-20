---
title: ccam im Vergleich mit Shell-Aliasen
description: Shell-Aliase mit CLAUDE_CONFIG_DIR und ccam im direkten Vergleich, um mehrere Claude Code-Konten zu betreiben.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: ea24bb091caf
---

Ein Shell-Alias pro Konto ist der einfachste Weg, mehrere Claude Code-Konten zu betreiben. ccam macht darunter dasselbe und ergänzt Teilen, Wechseln und das Verfolgen der Nutzung.

## Die Einrichtung mit Aliasen

Füge für jedes weitere Konto eine Zeile in `~/.bashrc` oder `~/.zshrc` ein, öffne ein neues Terminal und führe den Alias einmal aus, um dich anzumelden:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` bleibt dein erstes Konto und `claude-work` startet das zweite. Das ist die Methode, die [die Dokumentation von Claude Code beschreibt](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts).

## Die Einrichtung mit ccam

Installiere ccam und füge das Konto hinzu, hier mit einem eigenen Kurzbefehl:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` und `ccwork` starten das Konto, und beide funktionieren auch in Skripten.

## Im direkten Vergleich

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

[Mehrere Konten auf einem Computer nutzen](/guides/multiple-claude-code-accounts) erklärt beide Ansätze Schritt für Schritt, und [ccam und andere Tools zum Kontowechsel](/compare/account-switchers) vergleicht ccam mit Tools, die eine einzige Anmeldung austauschen.
