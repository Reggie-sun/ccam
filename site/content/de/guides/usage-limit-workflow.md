---
title: Was tun, wenn Claude Code sein Nutzungslimit erreicht
description: Prüfen, wie viel Spielraum jedes Claude Code-Konto noch hat, und dasselbe Gespräch in einem anderen Konto fortsetzen, wenn eines sein Limit erreicht.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 767a731b912e
---

Wenn ein Konto sein Limit erreicht, setzt du dasselbe Gespräch in einem anderen Konto fort, das dir gehört: Teile den Gesprächsverlauf einmal und führe dann `ccam next --continue` aus. Wenn du nur ein Konto hast, warte auf den Zeitpunkt des Zurücksetzens, den Claude Code anzeigt.

## Bevor es passiert

Richte die Konten so ein, dass der Wechsel nur einen Befehl braucht:

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Die Statuszeile zeigt jetzt beim Arbeiten das Konto und seine 5-Stunden-Nutzung, sodass dich ein Limit nicht überrascht.

## Wenn das Limit erreicht ist

Beende Claude Code (`/exit`), dann:

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

## Worauf du achten solltest

- Führe `ccam next --continue` im selben Ordner aus, in dem du gearbeitet hast; Gespräche werden pro Projektordner gespeichert.
- Vermeide es, zwei Konten gleichzeitig mit demselben Gespräch laufen zu lassen. ccam warnt dich, wenn du es versuchst.
- Das andere Konto muss den Gesprächsverlauf ebenfalls teilen, in derselben Gruppe.
- Ob du so zwischen Konten wechseln darfst, hängt von den Bedingungen jedes Kontos ab.

Alle Optionen von `ccam usage` und `ccam next` findest du unter [Nutzungslimits](/docs/usage-limits).
