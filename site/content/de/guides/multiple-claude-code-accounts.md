---
title: Zwei Claude Code-Konten auf einem Computer nutzen
description: Mehrere Claude Code-Konten parallel mit CLAUDE_CONFIG_DIR betreiben, von Hand mit Shell-Aliasen oder automatisch mit ccam.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e2fbee9f4f91
---

Gib jedem Konto einen eigenen Konfigurationsordner. Claude Code liest den Ordner aus der Umgebungsvariable `CLAUDE_CONFIG_DIR`, und jeder Ordner enthält seine eigene Anmeldung, Einstellungen und seinen eigenen Verlauf. Das kannst du von Hand mit einem Shell-Alias einrichten oder ccam erledigen lassen.

## Von Hand, mit einem Alias

Diesen Ansatz beschreibt die Dokumentation von Claude Code. Füge für jedes zusätzliche Konto eine Zeile zu `~/.bashrc` oder `~/.zshrc` hinzu:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Öffne ein neues Terminal und führe `claude-work` aus. Claude Code bittet dich, dich anzumelden, weil der Ordner neu ist. Ab dann ist `claude` dein erstes Konto und `claude-work` das zweite.

Für zwei Konten funktioniert das gut. Gut zu wissen:

- Jeder Ordner beginnt leer: keine Einstellungen, keine `CLAUDE.md`, keine MCP-Server, keine Agents oder Skills. Du kopierst oder verlinkst, was du brauchst, und hältst die Kopien selbst synchron.
- Aliase gibt es nur in interaktiven Shells. Skripte, Makefiles und andere Programme sehen sie nicht.
- MCP-Server liegen in `.claude.json`, neben der Anmeldung, und lassen sich deshalb nicht einfach zwischen Ordnern verlinken.
- Gespräche bleiben in dem Konto, in dem sie begonnen haben, du kannst also keines in einem anderen Konto fortsetzen.
- Nichts zeigt die Nutzung aller Konten auf einmal.

## Mit ccam

ccam automatisiert dieselbe Idee und ergänzt die fehlenden Teile:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

Jetzt starten `claude1` und `claude2` oder `ccam personal` und `ccam work` die beiden Konten. Weil das echte Programme und keine Aliase sind, funktionieren sie auch in Skripten. Einstellungen, `CLAUDE.md`, MCP-Server und Erweiterungen werden von Anfang an geteilt, und du kannst den geteilten Gesprächsverlauf einschalten, um ein Gespräch mit `ccam next --continue` im anderen Konto fortzusetzen.

Wenn du bereits `~/.claude` nutzt, führe `ccam import` aus, statt dieses Konto noch einmal hinzuzufügen.

## Was du wählen solltest

Für zwei Konten, bei denen nichts geteilt wird, reicht ein Alias. Wenn du geteilte Einstellungen und MCP-Server, mehr Konten, die Nutzung auf einen Blick oder fortgesetzte Gespräche willst, nimmt dir ccam die Buchführung ab. Der [Vergleich](/compare/aliases) stellt die Unterschiede nebeneinander.
