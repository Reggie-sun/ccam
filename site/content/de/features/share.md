---
title: Teilen zwischen Konten
description: Wie ccam share Einstellungen, CLAUDE.md, Erweiterungen und Verlauf mit Plan und Sicherungen in einen Gruppenordner verlinkt und wie off zurückkopiert.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

Jedes Konto hat seinen eigenen Ordner, aber das meiste darin muss nicht getrennt sein. `ccam share` schaltet das Teilen für ein Konto und ein Element nach dem anderen ein oder aus: Die Datei oder der Ordner des Kontos wird in einen Gruppenordner verschoben und durch einen symbolischen Link ersetzt, sodass jedes Konto der Gruppe dieselbe Kopie liest und schreibt. Verwende es nach `ccam add`, um zu ändern, was ein Konto teilt, oder um einige Konten in eine eigene Gruppe zu stecken.

## Verwendung

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Ohne Elementoptionen gibt `ccam share` den aktuellen Zustand aus:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

Mit Optionen gibt ccam einen Plan aus, fragt bei jeder Datei nach, die von der Kopie der Gruppe abweicht, und fragt `Proceed? [y/N]`, bevor es etwas ändert. Ein neues Konto teilt alles außer dem Gesprächsverlauf; [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing) beschreibt jedes Element aus Sicht der Nutzer.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | MCP-Server | unverändert |
| `--settings <on\|off\|group>` | `settings.json` | unverändert |
| `--memory <on\|off\|group>` | die globale `CLAUDE.md` | unverändert |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | unverändert |
| `--history <on\|off\|group>` | `projects/` und `history.jsonl` | unverändert |
| `--fresh` | Beim Ausschalten des Teilens leer beginnen, statt den Inhalt der Gruppe zu kopieren | aus |
| `--prefer <group\|mine\|merge>` | Wie Konflikte entschieden werden | nachfragen; mit `--yes` der Standard jedes Elements |
| `-y`, `--yes` | Ohne Nachfrage ausführen; Konflikte folgen `--prefer` | aus |
| `-h`, `--help` | Hilfe zu `share` | |
| `--dry-run` (global) | Wird beachtet: fragt bei Konflikten nach, gibt den Plan aus, ändert nichts | aus |

Der Wert `on` steht für die Gruppe `default`, `off` für lokal (nicht geteilt), und jeder andere Wert ist ein Gruppenname. Gruppennamen werden zu Ordnernamen unter `~/.local/share/ccam/groups/`, daher schlägt `ccam share side --history ../x` mit `invalid group name: "../x"` fehl. Das Konto wird über Namen oder Nummer angegeben, und nur die genannten Elemente ändern sich.

## So funktioniert es

Die beteiligten Dateien:

- `~/.config/ccam/config.toml`: die Tabelle `[profile.share]` jedes Kontos, ein Wert pro Element.
- `~/.local/share/ccam/profiles/<name>/`: das `CLAUDE_CONFIG_DIR` des Kontos.
- `~/.local/share/ccam/groups/<group>/`: die geteilten Kopien.

1. `internal/cli/share.go` liest die Konfiguration, wendet jeden `--<item>`-Wert mit `config.Share.Set` auf eine Kopie der Einstellungen des Kontos an und lehnt ab, wenn das Konto läuft.
2. `share.Build` vergleicht die aktuellen Einstellungen mit dem Ziel und listet die Schritte auf. Es liest nur von der Festplatte: `os.Lstat` und `os.Readlink` für jeden Pfad in `share.Layout`. Der Plan ist eine Liste von Operationen: in einen leeren Platz der Gruppe verschieben (übernehmen), die Dateien verschieben, die der Gruppe fehlen (Ordner zusammenführen), einen aufgelösten Konflikt in die Gruppe schreiben, sichern, Link entfernen, verlinken, zurückkopieren, Plugin-Pfade umschreiben und das Plugin-Register zusammenführen.
3. Konflikte werden entschieden, per Nachfrage oder durch `--prefer`.
4. `Plan.Describe` gibt den Plan aus, und ccam fragt nach.
5. Die Zieleinstellungen werden in `config.toml` geschrieben, bevor sich irgendetwas bewegt. Bricht `Plan.Apply` mittendrin ab, bauen `ccam doctor --fix` und ein erneuter Lauf ihren Plan auf das gespeicherte Ziel hin und schließen die Änderung ab.
6. `Plan.Apply` führt die Schritte der Reihe nach aus. Jeder Schritt verschiebt Daten in die Gruppe oder in den Sicherungsordner; nichts wird gelöscht.

`share.Layout` ist eine Positivliste der Pfade, die jedes Element abdeckt. Was nicht darauf steht, kann nie in eine Gruppe gelangen. `share.Never` nennt `.credentials.json` und `.claude.json`, und `Build` hält mit einem internen Fehler an, falls eine davon je in der Struktur auftaucht. Das Element `mcp` hat keine Pfade: Die Serverliste liegt in `.claude.json`, daher funktioniert das [Teilen von MCP-Servern](/features/mcp) stattdessen über die `mcp.json` der Gruppe, und `--mcp` ändert nur die Einstellung.

Das passiert mit einem Ordner, `projects/`, wenn `side` den Verlauf einschaltet und die Gruppe bereits Inhalt hat:

<!-- motion:share-on -->

Hat die Gruppe unter diesem Pfad noch nichts, teilt das Konto ihn als erstes, und sein Ordner wird unverändert und ohne Sicherung dorthin verschoben. Andernfalls verschiebt `mergeDir` jede Datei, die der Gruppe fehlt, in die Gruppe, und Dateien, die auf beiden Seiten existieren, bleiben im Ordner des Kontos, sodass die Version der Gruppe erhalten bleibt. Der Ordner des Kontos wandert dann nach `.ccam-backup/<time>/` im Kontoordner, und `os.Symlink` setzt einen Link an seine Stelle. Gesprächsdateien sind nach ihrer Sitzungs-ID benannt, daher lassen sich Verläufe aus mehreren Konten ohne Konflikte zusammenführen.

