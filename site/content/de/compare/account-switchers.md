---
title: "Claude Code-Konto wechseln: ccam und andere Tools"
description: Wie sich ccam von claude-swap, CCS, clauth, CAAM, cux und anderen Tools zum Wechseln von Claude Code-Konten unterscheidet und welcher Ansatz wozu passt.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

Tools zum Wechseln von Claude Code-Konten verfolgen einen von zwei Ansätzen. Die meisten tauschen die Anmeldung im einen `~/.claude`-Ordner aus, sodass immer nur ein Konto aktiv ist. Andere, darunter ccam, geben jedem Konto mit `CLAUDE_CONFIG_DIR` einen eigenen Ordner, sodass jedes Konto angemeldet bleibt und mehrere gleichzeitig laufen können. ccam ergänzt geteilte Einstellungen und MCP-Server und setzt ein Gespräch im Ordner eines anderen Kontos fort.

## Eine Anmeldung austauschen oder getrennte Ordner?

**Eine Anmeldung austauschen** behält einen einzigen Satz Einstellungen, MCP-Server und Verlauf, weil es nur einen Ordner gibt. Der Wechsel geht schnell, und eine laufende Sitzung kann die neue Anmeldung übernehmen. Aber es ist immer nur ein Konto aktiv, für alle Terminals und Editoren gleichzeitig.

**Getrennte Ordner** sind der Weg, den [die Dokumentation von Claude Code beschreibt](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Jedes Konto bleibt angemeldet, und du kannst eine berufliche und eine private Sitzung nebeneinander laufen lassen. Der Preis: Jeder Ordner beginnt mit eigenen Einstellungen, MCP-Servern und eigenem Verlauf, sofern kein Tool sie teilt.

## Die Tools im Vergleich

Entnommen aus der Dokumentation des jeweiligen Projekts, Stand 2026-10-19. „Nicht dokumentiert“ heißt, dass das Projekt nichts dazu sagt, nicht, dass die Funktion fehlt.

| Tool | Ansatz | Mehrere Konten gleichzeitig | 5-Stunden- und Wochennutzung | Wenn ein Limit erreicht ist | Plattformen |
| --- | --- | --- | --- | --- | --- |
| ccam | Getrennte Ordner | Ja | Ja, `ccam usage` | Du führst `ccam next --auto --continue` aus, um das Gespräch im Konto mit dem meisten Spielraum fortzusetzen | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | Tauscht die Anmeldung; getrennte Ordner mit `cswap run` | Ja, mit `cswap run` | Ja | `cswap auto` tauscht die Anmeldung ab einem Schwellenwert; die laufende Sitzung geht weiter | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Getrennte Ordner | Ja | Teilweise, in einer App für die macOS-Menüleiste | Konto-Pools reichen Anfragen über einen lokalen Proxy an das nächste Konto weiter | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Tauscht die Anmeldung; getrennte Ordner mit `clauth start` | Ja, mit `clauth start` | Ja | Wechselt entlang einer Fallback-Kette | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Tauscht die Anmeldung oder isolierte Profile | Ja, mit isolierten Profilen | Ja | `caam run` wiederholt nicht interaktive Befehle mit einem anderen Konto | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Tauscht die Anmeldung und umhüllt `claude` | Nicht dokumentiert | Ja | Tauscht die Anmeldung und setzt dasselbe Gespräch fort | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Tauscht die Anmeldung | Nein | Ja | Nicht dokumentiert | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Tauscht die Anmeldung | Nein | Ja | Optionaler Wechsel bei 100 % | macOS |

ccswap ist ein Fork von claude-swap, der auch Codex-Anmeldungen verwaltet. Mehrere dieser Tools unterstützen auch andere Coding-Agents wie Codex oder Gemini.

## Worauf ccam sich konzentriert

- Jedes Konto bleibt angemeldet. Ein Wechsel meldet nie eines ab, und `ccam work` und `ccam personal` können gleichzeitig in zwei Terminals laufen.
- MCP-Server, `settings.json`, `CLAUDE.md`, Agents, Skills und Plugins werden über Links zwischen Ordnern geteilt, Element für Element und Konto für Konto. Anmeldungen werden nie geteilt.
- Mit geteiltem Verlauf öffnet ein Befehl, `ccam next --auto --continue`, dasselbe Gespräch im Ordner des Kontos mit dem meisten Spielraum.
- ccam startet Claude Code und ersetzt sich selbst durch diesen Prozess, im Hintergrund läuft also nichts. Es überwacht deine Nutzung nicht und wechselt nicht von selbst; du entscheidest, wann.

## Wo andere Tools mehr können

- claude-swap, CCS, clauth, CAAM und cux laufen unter Windows ohne WSL. ccam nicht.
- claude-swap, clauth und cux wechseln von selbst, wenn ein Konto einen Nutzungsschwellenwert überschreitet. ccam überlässt den Wechsel dir.
- CCSwitcher, Claude Switcher und claude-swap zeigen die Nutzung in der macOS-Menüleiste.
- CCS, CAAM, ccswap und Claude Switcher verwalten auch Codex und andere Coding-Agents.

## Was passt zu dir?

Nimm ein Tool, das die Anmeldung austauscht, wenn du jeweils ein Konto nutzt, automatisch wechseln willst oder Windows ohne WSL verwendest. Nimm ccam, wenn du Konten nebeneinander, eine gemeinsame Einrichtung und dasselbe Gespräch im nächsten Konto willst. Wenn du nichts teilen musst, reichen [einfache Shell-Aliase](/compare/aliases).

Ob du mehrere Konten nutzen darfst, hängt von den Bedingungen des jeweiligen Kontos ab. Prüfe sie, bevor du loslegst.
