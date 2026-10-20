---
title: Befehlsreferenz
metaTitle: "ccam-Befehlsreferenz"
description: Alle ccam-Befehle und -Optionen, von ccam add und ccam <name> bis zu ccam usage, ccam next und ccam doctor.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 5d3275878535
---

Setze `--dry-run` vor einen Befehl, um zu sehen, was er tun würde, ohne es zu tun, zum Beispiel `ccam --dry-run work`. Das funktioniert beim Starten von Claude Code (`ccam <name>`, `ccam run`, `ccam next`) sowie mit `share`, `migrate`, `shortcut`, `update` und `doctor --fix`. Andere Befehle, die etwas ändern, lehnen die Option ab, statt zu laufen.

<!-- commands -->

## Claude Code starten

`ccam <name|number> [arguments]` startet Claude Code mit diesem Konto. ccam liest die Argumente nach dem Namen nicht; es gibt sie in derselben Reihenfolge an `claude` weiter, und Standardeingabe, Exit-Code und Signale verhalten sich genau wie bei `claude` selbst. Ein `--` nach dem Namen ist optional.

`ccam` mit Argumenten, aber ohne Kontonamen, etwa `ccam -p "summarise"`, verwendet das Konto, das für den aktuellen Ordner festgelegt ist, sonst das Standardkonto. Lege sie mit `ccam folder <name> [folder]` und `ccam use <name>` fest.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` vor dem Kontonamen ist für `-p`/`--print`-Durchläufe gedacht. Wenn das Konto sein Nutzungslimit erreicht, führt ccam denselben Befehl erneut mit dem angemeldeten Konto mit dem meisten Spielraum aus und gibt nur die Ausgabe des Durchlaufs aus, der fertig wurde. Siehe [Konten nach Ordner](/features/folder) und [Am Limit erneut ausführen](/features/next#run-again-at-the-limit).

Ein Kurzbefehl wie `ccwork` startet sein Konto auf dieselbe Weise: `ccwork --resume` macht dasselbe wie `ccam work --resume`. Siehe [Kurzbefehle](/docs/commands#shortcuts).

## Konten hinzufügen und entfernen

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` ohne Ordner sucht nach `~/.claude` und `~/.claude-*` und fragt bei jedem nach. Importierte Ordner werden dort registriert, wo sie liegen; ccam verschiebt oder kopiert nie eine Anmeldung.

`ccam rename`, `ccam remove` und `ccam share` ändern ein Konto nicht, solange Claude Code mit diesem Konto läuft. Beende zuerst die Sitzung.

`ccam list --json` gibt jedes Konto mit Slot-Nummer, Kurzbefehl, E-Mail-Adresse und Ordner aus. Das Objekt `share` ordnet jedem Element (`mcp`, `settings`, `memory`, `extensions`, `history`) einen Gruppennamen oder `local` zu.

## Kurzbefehle

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Ein Kurzbefehl ist ein Befehl, den du selbst benennst und der Claude Code mit einem bestimmten Konto startet. Wie `ccam work` gibt er jedes Argument unverändert an Claude Code weiter, sodass `ccwork --resume` und `ccwork -p "hi"` wie erwartet funktionieren. Er ist ein Link auf ccam in `~/.local/bin`.

Kurzbefehle sind optional. Der Assistent von `ccam add` fragt, ob einer angelegt werden soll, `ccam add <name> --shortcut <command>` legt einen ohne Nachfrage an, und `ccam list` zeigt den Kurzbefehl jedes Kontos oder `-`, wenn es keinen hat.

Namen bestehen aus den Zeichen a-z, 0-9, `.`, `-` und `_`. Ein Name darf kein ccam-Befehl sein, nicht der Kurzbefehl eines anderen Kontos und nicht der Name eines anderen Programms in deinem `PATH`.

`ccam rename` behält den Namen des Kurzbefehls bei, und `ccam remove` löscht den Kurzbefehl zusammen mit dem Konto. `ccam shims install` erstellt die Links aller Konten neu, zum Beispiel nachdem du ccam verschoben hast.

## Teilen

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Jedes Element nimmt `on` (die Gruppe `default`), `off` oder einen Gruppennamen. Siehe [Teilen](/docs/sharing).

## Eine vorhandene Einrichtung kopieren

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` kopiert `settings.json`, `CLAUDE.md`, Erweiterungen und Benutzer-MCP-Server aus einem Claude Code-Ordner in eine Gruppe, sodass jedes Konto, das diese Elemente teilt, sie bekommt. Der Ordner wird nur gelesen. Die Optionen `--settings`, `--memory`, `--extensions`, `--mcp` und `--history` kopieren nur die genannten Elemente. Ohne sie kannst du in einem Terminal in einer Checkliste auswählen, in der alles außer dem Gesprächsverlauf angehakt ist; anderswo wird alles außer dem Gesprächsverlauf kopiert. `-g` oder `--group` wählt die Gruppe, und `-y` oder `--yes` kopiert ohne Nachfrage und behält die Kopie der Gruppe, wenn eine Datei abweicht. Außerdem kopiert es die gespeicherten Anmeldungen der Server in die eigene Anmeldedatei jedes Kontos. Siehe [Eine vorhandene Einrichtung kopieren](/docs/sharing#copy-an-existing-setup).

## MCP-Server

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` verschiebt die Server, die nur dieses Konto hat, in seine Gruppe. Ein Server, dessen Namen die Gruppe bereits für eine andere Definition verwendet, wird übersprungen; mit `--force` ersetzt du die Kopie der Gruppe, die ccam vorher sichert.

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
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` prüft die SHA-256-Prüfsumme des Releases, bevor es ccam ersetzt, und `--version v0.2.0` installiert ein bestimmtes Release. Ein aus dem `.deb`-Paket installiertes ccam aktualisierst du stattdessen mit apt. Siehe [Aktualisieren](/docs/install#update).

## Dateien

| Pfad | Inhalt |
| --- | --- |
| `~/.config/ccam/config.toml` | Konten und was sie teilen |
| `~/.local/share/ccam/profiles/<name>/` | Das `CLAUDE_CONFIG_DIR` jedes Kontos |
| `~/.local/share/ccam/groups/<group>/` | Geteilte Einstellungen, MCP-Server, Erweiterungen und Verlauf |
| `~/.local/state/ccam/` | Nutzungscache und Zeitpunkte der letzten Verwendung |
| `~/.local/bin/<shortcut>` | Kurzbefehle wie `ccwork`, Links auf ccam |

Mit den Variablen `XDG_CONFIG_HOME`, `XDG_DATA_HOME` und `XDG_STATE_HOME` lassen sich diese Ordner verlegen. Unter Windows ist `~` gleich `%USERPROFILE%`, und Kurzbefehle sind `<shortcut>.exe`.
