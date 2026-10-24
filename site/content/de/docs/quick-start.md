---
title: Schnellstart
metaTitle: "Schnellstart: Claude Code-Konten anlegen und wechseln"
description: Zwei Claude Code-Konten mit ccam add hinzufügen, ihnen Kurzbefehle wie ccwork geben, zwischen ihnen wechseln und festlegen, was sie teilen.
updated: 2026-10-24
order: 2
translation: machine
sourceHash: e98d0637356a
---

Füge jedes Konto mit `ccam add` hinzu und starte Claude Code dann mit `ccam <name>` oder mit einem Kurzbefehl, den du selbst benennst, etwa `ccwork`. Mehr brauchen die meisten nicht.

## 1. Erstes Konto hinzufügen

```sh
ccam add
```

Der Assistent fragt nach einem Namen wie `personal` oder `work`, danach, welche Elemente das Konto mit deinen anderen Konten teilen soll, und ob ein Kurzbefehl angelegt werden soll (siehe Schritt 4). Dann führt er `claude auth login` aus, das den Browser öffnet, damit du dich anmelden kannst. Auf einem entfernten Rechner landet der Browser eventuell auf einer Seite `http://localhost:…/callback`, die nicht lädt; füge diese Adresse im Terminal ein, und ccam schließt die Anmeldung ab. Wenn er fertig ist, siehst du:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

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

Wenn du `ccam` allein ausführst, öffnet sich das [Dashboard](/features/dashboard): Gehe mit den Pfeiltasten oder der Nummer zu einem Konto und drücke dann Enter.

Alles nach dem Kurzbefehl oder dem Kontonamen geht unverändert an Claude Code, du nutzt also weiter die Optionen, die du kennst:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Festlegen, was geteilt wird

Standardmäßig werden MCP-Server, `settings.json`, die globale `CLAUDE.md` sowie Agents, Commands, Skills und Plugins über die Gruppe `default` zwischen Konten geteilt. Der Gesprächsverlauf nicht. Füge ihn der Gruppe hinzu, wenn du Gespräche kontoübergreifend fortsetzen willst:

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

**Team- und Enterprise-Plätze.** Die automatische Wahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation. Im Hilfe-Center von Anthropic ([Team-Tarife](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); ebenso für [Enterprise-Tarife](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)) heißt es: „Each seat on a Team plan is assigned to one named individual, and only that person may use it“ (jeder Platz in einem Team-Tarif ist einer namentlich benannten Person zugewiesen, und nur diese Person darf ihn nutzen) und „Users can't route their work through another user's account to use that user's remaining usage“ (Benutzer dürfen ihre Arbeit nicht über das Konto eines anderen Benutzers leiten, um dessen verbleibendes Nutzungskontingent zu verwenden). Wenn du die automatische Wahl zwischen Plätzen nutzt, kann das daher gegen die Bedingungen von Anthropic verstoßen. Kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt. Für private Tarife (Pro, Max) gilt der [Hinweis zu den Bedingungen in den FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Du nutzt bereits ~/.claude?

Registriere es, statt dich erneut anzumelden:

```sh
ccam import
```

Der Ordner bleibt, wo er ist, und einfaches `claude` funktioniert weiter wie bisher.