Das Verschieben verwendet `rename(2)` und greift über Dateisystemgrenzen hinweg auf Kopieren und Löschen zurück. Ein relativer Link im verschobenen Baum, der nach außerhalb zeigte, wird als absoluter Link neu geschrieben (`fsutil.PortableLink`), damit er vom neuen Ort aus weiter aufgelöst wird. Die `plugins/`-Register (`installed_plugins.json`, `known_marketplaces.json`) speichern absolute Pfade durch den Kontoordner; ccam richtet sie auf den Ordner der Gruppe aus und führt, wenn die Gruppe bereits ein Register hat, die Einträge des Kontos darin zusammen.

Das Ausschalten eines Elements kehrt den Link um:

<!-- motion:share-off -->

ccam entfernt den Link nur, wenn er in den Gruppenordner zeigt, und kopiert dann den aktuellen Inhalt der Gruppe mit `CopyTree` in das Konto. Mit `--fresh` entfällt die Kopie, und das Konto beginnt an diesem Pfad mit nichts. Die Gruppe behält ihren Inhalt für die anderen Konten. Wird ein Element von einer Gruppe in eine andere verschoben, wird es entlinkt und mit der neuen Gruppe verlinkt, während die alte Gruppe bleibt, wie sie ist; der Plan weist in einer `Note:`-Zeile darauf hin.

Eine einzelne Datei, die auf beiden Seiten mit unterschiedlichem Inhalt existiert, ist ein Konflikt:

<!-- motion:share-conflict -->

`share.NewConflict` listet die Wahlmöglichkeiten auf, die zur Datei passen. `history.jsonl` bietet `merge by line`, `keep group` und `use mine`, mit dem Zusammenführen an erster Stelle. Andere `.json`-Dateien bieten `keep group`, `use mine` und `merge JSON`. Alles andere bietet `keep group` und `use mine`. Die erste Wahl ist der Standard:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) führt Objekte rekursiv zusammen, verbindet Arrays ohne Duplikate und behält den Wert der Gruppe, wenn beide Seiten denselben Schlüssel unterschiedlich setzen, sodass andere Konten keine Verhaltensänderung bemerken. `merge by line` behält die Zeilen der Gruppe und fügt dann deine Zeilen hinzu, die ihr fehlen. Jede Wahl außer `keep group` sichert zuerst die Kopie der Gruppe in `.ccam-backup/<time>/group/` im Kontoordner und schreibt dann das Ergebnis mit `fsutil.AtomicWrite`. Deine eigene Datei kommt immer in die Sicherung.

`--prefer group`, `mine` oder `merge` entscheidet jeden Konflikt ohne Nachfrage; `merge` bedeutet `merge JSON` oder `merge by line`, je nachdem, was die Datei anbietet, und eine Datei, die keines davon anbietet, behält den Standard. `--yes` allein nimmt den Standard jeder Datei.

Bei jedem Start vergleicht `share.Check` jeden Pfad der Struktur mit der Einstellung und verwendet dabei nur `lstat` und `readlink`, damit es für den Startpfad günstig genug bleibt. Ein fehlender Link wird sofort angelegt, weil nichts verloren gehen kann. Eine echte Datei, wo ein Link hingehört, ein Link an die falsche Stelle oder ein veralteter Link hält den Start an:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` führt `share.Build` mit den gespeicherten Einstellungen auf beiden Seiten aus. Es leitet jeden Schritt aus dem ab, was auf der Festplatte liegt, behält bei Konflikten den Inhalt der Gruppe und sichert alles, was es ersetzt.

## Sicherheit und Sonderfälle

- **Laufende Konten.** ccam ändert das Teilen nicht, solange Claude Code als dieses Konto läuft (`account work is running; quit it first, then change its sharing`), weil die Sitzung weiter in die Dateien schreibt, die verschoben werden, und Zeilen, die sie in der Zwischenzeit anhängt, verloren gingen. `--dry-run` gibt den Plan trotzdem aus.
- **Sicherungen.** Jeder Plan erhält einen eigenen Ordner `.ccam-backup/<time>/` im Kontoordner, mit einer nummerierten Endung, wenn zwei Pläne in derselben Sekunde beginnen. ccam löscht sie nie.
- **Teilweise Fehlschläge.** Schlägt ein Schritt fehl, gibt ccam den Fehler und den Ort der Originale aus und rät dir, mit `ccam doctor --fix` abzuschließen. Das Ziel ist bereits gespeichert, daher setzt die Reparatur dort fort.
- **Deine eigenen Links.** Ist ein Pfad bereits ein Link, den du angelegt hast, zum Beispiel auf ein Dotfiles-Repository, und ist die Gruppe leer, kopiert ccam das Ziel des Links in die Gruppe und verschiebt den Link selbst in die Sicherung. Das Ziel wird nicht angefasst.
- **Anmeldedateien.** `.credentials.json` und `.claude.json` werden nie verlinkt, verschoben oder in eine Gruppe kopiert.
- **Zwei Konten in einem Ordner.** Konten, die den Verlauf teilen, schreiben in dieselben Dateien. Wenn ein zweites Konto derselben Verlaufsgruppe in einem Ordner startet, in dem bereits eines läuft, warnt ccam und fragt nach, bevor es fortfährt.

## Verwandte Themen

- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Geteilte MCP-Server](/features/mcp)
- [Eine vorhandene Einrichtung kopieren](/features/migrate)
- [Alle Befehle](/docs/commands)
