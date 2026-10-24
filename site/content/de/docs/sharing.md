---
title: MCP-Server, Einstellungen und Verlauf teilen
description: Wie ccam-Gruppen MCP-Server, settings.json, CLAUDE.md, Skills und den Gesprächsverlauf zwischen Claude Code-Konten teilen; jede Anmeldung bleibt getrennt.
updated: 2026-10-23
order: 4
translation: machine
sourceHash: 27e0a0e5abbd
---

Jedes Konto behält seine eigene Anmeldung. Alles andere kann über **Gruppen** geteilt werden: Eine Gruppe nennt die Elemente, die sie teilt, und die Konten in ihr, und ccam verlinkt die Dateien jedes dieser Konten mit dem Ordner der Gruppe, sodass eine Änderung in einem Konto in allen erscheint.

## Was geteilt werden kann

| Element | Was geteilt wird |
| --- | --- |
| `mcp` | MCP-Server: die `mcp.json` der Gruppe, die beim Start an Claude Code übergeben wird |
| `settings` | `settings.json`: Berechtigungen, Hooks, Modell, Statuszeile |
| `memory` | die globale `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | Gesprächsverlauf: `projects/` und `history.jsonl` |

Dein erstes Konto legt die Gruppe `default` an, die alles außer dem Gesprächsverlauf teilt, und spätere Konten treten ihr bei.

Die Anmeldedateien `.credentials.json` und `.claude.json` werden nie geteilt. ccam verlinkt nur Pfade aus einer festen Liste, und diese beiden stehen nicht darauf.

Diese Seite behandelt die alltäglichen Aufgaben. Jede Option ist auf den Funktionsseiten beschrieben: [Teilen mit Gruppen](/features/share), [geteilte MCP-Server](/features/mcp) und [eine vorhandene Einrichtung kopieren](/features/migrate).

## Gruppen

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

Ein Konto kann in mehreren Gruppen sein, aber jedes Element kommt nur aus einer Gruppe. Um also den Verlauf zwischen `work` und `side` zu teilen, aber nicht mit deinen anderen Konten, halte den Verlauf aus `default` heraus und gib den beiden eine eigene Gruppe, wie oben. `ccam share work` zeigt, was ein Konto teilt und über welche Gruppe. Das [Dashboard](/features/dashboard) zeigt dasselbe auf seinem Tab Groups.

## Wenn ein Konto anfängt zu teilen

Wenn ein Konto anfängt, ein Element zu teilen (es tritt einer Gruppe bei, oder seine Gruppe nimmt das Element auf), macht ccam Folgendes:

1. Es gibt einen Plan mit jeder Datei aus, die es für jedes Konto verschieben, kopieren oder verlinken wird, und fragt, bevor es beginnt.
2. Es verschiebt die Dateien des Kontos in den Gruppenordner, wenn die Gruppe noch leer ist, oder fügt die Dateien hinzu, die die Gruppe noch nicht hat.
3. Wenn sich dieselbe Datei unterscheidet, fragt es, ob die Version der Gruppe behalten, deine verwendet oder das JSON zusammengeführt werden soll.
4. Es bewahrt dein Original in `.ccam-backup/<time>/` im Kontoordner auf.
5. Es ersetzt das Original durch einen Link auf die Gruppe.

Gesprächsdateien sind nach ihrer Sitzungs-ID benannt, daher lassen sich Verläufe aus mehreren Konten ohne Konflikte zusammenführen.

ccam ändert ein Konto nicht, solange Claude Code mit ihm läuft, egal wie Claude Code gestartet wurde (auch Hintergrundsitzungen), weil die Sitzung weiter in die Dateien schreibt, die verschoben werden. Beende sie zuerst. ccam prüft unmittelbar vor dem Speichern erneut, sodass auch eine Sitzung erkannt wird, die gestartet wurde, während der Plan wartete.

## Wenn ein Konto aufhört zu teilen

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

Der Link wird entfernt, und das Konto erhält eine eigene Kopie des aktuellen Inhalts der Gruppe; die MCP-Server der Gruppe kommen in seine eigene `.claude.json`, neben seine eigenen Server. Füge `--fresh` hinzu, um stattdessen leer zu beginnen; der Plan listet dann die MCP-Server auf, die nicht mehr geladen werden. Ein importierter Ordner wie `~/.claude` behält immer eine Kopie. Die Gruppe behält ihren Inhalt für die anderen Konten.

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
- der Gesprächsverlauf (`projects/` und `history.jsonl`), aber nur, wenn du `--history` angibst.

Der Quellordner wird nur gelesen: Nichts darin wird verschoben, verlinkt oder geändert, und einfaches `claude` funktioniert weiter wie bisher. Die Anmeldedateien `.credentials.json` und `.claude.json` werden nie in die Gruppe kopiert.

MCP-Server, die sich per OAuth anmelden, speichern ihre Tokens in der Anmeldung, nicht in der Serverliste. ccam kopiert diese Anmeldungen nicht: Eine kopierte Anmeldung enthält dasselbe Refresh-Token, und die erste Erneuerung meldet jeden anderen Inhaber bei diesem Server ab, `~/.claude` eingeschlossen. Der Plan und die Zusammenfassung nennen die Konten, die sich selbst mit `/mcp` bei den Servern anmelden.

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
Group default: personal
To share them:
  ccam group join default work
```

