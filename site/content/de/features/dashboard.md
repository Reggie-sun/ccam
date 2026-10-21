---
title: Das Dashboard
description: ccam öffnet im Terminal ein Vollbild-Dashboard mit Konten, Gruppen, Ordnerregeln, Erweiterungen und Prüfungen; jede Änderung zeigt zuerst ihren Plan.
updated: 2026-10-21
order: 15
translation: machine
sourceHash: "27bd153b48b2"
---

`ccam` ohne Argumente öffnet ein Vollbild-Dashboard, wenn Standardeingabe und -ausgabe ein Terminal sind; andernfalls gibt es `ccam list` aus. Das Dashboard zeigt jedes Konto mit seiner Nutzung und seinen Gruppen und erledigt, was die Befehle tun: starten, hinzufügen, anmelden, Gruppen ändern, Ordnerregeln setzen, prüfen und reparieren. Ein Start schließt es; beim Hinzufügen eines Kontos oder bei der Anmeldung überlässt es das Terminal kurz der Aufgabe und kehrt dann zurück.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

Die Zeile unter den Tab-Namen sagt in einem Satz, wofür der aktuelle Tab da ist. In einem breiteren Terminal zeigt der Tab Accounts außerdem die E-Mail-Adresse jedes Kontos.

Das erste Bild entsteht allein aus `config.toml`, in deutlich unter 100 ms; die Nutzung erscheint, sobald das jeweilige Konto antwortet. Jedes Limit hat seinen eigenen Balken: das 5-Stunden-Fenster, das Wochenfenster für alle Modelle und das eigene Wochenfenster von Fable; `-` bedeutet, dass das Konto kein solches Fenster hat. Ein Prozentwert wird ab 80 % gelb und bei 100 % rot, und Werte aus einem Cache nach einer fehlgeschlagenen Abfrage werden blass dargestellt.

## Tabs

| Tab | Was er zeigt | Tasten |
| --- | --- | --- |
| Accounts | Nummer, Name, E-Mail-Adresse, Balken für 5 Stunden, Woche und Fable, letzte Verwendung, Gruppen; `✔` das Standardkonto, `▸` das Konto dieses Ordners, `•` läuft | siehe unten |
| Groups | jede Gruppe, was sie teilt und ihre Konten, sowie die Konten ohne Gruppe | `enter` bearbeiten, `n` neu, `e` umbenennen, `x` entfernen |
| Folders | die Regeln von [`ccam folder`](/features/folder) und die, die hier gilt | `n` ein Konto hier verwenden, `enter` ändern, `x` entfernen |
| Extensions | was [`ccam ext`](/features/ext) zeigt, jede Erweiterung einmal mit den Gruppen, aus denen sie kommt | `r` neu laden |
| Health | was [`ccam doctor`](/features/doctor) meldet, als Kasten pro Konto (siehe unten) | `f` den ausgewählten Kasten reparieren, `F` alles reparieren (`doctor --fix`), `↑`/`↓` zwischen Kästen wechseln |

`←`/`→` oder `Tab` wechseln den Tab, `?` listet die Tasten des aktuellen Tabs auf, `r` fragt die Nutzung ohne den kurzen Cache erneut ab, und `q` oder `Esc` beendet.

## Tasten auf dem Tab Accounts

| Taste | Wirkung |
| --- | --- |
| `Enter` | Claude Code als das ausgewählte Konto starten |
| `1`–`9` | Zum Konto mit dieser Nummer gehen; tippe `1` `2` schnell hintereinander für 12. Nur `Enter` startet |
| `↑`/`↓`, `k`/`j` | Bewegen |
| `n` | Das Konto starten, das [`ccam next --auto`](/features/next) wählen würde |
| `o` | Mit Optionen starten: `--continue`, `--resume`, andere Argumente oder ein Probelauf |
| `u` | Nutzungsdetails: 5-Stunden-, Wochen- und Modellfenster mit ihren Zurücksetzungen und der Herkunft der Zahlen |
| `a` | Ein Konto hinzufügen (der Assistent unten) |
| `l` | Erneut anmelden |
| `s` | Die Gruppen des Kontos wählen |
| `d` | Es zum Standardkonto machen |
| `f` | Es im aktuellen Ordner verwenden |
| `e` | Es umbenennen |
| `c` | Seinen Kurzbefehl setzen oder entfernen |
| `x` | Es entfernen (fragt zuerst nach, mit No vorausgewählt) |

## Der Tab Health

