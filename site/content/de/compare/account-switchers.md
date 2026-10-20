---
title: "Claude Code-Konto wechseln: ccam und andere Tools"
description: Wie sich ccam von claude-swap, CCS, clauth, CAAM, cux und anderen Tools zum Wechseln von Claude Code-Konten unterscheidet und welcher Ansatz wozu passt.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: d058da937218
---

Tools zum Wechseln von Claude Code-Konten verfolgen einen von zwei Ansätzen. Die meisten tauschen die Anmeldung im einen `~/.claude`-Ordner aus, sodass immer nur ein Konto aktiv ist. Andere, darunter ccam, geben jedem Konto mit `CLAUDE_CONFIG_DIR` einen eigenen Ordner, sodass jedes Konto angemeldet bleibt und mehrere gleichzeitig laufen können. ccam ergänzt geteilte Einstellungen und MCP-Server und setzt ein Gespräch im Ordner eines anderen Kontos fort.

## Eine Anmeldung austauschen oder getrennte Ordner?

**Eine Anmeldung austauschen** behält einen einzigen Satz Einstellungen, MCP-Server und Verlauf, weil es nur einen Ordner gibt. Der Wechsel geht schnell, und eine laufende Sitzung kann die neue Anmeldung übernehmen. Aber es ist immer nur ein Konto aktiv, für alle Terminals und Editoren gleichzeitig.

Zum Austauschen bewahrt das Tool eine Kopie der Anmeldung jedes anderen Kontos auf und schreibt sie an ihren Platz. Claude Code ersetzt beim Erneuern einer Anmeldung deren Refresh-Token, deshalb muss jede gespeicherte Kopie mit den laufenden Sitzungen Schritt halten. Konten, die nach einem Wechsel abgemeldet sind oder erneut `/login` brauchen, gehören zu den Problemen, die diese Projekte am häufigsten beheben.

**Getrennte Ordner** sind der Weg, den [die Dokumentation von Claude Code beschreibt](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Jedes Konto bleibt angemeldet, und du kannst eine berufliche und eine private Sitzung nebeneinander laufen lassen. Jede Anmeldung bleibt in ihrem eigenen Ordner, wo Claude Code sie wie gewohnt erneuert. Der Preis: Jeder Ordner beginnt mit eigenen Einstellungen, MCP-Servern und eigenem Verlauf, sofern kein Tool sie teilt.

## Die Tools im Vergleich

Entnommen aus der Dokumentation und dem Quellcode des jeweiligen Projekts, Stand 2026-10-20. „Nicht dokumentiert“ heißt, dass das Projekt nichts dazu sagt, nicht, dass die Funktion fehlt.

| Tool | Ansatz | Mehrere Konten gleichzeitig | 5-Stunden- und Wochennutzung | Wenn ein Limit erreicht ist | Plattformen |
| --- | --- | --- | --- | --- | --- |
| ccam | Getrennte Ordner | Ja | Ja, `ccam usage` | Du führst `ccam next --auto --continue` aus, um das Gespräch im Konto mit dem meisten Spielraum fortzusetzen; `ccam --fallback` führt `-p`-Jobs dort von selbst erneut aus | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Tauscht die Anmeldung; getrennte Ordner mit `cswap run` | Ja, mit `cswap run` | Ja | `cswap auto` tauscht die Anmeldung ab einem Schwellenwert; die laufende Sitzung geht weiter | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Getrennte Ordner | Ja | Teilweise | Konto-Pools reichen Anfragen über einen lokalen Proxy an das nächste Konto weiter | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Tauscht die Anmeldung; getrennte Ordner mit `clauth start` | Ja, mit `clauth start` | Ja | Wechselt entlang einer Fallback-Kette | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Tauscht die Anmeldung oder isolierte Profile | Ja, mit isolierten Profilen | Ja | `caam run` wiederholt nicht interaktive Befehle mit einem anderen Konto | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Tauscht die Anmeldung und umhüllt `claude` | Nein, ein aktives Konto | Ja | Tauscht die Anmeldung und setzt dasselbe Gespräch fort | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Tauscht die Anmeldung | Nein | Ja | Optionaler Wechsel ab einem Schwellenwert | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Tauscht die Anmeldung | Nein | Ja | Optionaler Wechsel bei 100 % | macOS |

ccswap ist ein Fork von claude-swap, der auch Codex-Anmeldungen verwaltet. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) und [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) geben ebenfalls jedem Konto einen eigenen Ordner. Mehrere dieser Tools unterstützen auch andere Coding-Agents wie Codex oder Gemini.

## Wie sie aufgebaut sind

