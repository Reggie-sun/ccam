---
title: Mehrere Claude Code-Konten auf einem Rechner nutzen
description: Mehrere Claude Code-Konten parallel mit CLAUDE_CONFIG_DIR betreiben, von Hand mit Shell-Aliasen oder automatisch mit ccam, und zwischen ihnen wechseln.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: f2982e4317a2
---

Gib jedem Konto einen eigenen Konfigurationsordner. Claude Code liest den Ordner aus der Umgebungsvariable `CLAUDE_CONFIG_DIR`, und jeder Ordner enthält seine eigene Anmeldung, Einstellungen und seinen eigenen Verlauf. So bleibt jedes Konto angemeldet, und du kannst sie gleichzeitig nutzen. Das kannst du von Hand mit einem Shell-Alias einrichten oder ccam erledigen lassen.

## Wie richte ich das von Hand ein?

Diesen Ansatz [beschreibt die Dokumentation von Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Füge für jedes zusätzliche Konto eine Zeile zu `~/.bashrc` oder `~/.zshrc` hinzu:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Öffne ein neues Terminal und führe `claude-work` aus. Claude Code bittet dich, dich anzumelden, weil der Ordner neu ist. Ab dann ist `claude` dein erstes Konto und `claude-work` das zweite. Lege für jedes Konto einen Alias mit eigenem Ordner an.

Unter Windows ohne WSL erledigt eine PowerShell-Funktion in deinem Profil (`notepad $PROFILE`) dasselbe:

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## Was fehlt beim Alias-Ansatz?

Für zwei Konten funktioniert er gut. Gut zu wissen:

- Jeder Ordner beginnt leer: keine Einstellungen, keine `CLAUDE.md`, keine MCP-Server, keine Agents oder Skills. Du kopierst oder verlinkst, was du brauchst, und hältst die Kopien selbst synchron.
- Aliase gibt es nur in interaktiven Shells. Skripte, Makefiles und andere Programme sehen sie nicht.
- MCP-Server liegen in `.claude.json`, neben der Anmeldung, und lassen sich deshalb nicht einfach zwischen Ordnern verlinken.
- Gespräche bleiben in dem Konto, in dem sie begonnen haben, du kannst also keines in einem anderen Konto fortsetzen.
- Nichts zeigt die Nutzung aller Konten auf einmal.

## Kann ich stattdessen mit /login das Konto wechseln?

Ja, wenn du nur ab und zu wechselst. `/login` in Claude Code [wechselt das Konto](https://code.claude.com/docs/en/quickstart) des Ordners, in dem du gerade bist. Es ist immer nur ein Konto angemeldet, du meldest dich also bei jedem Wechsel neu an und kannst nicht zwei Konten nebeneinander betreiben. Getrennte Ordner vermeiden beides.

## Wie richte ich das mit ccam ein?

ccam automatisiert dieselbe Idee und ergänzt die fehlenden Teile:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Jetzt starten `ccam personal` und `ccam work` die beiden Konten, ebenso `ccwork`, der Kurzbefehl, den du dem Arbeitskonto gegeben hast. Weil das echte Programme und keine Aliase sind, funktionieren sie auch in Skripten. Führe `ccam` ohne Argumente aus, um ein Konto aus einer Liste zu wählen, und `ccam use work`, um eines zum Standard zu machen.

Einstellungen, `CLAUDE.md`, MCP-Server und Erweiterungen werden von Anfang an geteilt, und du kannst den geteilten Gesprächsverlauf einschalten, um ein Gespräch mit `ccam next --continue` im anderen Konto fortzusetzen. `ccam usage` zeigt die 5-Stunden- und Wochennutzung aller Konten auf einmal.

Wenn du bereits `~/.claude` nutzt, führe `ccam import` aus, statt dieses Konto noch einmal hinzuzufügen, oder `ccam migrate`, um seine Einstellungen, MCP-Server und Erweiterungen in die geteilten Einstellungen zu kopieren.

## Woran erkenne ich, welches Konto läuft?

In Claude Code zeigt `/status` das Konto, bei dem du angemeldet bist. Mit ccam führst du einmal `ccam statusline install` aus, dann zeigt die Statuszeile ständig den Kontonamen und seine Nutzung, zum Beispiel `[work] · Opus · 5h 24% · 7d 41%`.

## Funktioniert das unter Windows und WSL?

Claude Code läuft [nativ unter Windows und in WSL](https://code.claude.com/docs/en/setup). `CLAUDE_CONFIG_DIR` funktioniert in beiden, unter nativem Windows reichen also der Alias oder die PowerShell-Funktion von oben. ccam läuft ebenfalls nativ unter Windows: Installiere es in PowerShell mit `irm https://ccam.dev/install.ps1 | iex` oder in WSL mit dem Linux-Befehl. Kurzbefehle werden zu `ccwork.exe`, und das Teilen von Dateien wie `settings.json` braucht den Entwicklermodus ([ccam installieren](/docs/install)).

## Fehlerbehebung

**Claude Code fragt jedes Mal nach der Anmeldung.** Beim ersten Start in einem neuen Ordner fragt es immer. Fragt es weiterhin, ist die Anmeldung abgelaufen oder wurde entfernt: Führe `ccam login work` aus (oder starte den Alias einmal und melde dich an).

**Ein MCP-Server erscheint nur in einem Konto.** Server, die in einer Sitzung mit `claude mcp add` hinzugefügt wurden, gehören zu diesem Konto. Mit `ccam mcp import work <server>` verschiebst du einen in alle Konten.

**Die Einstellungen unterscheiden sich zwischen Konten.** `ccam share work` zeigt, was das Konto teilt; `ccam share work --settings on` teilt seine Einstellungen, nachdem es gezeigt hat, was sich ändert. `ccam doctor --fix` repariert kaputte Links und legt von jeder geänderten Datei eine Sicherung an.

**Zwei Konten öffnen dasselbe Gespräch.** Arbeite nicht gleichzeitig in zwei Konten am selben Gespräch. ccam warnt dich, wenn du es versuchst.

## Was solltest du wählen?

Für zwei Konten, bei denen nichts geteilt wird, reicht ein Alias. Wenn du geteilte Einstellungen und MCP-Server, mehr Konten, die Nutzung auf einen Blick oder fortgesetzte Gespräche willst, nimmt dir ccam die Buchführung ab. Der [Vergleich mit Aliasen](/compare/aliases) stellt die Unterschiede nebeneinander, und [ccam und andere Tools zum Kontowechsel](/compare/account-switchers) behandelt Tools, die stattdessen eine einzige Anmeldung austauschen.

Ob du mehrere Konten nutzen darfst, hängt von den Bedingungen des jeweiligen Kontos ab. Prüfe sie, bevor du loslegst.