Der Tab Health führt beim Öffnen die Prüfungen von [`ccam doctor`](/features/doctor) aus und zeigt dieselbe Ansicht: eine erste Zeile mit einem Zeichen pro Konto, dann einen Kasten pro Konto. Jede Zeile in einem Kasten ist ein geprüfter Bereich, zuerst `Folder`, `Login` und `Share links`, markiert mit `✔` in Ordnung, `ℹ` Information, `⚠` eine Warnung oder `✘` ein Fehler, sodass du auf einen Blick siehst, welches Konto wo ein Problem hat. Befunde, die zu keinem Konto gehören, kommen in einen Kasten `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ──────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                 │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)      │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link                │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` wechseln von Kasten zu Kasten, und der ausgewählte Kasten hat einen Rahmen in der Akzentfarbe. `f` repariert nur den ausgewählten Kasten, so wie `ccam doctor --fix` es tun würde: im Kasten eines Kontos nur dieses Konto, im Kasten `General` die Befunde, die zu keinem Konto gehören (Plugin-Pfade und Kurzbefehle); alle anderen Konten werden nur erneut geprüft. Vorher kommt eine Rückfrage, und ein Kasten, in dem es nichts zu reparieren gibt, meldet nur `nothing to repair in <name>`. `F` repariert alles auf einmal, wie `ccam doctor --fix`, und `r` prüft erneut.

## Gruppen ändern

Auf dem Tab Groups öffnet `Enter` den Editor: links die Elemente, rechts die Konten, `Tab` zum Wechseln, `Space` zum Umschalten. Ein Element, das eines der Konten bereits von einer anderen Gruppe bekommt, ist mit dem Grund ausgegraut, ebenso ein Konto, das ein Element doppelt bekäme; ein Element kann nur aus einer Gruppe kommen. `Enter` zeigt den Plan:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Zu einer Datei, die von der Kopie der Gruppe abweicht, wird zuerst nachgefragt (die der Gruppe behalten, meine verwenden oder zusammenführen). Dann läuft die Änderung im Hintergrund mit demselben Code wie [`ccam group`](/features/share): Die Konfiguration wird zuerst gespeichert, die Daten jedes Kontos werden der Reihe nach verschoben, und die Originale kommen nach `.ccam-backup/`. Eine Änderung, die die Daten eines laufenden Kontos verschieben würde, wird mit dem Namen des Kontos abgelehnt.

## Ein Konto hinzufügen

`a` öffnet den Assistenten, den auch `ccam add` in einem Terminal verwendet:

1. Der Kontoname.
2. Wie es sich anmeldet: jetzt mit einem Claude-Abo (im Browser), mit einem [API-Schlüssel](/features/api) oder später mit einem Claude-Abo.
3. Bei einem API-Schlüssel: woher ccam den Schlüssel bei jedem Start liest (`env:NAME` oder `command:COMMAND`), der Endpunkt, wie der Schlüssel gesendet wird, und das Modell.
4. Die Gruppen zum Beitreten, mit angehaktem `default`; solange es noch keine Gruppen gibt, was die Gruppe `default` teilt.
5. Ein Kurzbefehl, wenn du einen willst.
6. Eine Zusammenfassung zum Bestätigen.

`Esc` geht einen Schritt zurück und schließt den Assistenten im ersten Schritt.

## Farben

Die Farben folgen den Rollen der Terminal-Oberfläche von Claude Code, sodass ccam daneben vertraut aussieht; Gruppennamen erhalten eine Farbe, die sich aus ihrem Namen ergibt. ccam fragt das Terminal nach seinem Hintergrund und wählt ein helles oder dunkles Design. Du kannst eines selbst in `config.toml` festlegen:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Die Designs mit `-daltonized` sind für Rot-Grün-Farbenblindheit gedacht, und die Designs mit `-ansi` verwenden die 16 eigenen Farben des Terminals. `NO_COLOR` schaltet die Farben aus.

## Sicherheit und Sonderfälle

- Das Dashboard ruft dieselben Funktionen auf wie die Befehle, daher gilt jede Regel: Pläne, bevor Daten bewegt werden, Ablehnungen bei laufenden Konten, Sicherungen und Anmeldetokens, die ein Konto nie verlassen.
- Das Dashboard hat keinen Probelauf: `ccam --dry-run` ohne Namen gibt aus, wie das Konto dieses Ordners starten würde, und `o` bei einem Konto bietet einen Probelauf genau dieses einen Starts an.
- Das Terminal wird wiederhergestellt, bevor Claude Code startet oder sich die Anmeldung öffnet.
- Bei einem Bildschirm, der für einen Dialog zu klein ist, fallen zuerst die Zeilen des Tabs weg; Zeilen werden mit `…` abgeschnitten statt umbrochen.

## Verwandte Themen

- [Teilen mit Gruppen](/features/share)
- [Claude Code starten](/features/launch)
- [Konten verwalten](/features/accounts)
- [Nutzungslimits](/features/usage)
