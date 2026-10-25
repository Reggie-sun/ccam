---
title: Das Dashboard
description: ccam öffnet im Terminal ein Vollbild-Dashboard mit Konten, Gruppen, Ordnerregeln, Erweiterungen und Prüfungen; jede Änderung zeigt zuerst ihren Plan.
updated: 2026-10-25
order: 15
translation: machine
sourceHash: "98ccf6cd81e9"
---

`ccam` ohne Argumente öffnet ein Vollbild-Dashboard, wenn Standardeingabe und -ausgabe ein Terminal sind; andernfalls gibt es `ccam list` aus. Das Dashboard zeigt jedes Konto mit seiner Nutzung und seinen Gruppen und erledigt, was die Befehle tun: starten, hinzufügen, anmelden, Gruppen ändern, Ordnerregeln setzen, Hintergrundsitzungen auflisten, prüfen und reparieren. Ein Start oder das Verbinden mit einer Hintergrundsitzung schließt es; beim Hinzufügen eines Kontos oder bei der Anmeldung überlässt es das Terminal kurz der Aufgabe und kehrt dann zurück.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Agents   Groups   Folders   Extensions   Health
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

Die Zeile unter den Tab-Namen sagt in einem Satz, wofür der aktuelle Tab da ist. In einem breiteren Terminal zeigt der Tab Accounts außerdem die E-Mail-Adresse jedes Kontos.

Bei jedem Öffnen fragt das Dashboard im Hintergrund bei GitHub nach, ob es ein neueres ccam-Release gibt; das erste Bild wartet nie darauf, und offline wird nichts angezeigt. Gibt es eines, empfiehlt eine Zeile unter dem Titel das Update:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` fragt zuerst nach und aktualisiert dann so wie [`ccam update`](/features/update): Das Release wird gegen seine SHA-256-Prüfsummen geprüft und ersetzt ccam in einem Schritt. Das Dashboard, in dem du gerade bist, ist noch die alte Version, also beende es und starte `ccam` neu. Wenn das Update hier nicht möglich ist, etwa bei einem aus dem `.deb`-Paket installierten ccam, zeigt das Dashboard denselben Hinweis wie `ccam update` und behält das Angebot bei.

Das erste Bild entsteht allein aus `config.toml`, in deutlich unter 100 ms; die Nutzung erscheint, sobald das jeweilige Konto antwortet. Jedes Limit hat seinen eigenen Balken: das 5-Stunden-Fenster, das Wochenfenster für alle Modelle und das eigene Wochenfenster von Fable; `-` bedeutet, dass das Konto kein solches Fenster hat. Ein Prozentwert wird ab 80 % gelb und bei 100 % rot, und Werte aus einem Cache nach einer fehlgeschlagenen Abfrage werden blass dargestellt. Ein Zeitfenster, das Claude Code als voll gemeldet hat, zeigt bis zu seinem Zurücksetzen 100 %. Ein abgelaufenes Zugriffstoken zeigt ein blasses `⚠ token expired` (Claude Code erneuert es beim nächsten Start), während `⚠ login expired` bedeutet, dass das Konto `ccam login` braucht.

## Tabs

| Tab | Was er zeigt | Tasten |
| --- | --- | --- |
| Accounts | Nummer, Name, E-Mail-Adresse, Balken für 5 Stunden, Woche und Fable, letzte Verwendung, Gruppen; `✔` das Standardkonto, `▸` das Konto dieses Ordners, `•` läuft (auch Claude Code, das ohne ccam gestartet wurde: einfaches `claude`, eine IDE, Hintergrundsitzungen und ihr Daemon) | siehe unten |
| Agents | was [`ccam agents`](/features/agents) auflistet: die Hintergrundsitzungen und offenen Sitzungen jedes Kontos, mit Zustand, Ordner und Startzeit (siehe unten) | `enter` verbinden (wie `ccam attach`), `x` stoppen, `r` neu laden |
| Groups | jede Gruppe, was sie teilt und ihre Konten, sowie die Konten ohne Gruppe | `enter` bearbeiten, `n` neu, `e` umbenennen, `x` entfernen |
| Folders | die Regeln von [`ccam folder`](/features/folder) und die, die hier gilt, auch wenn du den Ordner über einen symbolischen Link erreicht hast (neue Regeln werden mit aufgelösten Links gespeichert) | `n` ein Konto hier verwenden, `enter` ändern, `x` entfernen |
| Extensions | was [`ccam ext`](/features/ext) zeigt, jede Erweiterung einmal mit den Gruppen, aus denen sie kommt, `off:` die Konten, in denen ein Plugin nicht eingeschaltet ist, Server, die ein Konto unter seinen eigenen nur für diesen Ordner hat, `This folder's project` und Probleme, die zu keinem Konto gehören | `r` neu laden |
| Health | was [`ccam doctor`](/features/doctor) meldet, als Kasten pro Konto (siehe unten) | `f` den ausgewählten Kasten reparieren, `F` alles reparieren (`doctor --fix`), `↑`/`↓` zwischen Kästen wechseln |

