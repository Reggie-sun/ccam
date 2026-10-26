---
title: Schnellstart
metaTitle: "Schnellstart: Claude Code-Konten anlegen und wechseln"
description: Zwei Claude Code-Konten mit ccam add hinzufügen, ihnen Kurzbefehle wie ccwork geben, zwischen ihnen wechseln und festlegen, was sie teilen.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: "bd893bbb83ab"
---

Füge jedes Konto mit `ccam add` hinzu und starte Claude Code dann mit `ccam <name>` oder mit einem Kurzbefehl, den du selbst benennst, etwa `ccwork`. Mehr brauchen die meisten nicht.

Du nutzt Claude Code bereits? Lies vorher [Du nutzt bereits ~/.claude?](#already-using-claude) am Ende dieser Seite.

## 1. Erstes Konto hinzufügen

```sh
ccam add
```

Der Assistent fragt nach einem Namen wie `personal` oder `work`, danach, welche Elemente das Konto mit deinen anderen Konten teilen soll, und ob ein Kurzbefehl angelegt werden soll (siehe Schritt 3). Dann führt er `claude auth login` aus, das den Browser öffnet, damit du dich anmelden kannst. Auf einem entfernten Rechner landet dein Browser eventuell auf einer Seite `http://localhost:…/callback`, die nicht lädt; füge diese Adresse im Terminal ein, und ccam schließt die Anmeldung ab. Wenn die Anmeldung abgeschlossen ist, siehst du:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Vorhandene Einstellungen übernehmen

Wenn du Claude Code schon benutzt hast, liegen deine Einstellungen, `CLAUDE.md` und Erweiterungen in `~/.claude` und deine MCP-Server in `~/.claude.json`. Kopiere sie in die geteilte Gruppe `default`, damit jedes Konto sie bekommt:

```sh
ccam migrate
```

ccam zeigt, was es kopieren wird, und fragt vorher nach. Es liest `~/.claude` nur, kopiert nie eine Anmeldung und lässt den Gesprächsverlauf weg, außer du gibst `--history` an. Details findest du unter [Eine vorhandene Einrichtung kopieren](/docs/sharing#copy-an-existing-setup).

## 3. Zweites Konto hinzufügen

```sh
ccam add work
```

Der Assistent überspringt den Namen, den du schon auf der Befehlszeile angegeben hast, und zeigt die Gruppen zum Beitreten, mit `default` angehakt. Gib dem Konto diesmal einen eigenen Kurzbefehl, etwa `ccwork`, damit du es mit einem Wort starten kannst. Nach den Gruppen fragt der Assistent:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Drücke `y` und tippe dann einen Namen ein, oder drücke Enter, um den Vorschlag zu übernehmen:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Melde dich dann im Browser mit dem anderen Konto an. ccam bestätigt das Konto und seinen Kurzbefehl:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. Konten wechseln

`ccam work` startet Claude Code mit dem Arbeitskonto. Wenn du die Frage nach dem Kurzbefehl übersprungen hast, kannst du jederzeit einen hinzufügen:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Jeder dieser Befehle startet jetzt Claude Code mit dem Arbeitskonto (`2` ist seine Nummer in `ccam list`):

```sh
ccwork
ccam work
ccam 2
```

Wenn du `ccam` allein ausführst, öffnet sich das [Dashboard](/features/dashboard): Gehe mit den Pfeiltasten oder der Nummer zu einem Konto und drücke dann Enter.

Alles nach dem Kurzbefehl oder dem Kontonamen geht unverändert an Claude Code, du nutzt also weiter die Optionen, die du kennst:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Festlegen, was geteilt wird

Standardmäßig werden MCP-Server, `settings.json`, die globale `CLAUDE.md` und Erweiterungen (Agents, Commands, Skills und Plugins) über die Gruppe `default` zwischen Konten geteilt. Der Gesprächsverlauf nicht. Füge ihn der Gruppe hinzu, wenn du Gespräche kontoübergreifend fortsetzen willst:

```sh
ccam group set default --add history
```

ccam zeigt für jedes Konto, was es verschieben wird, bevor es etwas ändert. Details findest du unter [Teilen mit Gruppen](/features/share).

## 6. Wenn ein Konto sein Limit erreicht

```sh
ccam usage
ccam next --continue
```

`ccam usage` zeigt die 5-Stunden- und die Wochennutzung aller Konten. `ccam next --continue` öffnet dasselbe Gespräch im nächsten Konto. In Skripten wechselt `ccam --fallback work -p "…"` von selbst. Siehe [Nutzungslimits](/docs/usage-limits).

**Team- und Enterprise-Plätze.** Die automatische Kontoauswahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und, sobald du es einschaltest, [`auto_switch`](/features/next#auto-switch)) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation, was gegen die Bedingungen von Anthropic verstoßen kann. Lies die [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) und kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt.

## Du nutzt bereits ~/.claude?

Registriere es, statt dich erneut anzumelden:

```sh
ccam import
```

Der Ordner bleibt, wo er ist, und einfaches `claude` funktioniert weiter wie bisher.
