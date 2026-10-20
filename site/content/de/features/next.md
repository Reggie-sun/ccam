---
title: Zum nächsten Konto wechseln
description: ccam next startet Claude Code als anderes Konto; --auto wählt das mit dem meisten Spielraum, --continue nimmt das Gespräch mit, --fallback wiederholt -p-Jobs.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: "6a5dcb1d5c01"
---

`ccam next` öffnet Claude Code mit einem anderen Konto, ohne dass du es benennen musst. Verwende es, wenn ein Konto sein Nutzungslimit erreicht: Einfaches `ccam next` nimmt das Konto nach dem zuletzt verwendeten, `--auto` nimmt das mit dem meisten Spielraum, und `--continue` nimmt das letzte Gespräch dieses Ordners in das neue Konto mit, ob die beiden den Verlauf teilen oder nicht. In Claude Code meldet dir das [ccam-Plugin](#inside-claude-code), wenn ein Limit eine Sitzung anhält, und `/ccam:switch` gibt das Gespräch weiter. Für Skripte, die `claude -p` ausführen, erledigt `ccam --fallback` den Wechsel von selbst (siehe [Am Limit erneut ausführen](#run-again-at-the-limit)).

## Verwendung

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam gibt den Wechsel als `→ <from> → <to>` auf der Standardfehlerausgabe aus, bei einer Wahl mit `--auto` mit dem Grund, und startet dann Claude Code mit dem neuen Konto auf dieselbe Weise wie `ccam <name>`. `<from>` ist das zuletzt verwendete Konto oder mit `--continue` das Konto, zu dem das Gespräch gehört. Argumente nach `--` gehen unverändert an `claude`.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--auto` | Das angemeldete Konto nach Nutzung wählen (siehe unten) und Konten bei 100 % überspringen | aus |
| `--max-usage <1-100>` | Mit `--auto`: Konten unter dieser Nutzung bevorzugen | `auto_max_usage`, sonst 100 |
| `--strategy headroom\|expiring\|balance` | Mit `--auto`: niedrigste 5-Stunden-Nutzung, die früheste 5-Stunden-Zurücksetzung oder das Wochenkontingent, das als Erstes ungenutzt verfallen würde | `auto_strategy`, sonst `headroom` |
| `--cross-org` | Mit `--auto`: auch Konten einer anderen Organisation wählen (siehe unten) | `cross_org`, sonst aus |
| `--wait[=<duration>]` | Mit `--auto`: Wenn kein Konto Spielraum hat, auf die Zurücksetzung des ersten warten, höchstens so lange | aus; `5h`, wenn ohne Wert angegeben |
| `-c`, `--continue` | Das letzte Gespräch dieses Ordners im neuen Konto weiterführen (siehe unten) | aus |
| `-- <claude args>` | Argumente, die nach ccams eigenen an `claude` übergeben werden | keine |
| `--dry-run` | Wird beachtet: gibt `→ <from> → <to>` und die Befehlszeile aus, die ausgeführt würde, ohne Claude Code zu starten oder einen Start zu erfassen | aus |

`--auto` und `--continue` lassen sich kombinieren. Mit `--continue` werden die `--resume`-Argumente vor deine eigenen gesetzt.

## So funktioniert es

Jeder Start über ccam schreibt die Uhrzeit nach `~/.local/state/ccam/last-used/<name>`. `ccam next` liest diese Dateien, prüft, welche Konten angemeldet sind, und wählt.

1. `cli.candidates` baut aus `config.toml` pro Konto einen `next.Candidate`: das Konto, seine letzte Startzeit (`launcher.LastUsed`) und ob es eine Anmeldung hat (`creds.Has`: eine lesbare `.credentials.json` oder unter macOS ein Eintrag im Schlüsselbund).
2. Ohne `--auto` sortiert `next.Rotate` die Konten nach Slot-Nummer und findet das zuletzt gestartete. Von diesem Konto wechselst du weg. Wurde noch kein Konto gestartet, ist es der niedrigste Slot.
3. Von dort geht es in Slot-Reihenfolge weiter, springt vom letzten Slot zurück zum ersten und gibt das erste angemeldete Konto zurück.
4. `launcher.LaunchWith` bereitet das gewählte Konto genau so vor und startet es, wie es `ccam <name>` tun würde: Es prüft die Links des Kontos, wendet das MCP-Teilen an, registriert die Sitzung und ersetzt ccam über `syscall.Exec` durch `claude`.

Der Durchlauf beginnt beim zuletzt verwendeten Konto und endet beim ersten angemeldeten Konto danach:

<!-- motion:next-rotate -->

`work` wurde zuletzt gestartet, daher beginnt der Durchlauf danach. `side` hat keine Anmeldung und wird übersprungen, und `lab` wird gewählt. Wäre auch `lab` abgemeldet, würde der Durchlauf zu `personal` zurückspringen (`% len` im Diagramm). Die Rotation gibt nie das Konto zurück, von dem du gestartet bist: Gibt es kein anderes angemeldetes Konto, endet sie mit einem Fehler.

### --auto

`--auto` entscheidet nach Nutzung statt nach Slot-Reihenfolge. `cli.pickAuto` behält die angemeldeten Konten, fragt ihre Nutzung mit demselben `Fetcher.FetchAll`, das `ccam usage` verwendet, gleichzeitig ab und lässt `usage.Policy.Choose()` wählen:

<!-- motion:next-auto -->

Ein Konto wird nach seinem 5-Stunden- und seinem Wochenfenster beurteilt und, wenn `--model` unter den Argumenten für Claude Code steht, auch nach dem Wochenfenster dieses Modells (`opus`, `sonnet`, `haiku` oder `fable` im Namen). `Policy.Choose()` teilt die Konten in drei Gruppen ein:

1. **Spielraum**: jedes Fenster unter `auto_max_usage` (`config.toml`, Standard 100). Die Strategie wählt unter ihnen: `headroom` (Standard) nimmt die niedrigste 5-Stunden-Nutzung; `expiring` nimmt das 5-Stunden-Fenster, das am frühesten zurückgesetzt wird, damit sein restliches Kontingent genutzt wird, bevor es verfällt; `balance` nimmt das Konto mit dem meisten verbleibenden Wochenkontingent pro Tag bis zur Zurücksetzung seiner Woche und rechnet jede laufende ccam-Sitzung dagegen, damit kein Wochenkontingent ungenutzt bleibt und sich parallele Sitzungen verteilen.
2. **Darüber**: unter 100 %, aber bei oder über `auto_max_usage`. Wird nur verwendet, wenn kein Konto Spielraum hat, die niedrigste 5-Stunden-Nutzung zuerst.
3. **Voll**: ein Fenster bei 100 % oder unbekannte 5-Stunden-Nutzung. Wird nie gewählt.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organisationen.** Automatische Wahlen bleiben bei Konten wie dem, das du verlässt: dieselbe Team- oder Enterprise-Organisation oder private Tarife (Pro, Max) untereinander, sodass Arbeit aus einem Firmenkonto nicht in einem privaten weitergeht oder umgekehrt. ccam kennt die Organisation jedes Kontos aus seiner Anmeldung (`ccam login` hält sie fest; ältere Anmeldungen werden aus `.claude.json` gelesen); ein Konto, bei dem es das nicht erkennen kann, passt überall. Zwei Ordner, die beim selben Claude-Konto angemeldet sind, teilen sich dessen Limits, daher übernimmt nie einer vom anderen. `--cross-org` oder `cross_org = true` in `config.toml` hebt die Organisationsregel auf.

Wenn nichts gewählt werden kann, beendet sich ccam mit Code 75 (`EX_TEMPFAIL`), damit ein Skript „später erneut versuchen“ von einem Fehler unterscheiden kann. Mit `--wait` schläft es stattdessen, bis das blockierende Fenster des ersten vollen Kontos zurückgesetzt wird (30 Sekunden nach dem Zurücksetzungszeitpunkt, wenn der Endpunkt das neue Fenster zeigt), gibt `ccam: no account has room; waiting until 15:10` aus und wählt erneut; wenn kein Konto angibt, wann es zurückgesetzt wird, sieht es nach 10 Minuten wieder nach. Wird vor dem Ende der Wartezeit nichts frei, beendet es sich mit 75. Die Abfrage folgt den Regeln von `ccam usage`: Zahlen, die jünger als 60 Sekunden sind, werden wiederverwendet, jedes Konto hat 3 Sekunden, und ein Konto, das nicht antwortet oder dessen Anmeldung abgelaufen ist, wird nach seinen gespeicherten Zahlen beurteilt. Ein gespeichertes Zeitfenster, dessen Zurücksetzungszeitpunkt vorbei ist, zählt als 0 %, damit ein Konto, dessen Fenster inzwischen zurückgesetzt wurde, nicht übergangen wird. `--auto` erneuert nie Tokens. Anders als die einfache Rotation kann `--auto` das zuletzt verwendete Konto wählen, wenn es noch den meisten Spielraum hat.

### --continue

1. ccam sucht in jedem Konto nach dem letzten Gespräch dieses Ordners: der neuesten `<session>.jsonl` in `projects/<folder>/` (der Ordnername ist der Arbeitsordner, in dem jedes Zeichen außer Buchstaben und Ziffern zu `-` wird, so wie Claude Code ihn benennt). Hat die Sitzung eine Notiz hinterlassen (ein Nutzungslimit oder `/ccam:switch`, siehe unten) und wurde seitdem nichts Neueres geschrieben, wird stattdessen genau diese Sitzung verwendet.
2. Das Konto, zu dem es gehört, ist das, von dem du wechselst, und das neue Konto wird wie oben gewählt, innerhalb derselben Organisation.
3. Teilt das neue Konto diesen Verlauf (derselbe Gruppenordner), startet ccam `claude --resume <session id>`, also genau dieses Gespräch (`--continue` würde das jeweils neueste Gespräch im Ordner öffnen).
4. Andernfalls startet es `claude --resume <the conversation's file> --fork-session`: Claude Code kopiert das ganze Gespräch unter einer neuen Sitzungs-ID in den eigenen Verlauf des neuen Kontos, und das Original bleibt, wie es war. Von da an findet `--continue` des neuen Kontos es.

ccam sagt, wie groß das Gespräch ist: Das neue Konto liest es mit seiner ersten Antwort komplett erneut. Prompt-Caches werden nie zwischen Organisationen geteilt, daher kostet diese erste Antwort mit `--cross-org` die volle Größe. Wenn eine Gruppe `history` teilt, entfällt die Kopie, es bleibt ein Gespräch, und auch die Checkpoints, aus denen `/rewind` Code wiederherstellt (`file-history/`), werden geteilt.

## In Claude Code

Das [ccam-Plugin](/docs/install#claude-code-plugin) fügt zwei Hooks hinzu (Claude Code 2.1.78 oder neuer):

- Wenn eine Runde an einem Nutzungslimit anhält (`You've hit your session limit`, das Wochen- oder ein Modelllimit), merkt es das Gespräch vor und zeigt in Terminals, die das unterstützen (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty), eine Desktop-Benachrichtigung: `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` Eine kurze Drosselung durch den Server oder das Ratenlimit eines API-Schlüssels ist kein Nutzungslimit und bleibt unbeachtet. Claude Code 2.1.234 und neuer wartet außerdem in der Sitzung und macht nach einer 5-Stunden-Zurücksetzung von selbst weiter; die Benachrichtigung gibt dir die Wahl, stattdessen zu wechseln.
- `/ccam:switch [account]` merkt das Gespräch vor und beendet die Runde, bevor eine Anfrage gestellt wird: Tippe `/exit`, dann führt `ccam next --continue` es weiter, mit dem Konto, das du genannt hast, oder mit dem nächsten.

Die Notiz liegt 12 Stunden lang in `~/.local/state/ccam/handoff/<account>.json` und wird einmal verwendet, wenn die nächste Sitzung gerade starten soll; ein Start, der fehlschlägt, lässt sie für den nächsten Versuch liegen. Sie gehört zu dem Ordner, in dem die Sitzung gestartet wurde, auch nach einem `cd` in der Shell von Claude Code. Wurde das Konto, das sie nennt, inzwischen umbenannt oder entfernt, sagt ccam das und wählt wie gewohnt. Mit `relay = true` in `config.toml` bleibt ccam neben interaktiven Sitzungen, statt sich selbst durch Claude Code zu ersetzen, und führt nach `/ccam:switch` das Gespräch dieser Sitzung von selbst weiter (auf eine Notiz, die eine andere Sitzung hinterlassen hat, reagiert es nicht). Das kostet einen wartenden ccam-Prozess pro Sitzung, deshalb ist es standardmäßig aus. Die weitergeführte Sitzung startet mit `ccam next --continue`, ohne die Argumente, mit denen die erste gestartet wurde, etwa `--model`.

## Am Limit erneut ausführen

`ccam next` ist für eine interaktive Sitzung gedacht, die angehalten hat. Ein nicht interaktiver Durchlauf, `claude -p` in einem Skript oder einem geplanten Job, kann nicht auf dich warten, deshalb lässt `--fallback` vor dem Kontonamen ccam für ihn wechseln:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` funktioniert nur zusammen mit `-p` oder `--print`; andernfalls hält ccam an und verweist auf `ccam next --auto --continue`.
2. Die Standardeingabe wird einmal gelesen und jedem Versuch übergeben.
3. ccam startet `claude` als Kindprozess statt mit `exec`, hält seine Standardausgabe zurück und gibt die Standardfehlerausgabe sofort weiter.
4. Wenn `claude` mit 0 endet oder aus einem anderen Grund fehlschlägt, werden die zurückgehaltene Ausgabe und der Exit-Code unverändert weitergegeben.
5. Ein Fehlschlag zählt als Nutzungslimit, wenn die Meldung von Claude Code das sagt (`You've hit your limit`, `usage limit reached`) oder wenn die Nutzungsabfrage des Kontos sein 5-Stunden- oder Wochenfenster bei 100 % zeigt.
6. Dann wird die Ausgabe dieses Versuchs verworfen, und derselbe Befehl läuft mit einem anderen angemeldeten Konto: dem mit der niedrigsten 5-Stunden-Nutzung, danach Konten, deren Nutzung sich nicht abfragen lässt, in der Reihenfolge ihrer Nummern. Konten bei 100 % werden übersprungen, und jedes Konto wird höchstens einmal versucht.
7. Mit `-c`, `--continue`, `-r` oder `--resume` unter den Argumenten werden nur Konten in derselben Verlaufsgruppe versucht, damit das Gespräch gefunden wird.
8. Es werden nur Konten derselben Organisation versucht, wie bei `--auto` (`cross_org = true` hebt das auf).

Das nächste Konto wird nach denselben Regeln wie bei `--auto` gewählt (`auto_max_usage`, `auto_strategy`, das `--model`-Fenster). Wenn kein Konto mehr übrig ist, gibt ccam die Ausgabe des letzten Versuchs aus, meldet, dass kein angemeldetes Konto mehr Spielraum hat, und beendet sich mit Code 75. Mit `--wait` vor dem Kontonamen (`ccam --fallback --wait=2h work -p …`) wartet es stattdessen, bis das erste Konto zurückgesetzt wird, und kann dann jedes Konto erneut versuchen; es schläft immer zuerst mindestens einmal, damit ein Konto, das gerade sein Limit erreicht hat, nicht sofort wieder läuft. `ccam --dry-run --fallback work -p …` gibt den Befehl des ersten Versuchs aus.

Interaktive Sitzungen werden nicht von selbst gewechselt. Stattdessen ergänzt die Statuszeile mit `ccam statusline install` den Hinweis `→ ccam next --auto --continue`, sobald die 5-Stunden- oder Wochennutzung 90 % erreicht ([Nutzungslimits](/features/usage)).

## Sicherheit und Sonderfälle

- Konten ohne Anmeldung werden nie gewählt. Gibt es kein anderes angemeldetes Konto, hält einfaches `ccam next` mit `no logged-in account other than <name>` an.
- Mit `--continue` und keinem Gespräch im Ordner: `no account has a conversation in <folder> to continue`.
- Wenn die Organisationsregel kein Konto übrig lässt: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Mit `--auto` und keinem angemeldeten Konto: `no logged-in account to pick`. Wenn jedes Konto voll oder seine Nutzung unbekannt ist: `` no account with known usage and headroom (see `ccam usage`) ``. Beide enden mit Exit-Code 75.
- `--max-usage`, `--strategy`, `--wait` und `--cross-org` funktionieren nur mit `--auto`, und ein Wert außerhalb von 1-100 oder eine andere Strategie wird abgelehnt.
- Ein Fork kopiert das Gespräch, nicht den übrigen Zustand des Kontos: Checkpoints aus `file-history/` des ursprünglichen Kontos kommen nicht mit, es sei denn, der Verlauf wird geteilt.
- `--auto` stellt auch mit `--dry-run` Netzwerkanfragen. Es liest nur und schreibt nie Token-Dateien.
- Das gewählte Konto startet über denselben Pfad wie `ccam <name>`, daher gelten dieselben Prüfungen, einschließlich der Warnung, wenn in diesem Ordner bereits ein anderes Claude Code läuft, das denselben Verlauf teilt.
- `last-used` wird nur geschrieben, wenn Claude Code tatsächlich startet. `--dry-run` lässt es unverändert, daher rückt ein Probelauf die Rotation nicht weiter.
- `--fallback` hält die Standardausgabe zurück, bis ein Versuch endet, daher kommt eine fortlaufende Ausgabe wie `--output-format stream-json` am Ende in einem Stück an. Dateien, die der erste Versuch vor dem Erreichen des Limits geändert hat, bleiben geändert; der nächste Versuch beginnt den Prompt von vorn.

## Verwandte Themen

- [Nutzungslimits](/features/usage)
- [Nutzungslimits und Kontowechsel](/docs/usage-limits)
- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Umgang mit Nutzungslimits](/guides/usage-limit-workflow)