`←`/`→` oder `Tab` wechseln den Tab, `?` listet die Tasten des aktuellen Tabs auf, `r` fragt die Nutzung ohne den kurzen Cache erneut ab und lädt die Tabs Agents, Extensions und Health neu, und `q` oder `Esc` beendet.

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

**Team- und Enterprise-Plätze.** Die automatische Wahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und [`auto_switch`](/features/next#auto-switch), wenn du es einschaltest) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation. Im Hilfe-Center von Anthropic ([Team-Tarife](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); ebenso für [Enterprise-Tarife](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)) heißt es: „Each seat on a Team plan is assigned to one named individual, and only that person may use it“ (jeder Platz in einem Team-Tarif ist einer namentlich benannten Person zugewiesen, und nur diese Person darf ihn nutzen) und „Users can't route their work through another user's account to use that user's remaining usage“ (Benutzer dürfen ihre Arbeit nicht über das Konto eines anderen Benutzers leiten, um dessen verbleibendes Nutzungskontingent zu verwenden). Wenn du die automatische Wahl zwischen Plätzen nutzt, kann das daher gegen die Bedingungen von Anthropic verstoßen. Ob du sie dort nutzt, ist deine Entscheidung: Kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt. Für private Tarife (Pro, Max) gilt der [Hinweis zu den Bedingungen in den FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Der Tab Health

Der Tab Health führt beim Öffnen die Prüfungen von [`ccam doctor`](/features/doctor) aus und zeigt dieselbe Ansicht: eine erste Zeile mit einem Zeichen pro Konto, dann einen Kasten pro Konto. Jede Zeile in einem Kasten ist ein geprüfter Bereich, zuerst `Folder`, `Login` und `Share links`, markiert mit `✔` in Ordnung, `ℹ` Information, `⚠` eine Warnung oder `✘` ein Fehler, sodass du auf einen Blick siehst, welches Konto wo ein Problem hat. Befunde, die zu keinem Konto gehören, kommen in einen Kasten `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` wechseln von Kasten zu Kasten, und der ausgewählte Kasten hat einen Rahmen in der Akzentfarbe. `f` repariert nur den ausgewählten Kasten, so wie `ccam doctor --fix` es tun würde: im Kasten eines Kontos nur dieses Konto, im Kasten `General` die Befunde, die zu keinem Konto gehören (Plugin-Pfade und Kurzbefehle); alle anderen Konten werden nur erneut geprüft. Vorher kommt eine Rückfrage (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), und ein Kasten, in dem es nichts zu reparieren gibt, meldet nur `nothing to repair in <name>`. `F` repariert alles auf einmal, wie `ccam doctor --fix`, und `r` prüft erneut.

## Der Tab Agents

Der Tab Agents zeigt, was [`ccam agents`](/features/agents) ausgibt: die Hintergrundsitzungen und offenen Sitzungen jedes Kontos, jeweils mit dem Konto, zu dem sie gehören. Er fragt das Claude Code jedes Kontos gleichzeitig ab (`claude agents --json`, jeweils bis zu 15 Sekunden), wenn du den Tab zum ersten Mal öffnest, bei `r`, nach einem Stopp und nach einer Änderung der Konfiguration, und zeigt `Asking every account's Claude Code for its sessions…`, bis alle geantwortet haben.

```console
  Accounts   Agents   Groups   Folders   Extensions   Health
 Background agents and open sessions of every account; enter attaches to one.

    Account   ID        Name   State          Folder      Started
 ❯  work      -         -      open (busy)    ~/code/web  now
    personal  49ec0028  probe  blocked, idle  ~/code/app  12m

 ✘ lab: could not list its sessions: exit status 1 'claude agents --json' is
   not available in this environment.

 enter to attach · x to stop · r to reload · ? for shortcuts
```

`State` ist der Zustand einer Hintergrundsitzung und das, was ihr Prozess gerade tut, oder `open` für eine interaktive Sitzung, wie in `ccam agents`. Ein Konto, das sich nicht abfragen ließ, bekommt unter der Liste eine `✘`-Zeile mit dem, was Claude Code gesagt hat, wie `lab` oben, dessen Agent-Ansicht ausgeschaltet ist. Beendete Hintergrundsitzungen werden ausgelassen (`ccam agents --all` listet sie auf), und wenn alle Konten geantwortet haben und keines eine Sitzung hat, zeigt der Tab ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` In einem schmalen Terminal lässt der Tab zuerst `Started` weg, dann `Name`, dann `Folder`, und kürzt einen langen Ordner von links.

`Enter` auf einer Hintergrundsitzung schließt das Dashboard und öffnet die Sitzung in diesem Terminal als ihr Konto, so wie `ccam attach <id>`; das Dashboard kommt nicht zurück (in Claude Code führt `←` zu seiner Agent-Ansicht und `Ctrl+Z` zur Shell). Eine offene Sitzung hat keine ID, daher sagen `Enter` und `x` dort nur, dass eine Hintergrundsitzung nötig ist. `x` fragt vorher nach, beginnend bei `No, keep it`: ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Bei Ja führt ccam `claude stop <id>` als dieses Konto aus, meldet ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` oder warum es nicht ging, und liest die Liste neu ein.

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

Zu einer Datei, die von der Kopie der Gruppe abweicht, wird zuerst nachgefragt (die der Gruppe behalten, meine verwenden oder zusammenführen). Dann läuft die Änderung im Hintergrund mit demselben Code wie [`ccam group`](/features/share): Die Konfiguration wird zuerst gespeichert, die Daten jedes Kontos werden der Reihe nach verschoben, und die Originale kommen nach `.ccam-backup/`; ein Konto, das das Teilen von MCP verlässt, behält die Server der Gruppe in seiner `.claude.json`. Eine Änderung, die die Daten eines laufenden Kontos verschieben würde, wird mit dem Namen des Kontos abgelehnt, auch wenn das Konto gestartet wurde, während der Plan offen war, und ein Plan, den die Gruppenänderung eines anderen ccam veraltet gemacht hat, wird mit `the groups changed while this plan waited; nothing was moved, make the change again` abgelehnt. Während ein anderes ccam die Gruppen ändert, melden eine neue Änderung, ein Umbenennen oder Entfernen `another ccam is changing the groups right now`; nach einer Änderung, die auf halbem Weg stehen geblieben ist, warten sie, bis `F` auf dem Tab Health (oder `ccam doctor --fix`) sie abgeschlossen hat.

## Ein Konto hinzufügen

`a` öffnet den Assistenten, den auch `ccam add` in einem Terminal verwendet:

1. Der Kontoname.
2. Wie es sich anmeldet: jetzt mit einem Claude-Abo (im Browser), mit einem [API-Schlüssel](/features/api) oder später mit einem Claude-Abo.
3. Bei einem API-Schlüssel: woher ccam den Schlüssel bei jedem Start liest (`env:NAME` oder `command:COMMAND`), der Endpunkt, wie der Schlüssel gesendet wird, und das Modell.
4. Die Gruppen zum Beitreten, mit angehaktem `default`; solange es noch keine Gruppen gibt, was die Gruppe `default` teilt.
5. Ein Kurzbefehl, wenn du einen willst.
6. Eine Zusammenfassung zum Bestätigen.

`Esc` geht einen Schritt zurück und schließt den Assistenten im ersten Schritt.

## Ein Konto entfernen

`x` fragt nach, bevor irgendetwas passiert, mit No vorausgewählt, und sagt, was das Entfernen bewirkt: `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` Bei einem Konto mit API-Schlüssel gibt es keine Anmeldung zum Widerrufen, und wenn ein anderes Konto dieselbe Anmeldung hat, wird es nicht abgemeldet. Ein importiertes Konto behält seinen Ordner: `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Farben

Die Farben folgen den Rollen der Terminal-Oberfläche von Claude Code, sodass ccam daneben vertraut aussieht; Gruppennamen erhalten eine Farbe, die sich aus ihrem Namen ergibt. ccam fragt das Terminal nach seinem Hintergrund und wählt ein helles oder dunkles Design. Du kannst eines selbst am Anfang von `config.toml` festlegen, vor dem ersten `[[group]]`, `[[profile]]` oder `[[folder]]`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Die Designs mit `-daltonized` sind für Rot-Grün-Farbenblindheit gedacht, und die Designs mit `-ansi` verwenden die 16 eigenen Farben des Terminals. `NO_COLOR` schaltet die Farben aus.

## Sicherheit und Sonderfälle

- Das Dashboard ruft dieselben Funktionen auf wie die Befehle, daher gilt jede Regel: Pläne, bevor Daten bewegt werden, Ablehnungen bei laufenden Konten, Sicherungen und Anmeldetokens, die ein Konto nie verlassen.
- Das Dashboard hat keinen Probelauf: `ccam --dry-run` ohne Namen gibt aus, wie das Konto dieses Ordners starten würde, und `o` bei einem Konto bietet einen Probelauf genau dieses einen Starts an.
- Das Terminal wird wiederhergestellt, bevor Claude Code startet oder sich die Anmeldung öffnet.
- Bei einem Bildschirm, der für einen Dialog zu klein ist, fallen zuerst die Zeilen des Tabs weg; Zeilen werden mit `…` abgeschnitten statt umbrochen. Auf einem niedrigen Bildschirm fallen als Nächstes die Leerzeilen und die Tab-Beschreibung weg, sodass die Tastenzeile bleibt. Unter 40 Spalten oder 8 Zeilen zeigt das Dashboard nur `Too small: ccam needs 40x8 (now 36x7)`.
- `Ctrl+C`, während eine Änderung, ein Entfernen, ein Update oder eine Reparatur angewendet wird, wartet darauf: Die Statuszeile sagt `ccam quits when it is done.`
- Nachdem `U` ccam aktualisiert hat, kommt das Angebot nicht wieder, wenn das Dashboard vom Hinzufügen eines Kontos oder einer Anmeldung zurückkehrt.

## Verwandte Themen

- [Teilen mit Gruppen](/features/share)
- [Claude Code starten](/features/launch)
- [Konten verwalten](/features/accounts)
- [Nutzungslimits](/features/usage)
