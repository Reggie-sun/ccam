---
title: Zum nächsten Konto wechseln
description: ccam next öffnet Claude Code mit dem Konto nach dem letzten; --auto wählt das mit dem meisten Spielraum, --continue setzt fort, --fallback wiederholt -p-Jobs.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: 784ae4d1f6d9
---

`ccam next` öffnet Claude Code mit einem anderen Konto, ohne dass du es benennen musst. Verwende es, wenn ein Konto sein Nutzungslimit erreicht: Einfaches `ccam next` nimmt das Konto nach dem zuletzt verwendeten, `--auto` nimmt das mit dem meisten Spielraum im 5-Stunden-Fenster, und `--continue` setzt das letzte Gespräch im neuen Konto fort. Für Skripte, die `claude -p` ausführen, erledigt `ccam --fallback` diesen Wechsel von selbst (siehe [Am Limit erneut ausführen](#run-again-at-the-limit)).

## Verwendung

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam gibt den Wechsel als `→ <from> → <to>` auf der Standardfehlerausgabe aus und startet dann Claude Code mit dem neuen Konto auf dieselbe Weise wie `ccam <name>`. `<from>` ist das zuletzt verwendete Konto. Argumente nach `--` gehen unverändert an `claude`.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--auto` | Das angemeldete Konto nach Nutzung wählen (siehe unten) und Konten bei 100 % überspringen | aus |
| `--max-usage <1-100>` | Mit `--auto`: Konten unter dieser Nutzung bevorzugen | `auto_max_usage`, sonst 100 |
| `--strategy headroom\|expiring` | Mit `--auto`: niedrigste 5-Stunden-Nutzung oder die früheste 5-Stunden-Zurücksetzung | `auto_strategy`, sonst `headroom` |
| `--wait[=<duration>]` | Mit `--auto`: Wenn kein Konto Spielraum hat, auf die Zurücksetzung des ersten warten, höchstens so lange | aus; `5h`, wenn ohne Wert angegeben |
| `-c`, `--continue` | Nur Konten berücksichtigen, die den Gesprächsverlauf mit dem zuletzt verwendeten Konto teilen, und `claude --continue` starten | aus |
| `-- <claude args>` | Argumente, die nach ccams eigenen an `claude` übergeben werden | keine |
| `--dry-run` | Wird beachtet: gibt `→ <from> → <to>` und die Befehlszeile aus, die ausgeführt würde, ohne Claude Code zu starten oder einen Start zu erfassen | aus |

`--auto` und `--continue` lassen sich kombinieren. Mit `--continue` wird `--continue` vor deine eigenen Argumente gesetzt.

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

1. **Spielraum**: jedes Fenster unter `auto_max_usage` (`config.toml`, Standard 100). Die Strategie wählt unter ihnen: `headroom` (Standard) nimmt die niedrigste 5-Stunden-Nutzung, `expiring` nimmt das 5-Stunden-Fenster, das am frühesten zurückgesetzt wird, damit sein restliches Kontingent genutzt wird, bevor es verfällt.
2. **Darüber**: unter 100 %, aber bei oder über `auto_max_usage`. Wird nur verwendet, wenn kein Konto Spielraum hat, die niedrigste 5-Stunden-Nutzung zuerst.
3. **Voll**: ein Fenster bei 100 % oder unbekannte 5-Stunden-Nutzung. Wird nie gewählt.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

Wenn nichts gewählt werden kann, beendet sich ccam mit Code 75 (`EX_TEMPFAIL`), damit ein Skript „später erneut versuchen“ von einem Fehler unterscheiden kann. Mit `--wait` schläft es stattdessen, bis das blockierende Fenster des ersten vollen Kontos zurückgesetzt wird (30 Sekunden nach dem Zurücksetzungszeitpunkt, wenn der Endpunkt das neue Fenster zeigt), gibt `ccam: no account has room; waiting until 15:10` aus und wählt erneut; wenn kein Konto angibt, wann es zurückgesetzt wird, sieht es nach 10 Minuten wieder nach. Wird vor dem Ende der Wartezeit nichts frei, beendet es sich mit 75. Die Abfrage folgt den Regeln von `ccam usage`: Zahlen, die jünger als 60 Sekunden sind, werden wiederverwendet, jedes Konto hat 3 Sekunden, und ein Konto, das nicht antwortet oder dessen Anmeldung abgelaufen ist, wird nach seinen gespeicherten Zahlen beurteilt. Ein gespeichertes Zeitfenster, dessen Zurücksetzungszeitpunkt vorbei ist, zählt als 0 %, damit ein Konto, dessen Fenster inzwischen zurückgesetzt wurde, nicht übergangen wird. `--auto` erneuert nie Tokens. Anders als die einfache Rotation kann `--auto` das zuletzt verwendete Konto wählen, wenn es noch den meisten Spielraum hat.

### --continue

`--continue` von Claude Code öffnet das neueste Gespräch im aktuellen Ordner, was nur funktioniert, wenn das neue Konto dieses Gespräch sehen kann. Deshalb berücksichtigt ccam mit `--continue` nur Konten in derselben `history`-Gruppe wie das zuletzt verwendete Konto und setzt dann `--continue` vor deine Argumente. Konten, deren Verlauf `local` ist, teilen mit niemandem, daher schlägt `--continue` von einem solchen Konto aus immer fehl. Schalte das Teilen zuerst mit `ccam share <name> --history on` ein (siehe [Teilen](/docs/sharing)).

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

Das nächste Konto wird nach denselben Regeln wie bei `--auto` gewählt (`auto_max_usage`, `auto_strategy`, das `--model`-Fenster). Wenn kein Konto mehr übrig ist, gibt ccam die Ausgabe des letzten Versuchs aus, meldet, dass kein angemeldetes Konto mehr Spielraum hat, und beendet sich mit Code 75. Mit `--wait` vor dem Kontonamen (`ccam --fallback --wait=2h work -p …`) wartet es stattdessen, bis das erste Konto zurückgesetzt wird, und kann dann jedes Konto erneut versuchen; es schläft immer zuerst mindestens einmal, damit ein Konto, das gerade sein Limit erreicht hat, nicht sofort wieder läuft. `ccam --dry-run --fallback work -p …` gibt den Befehl des ersten Versuchs aus.

Interaktive Sitzungen werden nicht von selbst gewechselt. Stattdessen ergänzt die Statuszeile mit `ccam statusline install` den Hinweis `→ ccam next --auto --continue`, sobald die 5-Stunden- oder Wochennutzung 90 % erreicht ([Nutzungslimits](/features/usage)).

## Sicherheit und Sonderfälle

- Konten ohne Anmeldung werden nie gewählt. Gibt es kein anderes angemeldetes Konto, hält einfaches `ccam next` mit `no logged-in account other than <name>` an.
- Mit `--continue` und keinem Konto, das den Verlauf teilt: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Mit `--auto` und keinem angemeldeten Konto: `no logged-in account to pick`. Wenn jedes Konto voll oder seine Nutzung unbekannt ist: `` no account with known usage and headroom (see `ccam usage`) ``. Beide enden mit Exit-Code 75.
- `--max-usage`, `--strategy` und `--wait` funktionieren nur mit `--auto`, und ein Wert außerhalb von 1-100 oder eine andere Strategie wird abgelehnt.
- `--auto` stellt auch mit `--dry-run` Netzwerkanfragen. Es liest nur und schreibt nie Token-Dateien.
- Das gewählte Konto startet über denselben Pfad wie `ccam <name>`, daher gelten dieselben Prüfungen, einschließlich der Warnung, wenn in diesem Ordner bereits ein anderes Claude Code läuft, das denselben Verlauf teilt.
- `last-used` wird nur geschrieben, wenn Claude Code tatsächlich startet. `--dry-run` lässt es unverändert, daher rückt ein Probelauf die Rotation nicht weiter.
- `--fallback` hält die Standardausgabe zurück, bis ein Versuch endet, daher kommt eine fortlaufende Ausgabe wie `--output-format stream-json` am Ende in einem Stück an. Dateien, die der erste Versuch vor dem Erreichen des Limits geändert hat, bleiben geändert; der nächste Versuch beginnt den Prompt von vorn.

## Verwandte Themen

- [Nutzungslimits](/features/usage)
- [Nutzungslimits und Kontowechsel](/docs/usage-limits)
- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Umgang mit Nutzungslimits](/guides/usage-limit-workflow)
