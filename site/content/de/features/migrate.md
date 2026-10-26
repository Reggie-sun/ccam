---
title: Eine vorhandene Einrichtung kopieren
description: Mit ccam migrate die Einstellungen, CLAUDE.md, Erweiterungen und MCP-Server eines vorhandenen ~/.claude in eine ccam-Gruppe kopieren.
updated: 2026-10-26
order: 6
translation: machine
sourceHash: f835fbe84ade
---

Wenn du Claude Code schon vor ccam benutzt hast, liegt deine Einrichtung in `~/.claude` und `~/.claude.json`. `ccam migrate` kopiert sie in eine Gruppe, sodass jedes Konto, das diese Elemente über die Gruppe teilt, mit der Einrichtung startet, die du schon hast. Der Quellordner wird nur gelesen, und einfaches `claude` funktioniert weiter wie bisher. Führe es einmal nach deinem ersten `ccam add` aus; wenn `~/.claude.json` MCP-Server enthält, nennt das erste `ccam add` sie in seinem Tipp.

## Verwendung

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

In einem Terminal, ohne Elementoptionen und ohne `--yes`, zeigt ccam zuerst eine Checkliste:

```console
 Copy into group default

 From /home/me/.claude
 Choose what to copy (Space to toggle)
 ❯ [✔] MCP servers
   [✔] settings.json
   [✔] global CLAUDE.md
   [✔] skills, agents, commands, plugins
   [ ] conversation history

 space to toggle · enter to show the plan · esc to cancel
```

Die Leertaste (oder `x`) schaltet das Element unter dem Cursor um, die Pfeiltasten oder `j`/`k` bewegen die Auswahl, Enter zeigt den Plan, und Esc, `q` oder Ctrl+C gibt `Cancelled.` aus. Elemente, für die die Quelle nichts enthält, sind mit `not found` markiert. Wenn du alles abwählst, gibt ccam `Nothing selected.` aus und hört auf.

Dann gibt ccam den Plan aus und fragt nach, bevor es kopiert:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

Eine Gruppe wirkt sich nur auf ihre eigenen Konten aus und nur auf die Elemente, die sie teilt, daher listet die Zusammenfassung die Konten der Gruppe auf. Wenn die Gruppe eines der kopierten Elemente nicht teilt oder manche Konten nicht in ihr sind, gibt ccam die Befehle `ccam group set <group> --add <items>` und `ccam group join <group> <accounts>` aus, die das ändern.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `[folder]` | Der Claude Code-Ordner, aus dem kopiert wird | `~/.claude` |
| `--settings` | `settings.json` kopieren | aus |
| `--memory` | Die globale `CLAUDE.md` kopieren | aus |
| `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` kopieren | aus |
| `--mcp` | Benutzer-MCP-Server kopieren (`mcpServers` in `.claude.json`) | aus |
| `--history` | Den Gesprächsverlauf kopieren (`projects/`, `history.jsonl`) | aus |
| `-g`, `--group <group>` | Die Gruppe, in die kopiert wird | `default` |
| `-y`, `--yes` | Ohne Nachfrage kopieren; bei einem Konflikt bleibt die Kopie der Gruppe | aus |
| `--dry-run` (global) | Den Plan ausgeben und vor jeder Frage und jedem Kopieren anhalten | aus |

Gibst du in einem Terminal keine Elementoptionen an, wählst du die Elemente in der Checkliste aus. Außerhalb eines Terminals oder mit `--yes` kopiert ccam alles außer dem Gesprächsverlauf. Der Verlauf ist standardmäßig ausgenommen, weil `projects/` groß sein kann.

`--group on` bedeutet `default`. `--group off` und `--group local` werden mit `--group takes a group name, not "off"` abgelehnt. `ccam migrate` unterstützt `--dry-run`, und die Checkliste erscheint auch bei einem Probelauf, damit du dir eine Auswahl vorab ansehen kannst.

## So funktioniert es

ccam erstellt einen vollständigen Plan, bevor es irgendetwas kopiert. Es wählt die Elemente aus (`internal/cli/migrate.go`), liest dann die Quelle und die Gruppe und erstellt den Plan (`migrate.Build`). Bis `Plan.Apply` wird nichts geschrieben, und das erst nach deiner Bestätigung.

1. **Die Quelle finden.** ccam lehnt einen Pfad ab, der fehlt, kein Ordner ist oder keinen bekannten Claude Code-Eintrag enthält (`checkSource`). Nur für MCP-Server ist ein leeres oder fehlendes `~/.claude` in Ordnung, solange `~/.claude.json` Server enthält. ccam findet die `.claude.json`, die zum Ordner gehört (`migrate.ClaudeJSON`): darin, oder bei einem `~/.claude` ohne eigene Datei die `~/.claude.json` daneben. Dort speichert Claude Code sie, wenn `CLAUDE_CONFIG_DIR` nicht gesetzt ist.
2. **Elemente Pfaden zuordnen.** Jeder Pfad stammt aus `share.Layout`, derselben Positivliste, die Gruppen verwenden. Die MCP-Server werden mit `mcp.LocalServers` aus `.claude.json` gelesen, das den Scanner `profile.TopLevelValue` verwendet, statt die ganze Datei zu parsen.
3. **Kopieren, nie verschieben.** Dateien gehen mit `fsutil.AtomicWrite` oder `fsutil.CopyFile` nach `~/.local/share/ccam/groups/<group>/`; dabei wird eine temporäre Datei neben das Ziel geschrieben und umbenannt, sodass ein unterbrochener Kopiervorgang nie eine unvollständige Datei hinterlässt. MCP-Server gehen unter deren Sperre in die `mcp.json` der Gruppe: Ein Server, den ein anderer ccam-Prozess in der Zwischenzeit hinzufügt, bleibt erhalten, und die für [Remote Control](/guides/remote-control) angelegten Kopien werden danach aktualisiert.

