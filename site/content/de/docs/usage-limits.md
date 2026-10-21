---
title: Nutzungslimits und Kontowechsel
description: Die 5-Stunden-, Wochen- und Fable-Nutzung aller Claude Code-Konten mit ccam usage als Balken sehen und mit ccam next --auto das Konto mit dem meisten Spielraum öffnen.
updated: 2026-10-21
order: 5
translation: machine
sourceHash: 0e64e3251c69
---

`ccam usage` zeigt, wie viel jedes Konto von seinem 5-Stunden-Limit, seinem Wochenlimit und dem eigenen Wochenlimit jedes Modells (etwa Fable) verbraucht hat. `ccam next --auto` öffnet das Konto mit dem meisten Spielraum, und `--continue` nimmt das aktuelle Gespräch mit.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Jedes Konto bekommt eine Zeile, die zeigt, wie alt seine Zahlen sind, darunter einen Balken pro Limit mit der Zeit bis zum Zurücksetzen: das 5-Stunden-Fenster, das Wochenfenster für alle Modelle und das eigene Wochenfenster jedes Modells, etwa Fable. Ab 80 % wird ein Balken gelb, bei 100 % rot. Die letzte Zeile nennt das Konto, das `ccam next --auto` öffnen würde.

## Woher die Zahlen kommen

ccam fragt den Nutzungs-Endpunkt von Anthropic für alle Konten gleichzeitig ab, jeweils mit der eigenen Anmeldung des Kontos. Es wartet höchstens 3 Sekunden pro Konto und zeigt jedes Konto an, sobald es antwortet.

Dieser Endpunkt ist nicht öffentlich dokumentiert und kann sich ändern. Wenn eine Anfrage fehlschlägt, zeigt ccam die zuletzt bekannten Zahlen und markiert sie als zwischengespeichert. Ein Zeitfenster, das seit dem Speichern dieser Zahlen zurückgesetzt wurde, zeigt 0 %. `ccam doctor` sagt dir, ob sich das Antwortformat geändert hat.

Zahlen, die vor weniger als 60 Sekunden abgerufen wurden, werden wiederverwendet. Optionen:

| Option | Wirkung |
| --- | --- |
| `--watch [seconds]` | Laufend aktualisieren, standardmäßig alle 60 Sekunden |
| `--json` | Maschinenlesbare Ausgabe |
| `-v` | Zusätzlich zeigen, woher die Zahlen jedes Kontos stammen |
| `--no-cache` | Immer neu abfragen |
| `--source cache` | Nie abfragen; nur gespeicherte Zahlen zeigen |

## Abgelaufene Anmeldungen

Wenn das Zugriffstoken eines Kontos abgelaufen ist, erneuert ccam es nicht, denn eine Erneuerung, während eine Claude Code-Sitzung dieselbe Anmeldung nutzt, könnte diese Sitzung abmelden. Starte das Konto einmal mit `ccam <name>`, dann erneuert Claude Code es.

`ccam usage --refresh` erneuert abgelaufene Tokens selbst, aber nur für Konten, die nicht laufen. Anmeldungen im macOS-Schlüsselbund oder in einem importierten `~/.claude` erneuert es nie.

## Nutzung in der Statuszeile

```sh
ccam statusline install
```

Claude Code zeigt dann den Kontonamen und seine Nutzung in der Statuszeile an, zum Beispiel `[work] · Opus · 5h 24% · 7d 41%`. Ab 90 % schlägt die Zeile außerdem `ccam next --auto --continue` vor. Wenn du bereits einen Befehl für die Statuszeile hast, behält ccam ihn bei und stellt den Kontonamen voran. Die Statuszeile speichert außerdem ihre Zahlen, die `ccam usage` verwendet, wenn der Endpunkt nicht antwortet.

## Wechseln

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # carry this folder's conversation on as the next account
```

`ccam next` überspringt Konten, die nicht angemeldet sind, und `--auto` überspringt Konten bei 100 % sowie Konten einer anderen Organisation. Wenn die beiden Konten den Gesprächsverlauf teilen, öffnet `--continue` dasselbe Gespräch; andernfalls kopiert Claude Code es in das neue Konto. Siehe [Teilen](/docs/sharing) und [zum nächsten Konto wechseln](/features/next#continue).

Für Skripte und geplante Jobs, die `claude -p` ausführen, setze `--fallback` vor den Kontonamen. Wenn dieses Konto sein Limit erreicht, führt ccam denselben Befehl erneut mit dem Konto mit dem meisten Spielraum aus, und nur die Ausgabe des Durchlaufs, der fertig wurde, wird ausgegeben:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Wie ccam entscheidet, steht unter [Am Limit erneut ausführen](/features/next#run-again-at-the-limit).

Die [Anleitung zum Umgang mit Nutzungslimits](/guides/usage-limit-workflow) geht ein vollständiges Beispiel durch.
