---
title: "Claude Code-Konto wechseln: ccam und andere Tools"
description: Wie sich ccam von claude-swap, CCS, clausona, clauth, CAAM, cux und anderen Tools für Claude Code-Kontowechsel unterscheidet und welcher Ansatz wozu passt.
updated: 2026-10-25
order: 2
translation: machine
sourceHash: "3d87ba7c975a"
---

Tools zum Wechseln von Claude Code-Konten verfolgen einen von zwei Ansätzen. Die meisten tauschen die Anmeldung im einen `~/.claude`-Ordner aus, sodass immer nur ein Konto aktiv ist. Andere, darunter ccam, geben jedem Konto mit `CLAUDE_CONFIG_DIR` einen eigenen Ordner, sodass jedes Konto angemeldet bleibt und mehrere gleichzeitig laufen können. ccam ergänzt geteilte Einstellungen und MCP-Server und setzt ein Gespräch im Ordner eines anderen Kontos fort.

## Eine Anmeldung austauschen oder getrennte Ordner?

**Eine Anmeldung austauschen** behält einen einzigen Satz Einstellungen, MCP-Server und Verlauf, weil es nur einen Ordner gibt. Der Wechsel geht schnell, und eine laufende Sitzung kann die neue Anmeldung übernehmen. Aber es ist immer nur ein Konto aktiv, für alle Terminals und Editoren gleichzeitig.

Zum Austauschen bewahrt das Tool eine Kopie der Anmeldung jedes anderen Kontos auf und schreibt sie an ihren Platz. Claude Code ersetzt beim Erneuern einer Anmeldung deren Refresh-Token, deshalb muss jede gespeicherte Kopie mit den laufenden Sitzungen Schritt halten. Konten, die nach einem Wechsel abgemeldet sind oder erneut `/login` brauchen, gehören zu den Problemen, die diese Projekte am häufigsten beheben.

