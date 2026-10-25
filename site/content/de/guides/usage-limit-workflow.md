---
title: Was tun, wenn Claude Code sein Nutzungslimit erreicht
description: Auf den Reset warten oder dasselbe Gespräch sofort in einem anderen Claude Code-Konto fortsetzen. So funktionieren das 5-Stunden- und das Wochenlimit.
updated: 2026-10-25
order: 3
translation: machine
sourceHash: "18ef61a68c55"
---

Wenn ein Konto sein Nutzungslimit erreicht (oft Rate Limit genannt), hast du zwei Möglichkeiten: warten, bis das Limit zurückgesetzt wird, oder dasselbe Gespräch sofort in einem anderen Konto fortsetzen, das dir gehört. Mit ccam braucht das Zweite nur einen Befehl, `ccam next --auto --continue`.

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
ccam group set default --add history
ccam statusline install
```

Die Statuszeile zeigt jetzt beim Arbeiten das Konto und seine 5-Stunden-Nutzung, sodass dich ein Limit nicht überrascht.

Wenn das Limit erreicht ist, beende Claude Code (`/exit`) und führe dann aus:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Öffne das Gespräch im Konto mit dem meisten Spielraum:

```sh
ccam next --auto --continue
```

Claude Code startet als `work` mit dem neuesten Gespräch des aktuellen Ordners, einschließlich allem, was bisher gesagt wurde. Geteilter Verlauf (`ccam group set default --add history` oben) hält es bei einem Gespräch; ohne ihn kopiert Claude Code das Gespräch nach `work`, und ccam sagt dir, wie viele Tokens diese erste Antwort erneut liest.

Mit dem [ccam-Plugin](/docs/install#claude-code-plugin) in Claude Code zeigt eine Sitzung, die an ihrem Limit anhält, eine Benachrichtigung und vermerkt das volle Fenster, sodass `ccam next --auto` dieses Konto vor seinem Zurücksetzen nicht erneut wählt, und `/ccam:switch` (oder `/ccam:switch work`) merkt das Gespräch für `ccam next --continue` vor, sodass nichts zu raten bleibt.

**Team- und Enterprise-Plätze.** Die automatische Wahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und [`auto_switch`](/features/next#auto-switch), wenn du es einschaltest) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation. Im Hilfe-Center von Anthropic ([Team-Tarife](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); ebenso für [Enterprise-Tarife](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)) heißt es: „Each seat on a Team plan is assigned to one named individual, and only that person may use it“ (jeder Platz in einem Team-Tarif ist einer namentlich benannten Person zugewiesen, und nur diese Person darf ihn nutzen) und „Users can't route their work through another user's account to use that user's remaining usage“ (Benutzer dürfen ihre Arbeit nicht über das Konto eines anderen Benutzers leiten, um dessen verbleibendes Nutzungskontingent zu verwenden). Wenn du die automatische Wahl zwischen Plätzen nutzt, kann das daher gegen die Bedingungen von Anthropic verstoßen. Ob du sie dort nutzt, ist deine Entscheidung: Kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt. Für private Tarife (Pro, Max) gilt der [Hinweis zu den Bedingungen in den FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Wechselt ccam das Konto automatisch?

Nur wenn du es einschaltest. Standardmäßig läuft ccam nicht, während du arbeitest: Es startet Claude Code und ersetzt sich selbst durch diesen Prozess, also entscheidest du, wann gewechselt wird, und `ccam next --auto` wählt das Konto mit der geringsten 5-Stunden-Nutzung, überspringt Konten, die bei 100 % stehen oder nicht angemeldet sind, und bleibt innerhalb einer Team- oder Enterprise-Organisation oder unter privaten Tarifen. Mit `relay = true` in `config.toml` wartet ccam stattdessen neben der Sitzung, sodass `/ccam:switch` von selbst weitermacht und dabei die Optionen, den Berechtigungsmodus und die Effort-Stufe der Sitzung beibehält. Mit `auto_switch = true` (Linux und macOS; experimentell) nimmt ccam den Wechsel auch selbst vor: Wenn Claude Code die Sitzung an ihrem Limit stoppt und ein anderes gleichartiges Konto freies Kontingent hat, beendet ccam die Sitzung und setzt das Gespräch dort fort ([Bedingungen und Kosten](/features/next#auto-switch)). Manche anderen Tools tauschen stattdessen die Anmeldung einer laufenden Sitzung aus; [ccam und andere Tools zum Kontowechsel](/compare/account-switchers) erklärt den Unterschied.

## Was du beachten solltest

- Führe `ccam next --continue` im selben Ordner aus, in dem du gearbeitet hast; Gespräche werden pro Projektordner gespeichert.
- Arbeite nicht gleichzeitig mit zwei Konten am selben Gespräch. ccam warnt dich, wenn du es versuchst.
- Ohne geteilten Verlauf wird das Gespräch in das andere Konto kopiert, und die beiden Kopien gehen getrennte Wege; die Checkpoints, aus denen `/rewind` Code wiederherstellt, bleiben beim ersten Konto.
- Die erste Nachricht im neuen Konto kann mehr von dessen Limit verbrauchen als üblich, weil dieses Konto das Gespräch noch nicht im Cache hat. `ccam next --continue --brief` beginnt es stattdessen mit einer kurzen Zusammenfassung.
- Ob du so zwischen Konten wechseln darfst, hängt von den Bedingungen des jeweiligen Kontos ab.

Alle Optionen von `ccam usage` und `ccam next` findest du unter [Nutzungslimits](/docs/usage-limits).
