---
title: Geteilte MCP-Server
description: Mit ccam mcp add, add-json, list, remove und import eine MCP-Serverliste pro Gruppe führen, per --mcp-config übergeben oder in jedes Konto synchronisiert.
updated: 2026-10-26
order: 5
translation: machine
sourceHash: 27a06ad35295
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
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
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
only in account side, for /home/me/code/app only: notes
```

Unter den Gruppen nennt `ccam mcp list` jeden Server, der die anderen Konten nicht erreicht:

- `only in account work: pg` für einen Server, den ein Konto für sich selbst hinzugefügt hat;
- `only in account side, for /home/me/code/app only: notes` für einen, der für einen einzelnen Ordner hinzugefügt wurde (einfaches `claude mcp add`, siehe unten);
- `account work has its own …, used in its sessions in place of the group's`, wenn ein eigener Server eines Kontos den Namen eines Gruppenservers trägt;
- die `.mcp.json` dieses Ordners, mit den Konten, deren Freigabe sie `waiting for approval in` ist;
- Server, die `in no account, only in ~/.claude.json (plain claude's)` sind, zusammen mit dem Befehl [`ccam migrate --mcp`](/features/migrate), der sie kopiert.

Fügst du einen Namen hinzu, den die Gruppe bereits hat, wird dessen Definition ersetzt, und es erscheint `updated` statt `added`. Server erreichen ein Konto nur, wenn es in dieser Gruppe ist und die Gruppe MCP teilt (`ccam group set <group> --add mcp`); tut die Gruppe das nicht, weist `add` darauf hin.

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
| `--dry-run` (global) | Nur von `mcp list` unterstützt; `add`, `add-json`, `remove` und `import` lehnen die Option ab und ändern nichts | aus |

`add` braucht entweder einen Befehl nach `--` oder `--url`; sonst schlägt es mit `a command to run (after --) or --url is required` fehl. `--header` gilt für `--url`-Server und `--env` für Befehlsserver. `add-json` nimmt dasselbe JSON wie `claude mcp add-json`, und die Definition muss ein Objekt mit einem `command` oder einer `url` sein. `import` verwendet immer die Gruppe, über die das Konto MCP teilt, und `list` zeigt jede Gruppe, daher ändert `-g` daran nichts.

## So funktioniert es

Die beteiligten Dateien:

- `~/.local/share/ccam/groups/<group>/mcp.json`: die geteilten Server, `{"mcpServers": {…}}`, das Format, das `--mcp-config` liest. Sie wird mit `fsutil.AtomicWrite` und Modus `0600` geschrieben, weil Servereinträge oft API-Schlüssel enthalten, und nur unter der Sperre `<group>/.mcp.json.lock`, sodass zwei gleichzeitig ausgeführte `ccam mcp add`-Befehle beide ihre Server behalten.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: die eigene Datei von Claude Code mit den Servern des Kontos. ccam liest nur ihren Wert `mcpServers` mit dem Scanner `profile.TopLevelValue` und parst auf dem Startpfad nie die ganze Datei.
- `~/.local/state/ccam/mcp/<name>/`: die Datei `seen`, der `sync-stamp` (mit einem Hash jeder Kopie, die ccam geschrieben hat, damit es eine unveränderte Kopie von einer unterscheiden kann, die du bearbeitet hast), `sync-checked` (was die Strategie sync zuletzt gesehen hat), `kept` (die Server, die das Konto behalten hat, als es das Teilen von MCP verlassen hat) und Sicherungen von `.claude.json`.

Es gibt zwei Strategien, festgelegt durch `mcp_strategy` in `~/.config/ccam/config.toml`. Standard ist die Flag-Strategie:

<!-- motion:mcp-flag -->

1. ccam findet die MCP-Gruppe des Kontos (`launcher.applyMCP`) und prüft mit einem einzigen `stat`, ob die `mcp.json` der Gruppe existiert, ohne sie zu parsen (`mcp.FlagArgs`).
2. ccam setzt `--mcp-config=<path>` direkt hinter `claude`, vor deine eigenen Argumente, sodass dein eigenes `--mcp-config` oder `--strict-mcp-config` zusätzlich weiterhin gilt (`Plan.InsertArgs`).
3. ccam prüft, ob sich die `.claude.json` des Kontos seit dem letzten Start geändert hat, und parst nur dann ihre Server (`mcp.NewLocal`).
4. ccam ersetzt sich selbst durch Claude Code (`syscall.Exec`).