**Getrennte Ordner** sind der Weg, den [die Dokumentation von Claude Code beschreibt](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Jedes Konto bleibt angemeldet, und du kannst eine berufliche und eine private Sitzung nebeneinander laufen lassen. Jede Anmeldung bleibt in ihrem eigenen Ordner, wo Claude Code sie wie gewohnt erneuert. Der Preis: Jeder Ordner beginnt mit eigenen Einstellungen, MCP-Servern und eigenem Verlauf, sofern kein Tool sie teilt.

## Die Tools im Vergleich

Entnommen aus der Dokumentation und dem Quellcode des jeweiligen Projekts, Stand 2026-10-20 (clausona: v0.12.0-beta). „Nicht dokumentiert“ heißt, dass das Projekt nichts dazu sagt, nicht, dass die Funktion fehlt.

| Tool | Ansatz | Mehrere Konten gleichzeitig | 5-Stunden- und Wochennutzung | Wenn ein Limit erreicht ist | Plattformen |
| --- | --- | --- | --- | --- | --- |
| ccam | Getrennte Ordner | Ja | Ja, `ccam usage` | Eine Benachrichtigung in der Sitzung; `ccam next --auto --continue` oder `/ccam:switch` bringt das Gespräch zum Konto mit dem meisten Spielraum, mit oder ohne geteilten Verlauf, und mit eingeschaltetem `auto_switch` (Linux, macOS) wandert die Sitzung von selbst dorthin; `ccam --fallback` führt `-p`-Jobs dort von selbst erneut aus, und `ccam pick` verteilt Worker | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Tauscht die Anmeldung; getrennte Ordner mit `cswap run` | Ja, mit `cswap run` | Ja | `cswap auto` tauscht die Anmeldung ab einem Schwellenwert; die laufende Sitzung geht weiter | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Getrennte Ordner | Ja | Teilweise | Konto-Pools reichen Anfragen über einen lokalen Proxy an das nächste Konto weiter | Linux, macOS, Windows |
| [clausona](https://github.com/larcane97/clausona) | Getrennte Ordner | Ja | Ja | Routen wählen vor jedem Lauf ein Konto mit Kontingent, können auf eine Zurücksetzung warten und verteilen verschiedene Konten an Worker | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Tauscht die Anmeldung; getrennte Ordner mit `clauth start` | Ja, mit `clauth start` | Ja | Wechselt entlang einer Fallback-Kette | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Tauscht die Anmeldung oder isolierte Profile | Ja, mit isolierten Profilen | Ja | `caam run` wiederholt nicht interaktive Befehle mit einem anderen Konto | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Tauscht die Anmeldung und umhüllt `claude` | Nein, ein aktives Konto | Ja | Tauscht die Anmeldung und setzt dasselbe Gespräch fort | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Tauscht die Anmeldung | Nein | Ja | Optionaler Wechsel ab einem Schwellenwert | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Tauscht die Anmeldung | Nein | Ja | Optionaler Wechsel bei 100 % | macOS |

ccswap ist ein Fork von claude-swap, der auch Codex-Anmeldungen verwaltet. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) und [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) geben ebenfalls jedem Konto einen eigenen Ordner. Mehrere dieser Tools unterstützen auch andere Coding-Agents wie Codex oder Gemini.

## Wie sie aufgebaut sind

| Tool | Geschrieben in | Während Claude Code läuft | Kopiert Anmeldungen hin und her | Was Konten teilen können |
| --- | --- | --- | --- | --- |
| ccam | Go, eine statische Binärdatei | Standardmäßig nichts: ccam ersetzt sich selbst durch `claude` (unter Windows, wo das nicht geht, wartet es darauf); mit eingeschaltetem `relay` oder `auto_switch` wartet ccam neben interaktiven Sitzungen | Nein | Gruppen: Jede nennt die Elemente, die sie teilt, und ihre Konten; ein Konto kann in mehreren sein |
| claude-swap | Python 3.12+, installiert mit uv oder pipx | Nichts mit `cswap run`; `cswap auto` und die Menüleisten-App laufen weiter | Ja, wenn es die Anmeldung tauscht | Mit `cswap run`: Einstellungen, `CLAUDE.md` und Erweiterungen zusammen, Verlauf optional; MCP-Server werden bei jedem Start kopiert |
| CCS | Node.js 18+ oder Bun | CCS bleibt als Elternprozess; Pools fügen einen lokalen Proxy hinzu | Nein für seine Kontoordner; Pools haben eigene Anmeldungen | Geteilt oder lokal pro Konto, Verlauf in Kontextgruppen; MCP-Server werden kopiert |
| clausona | TypeScript, Node.js 20+ | Einfaches `claude` läuft über eine Shell-Funktion; `clausona run` behält Node als Elternprozess; nach jedem Beenden ein Node-Prozess im Hintergrund | Nein, erneuert aber selbst die Tokens ungenutzter Konten, um deren Nutzung zu lesen | Alles außer kontoeigenen Dateien, von einem primären Konto aus; Verlauf optional; mit `claude mcp add` hinzugefügte Server werden nicht geteilt |
| clauth | Rust | clauth bleibt als Elternprozess | Ja, wenn es die Anmeldung tauscht | Das ganze `~/.claude` oder nichts |
| CAAM | Go | Hängt vom Modus ab; ein optionaler Daemon | Ja, im Vault-Modus | Durch den Modus festgelegt; MCP-Server werden bei jedem Start kopiert |
| cux | Go | Ein Wrapper bleibt um `claude` | Ja | Alles: ein Ordner |
| CCSwitcher, Claude Switcher | macOS-Apps (Swift, Python) | Die Menüleisten-App | Ja | Alles: ein Ordner |

## Was ccam gut macht

### Anmeldungen bleiben, wo Claude Code sie abgelegt hat

ccam kopiert nie eine Claude-Anmeldung von einem Ort an einen anderen. Die Anmeldung jedes Kontos bleibt in seinem eigenen Ordner, unter macOS im eigenen Schlüsselbundeintrag dieses Ordners, und Claude Code ist das Programm, das sie erneuert. Es gibt keine gespeicherte Kopie, die aus dem Takt geraten kann. Ein bestehendes `~/.claude`, das mit `ccam import` registriert wurde, behält den Schlüsselbundeintrag, den das einfache `claude` nutzt, und bleibt so in beiden Fällen angemeldet.

ccam erneuert ein abgelaufenes Token nur dann selbst, wenn du `ccam usage --refresh` ausführst. Es setzt vorher eine Sperre und verweigert den Vorgang, solange das Konto läuft, egal wie Claude Code gestartet wurde (auch Hintergrundsitzungen und ihr Daemon), und solange Claude Code selbst diese Anmeldung erneuert, rotiert also nie ein Token unter einer laufenden Sitzung.

### Standardmäßig läuft nichts zwischen dir und Claude Code

ccam liest eine Datei, prüft seine Links und ersetzt sich selbst durch `claude` (`exec`). Standardmäßig gibt es, während du arbeitest, keinen ccam-Prozess, keinen Wrapper, keinen Proxy und keinen Daemon: Standardeingabe, Signale und Exit-Code gehören Claude Code, und seine Anfragen gehen direkt an Anthropic. Windows kann einen Prozess nicht ersetzen, deshalb startet ccam dort Claude Code mit derselben Konsole und wartet auf dessen Exit-Code, ohne sonst etwas zu tun. Mit `relay = true` oder `auto_switch = true` wartet ccam auf dieselbe Weise neben interaktiven Sitzungen, damit es ein Gespräch als ein anderes Konto weiterführen kann; die Anfragen von Claude Code gehen weiterhin direkt an Anthropic.

Ein Konto so zu starten kostet etwa 4 ms, unter 6 ms beim 95. Perzentil, bei 8 MB Spitzenspeicher ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) mit hyperfine, 300 Durchläufe). Auf demselben Laptop braucht allein der Start des Python-Interpreters 13 ms und der von Node.js 31 ms, bevor ein Programm überhaupt etwas tut.

