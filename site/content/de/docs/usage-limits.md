---
title: Nutzungslimits und Kontowechsel
description: Die 5-Stunden- und Wochennutzung aller Claude Code-Konten mit ccs usage sehen und mit ccs next --auto das Konto mit dem meisten Spielraum öffnen.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: a9ee77936668
---

`ccs usage` zeigt, wie viel jedes Konto von seinem 5-Stunden- und seinem Wochenlimit verbraucht hat. `ccs next --auto` öffnet das Konto mit dem meisten Spielraum, und `--continue` nimmt das aktuelle Gespräch mit.

```console
$ ccs usage
 계정        5시간         초기화      주간          초기화      기준
 1 personal  ██████░░  72% 1시간 40분  ███░░░░░  41% 3일 6시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
 3 side      ████████ 100% 22분        █████░░░  68% 1일 9시간   12분 전 (캐시)
→ 여유가 가장 많은 계정: work   (ccs next --auto 로 바로 열기)
```

Die Spalten sind: Konto, 5-Stunden-Nutzung, Zeit bis zum Zurücksetzen, Wochennutzung, Zeit bis zum Zurücksetzen und wie alt die Zahlen sind.

## Woher die Zahlen kommen

ccam fragt den Nutzungs-Endpunkt von Anthropic für alle Konten gleichzeitig ab, jeweils mit der eigenen Anmeldung des Kontos. Es wartet höchstens 3 Sekunden pro Konto und zeigt jedes Konto an, sobald es antwortet.

Dieser Endpunkt ist nicht öffentlich dokumentiert und kann sich ändern. Wenn eine Anfrage fehlschlägt, zeigt ccam die zuletzt bekannten Zahlen und markiert sie als zwischengespeichert. `ccam doctor` sagt dir, ob sich das Antwortformat geändert hat.

Zahlen, die vor weniger als 60 Sekunden abgerufen wurden, werden wiederverwendet. Optionen:

| Option | Wirkung |
| --- | --- |
| `--watch [seconds]` | Laufend aktualisieren, standardmäßig alle 60 Sekunden |
| `--json` | Maschinenlesbare Ausgabe |
| `-v` | Zusätzlich modellspezifische Wochenlimits zeigen und woher jede Zahl stammt |
| `--no-cache` | Immer neu abfragen |
| `--source cache` | Nie abfragen; nur gespeicherte Zahlen zeigen |

## Abgelaufene Anmeldungen

Wenn das Zugriffstoken eines Kontos abgelaufen ist, erneuert ccam es nicht, denn eine Erneuerung, während eine Claude Code-Sitzung dieselbe Anmeldung nutzt, könnte diese Sitzung abmelden. Starte das Konto einmal mit `ccs <name>`, dann erneuert Claude Code es.

`ccs usage --refresh` erneuert abgelaufene Tokens selbst, aber nur für Konten, die nicht laufen. Anmeldungen im macOS-Schlüsselbund oder in einem importierten `~/.claude` erneuert es nie.

## Nutzung in der Statuszeile

```sh
ccam statusline install
```

Claude Code zeigt dann den Kontonamen und seine Nutzung in der Statuszeile an, zum Beispiel `[work] · Opus · 5h 24% · 7d 41%`. Wenn du bereits einen Befehl für die Statuszeile hast, behält ccam ihn bei und stellt den Kontonamen voran. Die Statuszeile speichert außerdem ihre Zahlen, die `ccs usage` verwendet, wenn der Endpunkt nicht antwortet.

## Wechseln

```sh
ccs next              # the account after the one you used last
ccs next --auto       # the account with the lowest 5-hour usage
ccs next --continue   # the next account that shares your conversation history
```

`ccs next` überspringt Konten, die nicht angemeldet sind, und `--auto` überspringt Konten bei 100 %. Um ein Gespräch fortzusetzen, müssen beide Konten den Gesprächsverlauf teilen; siehe [Teilen](/docs/sharing).

Die [Anleitung zum Umgang mit Nutzungslimits](/guides/usage-limit-workflow) geht ein vollständiges Beispiel durch.