Der Pfad wird als ein Argument übergeben, `--mcp-config=<path>`, nicht als zwei. `--mcp-config` von Claude Code nimmt mehrere Werte, daher würde `claude --mcp-config a.json "prompt"` den Prompt als zweite Konfigurationsdatei lesen und mit `MCP config file not found` fehlschlagen. Die Form mit einem Argument vermeidet das, sodass ein Prompt oder Unterbefehl danach funktioniert. `ccam work mcp list`, `get` und `remove` erhalten die Option nicht, weil Claude Code sie dort ablehnt; ccam nennt stattdessen die Server der Gruppe auf der Standardfehlerausgabe ([Details](/features/launch#shared-mcp-servers)). Claude Code lädt Server aus `--mcp-config` zusätzlich zu den eigenen des Kontos; bei gleichem Namen gewinnt der Eintrag in `.claude.json`. Die Flag-Strategie schreibt nur dann in `.claude.json`, wenn es um Remote Control geht, wenn ein Konto das Teilen von MCP verlässt oder ihm wieder beitritt ([unten](#leaving-and-joining)), und bei `ccam mcp import`.

`claude remote-control` (`rc`) nimmt kein `--mcp-config`, und seine Sitzungen lesen Server nur aus `.claude.json`. Bei `ccam work rc` lässt ccam die Option weg und fragt in einem Terminal, ob die Server der Gruppe in `.claude.json` kopiert werden sollen, mit Sicherung, so wie es die Strategie `sync` tut. Es fragt, bis eine Kopie existiert, und danach erst wieder, wenn sich die `mcp.json` der Gruppe ändert oder ein kopierter Server fehlt ([Remote Control](/features/launch#remote-control)). Weil eine solche Kopie bei gleichem Namen gegen `--mcp-config` gewinnt, aktualisieren `ccam mcp add`, `add-json`, `remove` und `import` die Kopien, die andere Konten der Gruppe haben:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Die Zeile `✓ Updated` erscheint nur, wenn sich eine Kopie geändert hat. Eine Kopie, die eine ältere Version von ccam angelegt hat, wird auf dieselbe Weise aktualisiert, ebenso eine nach `ccam migrate`.

Ein Server, der innerhalb einer Sitzung mit `claude mcp add` hinzugefügt wird, landet nur in der `.claude.json` dieses Kontos. Einfaches `claude mcp add` verwendet den lokalen Geltungsbereich: Der Server wird für dieses eine Konto in diesem einen Ordner geladen (`-s user` macht ihn für das Konto in jedem Ordner verfügbar). ccam kann das nicht sehen, wenn es passiert, weil Claude Code den ccam-Prozess bereits ersetzt hat, daher prüft es beim nächsten Start. Die Datei `seen` enthält die Änderungszeit von `.claude.json` und die bereits gemeldeten Servernamen. Weicht die Zeit ab, vergleicht ccam die Namen mit der Gruppe und gibt für jeden neuen Server einmal einen Hinweis aus, ohne auf eine Eingabe zu warten, bei beiden Strategien:

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

Der Hinweis nennt auch einen eigenen Server eines Kontos, der den Namen eines Gruppenservers trägt, da der des Kontos in dessen Sitzungen gewinnt.

Mit `mcp_strategy = "sync"` erscheinen die Server in `claude mcp list` und `/mcp` als gewöhnliche Benutzer-Server, und ccam schreibt sie vor dem Start in die `.claude.json` jedes Kontos:

<!-- motion:mcp-sync -->

1. ccam vergleicht die Änderungszeit der `mcp.json` der Gruppe mit `sync-stamp` und die der `.claude.json` des Kontos mit dem, was es zuletzt gesehen hat (`mcp.Sync`). Hat sich keine von beiden geändert, passiert nichts weiter. Andernfalls prüft es zusätzlich, ob jeder Gruppenserver noch in `.claude.json` steht, und stellt eine fehlende Kopie wieder her (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Es nimmt die eigene Sperre von Claude Code, den Ordner `.claude.json.lock` neben der Datei, angelegt mit `mkdir`. Es wartet bis zu 2 Sekunden und behandelt eine Sperre, die älter als 10 Sekunden ist, als veraltet, wie es die Sperrbibliothek von Claude Code tut.
3. Es sichert die aktuelle `.claude.json` in `backups/` im Statusordner und behält die neuesten 5.
4. Es setzt jeden Gruppenserver in `mcpServers`, entfernt Server, die es früher geschrieben hat und die die Gruppe nicht mehr hat, und schreibt die Datei atomar mit ihrem ursprünglichen Modus, nur wenn sich etwas geändert hat. Es ändert oder entfernt nur die Kopien, die es selbst geschrieben hat und die seitdem niemand geändert hat. Hat das Konto einen eigenen Server mit demselben Namen wie ein Gruppenserver, behält ccam den Server des Kontos und sagt das einmal: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`.
5. Es speichert den neuen Zeitstempel mit einem Hash jeder Kopie, und Claude Code startet ohne `--mcp-config`.

Die Synchronisierung wird mit `ccam: account work is already running, skipping MCP sync` übersprungen, wenn eine Sitzung des Kontos läuft, egal wie Claude Code gestartet wurde, weil Claude Code `.claude.json` während der Laufzeit neu schreibt. Eine fehlgeschlagene Synchronisierung gibt `ccam: MCP sync failed (launching anyway): …` aus und startet trotzdem.

### Verlassen und Beitreten

Ein Konto, das das Teilen von MCP beendet (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`), behält die Server der Gruppe: Der Plan kopiert sie mit Sicherung in seine `.claude.json`, und seine eigenen Server gleichen Namens bleiben. Mit `--fresh` erhält es keine, und der Plan listet die Server auf, die nicht mehr geladen werden. Tritt es wieder einer Gruppe bei, die diese Server hat, werden die behaltenen Kopien entfernt, die niemand geändert hat, sodass spätere Änderungen der Gruppe es erreichen.

`ccam mcp import` verschiebt die eigenen Server eines Kontos in seine Gruppe:

<!-- motion:mcp-import -->

ccam liest die Server des Kontos, seine eigenen und die, die es nur für einen Ordner hat, und vergleicht jeden mit dem gleichnamigen Server der Gruppe über `mcp.SameDef`, das Schlüsselreihenfolge und Leerraum ignoriert. Ein Ordnerserver wird zu einem Gruppenserver für jeden Ordner (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); ein Name mit zwei unterschiedlichen Ordnerdefinitionen wird abgelehnt. Server aus Plugins werden nicht angeboten, da das Plugin sie mitbringt. Ein Server, den die Gruppe nicht hat oder mit derselben Definition hat, wird kopiert. Ein abweichender wird übersprungen, weil die Definition der Gruppe diejenige ist, die jedes andere Konto verwendet:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Mit `--force` wird der Server der Gruppe ersetzt, und die vorherige `mcp.json` wird zuerst in `.ccam-backup/<time>/mcp.json` im Gruppenordner gesichert. Dann entfernt ccam die kopierten Einträge aus der `.claude.json` des Kontos, die eigenen und die des Ordners in einem Schreibvorgang, unter derselben Sperre und nach einer Sicherung im Statusordner (`mcp.Remove`). ccam hält die `mcp.json` der Gruppe von dem Moment an gesperrt, in dem es die Datei liest, bis diese Einträge entfernt sind, sodass niemand den Server zwischendurch aus der Gruppe entfernt. Das ist nötig, weil der Eintrag des Kontos bei gleichem Namen über den der Gruppe gewinnen würde. Mit `--keep-local` bleiben die Einträge in der `.claude.json` des Kontos.

## Sicherheit und Sonderfälle

- **Anmeldungen.** Server, die sich per OAuth anmelden, speichern ihre Tokens in der `.credentials.json` jedes Kontos (`mcpOAuth`), nicht in der Serverliste, daher teilt das Teilen eines Servers nicht seine Anmeldung. Melde dich in jedem Konto einmal mit `/mcp` an. ccam kopiert nie Anmeldungen, und [`ccam migrate`](/features/migrate) tut es auch nicht: Ein kopiertes Refresh-Token meldet jeden anderen Inhaber ab. [`ccam ext`](/features/ext) und `ccam doctor` nennen die Server, bei denen ein Konto nicht angemeldet ist, ein anderes aber schon, und `ccam next --continue` sagt, welche dem neuen Konto fehlen. Die Gruppe erhält nie eine Anmeldedatei.
- **Laufende Konten.** `import` kopiert in die Gruppe, lässt aber die eigenen Einträge eines laufenden Kontos stehen und weist darauf hin; führe es nach dem Ende der Sitzung erneut aus. Die Synchronisierung überspringt laufende Konten. Beide erkennen auch Claude Code-Sitzungen, die ohne ccam gestartet wurden.
- **Kopien für Remote Control.** Die Kopie eines laufenden Kontos bleibt unberührt: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Mit `mcp_strategy = "sync"` muss nichts aktualisiert werden, weil jeder Start synchronisiert.
- **Beschädigte Dateien.** Ist `.claude.json` kein gültiges JSON, lässt die Synchronisierung die Datei unberührt und meldet den Fehler. Lässt sich die Sperre nicht innerhalb von 2 Sekunden nehmen, meldet ccam `could not lock … (Claude Code is writing the file)`.
- **Konten, die MCP nicht teilen.** `import` lehnt ab und rät dir, das Konto zuerst in eine Gruppe aufzunehmen, die MCP teilt (`ccam group join <group> side`).
- **Gruppennamen.** `-g` nimmt nur einen Gruppennamen; er wird zu einem Ordner unter `groups/`, daher werden Namen wie `../x` abgelehnt.

## Verwandte Themen

- [Teilen mit Gruppen](/features/share)
- [Remote Control mit ccam nutzen](/guides/remote-control)
- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Eine vorhandene Einrichtung kopieren](/features/migrate)
- [Befehlsreferenz](/docs/commands)
