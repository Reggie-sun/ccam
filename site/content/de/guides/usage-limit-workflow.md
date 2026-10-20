---
title: Was tun, wenn Claude Code sein Nutzungslimit erreicht
description: Auf den Reset warten oder dasselbe Gespräch sofort in einem anderen Claude Code-Konto fortsetzen. So funktionieren das 5-Stunden- und das Wochenlimit.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: ff765a189332
---

Wenn ein Konto sein Nutzungslimit erreicht (oft Rate Limit genannt), hast du zwei Möglichkeiten: warten, bis das Limit zurückgesetzt wird, oder dasselbe Gespräch sofort in einem anderen Konto fortsetzen, das dir gehört. Mit ccam braucht das Zweite nur einen Befehl, `ccam next --auto --continue`, sobald der Gesprächsverlauf zwischen den Konten geteilt wird.

## Welche Limits hat Claude Code?

In den Pro- und Max-Tarifen hat Claude Code ein [5-Stunden-Sitzungslimit und ein Wochenlimit](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Beide gelten auch für Claude im Web und in den Apps, und die Nutzung zählt gleichzeitig auf beide. Wenn eines aufgebraucht ist, [zeigt Claude Code eine Meldung](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) wie „You've hit your session limit · resets 3:45pm“. Ein anderes Modell zu wählen hilft nicht, weil die Limits für alle Modelle gelten.

## Option 1: auf den Reset warten

Wenn du nur ein Konto hast, warte. Aktuelle Versionen von Claude Code [warten in der offenen Sitzung](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) und setzen die Aufgabe von selbst fort, sobald das Limit zurückgesetzt ist. `/usage` in Claude Code zeigt, wie viel du von den Limits deines Tarifs verbraucht hast.

## Option 2: sofort in einem anderen Konto weitermachen

Wenn du ein zweites Konto hast, kannst du direkt mit demselben Gespräch weitermachen. Richte das einmal ein:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Die Statuszeile zeigt jetzt beim Arbeiten das Konto und seine 5-Stunden-Nutzung, sodass dich ein Limit nicht überrascht.

Wenn das Limit erreicht ist, beende Claude Code (`/exit`) und führe dann aus:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

Öffne das Gespräch im Konto mit dem meisten Spielraum:

```sh
ccam next --auto --continue
```

Claude Code startet als `work` mit dem neuesten Gespräch des aktuellen Ordners, einschließlich allem, was bisher gesagt wurde.

## Wechselt ccam das Konto automatisch?

Nein. ccam läuft nicht, während du arbeitest: Es startet Claude Code und ersetzt sich selbst durch diesen Prozess. Du entscheidest, wann gewechselt wird, und `ccam next --auto` wählt das Konto mit der geringsten 5-Stunden-Nutzung und überspringt Konten, die bei 100 % oder nicht angemeldet sind. Manche anderen Tools tauschen stattdessen die Anmeldung einer laufenden Sitzung aus; [ccam und andere Tools zum Kontowechsel](/compare/account-switchers) erklärt den Unterschied.

## Was du beachten solltest

- Führe `ccam next --continue` im selben Ordner aus, in dem du gearbeitet hast; Gespräche werden pro Projektordner gespeichert.
- Arbeite nicht gleichzeitig mit zwei Konten am selben Gespräch. ccam warnt dich, wenn du es versuchst.
- Das andere Konto muss den Gesprächsverlauf ebenfalls teilen, in derselben Gruppe.
- Die erste Nachricht im neuen Konto kann mehr von dessen Limit verbrauchen als üblich, weil dieses Konto das Gespräch noch nicht im Cache hat.
- Ob du so zwischen Konten wechseln darfst, hängt von den Bedingungen des jeweiligen Kontos ab.

Alle Optionen von `ccam usage` und `ccam next` findest du unter [Nutzungslimits](/docs/usage-limits).
