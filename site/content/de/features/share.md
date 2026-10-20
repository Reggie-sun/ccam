---
title: Teilen mit Gruppen
description: Eine ccam-Gruppe teilt gewählte Elemente über einen Ordner zwischen ihren Konten; ccam group verschiebt die Daten mit Plan, Sicherungen und Konfliktwahl.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 0c0338362ba2
---

Jedes Konto hat seinen eigenen Ordner, aber das meiste darin muss nicht getrennt sein. Eine **Gruppe** legt fest, was geteilt wird und von wem: Jedes Konto der Gruppe verlinkt diese Elemente mit einem Ordner, `~/.local/share/ccam/groups/<group>/`, sodass alle dieselbe Kopie lesen und schreiben. Anmeldetokens werden nie geteilt.

Ein Konto kann in mehreren Gruppen sein, aber jedes Element kommt nur aus einer Gruppe. Eine typische Einrichtung ist eine Gruppe für Einstellungen und Erweiterungen und eine zweite für den Gesprächsverlauf zwischen zwei Konten:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

Das erste Konto, das du hinzufügst, legt `default` an, und spätere Konten treten ihr bei, sofern du nichts anderes angibst.

## Verwendung

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

Jeder Befehl, der Daten verschiebt, gibt einen Plan aus, fragt bei jeder Datei nach, die von der Kopie der Gruppe abweicht, und fragt `Proceed? [y/N]`, bevor er etwas ändert:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261020-153012
  link     projects, history.jsonl: profiles/work → groups/default
