---
title: ccam stellt sich vor
description: Warum es ccam gibt, was es für Menschen mit mehreren Claude Code-Konten tut und wie es dir nicht in die Quere kommt.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: b116c91bb8d4
---

ccam ist ein kleines Kommandozeilentool für alle, die mehr als ein Claude Code-Konto auf demselben Computer nutzen: ein privates und ein berufliches Konto oder mehrere Konten für ein langes Projekt.

## Das Problem

Claude Code unterstützt mehrere Konten über `CLAUDE_CONFIG_DIR`: ein Ordner pro Konto. In der Praxis ist jeder Ordner eine Insel. Einstellungen, `CLAUDE.md`, MCP-Server, Agents und Skills müssen in jedem Ordner neu eingerichtet und von Hand synchron gehalten werden. Gespräche bleiben dort, wo sie begonnen haben. Und wenn ein Konto sein Nutzungslimit erreicht, merkst du das erst beim Ausprobieren.

## Was ccam tut

- `ccam add` legt den Ordner an, führt die Anmeldung aus und gibt dem Konto auf Wunsch einen Kurzbefehl, den du selbst benennst, etwa `ccwork`.
- `ccam work` und `ccwork` starten Claude Code mit diesem Konto und geben jedes Argument unverändert weiter.
- MCP-Server, Einstellungen, `CLAUDE.md` und Erweiterungen werden über Links geteilt, Element für Element und Konto für Konto. Anmeldungen werden nie geteilt.
- `ccam usage` zeigt die 5-Stunden- und die Wochennutzung aller Konten; `ccam next --continue` setzt das Gespräch im nächsten Konto fort.

## Nicht im Weg stehen

Ein Tool, das du vor jeder Sitzung ausführst, sollte nichts kosten, was du bemerkst. ccam liest eine Konfigurationsdatei, prüft die geteilten Links mit jeweils einem einzigen Systemaufruf und ersetzt sich dann durch `claude`. Es kostet etwa 4 Millisekunden und läuft nicht, während du arbeitest.

Außerdem geht es sorgsam mit deinen Dateien um. Jede Änderung am Teilen gibt vor dem Start einen Plan aus und bewahrt die Originale in einem Sicherungsordner auf. `ccam doctor` prüft Links, Berechtigungen und Anmeldungen und repariert, was sich sicher reparieren lässt.

## Ausprobieren

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

ccam ist Open Source und inoffiziell: Es steht in keiner Verbindung zu Anthropic und wird nicht von Anthropic unterstützt. Ob du mehrere Konten nutzen darfst, hängt von den Bedingungen jedes Kontos ab. Feedback und Übersetzungen sind auf [GitHub](https://github.com/jhyunwoo/ccam) willkommen.
