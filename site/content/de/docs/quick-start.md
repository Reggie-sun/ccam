---
title: Schnellstart
description: Zwei Claude Code-Konten mit ccam add hinzufügen, mit claude1, claude2 oder ccs zwischen ihnen wechseln und festlegen, was sie teilen.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 9462b80d61bc
---

Füge jedes Konto mit `ccam add` hinzu und starte Claude Code dann mit `claude1`, `claude2` oder `ccs <name>`. Mehr brauchen die meisten nicht.

## 1. Erstes Konto hinzufügen

```sh
ccam add
```

Der Assistent fragt nach einem Namen wie `personal` oder `work` und danach, welche Elemente das Konto mit deinen anderen Konten teilen soll. Dann führt er `claude auth login` aus, das den Browser öffnet, damit du dich anmelden kannst. Wenn er fertig ist, siehst du:

```console
✓ personal (me@example.com) 추가됨 — 단축 명령: claude1
```

ccam gibt seine eigenen Meldungen vorerst auf Koreanisch aus; die Befehle und Optionen sind in jeder Sprache gleich.

## 2. Zweites Konto hinzufügen

```sh
ccam add work
```

Melde dich im Browser mit dem anderen Konto an. Es wird zu `claude2`.

## 3. Wechseln

Jeder dieser Befehle startet Claude Code mit dem Arbeitskonto:

```sh
claude2
ccs work
ccs 2
```

Alles nach dem Kontonamen geht unverändert an Claude Code, du nutzt also weiter die Optionen, die du kennst:

```sh
ccs work --resume
ccs work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

Wenn du `ccs` allein ausführst, öffnet sich eine Auswahl: Nutze die Pfeiltasten und Enter oder drücke die Nummer des Kontos.

## 4. Festlegen, was geteilt wird

Standardmäßig werden MCP-Server, `settings.json`, die globale `CLAUDE.md` sowie Agents, Commands, Skills und Plugins zwischen Konten geteilt. Der Gesprächsverlauf nicht. Schalte ihn für beide Konten ein, wenn du Gespräche kontoübergreifend fortsetzen willst:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam zeigt, was es verschieben wird, bevor es etwas ändert. Details findest du unter [Teilen](/docs/sharing).

## 5. Wenn ein Konto sein Limit erreicht

```sh
ccs usage
ccs next --continue
```

`ccs usage` zeigt die 5-Stunden- und die Wochennutzung aller Konten. `ccs next --continue` öffnet dasselbe Gespräch im nächsten Konto. Siehe [Nutzungslimits](/docs/usage-limits).

## Du nutzt bereits ~/.claude?

Registriere es, statt dich erneut anzumelden:

```sh
ccam import
```

Der Ordner bleibt, wo er ist, und einfaches `claude` funktioniert weiter wie bisher.