Proceed? [y/N]
```

Dieselben Änderungen gibt es auf dem Tab **Groups** des [Dashboards](/features/dashboard), und `s` auf einem Konto wählt seine Gruppen.

## Elemente

| Element | Was es umfasst |
| --- | --- |
| `mcp` | MCP-Server, über die `mcp.json` der Gruppe ([Geteilte MCP-Server](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | die globale `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` und `history.jsonl`, damit ein anderes Konto mit `--continue` weitermachen kann |

Ein Element, das keine Gruppe des Kontos teilt, bleibt im eigenen Ordner des Kontos.

## Befehle und Optionen

| Befehl | Wirkung |
| --- | --- |
| `ccam group` oder `ccam group list` | die Gruppen auflisten |
| `ccam group show <group>` | was die Gruppe teilt, ihre Konten, ihr Ordner und ihre MCP-Server |
| `ccam group add <group> --items <items> [--members <accounts>]` | eine Gruppe anlegen; die Elemente ihrer Konten wandern in sie |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | ändern, was die Gruppe teilt, für jedes Konto in ihr |
| `ccam group join <group> <account>...` | Konten in die Gruppe aufnehmen |
| `ccam group leave <group> <account>...` | Konten herausnehmen; jedes behält eine Kopie |
| `ccam group rename <group> <new-name>` | nur in der Konfiguration umbenennen |
| `ccam group remove <group> [--leave]` | eine Gruppe entfernen |

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: das Konto, dessen Kopie als Erste in die Gruppe kommt | das Standardkonto, dann nach Nummer |
| `--replace` | `join`: zuerst die Gruppen verlassen, die dieselben Elemente teilen | aus |
| `--fresh` | Konten, die ein Element nicht mehr teilen, beginnen ohne es statt mit einer Kopie | aus |
| `--leave` | `remove`: zuerst die Konten der Gruppe herausnehmen | nötig, wenn sie Konten hat |
| `--prefer <group\|mine\|merge>` | wie Konflikte entschieden werden | nachfragen; mit `--yes` der Standard jeder Datei |
| `-y`, `--yes` | ohne Nachfrage ausführen; Konflikte folgen `--prefer` | aus |
| `--dry-run` (global) | den Plan ausgeben, nichts ändern | aus |

Konten werden über Namen oder Nummer angegeben, Elemente als kommagetrennte Liste. `ccam add` nimmt `--group <group>` (mehrfach verwendbar) oder `--no-group`.

## So funktioniert es

Die beteiligten Dateien:

- `~/.config/ccam/config.toml`: die Einträge `[[group]]`: `name`, `items`, `members` und, wenn der Name des Ordners abweicht, `dir`.
- `~/.local/share/ccam/profiles/<name>/`: das `CLAUDE_CONFIG_DIR` des Kontos.
- `~/.local/share/ccam/groups/<folder>/`: die geteilten Kopien.

Beim Lesen von `config.toml` werden die Gruppen jedes Kontos zu einem Ordner pro Element aufgelöst. Das dauert bei fünf Konten etwa 25 µs, sodass ein Start weiterhin nur `config.toml` liest. Eine Konfiguration, in der ein Konto ein Element aus zwei Gruppen bekäme, ist ein Fehler, und nichts startet, bis er behoben ist.

Eine Änderung arbeitet auf einer Kopie der Konfiguration:

1. Der Befehl bearbeitet die Kopie (`AddGroup`, `SetGroupItems`, `JoinGroup` …), die die Regeln prüft.
2. `internal/group` vergleicht alt und neu und baut für jedes Konto, dessen Elemente sich ändern, mit `share.Build` einen Plan. Die Pläne entstehen nacheinander, und jeder sieht, was die früheren in die Gruppenordner legen werden, sodass nur das erste Konto seine Kopie hineinverschiebt. Die erste Kopie kommt aus dem Gruppenordner, wenn er schon eine hat, dann von `--from`, dann vom Standardkonto, dann von den übrigen nach Nummer.
3. ccam lehnt ab, wenn eines dieser Konten läuft, und unter Windows ohne Entwicklermodus, wenn eine Datei einen Link bräuchte.
4. Konflikte werden entschieden, per Nachfrage oder durch `--prefer`, der Plan wird ausgegeben, und ccam fragt nach.
5. Die neue Konfiguration wird gespeichert, bevor sich irgendetwas bewegt. Bricht ein Schritt mittendrin ab, bauen `ccam doctor --fix` und ein erneuter Lauf ihren Plan auf die gespeicherten Gruppen hin und schließen die Änderung ab.
6. Die Pläne laufen der Reihe nach. Jeder Schritt verschiebt Daten in die Gruppe oder in den Sicherungsordner; nichts wird gelöscht.

`share.Layout` ist eine Positivliste der Pfade, die jedes Element abdeckt. Was nicht darauf steht, kann nie in eine Gruppe gelangen. `share.Never` nennt `.credentials.json` und `.claude.json`, und `Build` hält mit einem internen Fehler an, falls eine davon je in der Struktur auftaucht. Das Element `mcp` hat keine Pfade: Die Serverliste liegt in `.claude.json`, daher funktioniert das [Teilen von MCP-Servern](/features/mcp) über die `mcp.json` der Gruppe.

Das passiert mit einem Ordner, `projects/`, wenn `side` anfängt, den Verlauf zu teilen, und die Gruppe bereits Inhalt hat:

<!-- motion:share-on -->

Hat die Gruppe unter diesem Pfad noch nichts, teilt das Konto ihn als erstes, und sein Ordner wird unverändert dorthin verschoben. Andernfalls wandert jede Datei, die der Gruppe fehlt, in die Gruppe, und Dateien, die auf beiden Seiten existieren, bleiben im Ordner des Kontos, sodass die Version der Gruppe erhalten bleibt. Der Ordner des Kontos wandert dann nach `.ccam-backup/<time>/` im Kontoordner, und ein Link tritt an seine Stelle. Gesprächsdateien sind nach ihrer Sitzungs-ID benannt, daher lassen sich Verläufe aus mehreren Konten ohne Konflikte zusammenführen. Die `plugins/`-Register speichern absolute Pfade durch den Kontoordner; ccam richtet sie auf den Ordner der Gruppe aus und führt die Einträge des Kontos in ein Register zusammen, das die Gruppe bereits hat.

Das Verlassen kehrt den Link um:

<!-- motion:share-off -->

ccam entfernt den Link nur, wenn er in den Gruppenordner zeigt, und kopiert dann den aktuellen Inhalt der Gruppe in das Konto. Mit `--fresh` entfällt die Kopie. Die Gruppe behält ihren Inhalt für die anderen Konten. Ein Konto, das ein Element in eine Gruppe verschiebt, deren Ordner noch leer ist, bringt eine Kopie dessen mit, was es gesehen hat, und behält so seine Dateien; MCP-Server kommen auf dieselbe Weise mit. Ein Konto, das MCP nicht mehr teilt, verliert die Server, die ccam in seine `.claude.json` kopiert hatte (eine Sicherung bleibt), da sie von der Gruppe stammten.

Eine einzelne Datei, die auf beiden Seiten mit unterschiedlichem Inhalt existiert, ist ein Konflikt:

<!-- motion:share-conflict -->

`history.jsonl` bietet `merge by line`, `keep group` und `use mine`, mit dem Zusammenführen an erster Stelle. Andere `.json`-Dateien bieten `keep group`, `use mine` und `merge JSON`. Alles andere bietet `keep group` und `use mine`. `merge JSON` führt Objekte rekursiv zusammen, verbindet Arrays ohne Duplikate und behält den Wert der Gruppe, wenn beide Seiten denselben Schlüssel unterschiedlich setzen, sodass die anderen Konten keine Änderung bemerken. Jede Wahl außer `keep group` sichert zuerst die Kopie der Gruppe in `.ccam-backup/<time>/group/` des Kontos. Deine eigene Datei kommt immer in die Sicherung.

Bei jedem Start vergleicht `share.Check` jeden Pfad der Struktur mit den Gruppen des Kontos und verwendet dabei nur `lstat` und `readlink`. Ein fehlender Link wird sofort angelegt. Eine echte Datei, wo ein Link hingehört, ein Link an die falsche Stelle oder ein veralteter Link hält den Start an und verweist auf `ccam doctor --fix`.

### Umbenennen, Entfernen und das alte Format

- `ccam group rename` ändert nur die Konfiguration. Der Ordner behält seinen Namen, und die Gruppe hält ihn als `dir` fest, sodass sich nichts bewegt und nichts verloren gehen kann.
- `ccam group remove` lehnt eine Gruppe mit Konten ab, sofern du nicht `--leave` angibst. Jedes Konto behält dann eine Kopie, und der Ordner kommt nach `~/.local/share/ccam/trash/`, sofern keine andere Gruppe ihn verwendet.
- Vor den Gruppen speicherte `config.toml` pro Konto eine Tabelle `[profile.share]`. ccam liest dieses Format als Gruppen, ohne etwas zu verschieben. Beim ersten Speichern schreibt es das neue Format und behält die alte Datei als `config.toml.bak`. Elemente, die dieselben Konten teilen, werden zu einer Gruppe, die nach dem Ordner benannt ist; die übrigen erhalten Namen wie `default-history` und verwenden weiter denselben Ordner.

### Unter Windows

Windows legt symbolische Links nur im Entwicklermodus (Einstellungen > System > Für Entwickler) oder für einen Administrator an. Ohne dieses Recht wird ein geteilter Ordner zu einer Verzeichnisverbindung (Junction), die für Claude Code genauso funktioniert. Eine geteilte Datei (`settings.json`, `CLAUDE.md`, `history.jsonl`) kann nur ein symbolischer Link sein, deshalb wird eine Gruppenänderung, die einen bräuchte, abgelehnt, bevor sich etwas bewegt.

## Sicherheit und Sonderfälle

- **Laufende Konten.** Eine Änderung wird abgelehnt, solange Claude Code als eines der Konten läuft, die sie verschiebt, weil die Sitzung weiter in diese Dateien schreibt. Verlassen und Beitreten betreffen nur das genannte Konto. `--dry-run` gibt den Plan trotzdem aus.
- **Sicherungen.** Jeder Plan erhält einen eigenen Ordner `.ccam-backup/<time>/` im Kontoordner. ccam löscht sie nie.
- **Deine eigenen Links.** Ist ein Pfad bereits ein Link, den du angelegt hast, zum Beispiel auf ein Dotfiles-Repository, und ist die Gruppe leer, kopiert ccam das Ziel des Links in die Gruppe und verschiebt den Link selbst in die Sicherung. Das Ziel wird nicht angefasst.
- **Importierte Ordner.** Ein Konto, das mit `ccam import` an Ort und Stelle registriert wurde (etwa `~/.claude`), ändert sich an Ort und Stelle, sodass auch einfaches `claude` die Änderung sieht; der Plan weist darauf hin.
- **Zwei Konten in einem Ordner.** Konten, die den Verlauf teilen, schreiben in dieselben Dateien. Wenn ein zweites in einem Ordner startet, in dem bereits eines läuft, warnt ccam und fragt nach.
- **Änderungen von Hand.** Ein Mitglied, das kein Konto ist, wird ignoriert und von `ccam doctor` gemeldet. Die Namen `local`, `on`, `off` und `none` können keine Gruppennamen sein.

## Verwandte Themen

- [Das Dashboard](/features/dashboard)
- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Geteilte MCP-Server](/features/mcp)
- [Eine vorhandene Einrichtung kopieren](/features/migrate)
- [Alle Befehle](/docs/commands)
