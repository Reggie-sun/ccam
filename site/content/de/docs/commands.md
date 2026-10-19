---
title: Befehlsreferenz
description: Alle ccam-Befehle und -Optionen, von ccam add und ccam <name> bis zu ccam usage, ccam next und ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

Setze `--dry-run` vor einen Befehl, um zu sehen, was er tun würde, ohne es zu tun, zum Beispiel `ccam --dry-run work`.

<!-- commands -->

## Claude Code starten

`ccam <name|number> [arguments]` startet Claude Code mit diesem Konto. ccam liest die Argumente nach dem Namen nicht; es gibt sie in derselben Reihenfolge an `claude` weiter, und Standardeingabe, Exit-Code und Signale verhalten sich genau wie bei `claude` selbst. Ein `--` nach dem Namen ist optional.

`ccam` mit Argumenten, aber ohne Kontonamen, etwa `ccam -p "summarise"`, verwendet das Standardkonto. Lege es mit `ccam use <name>` fest.

Die Kurzbefehle `claude1`, `claude2` und so weiter sind Links auf ccam in `~/.local/bin`. `ccam shims install --names` erstellt außerdem `claude-<name>` für jedes Konto.

## Konten hinzufügen und entfernen

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` ohne Ordner sucht nach `~/.claude` und `~/.claude-*` und fragt bei jedem nach. Importierte Ordner werden dort registriert, wo sie liegen; ccam verschiebt oder kopiert nie eine Anmeldung.

## Teilen

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Jedes Element nimmt `on` (die Gruppe `default`), `off` oder einen Gruppennamen. Siehe [Teilen](/docs/sharing).

## MCP-Server

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Nutzung und Wechseln

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Wartung

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## Dateien

| Pfad | Inhalt |
| --- | --- |
| `~/.config/ccam/config.toml` | Konten und was sie teilen |
| `~/.local/share/ccam/profiles/<name>/` | Das `CLAUDE_CONFIG_DIR` jedes Kontos |
| `~/.local/share/ccam/groups/<group>/` | Geteilte Einstellungen, MCP-Server, Erweiterungen und Verlauf |
| `~/.local/state/ccam/` | Nutzungscache und Zeitpunkte der letzten Verwendung |
| `~/.local/bin/claude1`, `claude2`… | Links auf ccam |

Mit den Variablen `XDG_CONFIG_HOME`, `XDG_DATA_HOME` und `XDG_STATE_HOME` lassen sich diese Ordner verlegen.
