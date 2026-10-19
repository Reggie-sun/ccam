---
title: MCP-Server, Einstellungen und Verlauf teilen
description: Wie ccam MCP-Server, settings.json, CLAUDE.md, Agents, Skills und den Gesprächsverlauf zwischen Claude Code-Konten teilt und jede Anmeldung getrennt hält.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: f266f8bfd55c
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

ccam ändert das Teilen nicht, solange Claude Code mit dem Konto läuft, weil die Sitzung weiter in die Dateien schreibt, die verschoben werden. Beende sie zuerst.

## Teilen ausschalten

```sh
ccam share work --history off
```

Der Link wird entfernt, und das Konto erhält eine eigene Kopie des aktuellen Inhalts der Gruppe. Füge `--fresh` hinzu, um stattdessen leer zu beginnen. Die Gruppe behält ihren Inhalt für die anderen Konten.

## Eine vorhandene Einrichtung kopieren

Wenn du Claude Code schon vor ccam benutzt hast, liegen deine Einstellungen in `~/.claude`. `ccam migrate` kopiert sie in eine Gruppe, sodass jedes Konto, das diese Elemente teilt, sie bekommt:

```sh
ccam migrate
```

In einem Terminal zeigt ccam zuerst eine Checkliste der Elemente. Alles außer dem Gesprächsverlauf ist angehakt; drücke die Leertaste, um ein Element umzuschalten, und Enter, um den Plan zu sehen. Elemente, für die der Ordner nichts enthält, sind mit `not found` markiert.

Was kopiert wird:

- `settings.json` und die globale `CLAUDE.md`;
- die Erweiterungen `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` und `rules/`;
- deine Benutzer-MCP-Server, also der Eintrag `mcpServers` in `~/.claude.json`, in die `mcp.json` der Gruppe;
- die gespeicherten Anmeldungen (OAuth) der kopierten MCP-Server, in die eigene Anmeldedatei jedes Kontos;
- der Gesprächsverlauf (`projects/` und `history.jsonl`), aber nur, wenn du `--history` angibst.

Der Quellordner wird nur gelesen: Nichts darin wird verschoben, verlinkt oder geändert, und einfaches `claude` funktioniert weiter wie bisher. Die Anmeldedateien `.credentials.json` und `.claude.json` werden nie in die Gruppe kopiert.

MCP-Server, die sich per OAuth anmelden, speichern ihre Tokens in der Anmeldung, nicht in der Serverliste. Damit du dich nicht erneut bei ihnen anmelden musst, kopiert ccam diese Anmeldungen in die Anmeldedatei jedes Kontos, das die Server über die Gruppe bekommt, sofern das Konto nicht schon eigene hat. Konten, die gerade laufen, noch nicht angemeldet sind, mit `ccam import` an Ort und Stelle registriert wurden oder ihre Anmeldung im macOS-Schlüsselbund speichern, überspringt es und nennt sie; melde dich dort mit `/mcp` bei den Servern an. Ein Konto, das MCP-Server erst später teilt, bekommt die Anmeldungen, wenn du `ccam migrate` erneut ausführst.

ccam gibt den Plan aus und fragt nach, bevor es etwas kopiert:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

Am Ende listet ccam die Konten auf, die jedes Element über die Gruppe teilen. Wenn einige Konten das nicht tun, gibt es die `ccam share`-Befehle aus, die das Teilen für sie einschalten.

Wenn die Gruppe bereits eine abweichende Kopie einer Datei hat, fragt ccam wie beim Einschalten des Teilens: die Kopie der Gruppe behalten oder deine verwenden. `settings.json` und andere JSON-Dateien lassen sich auch zusammenführen, `history.jsonl` Zeile für Zeile. Die vorherigen Kopien der Gruppe werden in `.ccam-backup/<time>/` im Gruppenordner gesichert.

| Option | Wirkung |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Nur die genannten Elemente kopieren, ohne Checkliste; ohne sie und außerhalb eines Terminals alles außer dem Gesprächsverlauf |
| `-g`, `--group <group>` | In diese Gruppe statt in `default` kopieren |
| `-y`, `--yes` | Ohne Nachfrage kopieren; wenn eine Datei abweicht, die Kopie der Gruppe behalten |

Setze `--dry-run` vor den Befehl, etwa `ccam --dry-run migrate`, um den Plan zu sehen, ohne etwas zu kopieren. Um aus einem anderen Ordner zu kopieren, gib ihn an: `ccam migrate ~/.claude-work --group work`.

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

Hat die Gruppe bereits einen anderen Server mit diesem Namen, überspringt ccam ihn und weist darauf hin. `--force` ersetzt den Server der Gruppe und behält die alte `mcp.json` im Ordner `.ccam-backup/` der Gruppe.

MCP-Server, die sich per OAuth anmelden, speichern ihre Tokens zusammen mit der Anmeldung des Kontos, autorisiere sie also in jedem Konto einmal oder kopiere die Anmeldungen mit `ccam migrate` aus `~/.claude`.

Wenn die Server in `claude mcp list` und `/mcp` als gewöhnliche Benutzer-Server erscheinen sollen, setze `mcp_strategy = "sync"` in `~/.config/ccam/config.toml`. ccam schreibt sie dann vor dem Start in die `.claude.json` jedes Kontos, unter der eigenen Dateisperre von Claude Code und mit einer Sicherung.

## Zwei Konten im selben Ordner

Konten, die den Gesprächsverlauf teilen, schreiben in dieselben Dateien. Wenn du ein zweites Claude Code in einem Ordner startest, in dem bereits ein anderes Konto derselben Gruppe läuft, warnt ccam dich und fragt, ob du fortfahren willst.

## Wenn etwas kaputtgeht

```sh
ccam doctor
ccam doctor --fix
```

Wenn ein geteilter Pfad zu einer gewöhnlichen Datei geworden ist, zum Beispiel nachdem du eine Sicherung von Hand wiederhergestellt hast, hält ccam vor dem Start von Claude Code an und verweist auf `ccam doctor --fix`. Dieser Befehl führt die Datei mit der Gruppe zusammen und behält eine Sicherung.