### Gruppen bestimmen, was geteilt wird

MCP-Server, `settings.json`, `CLAUDE.md`, Erweiterungen (Agents, Befehle, Skills, Plugins) und Gesprächsverlauf sind fünf getrennte Elemente. Eine Gruppe nennt die Elemente, die sie teilt, und die Konten in ihr. So können zwei berufliche Konten den Verlauf in einer Gruppe `work` teilen, während ein privates Konto seinen eigenen behält, und alle drei teilen über `default` dieselben MCP-Server. `ccam group set work --add memory` verschiebt jedes Konto der Gruppe in einem Plan; das [Dashboard](/features/dashboard) tut dasselbe mit Kontrollkästchen und graut ein Element aus, das ein Konto bereits von einer anderen Gruppe bekommt.

- MCP-Server werden beim Start mit `--mcp-config` übergeben, daher lässt ein Start `.claude.json`, die Datei, die Claude Code während des Laufs neu schreibt, standardmäßig unverändert. ccam schreibt dort nur in Schritten, die das ankündigen, unter der Sperre von Claude Code und mit einer Sicherung: für Remote Control nach einer Rückfrage, wenn ein Konto das Teilen von MCP verlässt oder wieder aufnimmt, und für `ccam mcp import`. `ccam mcp add` fügt einen Server einmal für die ganze Gruppe hinzu.
- Die anderen Elemente sind Links in den Gruppenordner, sodass eine Änderung in einem Konto sofort in den anderen erscheint. Beim Start wird nichts kopiert.
- Anmeldungen, `.claude.json` und der Laufzeitzustand werden nie geteilt.

### Jede Änderung wird vorher gezeigt und gesichert

`ccam group` und `ccam migrate` zeigen einen Plan und fragen nach, bevor sie etwas verschieben. Wenn sich eine Datei zwischen einem Konto und seiner Gruppe unterscheidet, entscheidest du: die der Gruppe behalten, deine verwenden oder das JSON zusammenführen; in Ordnern bleibt die Kopie der Gruppe, und der Plan nennt die Dateien, die sich unterscheiden. Originale landen in `.ccam-backup/<time>/`, ein entferntes Konto wird abgemeldet und sein Ordner kommt in einen Papierkorb-Ordner, und Dateien werden mit einer einzigen Umbenennung ersetzt, sodass ein unterbrochener Schreibvorgang die alte Datei hinterlässt. `--dry-run` zeigt, was Starten, Teilen, Migrieren, Verknüpfungen, Updates und Reparaturen tun würden. `ccam doctor --fix` repariert kaputte Links und Berechtigungen.

### Dasselbe Gespräch im nächsten Konto

