---
title: Zum nächsten Konto wechseln
description: ccam next öffnet Claude Code mit dem Konto nach dem zuletzt verwendeten; --auto wählt die niedrigste 5-Stunden-Nutzung, und --continue setzt das Gespräch fort.
updated: 2026-10-19
order: 8
translation: machine
sourceHash: 92352f532e74
---

`ccam next` öffnet Claude Code mit einem anderen Konto, ohne dass du es benennen musst. Verwende es, wenn ein Konto sein Nutzungslimit erreicht: Einfaches `ccam next` nimmt das Konto nach dem zuletzt verwendeten, `--auto` nimmt das mit dem meisten Spielraum im 5-Stunden-Fenster, und `--continue` setzt das letzte Gespräch im neuen Konto fort.

## Verwendung

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
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
| `--auto` | Das angemeldete Konto mit der niedrigsten 5-Stunden-Nutzung wählen und Konten bei 100 % überspringen | aus |
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

`--auto` entscheidet nach Nutzung statt nach Slot-Reihenfolge. `cli.pickAuto` behält die angemeldeten Konten, fragt ihre Nutzung mit demselben `Fetcher.FetchAll`, das `ccam usage` verwendet, gleichzeitig ab und lässt `usage.Best()` wählen:

<!-- motion:next-auto -->

`usage.Best()` nimmt den niedrigsten 5-Stunden-Prozentsatz und überspringt Konten bei 100 % oder mit unbekannter 5-Stunden-Nutzung. Die Abfrage folgt den Regeln von `ccam usage`: Zahlen, die jünger als 60 Sekunden sind, werden wiederverwendet, jedes Konto hat 3 Sekunden, und ein Konto, das nicht antwortet oder dessen Anmeldung abgelaufen ist, wird nach seinen gespeicherten Zahlen beurteilt. Ein gespeichertes Zeitfenster, dessen Zurücksetzungszeitpunkt vorbei ist, zählt als 0 %, damit ein Konto, dessen Fenster inzwischen zurückgesetzt wurde, nicht übergangen wird. `--auto` erneuert nie Tokens. Anders als die einfache Rotation kann `--auto` das zuletzt verwendete Konto wählen, wenn es noch den meisten Spielraum hat.

### --continue

`--continue` von Claude Code öffnet das neueste Gespräch im aktuellen Ordner, was nur funktioniert, wenn das neue Konto dieses Gespräch sehen kann. Deshalb berücksichtigt ccam mit `--continue` nur Konten in derselben `history`-Gruppe wie das zuletzt verwendete Konto und setzt dann `--continue` vor deine Argumente. Konten, deren Verlauf `local` ist, teilen mit niemandem, daher schlägt `--continue` von einem solchen Konto aus immer fehl. Schalte das Teilen zuerst mit `ccam share <name> --history on` ein (siehe [Teilen](/docs/sharing)).

## Sicherheit und Sonderfälle

- Konten ohne Anmeldung werden nie gewählt. Gibt es kein anderes angemeldetes Konto, hält einfaches `ccam next` mit `no logged-in account other than <name>` an.
- Mit `--continue` und keinem Konto, das den Verlauf teilt: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Mit `--auto` und keinem angemeldeten Konto: `no logged-in account to pick`. Wenn jedes Konto voll oder seine Nutzung unbekannt ist: `` no account with known usage and headroom (see `ccam usage`) ``.
- `--auto` stellt auch mit `--dry-run` Netzwerkanfragen. Es liest nur und schreibt nie Token-Dateien.
- Das gewählte Konto startet über denselben Pfad wie `ccam <name>`, daher gelten dieselben Prüfungen, einschließlich der Warnung, wenn in diesem Ordner bereits ein anderes Claude Code läuft, das denselben Verlauf teilt.
- `last-used` wird nur geschrieben, wenn Claude Code tatsächlich startet. `--dry-run` lässt es unverändert, daher rückt ein Probelauf die Rotation nicht weiter.

## Verwandte Themen

- [Nutzungslimits](/features/usage)
- [Nutzungslimits und Kontowechsel](/docs/usage-limits)
- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Umgang mit Nutzungslimits](/guides/usage-limit-workflow)
