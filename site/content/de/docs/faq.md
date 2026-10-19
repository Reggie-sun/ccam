---
title: Häufig gestellte Fragen
description: Antworten zu ccam, Claude Code-Konten, Nutzungslimits, Teilen, Datenschutz und unterstützten Systemen.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 5ff6e94786a2
---

## Ist ccam ein offizielles Tool von Anthropic?

Nein. ccam ist ein inoffizielles Open-Source-Projekt und steht in keiner Verbindung zu Anthropic und wird nicht von Anthropic unterstützt. Es verwendet die Einstellung `CLAUDE_CONFIG_DIR`, die Claude Code dafür dokumentiert, mehrere Konten nebeneinander zu betreiben.

## Darf ich mehrere Claude Code-Konten nutzen?

Das hängt von den Bedingungen jedes Kontos ab und bei Arbeitskonten von den Richtlinien deiner Organisation. Prüfe sie, bevor du mehrere Konten nutzt. ccam macht nur den Wechsel zwischen Konten einfacher, die du bereits hast.

## Umgeht ccam Nutzungslimits?

Nein. Jedes Konto behält seine eigenen Limits. ccam zeigt sie an und lässt dich in einem anderen Konto weitermachen, das dir gehört, wenn eines erschöpft ist.

## Wie nutze ich zwei Claude Code-Konten auf einem Computer?

Gib jedem Konto mit `CLAUDE_CONFIG_DIR` einen eigenen Konfigurationsordner. ccam erledigt das für dich: Führe `ccam add` für jedes Konto aus und starte sie dann mit `ccam personal` und `ccam work` oder mit Kurzbefehlen, die du selbst benennst, etwa `ccwork`. Die [Anleitung](/guides/multiple-claude-code-accounts) zeigt auch, wie es von Hand geht.

## Wie setze ich ein Gespräch in einem anderen Konto fort?

Teile den Gesprächsverlauf zwischen den beiden Konten mit `ccam share <name> --history on` und führe dann `ccam next --continue` aus. Claude Code öffnet das neueste Gespräch des aktuellen Ordners im nächsten Konto.

## Sieht oder sendet ccam meine Anmeldung?

ccam teilt deine Claude-Anmeldung nie zwischen Konten. Es liest das Zugriffstoken eines Kontos nur, um beim Nutzungs-Endpunkt von Anthropic die Nutzung dieses Kontos abzufragen, wenn du `ccam usage`, `ccam next --auto` oder `ccam doctor` ausführst. `ccam migrate` kopiert die gespeicherten Anmeldungen von MCP-Servern aus `~/.claude` in die eigenen Anmeldedateien deiner Konten auf demselben Computer. ccam hat keine Telemetrie und sendet an keine andere Stelle etwas.

## Wie melde ich mich auf einem entfernten Server an?

Führe `ccam add` oder `ccam login <name>` wie gewohnt in der SSH-Sitzung aus. Claude Code wartet auf dem Server auf die Anmeldung, daher landet der Browser auf deinem Computer auf einer Seite `http://localhost:<port>/callback?code=…`, die nicht lädt. Kopiere diese Adresse aus dem Browser und füge sie im Terminal ein: ccam gibt sie an Claude Code auf dem Server weiter, das die Anmeldung abschließt. Die Adresse enthält einen Einmalcode, gib sie also nicht weiter.

## Wo werden meine Konten gespeichert?

In `~/.local/share/ccam/profiles/<name>/`, ein Ordner pro Konto. Unter macOS speichert Claude Code die Anmeldung selbst im Schlüsselbund, unter einem Namen, der von diesem Ordner abgeleitet ist.

## Was passiert mit meinem vorhandenen ~/.claude?

Nichts, außer du führst `ccam import` oder `ccam migrate` aus. Der Import registriert den Ordner dort, wo er liegt, sodass einfaches `claude` und `ccam <name>` weiterhin dasselbe Konto und dieselben Einstellungen verwenden. `ccam migrate` kopiert seine Einstellungen, `CLAUDE.md`, Erweiterungen und MCP-Server in die geteilten Einstellungen und liest den Ordner dabei nur.

## Welche Systeme unterstützt ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine und andere) und macOS. Unter Windows nutzt du es innerhalb von WSL.

## Warum erscheint ein MCP-Server nur in einem Konto?

Server, die innerhalb einer Sitzung mit `claude mcp add` hinzugefügt werden, werden nur für dieses Konto gespeichert. ccam weist dich beim nächsten Start des Kontos darauf hin; verschiebe den Server mit `ccam mcp import <account> <server>` in die geteilte Gruppe.

## Macht ccam Claude Code langsamer?

Claude Code über ccam zu starten kostet etwa 4 Millisekunden, gemessen mit `scripts/perf-check.sh` auf einem Laptop. ccam liest eine kleine Datei, prüft die geteilten Links und ersetzt sich dann durch `claude`, läuft also nicht, während du arbeitest.

## Wie aktualisiere ich ccam?

Führe `ccam update` aus. Der Befehl lädt das neueste Release herunter, prüft seine Prüfsumme und ersetzt die Datei `ccam`; `ccam update --check` zeigt nur an, ob es ein neueres Release gibt. ccam aktualisiert sich nie von selbst. Wenn du das `.deb`-Paket installiert hast, installiere stattdessen das neue `.deb` mit apt. Siehe [Aktualisieren](/docs/install#update).

## Wie deinstalliere ich ccam?

Entferne den Kurzbefehl jedes Kontos mit `ccam shortcut <account> --remove`, lösche `ccam` aus `~/.local/bin` und danach `~/.config/ccam`, `~/.local/share/ccam` und `~/.local/state/ccam`, wenn du die Konten nicht mehr brauchst. Siehe [Installation](/docs/install#uninstall).
