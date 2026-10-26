---
title: Teilen mit Gruppen
description: Eine ccam-Gruppe teilt gewählte Elemente über einen Ordner zwischen ihren Konten; ccam group verschiebt die Daten mit Plan, Sicherungen und Konfliktwahl.
updated: 2026-10-26
order: 4
translation: machine
sourceHash: 3f94dd13dbf1
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
ccam group set default --add history              # all accounts in default start sharing history
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # the group starts from work's copies
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # first leave the group that shares those items now
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
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

Dieselben Änderungen kannst du auf dem Tab **Groups** des [Dashboards](/features/dashboard) vornehmen; drücke `s` auf einem Konto, um seine Gruppen zu wählen.

## Elemente

| Element | Was es umfasst |
| --- | --- |
| `mcp` | MCP-Server, über die `mcp.json` der Gruppe ([Geteilte MCP-Server](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | die globale `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` und `history.jsonl`, damit ein anderes Konto mit `--continue` weitermachen kann |

Teilt keine der Gruppen eines Kontos ein Element, bleibt dieses Element im eigenen Ordner des Kontos.

## Befehle und Optionen

| Befehl | Wirkung |
| --- | --- |
| `ccam group` oder `ccam group list` | die Gruppen auflisten |
| `ccam group show <group>` | was die Gruppe teilt, ihre Konten, ihr Ordner und ihre MCP-Server |
| `ccam group add <group> --items <items> [--members <accounts>]` | eine Gruppe anlegen; die Elemente ihrer Konten wandern in sie |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | ändern, was die Gruppe teilt, für jedes Konto in ihr |
| `ccam group join <group> <account>...` | Konten in die Gruppe aufnehmen |
| `ccam group leave <group> <account>...` | Konten herausnehmen; jedes behält eine Kopie, MCP-Server eingeschlossen |
| `ccam group rename <group> <new-name>` | nur in der Konfiguration umbenennen |
| `ccam group remove <group> [--leave]` | eine Gruppe entfernen |

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: das Konto, dessen Kopie als Erste in die Gruppe kommt | das Standardkonto, dann nach Nummer |
| `--replace` | `join`: zuerst die Gruppen verlassen, die dieselben Elemente teilen | aus |
| `--fresh` | Konten, die ein Element nicht mehr teilen, beginnen ohne es statt mit einer Kopie (ein importierter Ordner behält trotzdem eine) | aus |
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
2. ccam vergleicht alt und neu und baut für jedes Konto, dessen Elemente sich ändern, einen Plan (`internal/group`, `share.Build`). Die Pläne entstehen nacheinander, und jeder sieht, was die früheren in die Gruppenordner legen werden, sodass nur das erste Konto seine Kopie hineinverschiebt. Die erste Kopie der Gruppe kommt in dieser Reihenfolge aus: dem Gruppenordner (wenn er schon eine hat), dem Konto aus `--from`, dem Standardkonto, dann den übrigen Konten nach Nummer.
3. ccam lehnt ab, wenn eines dieser Konten läuft, egal wie Claude Code gestartet wurde (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``), und unter Windows ohne Entwicklermodus, wenn eine Datei einen Link bräuchte.
4. ccam löst Konflikte (per Nachfrage oder mit `--prefer`), gibt den Plan aus und bittet um Bestätigung.
5. Direkt vor dem Speichern prüft ccam erneut. Wurde ein Konto gestartet, während ccam auf deine Antworten wartete, hält die Änderung an; sie hält auch an, wenn ein anderer ccam-Prozess in der Zwischenzeit dieselben Gruppen geändert hat (`the groups changed while this plan waited; nothing was moved, make the change again`). Die neue Konfiguration wird gespeichert, bevor sich irgendetwas bewegt, und die Änderung wird in `~/.local/state/ccam/group-change.json` vermerkt. Bricht ein Schritt mittendrin ab, bleibt der Vermerk: `ccam doctor --fix` schließt die Änderung in ihrer Reihenfolge ab, mit den Antworten, die du gegeben hast, MCP-Schritte eingeschlossen, und bis dahin beginnt keine andere Gruppenänderung (``the last group change stopped halfway; finish it with `ccam doctor --fix` first``).
6. Die MCP-Schritte laufen zuerst. Wenn Konten Erweiterungen aus einer Gruppe herausnehmen, werden die Plugin-Listen der Gruppe überall dort auf den Gruppenordner ausgerichtet, wo ihre Pfade durch den Ordner eines Kontos führen (`group plugin lists:` im Plan). Dann laufen die Pläne der Reihe nach. Jeder Schritt verschiebt Daten in die Gruppe oder in den Sicherungsordner; nur Dateien, die mit der Kopie der Gruppe identisch sind, werden gelöscht statt gesichert.

ccam führt eine Positivliste der Pfade, die jedes Element abdeckt (`share.Layout`). Was nicht darauf steht, kann nie in eine Gruppe gelangen. ccam führt außerdem `.credentials.json` und `.claude.json` als nie geteilt (`share.Never`) und hält mit einem internen Fehler an, falls eine davon je in der Struktur auftaucht (`Build`). Das Element `mcp` hat keine Pfade: Die Serverliste liegt in `.claude.json`, daher funktioniert das [Teilen von MCP-Servern](/features/mcp) über die `mcp.json` der Gruppe.

Das passiert mit einem Ordner, `projects/`, wenn `side` anfängt, den Verlauf zu teilen, und die Gruppe bereits Inhalt hat:

<!-- motion:share-on -->

Hat die Gruppe unter diesem Pfad noch nichts, teilt das Konto ihn als erstes, und sein Ordner wird unverändert dorthin verschoben. Andernfalls wandert jede Datei, die der Gruppe fehlt, in die Gruppe, und Dateien, die auf beiden Seiten existieren, bleiben im Ordner des Kontos, sodass die Version der Gruppe erhalten bleibt. Dateien, die mit denen der Gruppe identisch sind, werden gelöscht, und der Plan nennt die abweichenden (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`).

Der Ordner des Kontos wandert dann nach `.ccam-backup/<time>/` im Kontoordner, und ein Link tritt an seine Stelle. Ist der Eintrag der Gruppe selbst ein Link auf den eigenen Eintrag des beitretenden Kontos, kommt stattdessen der Inhalt des Kontos in die Gruppe. Gesprächsdateien sind nach ihrer Sitzungs-ID benannt, daher lassen sich Verläufe aus mehreren Konten ohne Konflikte zusammenführen. Die `plugins/`-Register speichern absolute Pfade durch den Kontoordner; ccam richtet sie auf den Ordner der Gruppe aus und führt die Einträge des Kontos in ein Register zusammen, das die Gruppe bereits hat.

Das Verlassen kehrt den Link um:

<!-- motion:share-off -->

ccam kopiert den aktuellen Inhalt der Gruppe zuerst nach `.ccam-backup/incomplete-<name>` im Kontoordner und ersetzt den Link erst, wenn die Kopie vollständig ist. Bricht die Kopie ab, bleibt der Link stehen, und `ccam doctor --fix` kopiert erneut (der Plan eines solchen Abschlusses zeigt `restore <entry>: … (copied back from the group before a stop)`). Die Gruppe behält ihren Inhalt für die anderen Konten.

- ccam entfernt den Link nur, wenn er noch so in den Gruppenordner zeigt wie beim Erstellen des Plans. Ein Pfad, der sich seitdem geändert hat, etwa eine echte Datei, wo der Link war, bleibt, wie er ist, und die Änderung hält an.
- Mit `--fresh` entfällt die Kopie.
- Ein Konto, das ein Element in eine Gruppe verschiebt, deren Ordner noch leer ist, bringt eine Kopie dessen mit, was es gesehen hat, und behält so seine Dateien. Diese Kopie wird in `<group>/.ccam-incomplete-<name>` vorbereitet, bis sie vollständig ist.
- Ein Konto, das MCP nicht mehr teilt, behält die Server der Gruppe: Sie kommen in seine eigene `.claude.json` (eine Sicherung bleibt), außer wo das Konto bereits einen Server gleichen Namens hat, der unverändert bleibt. Mit `--fresh` listet der Plan stattdessen die Server auf, die nicht mehr geladen werden. Ein erneuter Beitritt entfernt die Kopien, die niemand geändert hat.
- Das Verlassen von `extensions` richtet außerdem die Plugin-Listen der Gruppe auf den Gruppenordner aus, wo ihre Pfade durch den Ordner des austretenden Kontos führten; für die anderen Konten führen die umgeschriebenen Pfade zu denselben Plugins.

Eine einzelne Datei, die auf beiden Seiten mit unterschiedlichem Inhalt existiert, ist ein Konflikt:

<!-- motion:share-conflict -->

Bei einzelnen Dateien wie `settings.json` wählst du; innerhalb von Ordnern bleibt die Kopie der Gruppe, und der Plan nennt die abweichenden Dateien. Die Auswahl hängt von der Datei ab; die erste ist der Standard:

| Datei | Auswahl |
| --- | --- |
| `history.jsonl` | `merge by line`, `keep group`, `use mine` |
| Andere `.json`-Dateien | `keep group`, `use mine`, `merge JSON` |
| Alles andere | `keep group`, `use mine` |

`merge JSON` führt Objekte rekursiv zusammen, verbindet Arrays ohne Duplikate und behält den Wert der Gruppe, wenn beide Seiten denselben Schlüssel unterschiedlich setzen, sodass die anderen Konten keine Änderung bemerken. Jede Wahl außer `keep group` sichert zuerst die Kopie der Gruppe in `.ccam-backup/<time>/group/` des Kontos. Deine eigene Datei kommt immer in die Sicherung.

Bei jedem Start vergleicht ccam jeden Pfad der Struktur mit den Gruppen des Kontos und verwendet dabei nur `lstat` und `readlink` (`share.Check`). Ein fehlender Link wird sofort angelegt. Eine echte Datei, wo ein Link hingehört, ein Link an die falsche Stelle oder ein veralteter Link hält den Start an und verweist auf `ccam doctor --fix`.

### Umbenennen, Entfernen und das alte Format

- `ccam group rename` ändert nur die Konfiguration. Der Ordner behält seinen Namen, und die Gruppe hält ihn als `dir` fest, sodass sich nichts bewegt und nichts verloren gehen kann.
- `ccam group remove` lehnt eine Gruppe mit Konten ab, sofern du nicht `--leave` angibst. Jedes Konto behält dann eine Kopie, und der Ordner kommt nach `~/.local/share/ccam/trash/`, sofern keine andere Gruppe ihn verwendet.
- Vor den Gruppen speicherte `config.toml` pro Konto eine Tabelle `[profile.share]`. ccam liest dieses Format als Gruppen, ohne etwas zu verschieben. Beim ersten Speichern schreibt es das neue Format und behält die alte Datei als `config.toml.bak`. Elemente, die dieselben Konten teilen, werden zu einer Gruppe, die nach dem Ordner benannt ist; die übrigen erhalten Namen wie `default-history` und verwenden weiter denselben Ordner.

### Unter Windows

Windows legt symbolische Links nur im Entwicklermodus (Einstellungen > System > Für Entwickler) oder für einen Administrator an. Ohne dieses Recht wird ein geteilter Ordner zu einer Verzeichnisverbindung (Junction), die für Claude Code genauso funktioniert. Eine geteilte Datei (`settings.json`, `CLAUDE.md`, `history.jsonl`) kann nur ein symbolischer Link sein, deshalb wird eine Gruppenänderung, die einen bräuchte, abgelehnt, bevor sich etwas bewegt. Aus demselben Grund legt das erste `ccam add` dort die Gruppe `default` nur mit `mcp` und `extensions` an und sagt das, und `ccam add --group` in eine Gruppe, die Einstellungen, Memory oder Verlauf teilt, wird abgelehnt.

## Sicherheit und Sonderfälle

- **Laufende Konten.** Eine Änderung wird abgelehnt, solange Claude Code als eines der Konten läuft, die sie verschiebt, weil die Sitzung weiter in diese Dateien schreibt, egal ob ccam sie gestartet hat oder nicht (einfaches `claude`, eine IDE, Hintergrundsitzungen). Die Prüfung wird direkt vor dem Speichern wiederholt. Verlassen und Beitreten betreffen nur das genannte Konto; das Verlassen von `extensions` schreibt außerdem Pfade in den Plugin-Listen der Gruppe um, ohne zu ändern, was sie für die anderen bedeuten. `--dry-run` gibt den Plan trotzdem aus.
- **Sicherungen.** Jeder Plan erhält einen eigenen Ordner `.ccam-backup/<time>/` im Kontoordner. ccam löscht sie nie.
- **Deine eigenen Links.** Ist ein Pfad bereits ein Link, den du angelegt hast, zum Beispiel auf ein Dotfiles-Repository, und ist die Gruppe leer, kopiert ccam das Ziel des Links in die Gruppe und verschiebt den Link selbst in die Sicherung. Das Ziel wird nicht angefasst.
- **Importierte Ordner.** Ein Konto, das mit `ccam import` an Ort und Stelle registriert wurde (etwa `~/.claude`), ändert sich an Ort und Stelle, sodass auch einfaches `claude` die Änderung sieht; der Plan weist darauf hin. Verlässt es eine Gruppe, behält es immer eine Kopie, selbst mit `--fresh`, und `ccam remove` nimmt es auf diese Weise aus seinen Gruppen, bevor es das Konto vergisst. Seine geteilten Dateien liegen im Gruppenordner, solange es in einer Gruppe ist, also entferne das Konto (oder verlasse die Gruppen), bevor du `~/.local/share/ccam` löschst.
- **Zwei Konten in einem Ordner.** Konten, die den Verlauf teilen, schreiben in dieselben Dateien. Wenn ein zweites in einem Ordner startet, in dem bereits eines läuft, warnt ccam und fragt nach.
- **Manuelle Änderungen.** Steht in `members` ein Name, der kein Konto ist, ignoriert ccam ihn, und `ccam doctor` meldet ihn. Die Namen `local`, `on`, `off` und `none` können keine Gruppennamen sein.

## Verwandte Themen

- [Das Dashboard](/features/dashboard)
- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Geteilte MCP-Server](/features/mcp)
- [Eine vorhandene Einrichtung kopieren](/features/migrate)
- [Befehlsreferenz](/docs/commands)