| Tool | Geschrieben in | Während Claude Code läuft | Kopiert Anmeldungen hin und her | Was Konten teilen können |
| --- | --- | --- | --- | --- |
| ccam | Go, eine statische Binärdatei | Nichts: ccam ersetzt sich selbst durch `claude` (unter Windows, wo das nicht geht, wartet es darauf) | Nein | Jedes Element einzeln: lokal oder eine benannte Gruppe, pro Konto |
| claude-swap | Python 3.12+, installiert mit uv oder pipx | Nichts mit `cswap run`; `cswap auto` und die Menüleisten-App laufen weiter | Ja, wenn es die Anmeldung tauscht | Mit `cswap run`: Einstellungen, `CLAUDE.md` und Erweiterungen zusammen, Verlauf optional; MCP-Server werden bei jedem Start kopiert |
| CCS | Node.js 18+ oder Bun | CCS bleibt als Elternprozess; Pools fügen einen lokalen Proxy hinzu | Nein für seine Kontoordner; Pools haben eigene Anmeldungen | Geteilt oder lokal pro Konto, Verlauf in Kontextgruppen; MCP-Server werden kopiert |
| clauth | Rust | clauth bleibt als Elternprozess | Ja, wenn es die Anmeldung tauscht | Das ganze `~/.claude` oder nichts |
| CAAM | Go | Hängt vom Modus ab; ein optionaler Daemon | Ja, im Vault-Modus | Durch den Modus festgelegt; MCP-Server werden bei jedem Start kopiert |
| cux | Go | Ein Wrapper bleibt um `claude` | Ja | Alles: ein Ordner |
| CCSwitcher, Claude Switcher | macOS-Apps (Swift, Python) | Die Menüleisten-App | Ja | Alles: ein Ordner |

## Was ccam gut macht

### Anmeldungen bleiben, wo Claude Code sie abgelegt hat

ccam kopiert nie eine Claude-Anmeldung von einem Ort an einen anderen. Die Anmeldung jedes Kontos bleibt in seinem eigenen Ordner, unter macOS im eigenen Schlüsselbundeintrag dieses Ordners, und Claude Code ist das Programm, das sie erneuert. Es gibt keine gespeicherte Kopie, die aus dem Takt geraten kann. Ein bestehendes `~/.claude`, das mit `ccam import` registriert wurde, behält den Schlüsselbundeintrag, den das einfache `claude` nutzt, und bleibt so in beiden Fällen angemeldet.

ccam erneuert ein abgelaufenes Token nur dann selbst, wenn du `ccam usage --refresh` ausführst. Es setzt vorher eine Sperre und verweigert den Vorgang, solange das Konto läuft, rotiert also nie ein Token unter einer laufenden Sitzung.

### Nichts läuft zwischen dir und Claude Code

ccam liest eine Datei, prüft seine Links und ersetzt sich selbst durch `claude` (`exec`). Während du arbeitest, gibt es keinen ccam-Prozess, keinen Wrapper, keinen Proxy und keinen Daemon: Standardeingabe, Signale und Exit-Code gehören Claude Code, und seine Anfragen gehen direkt an Anthropic. Windows kann einen Prozess nicht ersetzen, deshalb startet ccam dort Claude Code mit derselben Konsole und wartet auf dessen Exit-Code, ohne sonst etwas zu tun.

Ein Konto so zu starten kostet etwa 4 ms, unter 6 ms beim 95. Perzentil, bei 8 MB Spitzenspeicher ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) mit hyperfine, 300 Durchläufe). Auf demselben Laptop braucht allein der Start des Python-Interpreters 13 ms und der von Node.js 31 ms, bevor ein Programm überhaupt etwas tut.

### Teilen Element für Element, Konto für Konto

MCP-Server, `settings.json`, `CLAUDE.md`, Erweiterungen (Agents, Befehle, Skills, Plugins) und Gesprächsverlauf sind fünf getrennte Elemente. Jedes Konto setzt jedes Element auf `local` oder auf eine Gruppe. So können zwei berufliche Konten den Verlauf in einer Gruppe `work` teilen, während ein privates Konto seinen eigenen behält, und alle drei teilen dieselben MCP-Server.

- MCP-Server werden beim Start mit `--mcp-config` übergeben, daher bearbeitet ccam standardmäßig nie `.claude.json`, die Datei, die Claude Code während des Laufs neu schreibt. `ccam mcp add` fügt einen Server einmal für die ganze Gruppe hinzu.
- Die anderen Elemente sind Links in den Gruppenordner, sodass eine Änderung in einem Konto sofort in den anderen erscheint. Beim Start wird nichts kopiert.
- Anmeldungen, `.claude.json` und der Laufzeitzustand werden nie geteilt.