`ccam next --auto --continue` fragt alle Konten gleichzeitig nach ihrer Nutzung, öffnet das mit dem meisten Spielraum und setzt das Gespräch fort. Mit geteiltem Verlauf ist es dasselbe Gespräch; ohne kopiert Claude Code es in das neue Konto (`--resume <file> --fork-session`), und ccam sagt, wie viele Tokens die erste Antwort erneut liest. Wenn ein Limit eine Sitzung anhält, zeigt der Hook des ccam-Plugins eine Benachrichtigung, und `/ccam:switch` gibt das Gespräch weiter. Standardmäßig entscheidest du, wann gewechselt wird. Mit eingeschaltetem `auto_switch` (Linux, macOS) verschiebt ccam die Sitzung selbst, aber nur, nachdem Claude Code sie an einem Limit gestoppt hat, und nur zu einem Konto wie dem, das sie verlässt: einem anderen Platz derselben Organisation oder einem anderen privaten Tarif ([Details](/features/next#auto-switch)). Nichts überwacht deine Nutzung im Hintergrund, und standardmäßig verschieben automatische Wahlen Arbeit nie aus einer Team- oder Enterprise-Organisation heraus oder in eine hinein, sodass ein Firmengespräch nie versehentlich in einem privaten Konto landet. Die Statuszeile schlägt den Befehl vor, sobald ein Konto 90 % erreicht.

Skripte können nicht auf dich warten, deshalb wechselt `ccam --fallback work -p "…"` für sie: Wenn das Konto sein Limit erreicht, läuft derselbe Befehl erneut mit dem Konto mit dem meisten Spielraum, und nur die Ausgabe des Durchlaufs, der fertig wurde, wird ausgegeben. Ein Fehler, der kein Nutzungslimit ist, wird unverändert weitergegeben. Für parallele Worker gibt `ccam pick` jedem ein anderes Konto mit Spielraum, und beide beenden sich mit Code 75, wenn kein Konto mehr Spielraum hat.

**Team- und Enterprise-Plätze.** Die automatische Wahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und [`auto_switch`](/features/next#auto-switch), wenn du es einschaltest) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation. Im Hilfe-Center von Anthropic ([Team-Tarife](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); ebenso für [Enterprise-Tarife](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)) heißt es: „Each seat on a Team plan is assigned to one named individual, and only that person may use it“ (jeder Platz in einem Team-Tarif ist einer namentlich benannten Person zugewiesen, und nur diese Person darf ihn nutzen) und „Users can't route their work through another user's account to use that user's remaining usage“ (Benutzer dürfen ihre Arbeit nicht über das Konto eines anderen Benutzers leiten, um dessen verbleibendes Nutzungskontingent zu verwenden). Wenn du die automatische Wahl zwischen Plätzen nutzt, kann das daher gegen die Bedingungen von Anthropic verstoßen. Ob du sie dort nutzt, ist deine Entscheidung: Kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt. Für private Tarife (Pro, Max) gilt der [Hinweis zu den Bedingungen in den FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

### Das richtige Konto in jedem Ordner

`ccam folder work ~/code/company` sorgt dafür, dass jeder Start ohne Kontonamen in diesem Ordner `work` verwendet, während deine übrigen Projekte beim Standardkonto bleiben. Die Einstellung steht in derselben `config.toml`, die jeder Start liest, und kostet deshalb keine Zeit.

### Eine Binärdatei, geprüfte Downloads

ccam ist eine statische Binärdatei von etwa 10 MB, ohne Laufzeitumgebung, die installiert werden müsste, für Linux, macOS und Windows. Die Installationsskripte (`install.sh`, `install.ps1`) und `ccam update` prüfen jeden Download gegen die SHA-256-Prüfsummen des Releases, und ccam aktualisiert sich nie, ohne zu fragen. Es gibt keine Telemetrie: Die einzigen Anfragen sind Nutzungsabfragen bei Anthropic für deine Konten, parallel mit 3 Sekunden Zeitlimit und 60 Sekunden Cache, die Frage des Dashboards an GitHub nach dem neuesten Release, wenn es sich öffnet, und die Downloads, die du anforderst. Mit `usage_source = "claude"` laufen sogar diese über Claude Code, und ccam sendet kein Token.

## Wo andere Tools mehr können

- claude-swap, clauth, cux, CAAM und CCSwitcher können eine laufende interaktive Sitzung von selbst wechseln, wenn ein Konto einen Nutzungsschwellenwert überschreitet. ccam wechselt `-p`-Durchläufe mit `--fallback` und verschiebt eine interaktive Sitzung nur dann von selbst, wenn `auto_switch` eingeschaltet ist (Linux, macOS), und nur, nachdem Claude Code sie an einem Limit gestoppt hat; andernfalls verschiebst du sie, mit einer Benachrichtigung am Limit und `/ccam:switch`.
- Unter Windows ohne Entwicklermodus teilt ccam Ordner, kann aber keine einzelnen Dateien wie `settings.json` teilen (das erste `ccam add` teilt dann nur MCP-Server und Erweiterungen); claude-swap und CCS kopieren sie dort stattdessen.
- CCSwitcher, Claude Switcher und claude-swap zeigen die Nutzung in der macOS-Menüleiste.
- CCS, clausona, CAAM, clauth, ccswap und Claude Switcher verwalten auch Codex und andere Coding-Agents.
- clausona schaltet außerdem Erweiterungen pro Konto ein und aus (ccams `ccam ext` zeigt sie nur an), erfasst lokale Token-Kosten und liefert Claude Code-Plugins mit, die Flotten von Workern in tmux und anderen Terminals ausführen.
- CAAM, CCSwitcher und Claude Switcher lassen sich mit Homebrew installieren, und clauth und CAAM signieren ihre Release-Prüfsummen. ccam installierst du mit seinem Skript, einem `.deb` oder `go install`, und es prüft Prüfsummen ohne Signatur.

## Was passt zu dir?

Nimm ein Tool, das die Anmeldung austauscht, wenn du jeweils ein Konto nutzt oder laufende Sitzungen schon vor dem Limit oder durch Austausch der Anmeldung wechseln sollen. Nimm ccam, wenn du Konten nebeneinander, eine gemeinsame Einrichtung, dasselbe Gespräch im nächsten Konto und, standardmäßig, nichts willst, das neben Claude Code mitläuft, unter Linux, macOS oder Windows. Wenn du nichts teilen musst, reichen [einfache Shell-Aliase](/compare/aliases).

Ob du mehrere Konten nutzen darfst, hängt von den Bedingungen des jeweiligen Kontos ab. Prüfe sie, bevor du loslegst.
