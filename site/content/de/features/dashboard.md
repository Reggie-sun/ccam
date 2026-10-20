---
title: Das Dashboard
description: ccam öffnet im Terminal ein Vollbild-Dashboard mit Konten, Gruppen, Ordnerregeln, Erweiterungen und Prüfungen; jede Änderung zeigt zuerst ihren Plan.
updated: 2026-10-20
order: 15
translation: machine
sourceHash: "decdd5fa54c9"
---

`ccam` ohne Argumente öffnet ein Vollbild-Dashboard, wenn Standardeingabe und -ausgabe ein Terminal sind; andernfalls gibt es `ccam list` aus. Das Dashboard zeigt jedes Konto mit seiner Nutzung und seinen Gruppen und erledigt, was die Befehle tun: starten, hinzufügen, anmelden, Gruppen ändern, Ordnerregeln setzen, prüfen und reparieren. Ein Start schließt es; beim Hinzufügen eines Kontos oder bei der Anmeldung überlässt es das Terminal kurz der Aufgabe und kehrt dann zurück.

```console
 ccam v0.9.0 · 3 accounts · 2 groups         ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health

    # Account     Email             5-hour         7-day Last  Groups
 ❯ 1. personal ✔   me@example.com    █████▊░░  72%   41% 2h    default
   2. work      ▸• dev@company.com   ██░░░░░░  25%  100% now   default team
   3. side         -                 ✘ not logged in         -     no group

 ✔ personal is the default account
 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

Das erste Bild entsteht allein aus `config.toml`, in deutlich unter 100 ms; die Nutzung erscheint, sobald das jeweilige Konto antwortet. Ein Prozentwert wird ab 80 % gelb und bei 100 % rot, und Werte aus einem Cache nach einer fehlgeschlagenen Abfrage werden blass dargestellt.

## Tabs

| Tab | Was er zeigt | Tasten |
| --- | --- | --- |
| Accounts | Nummer, Name, E-Mail-Adresse, 5-Stunden-Balken, 7-Tage-Nutzung, letzte Verwendung, Gruppen; `✔` das Standardkonto, `▸` das Konto dieses Ordners, `•` läuft | siehe unten |
| Groups | jede Gruppe, was sie teilt und ihre Konten, sowie die Konten ohne Gruppe | `enter` bearbeiten, `n` neu, `e` umbenennen, `x` entfernen |
| Folders | die Regeln von [`ccam folder`](/features/folder) und die, die hier gilt | `n` ein Konto hier verwenden, `enter` ändern, `x` entfernen |
| Extensions | was [`ccam ext`](/features/ext) zeigt, jede Erweiterung einmal mit den Gruppen, aus denen sie kommt | `r` neu laden |
| Health | was [`ccam doctor`](/features/doctor) meldet | `f` reparieren (`doctor --fix`) |

`←`/`→` oder `Tab` wechseln den Tab, `?` listet die Tasten des aktuellen Tabs auf, `r` fragt die Nutzung ohne den kurzen Cache erneut ab, und `q` oder `Esc` beendet.

## Tasten auf dem Tab Accounts

| Taste | Wirkung |
| --- | --- |
| `Enter` | Claude Code als das ausgewählte Konto starten |
| `1`–`9` | Zum Konto mit dieser Nummer gehen; tippe `1` `2` schnell hintereinander für 12. Nur `Enter` startet |
| `↑`/`↓`, `k`/`j` | Bewegen |
| `n` | Das Konto starten, das [`ccam next --auto`](/features/next) wählen würde |
| `o` | Mit Optionen starten: `--continue`, `--resume`, andere Argumente oder ein Probelauf |
| `u` | Nutzungsdetails: 5-Stunden-, 7-Tage- und Modellfenster mit ihren Zurücksetzungen und der Herkunft der Zahlen |
| `a` | Ein Konto hinzufügen (der Assistent unten) |
| `l` | Erneut anmelden |
| `s` | Die Gruppen des Kontos wählen |
| `d` | Es zum Standardkonto machen |
| `f` | Es im aktuellen Ordner verwenden |
| `e` | Es umbenennen |
| `c` | Seinen Kurzbefehl setzen oder entfernen |
| `x` | Es entfernen (fragt zuerst nach, mit No vorausgewählt) |

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
