---
title: Schnellstart
description: Zwei Claude Code-Konten mit ccam add hinzufügen, ihnen Kurzbefehle wie ccwork geben, zwischen ihnen wechseln und festlegen, was sie teilen.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: b73914ff371f
---

Füge jedes Konto mit `ccam add` hinzu und starte Claude Code dann mit `ccam <name>` oder mit einem Kurzbefehl, den du selbst benennst, etwa `ccwork`. Mehr brauchen die meisten nicht.

## 1. Erstes Konto hinzufügen

```sh
ccam add
```

Der Assistent fragt nach einem Namen wie `personal` oder `work`, danach, welche Elemente das Konto mit deinen anderen Konten teilen soll, und ob ein Kurzbefehl angelegt werden soll (siehe Schritt 4). Dann führt er `claude auth login` aus, das den Browser öffnet, damit du dich anmelden kannst. Wenn er fertig ist, siehst du:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

ccam gibt seine eigenen Meldungen auf Englisch aus; die Befehle und Optionen sind in jeder Sprache gleich.

## 2. Vorhandene Einstellungen übernehmen

Wenn du Claude Code schon benutzt hast, liegen deine Einstellungen, `CLAUDE.md`, Erweiterungen und MCP-Server in `~/.claude`. Kopiere sie in die geteilten Einstellungen, damit jedes Konto sie bekommt:

```sh
ccam migrate
```

ccam zeigt, was es kopieren wird, und fragt vorher nach. Es liest `~/.claude` nur, kopiert nie eine Anmeldung und lässt den Gesprächsverlauf weg, außer du gibst `--history` an. Details findest du unter [Eine vorhandene Einrichtung kopieren](/docs/sharing#copy-an-existing-setup).

## 3. Zweites Konto hinzufügen

```sh
ccam add work
```

Melde dich im Browser mit dem anderen Konto an.

## 4. Konten wechseln

`ccam work` startet Claude Code mit dem Arbeitskonto. Kürzer geht es mit einem eigenen Kurzbefehl für das Konto, etwa `ccwork`. Nach den geteilten Elementen fragt der Assistent:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Drücke `y` und tippe dann einen Namen ein, oder drücke Enter, um den Vorschlag zu übernehmen:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Nach der Anmeldung bestätigt ccam ihn:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Wenn du die Frage übersprungen hast, kannst du jederzeit einen Kurzbefehl hinzufügen:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Jetzt startet jeder dieser Befehle Claude Code mit dem Arbeitskonto:

```sh
ccwork
ccam work
ccam 2
```

Wenn du `ccam` allein ausführst, öffnet sich eine Auswahl: Nutze die Pfeiltasten und Enter oder drücke die Nummer des Kontos.

Alles nach dem Kurzbefehl oder dem Kontonamen geht unverändert an Claude Code, du nutzt also weiter die Optionen, die du kennst:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Festlegen, was geteilt wird

Standardmäßig werden MCP-Server, `settings.json`, die globale `CLAUDE.md` sowie Agents, Commands, Skills und Plugins zwischen Konten geteilt. Der Gesprächsverlauf nicht. Schalte ihn für beide Konten ein, wenn du Gespräche kontoübergreifend fortsetzen willst:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam zeigt, was es verschieben wird, bevor es etwas ändert. Details findest du unter [Teilen](/docs/sharing).

## 6. Wenn ein Konto sein Limit erreicht

```sh
ccam usage
ccam next --continue
```

`ccam usage` zeigt die 5-Stunden- und die Wochennutzung aller Konten. `ccam next --continue` öffnet dasselbe Gespräch im nächsten Konto. Siehe [Nutzungslimits](/docs/usage-limits).

## Du nutzt bereits ~/.claude?

Registriere es, statt dich erneut anzumelden:

```sh
ccam import
```

Der Ordner bleibt, wo er ist, und einfaches `claude` funktioniert weiter wie bisher.