### Jede Änderung wird vorher gezeigt und gesichert

`ccam share` und `ccam migrate` zeigen einen Plan und fragen nach, bevor sie etwas verschieben. Wenn sich eine Datei zwischen einem Konto und seiner Gruppe unterscheidet, entscheidest du: die der Gruppe behalten, deine verwenden oder das JSON zusammenführen. Originale landen in `.ccam-backup/<time>/`, der Ordner eines entfernten Kontos kommt in einen Papierkorb-Ordner, und Dateien werden mit einer einzigen Umbenennung ersetzt, sodass ein unterbrochener Schreibvorgang die alte Datei hinterlässt. `--dry-run` zeigt, was Starten, Teilen, Migrieren, Verknüpfungen, Updates und Reparaturen tun würden. `ccam doctor --fix` repariert kaputte Links und Berechtigungen.

### Dasselbe Gespräch im nächsten Konto

Mit geteiltem Verlauf fragt `ccam next --auto --continue` alle Konten gleichzeitig nach ihrer Nutzung, öffnet das mit dem meisten Spielraum und setzt das Gespräch fort. Du entscheidest, wann gewechselt wird; ccam überwacht deine Nutzung nicht im Hintergrund. Die Statuszeile schlägt den Befehl vor, sobald ein Konto 90 % erreicht.

Skripte können nicht auf dich warten, deshalb wechselt `ccam --fallback work -p "…"` für sie: Wenn das Konto sein Limit erreicht, läuft derselbe Befehl erneut mit dem Konto mit dem meisten Spielraum, und nur die Ausgabe des Durchlaufs, der fertig wurde, wird ausgegeben. Ein Fehler, der kein Nutzungslimit ist, wird unverändert weitergegeben.

### Das richtige Konto in jedem Ordner

`ccam folder work ~/code/company` sorgt dafür, dass jeder Start ohne Kontonamen in diesem Ordner `work` verwendet, während deine übrigen Projekte beim Standardkonto bleiben. Die Einstellung steht in derselben `config.toml`, die jeder Start liest, und kostet deshalb keine Zeit.

### Eine Binärdatei, geprüfte Downloads

ccam ist eine statische Binärdatei von etwa 10 MB, ohne Laufzeitumgebung, die installiert werden müsste, für Linux, macOS und Windows. Die Installationsskripte (`install.sh`, `install.ps1`) und `ccam update` prüfen jeden Download gegen die SHA-256-Prüfsummen des Releases, und ccam sucht nie von sich aus nach Updates. Es gibt keine Telemetrie: Die einzigen Anfragen sind Nutzungsabfragen bei Anthropic für deine Konten, parallel mit 3 Sekunden Zeitlimit und 60 Sekunden Cache, und die Downloads, die du anforderst.

## Wo andere Tools mehr können

- claude-swap, clauth, cux, CAAM und CCSwitcher können eine laufende interaktive Sitzung von selbst wechseln, wenn ein Konto einen Nutzungsschwellenwert überschreitet. ccam wechselt `-p`-Durchläufe mit `--fallback` und überlässt interaktive Sitzungen dir.
- Unter Windows ohne Entwicklermodus teilt ccam Ordner, kann aber keine einzelnen Dateien wie `settings.json` teilen; claude-swap und CCS kopieren sie dort stattdessen.
- CCSwitcher, Claude Switcher und claude-swap zeigen die Nutzung in der macOS-Menüleiste.
- CCS, CAAM, clauth, ccswap und Claude Switcher verwalten auch Codex und andere Coding-Agents.
- CAAM, CCSwitcher und Claude Switcher lassen sich mit Homebrew installieren, und clauth und CAAM signieren ihre Release-Prüfsummen. ccam installierst du mit seinem Skript, einem `.deb` oder `go install`, und es prüft Prüfsummen ohne Signatur.

## Was passt zu dir?

Nimm ein Tool, das die Anmeldung austauscht, wenn du jeweils ein Konto nutzt oder laufende Sitzungen automatisch wechseln sollen. Nimm ccam, wenn du Konten nebeneinander, eine gemeinsame Einrichtung, dasselbe Gespräch im nächsten Konto und nichts willst, das neben Claude Code mitläuft, unter Linux, macOS oder Windows. Wenn du nichts teilen musst, reichen [einfache Shell-Aliase](/compare/aliases).

Ob du mehrere Konten nutzen darfst, hängt von den Bedingungen des jeweiligen Kontos ab. Prüfe sie, bevor du loslegst.