Am Ende listet ccam die Konten der Gruppe auf. Wenn die Gruppe eines der Elemente noch nicht teilt oder manche Konten nicht in ihr sind, gibt es die `ccam group`-Befehle aus, die das ändern.

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

Ein Server, den du innerhalb einer Sitzung mit `claude mcp add` hinzufügst, wird nur für dieses Konto gespeichert, und für den Ordner, in dem du warst (der lokale Bereich von Claude Code; `-s user` macht ihn zu einem Server des Kontos in jedem Ordner). Beim nächsten Start dieses Kontos nennt ccam den Server mit seinem Ordner und den Befehl, der ihn in die Gruppe verschiebt:

```sh
ccam mcp import work gh
```

Hat die Gruppe bereits einen anderen Server mit diesem Namen, überspringt ccam ihn und weist darauf hin. `--force` ersetzt den Server der Gruppe und behält die alte `mcp.json` im Ordner `.ccam-backup/` der Gruppe.

MCP-Server, die sich per OAuth anmelden, speichern ihre Tokens zusammen mit der Anmeldung des Kontos, autorisiere sie also in jedem Konto einmal mit `/mcp`. `ccam ext` und `ccam doctor` nennen die Server, bei denen ein Konto nicht angemeldet ist, ein anderes Konto aber schon.

Wenn die Server in `claude mcp list` und `/mcp` als gewöhnliche Benutzer-Server erscheinen sollen, setze `mcp_strategy = "sync"` am Anfang von `~/.config/ccam/config.toml`, vor dem ersten `[[group]]`, `[[profile]]` oder `[[folder]]`. ccam schreibt sie dann vor dem Start in die `.claude.json` jedes Kontos, unter der eigenen Dateisperre von Claude Code und mit einer Sicherung. Ein eigener Server des Kontos mit demselben Namen bleibt der des Kontos, und ccam weist darauf hin.

## Zwei Konten im selben Ordner

Konten, die den Gesprächsverlauf teilen, schreiben in dieselben Dateien. Wenn du ein zweites Claude Code in einem Ordner startest, in dem bereits ein anderes Konto derselben Gruppe läuft, warnt ccam dich und fragt, ob du fortfahren willst.

## Wenn etwas kaputtgeht

```sh
ccam doctor
ccam doctor --fix
```

Wenn ein geteilter Pfad zu einer gewöhnlichen Datei geworden ist, zum Beispiel nachdem du eine Sicherung von Hand wiederhergestellt hast, hält ccam vor dem Start von Claude Code an und verweist auf `ccam doctor --fix`, das sie wieder mit der Gruppe verlinkt: Die Datei wird zur Kopie der Gruppe, wenn die Gruppe keine hat, und andernfalls behält die Gruppe ihre eigene, und deine kommt in eine Sicherung, die `--fix` nennt. Eine Gruppenänderung, die auf halbem Weg stehen geblieben ist, schließt `ccam doctor --fix` ebenfalls ab.