Jedes Element ordnet einen Quellpfad einem Gruppenpfad zu:

<!-- motion:migrate-items -->

Die beiden Verlaufszeilen sind gestrichelt, weil sie nur mit `--history` kopiert werden. `.credentials.json` und `.claude.json` stehen in `share.Never` und erreichen die Gruppe nie: ccam hält mit einem internen Fehler an, falls eine davon je in der Struktur auftaucht (`Build`), und der Durchlauf durch den Ordner überspringt sie anhand des Namens.

ccam gibt jedem Pfad einen Zustand, und den beschreibt jede Zeile des Plans (`buildEntry`):

<!-- motion:migrate-states -->

- **Fehlend**: nicht in der Quelle, oder eine Datei, wo ein Ordner erwartet wird. Wird pro Element einmal als `not found, skipped` aufgeführt.
- **Verlinkt**: Der Quellpfad führt bereits auf die Kopie der Gruppe, zum Beispiel wenn die Quelle ein ccam-Profilordner ist.
- **Gleich**: Die Gruppe hat bereits identischen Inhalt.
- **Neu**: Der Gruppe fehlt er, daher wird er vollständig kopiert.
- **Zusammenführen**: ein Ordner auf beiden Seiten. ccam fügt nur die Dateien hinzu, die der Gruppe fehlen (`copyMissing`); abweichende Dateien bleiben so, wie die Gruppe sie hat, und der Plan meldet `group copies kept`. Bei `plugins/` werden die Register `installed_plugins.json` und `known_marketplaces.json` zusammengeführt, und die absoluten Plugin-Pfade werden auf den Ordner der Gruppe umgeschrieben.
- **Konflikt**: eine Datei auf beiden Seiten mit unterschiedlichem Inhalt. Du wählst pro Datei: `keep group`, `use mine`, dazu `merge JSON` bei `.json`-Dateien (eine tiefe Zusammenführung, bei der die Gruppe bei widersprüchlichen Werten gewinnt) oder `merge by line` bei `history.jsonl`. Ein MCP-Server mit abweichender Definition ist ebenfalls ein Konflikt, mit `keep group` oder `use mine`. Die erste aufgeführte Wahl ist der Standard, wenn du Enter drückst; `--yes` behält immer die Kopie der Gruppe.

Die Checkliste selbst ist `tui.MigratePicker`, gebaut auf der Komponente `checklist`, die auch `ccam add` verwendet; ccam markiert fehlende Elemente anhand von lstat-Aufrufen und der Liste der MCP-Server (`migrate.Found`).

### MCP-Anmeldungen

MCP-Server, die sich per OAuth anmelden, speichern ihre Tokens in der Anmeldedatei, `.credentials.json` unter `mcpOAuth`, nicht in der Serverliste. ccam kopiert sie nicht. Eine kopierte Anmeldung enthält dasselbe Refresh-Token wie das Original, und der Server wechselt dieses Token bei der ersten Erneuerung: Jeder andere Inhaber, auch das Quell-`~/.claude`, ist dann bei diesem Server abgemeldet. Deshalb nennt die Zeile `sign-ins` des Plans für jeden kopierten Server, bei dem die Quelle angemeldet ist, die Konten der Gruppe, die sich selbst mit `/mcp` anmelden müssen, und die Zusammenfassung wiederholt das. Ein Eintrag ohne Token ist keine Anmeldung und wird nicht aufgeführt, ebenso wenig ein Konto, das bereits eine eigene Anmeldung hat, oder ein Server, dessen Definition die Gruppe behält. Es wird keine Anmeldedatei geschrieben.

## Sicherheit und Sonderfälle

- Die Quelle wird nur gelesen. Keine Datei darin wird verschoben, verlinkt oder geändert, und die Anmeldung der Quelle wird nie geschrieben, auch nicht unter macOS, wo sie womöglich aus dem Schlüsselbund gelesen wird.
- Bevor eine Gruppendatei ersetzt wird, durch `use mine`, eine Zusammenführung oder neue MCP-Server in einer vorhandenen `mcp.json`, wird die Kopie der Gruppe in `<group>/.ccam-backup/<time>/` gesichert. Der Ordnername ist pro Lauf eindeutig, sodass zwei Läufe in derselben Sekunde ihre Sicherungen nicht gegenseitig überschreiben. Die Zusammenfassung gibt den Pfad aus.
- Mit `--yes` behält jeder Konflikt die Kopie der Gruppe, und ccam sagt, wie du stattdessen wählen kannst: Führe denselben Befehl ohne `--yes` aus.
- Würde sich nichts ändern, gibt ccam `Nothing to copy.` aus; behält jeder Konflikt die Kopie der Gruppe, `Nothing to copy: the group keeps its own copies.`
- Schlägt ein Kopiervorgang mittendrin fehl, hält ccam mit `The source folder was not changed. Group files replaced so far are saved in` an, gefolgt vom Sicherungsordner. Ein erneuter Aufruf kopiert die Dateien, die nicht vollständig kopiert wurden.
- Wenn mehrere Gruppen MCP teilen, lauten die Befehle, die `ccam doctor`, `ccam ext` und `ccam mcp list` vorschlagen, `ccam migrate --mcp --group <a|b>`.
- Eine fehlende Quelle endet mit `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Verwandte Themen

- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Prüfen und reparieren](/features/doctor)
- [Mehrere Claude Code-Konten auf einem Rechner nutzen](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR erklärt](/guides/claude-config-dir)
