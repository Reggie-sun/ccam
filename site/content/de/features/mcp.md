---
title: Geteilte MCP-Server
description: Mit ccam mcp add, add-json, list, remove und import eine Liste von MCP-Servern pro Gruppe führen, die mit --mcp-config übergeben oder in jedes Konto synchronisiert wird.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 69fd3233fd09
---

Claude Code speichert Benutzer-MCP-Server in `.claude.json`, derselben Datei, die die Anmeldung des Kontos enthält, daher kann diese Datei nicht zwischen Konten geteilt werden. ccam hält geteilte Server in einer Gruppendatei, `mcp.json`, und übergibt sie beim Start an jedes Konto der Gruppe. Mit `ccam mcp` fügst du einen Server einmal für alle deine Konten hinzu, und mit `ccam mcp import` verschiebst du einen Server, den du innerhalb einer Sitzung hinzugefügt hast, zu den anderen.

## Verwendung

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

Eine Änderung gilt ab dem nächsten Start jedes Kontos; laufende Sitzungen behalten die Server, mit denen sie gestartet sind:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

Fügst du einen Namen hinzu, den die Gruppe bereits hat, wird dessen Definition ersetzt, und es erscheint `updated` statt `added`. Server erreichen ein Konto nur, wenn es MCP mit dieser Gruppe teilt (`ccam share <name> --mcp on`); teilt noch kein Konto MCP, weist `add` darauf hin.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `-g`, `--group <group>` | Zielgruppe für `add`, `add-json` und `remove` | die einzige Gruppe, über die ein Konto MCP teilt, sonst `default` |
| `-- <command> [args…]` (`add`) | Diesen Befehl als stdio-Server ausführen | |
| `--url <url>` (`add`) | URL eines entfernten Servers (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | Transport eines `--url`-Servers | `http` |
| `--header <'Name: value'>` (`add`) | HTTP-Header, mehrfach verwendbar | keiner |
| `-e`, `--env <KEY=VALUE>` (`add`) | Umgebungsvariable für einen stdio-Server, mehrfach verwendbar | keine |
| `--force` (`import`) | Einen abweichenden Gruppenserver gleichen Namens ersetzen; die alte `mcp.json` wird gesichert | aus |
| `--keep-local` (`import`) | Die Einträge in der `.claude.json` des Kontos behalten | aus |
| `-h`, `--help` | Hilfe zu jedem Befehl | |
| `--dry-run` (global) | Nur von `mcp list` beachtet; `add`, `add-json`, `remove` und `import` lehnen die Option ab und ändern nichts | aus |

`add` braucht entweder einen Befehl nach `--` oder `--url`; sonst schlägt es mit `a command to run (after --) or --url is required` fehl. `--header` gilt für `--url`-Server und `--env` für Befehlsserver. `add-json` nimmt dasselbe JSON wie `claude mcp add-json`, und die Definition muss ein Objekt mit einem `command` oder einer `url` sein. `import` verwendet immer die Gruppe, über die das Konto MCP teilt, und `list` zeigt jede Gruppe, daher ändert `-g` daran nichts.

## So funktioniert es

Die beteiligten Dateien:

- `~/.local/share/ccam/groups/<group>/mcp.json`: die geteilten Server, `{"mcpServers": {…}}`, das Format, das `--mcp-config` liest. Sie wird mit `fsutil.AtomicWrite` und Modus `0600` geschrieben, weil Servereinträge oft API-Schlüssel enthalten.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: die eigene Datei von Claude Code mit den Servern des Kontos. ccam liest nur ihren Wert `mcpServers` mit dem Scanner `profile.TopLevelValue` und parst auf dem Startpfad nie die ganze Datei.
- `~/.local/state/ccam/mcp/<name>/`: die Datei `seen`, der `sync-stamp` und Sicherungen von `.claude.json`.

Es gibt zwei Strategien, festgelegt durch `mcp_strategy` in `~/.config/ccam/config.toml`. Standard ist die Flag-Strategie:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` findet die MCP-Gruppe des Kontos. `mcp.FlagArgs` prüft mit einem einzigen `stat`, ob die `mcp.json` der Gruppe existiert, ohne sie zu parsen.
2. `Plan.InsertArgs` setzt `--mcp-config=<path>` direkt hinter `claude`, vor deine eigenen Argumente, sodass dein eigenes `--mcp-config` oder `--strict-mcp-config` zusätzlich weiterhin gilt.
3. `mcp.NewLocal` prüft, ob sich die `.claude.json` des Kontos seit dem letzten Start geändert hat, und parst nur dann ihre Server.
4. `syscall.Exec` ersetzt ccam durch Claude Code.

Der Pfad wird als ein Argument übergeben, `--mcp-config=<path>`, nicht als zwei. `--mcp-config` von Claude Code nimmt mehrere Werte, daher würde `claude --mcp-config a.json "prompt"` den Prompt als zweite Konfigurationsdatei lesen und mit `MCP config file not found` fehlschlagen. Die Form mit einem Argument endet an dieser Stelle, und ein folgender Prompt oder Unterbefehl wie `mcp list` funktioniert. Claude Code lädt Server aus `--mcp-config` zusätzlich zu den eigenen des Kontos; bei gleichem Namen gewinnt der Eintrag in `.claude.json`. Die Flag-Strategie schreibt nur für Remote Control in `.claude.json`.

`claude remote-control` (`rc`) nimmt kein `--mcp-config`, und seine Sitzungen lesen Server nur aus `.claude.json`. Bei `ccam work rc` lässt ccam die Option weg und fragt in einem Terminal, ob die Server der Gruppe in `.claude.json` kopiert werden sollen, mit Sicherung, so wie es die Strategie `sync` tut. Es fragt, bis eine Kopie existiert, und danach erst wieder, wenn sich die `mcp.json` der Gruppe ändert ([Remote Control](/features/launch#remote-control)). Weil eine solche Kopie bei gleichem Namen gegen `--mcp-config` gewinnt, aktualisieren `ccam mcp add`, `add-json`, `remove` und `import` die Kopien, die andere Konten der Gruppe haben:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Ein Server, der innerhalb einer Sitzung mit `claude mcp add` hinzugefügt wird, landet nur in der `.claude.json` dieses Kontos. ccam kann ihn zu diesem Zeitpunkt nicht sehen, weil Claude Code ccam ersetzt hat, daher prüft es beim nächsten Start. Die Datei `seen` enthält die Änderungszeit von `.claude.json` und die bereits gemeldeten Servernamen. Weicht die Zeit ab, vergleicht ccam die Namen mit der Gruppe und gibt für jeden neuen Server einmal einen Hinweis aus, ohne auf eine Eingabe zu warten:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

Mit `mcp_strategy = "sync"` erscheinen die Server in `claude mcp list` und `/mcp` als gewöhnliche Benutzer-Server, und ccam schreibt sie vor dem Start in die `.claude.json` jedes Kontos:

<!-- motion:mcp-sync -->

1. `mcp.Sync` vergleicht die Änderungszeit der `mcp.json` der Gruppe mit `sync-stamp`. Hat sie sich nicht geändert, wird die `.claude.json` des Kontos nicht geöffnet.
2. Es nimmt die eigene Sperre von Claude Code, den Ordner `.claude.json.lock` neben der Datei, angelegt mit `mkdir`. Es wartet bis zu 2 Sekunden und behandelt eine Sperre, die älter als 10 Sekunden ist, als veraltet, wie es die Sperrbibliothek von Claude Code tut.
3. Es sichert die aktuelle `.claude.json` in `backups/` im Statusordner und behält die neuesten 5.
4. Es setzt jeden Gruppenserver in `mcpServers`, entfernt Server, die es früher geschrieben hat und die die Gruppe nicht mehr hat, behält die eigenen Server des Kontos und schreibt die Datei atomar mit ihrem ursprünglichen Modus.
5. Es speichert den neuen Zeitstempel und die Namen, die es verwaltet, und Claude Code startet ohne `--mcp-config`.

Die Synchronisierung wird mit `ccam: account work is already running, skipping MCP sync` übersprungen, wenn eine Sitzung des Kontos läuft, weil Claude Code `.claude.json` während der Laufzeit neu schreibt. Eine fehlgeschlagene Synchronisierung gibt `ccam: MCP sync failed (launching anyway): …` aus und startet trotzdem. Der Hinweis auf innerhalb einer Sitzung hinzugefügte Server gehört nur zur Flag-Strategie.

`ccam mcp import` verschiebt die eigenen Server eines Kontos in seine Gruppe:

<!-- motion:mcp-import -->

ccam liest die Server des Kontos mit `mcp.LocalServers` und vergleicht jeden mit dem gleichnamigen Server der Gruppe über `mcp.SameDef`, das Schlüsselreihenfolge und Leerraum ignoriert. Ein Server, den die Gruppe nicht hat oder mit derselben Definition hat, wird kopiert. Ein abweichender wird übersprungen, weil die Definition der Gruppe diejenige ist, die jedes andere Konto verwendet:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Mit `--force` wird der Server der Gruppe ersetzt, und die vorherige `mcp.json` wird zuerst in `.ccam-backup/<time>/mcp.json` im Gruppenordner gesichert. Dann entfernt `mcp.RemoveLocal` die kopierten Einträge aus der `.claude.json` des Kontos, unter derselben Sperre und nach einer Sicherung im Statusordner. Das ist nötig, weil der Eintrag des Kontos bei gleichem Namen über den der Gruppe gewinnen würde. `--keep-local` lässt sie stehen.

## Sicherheit und Sonderfälle

- **Anmeldungen.** Server, die sich per OAuth anmelden, speichern ihre Tokens in der `.credentials.json` jedes Kontos (`mcpOAuth`), nicht in der Serverliste, daher teilt das Teilen eines Servers nicht seine Anmeldung. Melde dich in jedem Konto einmal mit `/mcp` an, oder kopiere die Anmeldungen mit [`ccam migrate`](/features/migrate) aus `~/.claude`. Die Gruppe erhält nie eine Anmeldedatei.
- **Laufende Konten.** `import` kopiert in die Gruppe, lässt aber die eigenen Einträge eines laufenden Kontos stehen und weist darauf hin; führe es nach dem Ende der Sitzung erneut aus. Die Synchronisierung überspringt laufende Konten.
- **Kopien für Remote Control.** Die Kopie eines laufenden Kontos bleibt unberührt: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Mit `mcp_strategy = "sync"` muss nichts aktualisiert werden, weil jeder Start synchronisiert.
- **Beschädigte Dateien.** Ist `.claude.json` kein gültiges JSON, lässt die Synchronisierung die Datei unberührt und meldet den Fehler. Lässt sich die Sperre nicht innerhalb von 2 Sekunden nehmen, meldet ccam `could not lock … (Claude Code is writing the file)`.
- **Konten, die MCP nicht teilen.** `import` lehnt ab und rät dir, zuerst `ccam share side --mcp on` auszuführen.
- **Gruppennamen.** `-g` nimmt nur einen Gruppennamen; er wird zu einem Ordner unter `groups/`, daher werden Namen wie `../x` abgelehnt.

## Verwandte Themen

- [Teilen zwischen Konten](/features/share)
- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Eine vorhandene Einrichtung kopieren](/features/migrate)
- [Alle Befehle](/docs/commands)
