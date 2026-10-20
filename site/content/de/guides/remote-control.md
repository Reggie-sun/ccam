---
title: Remote Control mit ccam nutzen
description: "Die Remote Control von Claude Code mit ccam <name> rc als beliebiges Konto starten: geltende Optionen, warum ccam MCP-Server kopieren will, Fehlerbehebung."
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 5f0b90ea53d6
---

`ccam work rc` startet die Remote Control von Claude Code als das Konto `work`, sodass du Sitzungen auf diesem Computer von claude.ai/code oder der mobilen Claude-App aus steuern kannst. Es funktioniert wie `claude rc`, mit einem Unterschied: Teilt das Konto MCP-Server, fragt ccam einmal, ob sie in die `.claude.json` des Kontos kopiert werden sollen, denn Sitzungen der Remote Control lesen Server von nirgendwo sonst.

## Starten

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

Ein [Kurzbefehl](/features/shortcuts) funktioniert genauso: `ccwork rc` ist `ccam work rc`. Führe den Befehl in dem Ordner aus, in dem du arbeiten willst, und öffne dann claude.ai/code oder die Claude-App, angemeldet mit demselben Claude-Konto wie `work`.

## Optionen

Hier treffen zwei Arten von Optionen aufeinander: die von ccam, die vor dem Kontonamen stehen oder in den Einstellungen von ccam liegen, und die der Remote Control, die nach `rc` stehen und unverändert bei Claude Code ankommen.

### Optionen und Einstellungen von ccam

| Option oder Einstellung | Was sie tut | Standard |
| --- | --- | --- |
| `--dry-run` (vor dem Konto) | Gibt Umgebung und Befehl aus, dazu den MCP-Hinweis, falls es einen gibt. Fragt nicht und kopiert nichts | aus |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | Kommt, wenn das Konto MCP-Server teilt und seiner `.claude.json` die aktuelle Liste der Gruppe fehlt. `y` kopiert sie, `N` oder Enter startet ohne sie | `N` |
| `mcp_strategy = "sync"` in `~/.config/ccam/config.toml` | Kopiert die Server der Gruppe vor jedem Start in `.claude.json`, ob Remote Control oder nicht, sodass nie gefragt wird | `"flag"` |
| `ccam group leave <group> work` | Nimmt das Konto aus der Gruppe, die MCP mit ihm teilt; Remote Control startet dann nur mit den eigenen Servern des Kontos, und es wird nicht gefragt | in `default` |

### Optionen der Remote Control

Alles nach `rc` gehört Claude Code. Dies sind die Optionen von `claude remote-control` in Claude Code 2.1.304; `ccam work rc --help` zeigt die Liste für deine Version.

| Option | Was sie tut |
| --- | --- |
| `--name <name>` | Name der Sitzung, der in claude.ai/code angezeigt wird |
| `--remote-control-session-name-prefix <prefix>` | Präfix für automatisch erzeugte Sitzungsnamen (Standard: der Hostname) |
| `-c`, `--continue` | Wieder mit der Sitzung verbinden, die Remote Control zuletzt für diesen Ordner gespeichert hat, innerhalb von etwa 4 Stunden |
| `--session-id <id>` | Wieder mit einer bestimmten Sitzung verbinden |
| `--permission-mode <mode>` | Berechtigungsmodus für neue Sitzungen: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` oder `plan` |
| `--spawn <mode>` | `same-dir` (Standard), `worktree` (ein git worktree pro Sitzung) oder `session` (eine Sitzung, beendet sich mit ihr) |
| `--capacity <N>` | Höchstzahl gleichzeitiger Sitzungen im Modus `same-dir` und `worktree` (Standard 32) |
| `--[no-]create-session-in-dir` | Beim Start eine erste Sitzung im aktuellen Ordner anlegen (Standard an) |
| `--[no-]chrome` | Claude in Chrome für neue Sitzungen (Standard aus) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Debug- und ausführliche Ausgabe |

Schreib `rc` direkt hinter den Kontonamen. ccam nimmt das erste Argument, das nicht mit `-` beginnt, als Befehl von Claude Code, also sieht es in `ccam work rc --name laptop` das `rc`, in `ccam work --model opus rc` aber `opus`.

## Was ccam tut

Der Start nimmt den üblichen Weg ([Claude Code starten](/features/launch)) mit einer Änderung im MCP-Schritt:

1. ccam findet das Konto, prüft seine Freigabelinks und löst `claude` auf, wie bei jedem Start.
2. Es erkennt, dass der Befehl `rc` oder `remote-control` ist, und lässt `--mcp-config` weg.
3. Teilt das Konto MCP-Server, vergleicht es die `mcp.json` der Gruppe mit dem `sync-stamp` in `~/.local/state/ccam/mcp/work/`. Enthält `.claude.json` bereits die aktuellen Server oder hat die Gruppe keine, geht es direkt zu Schritt 6.
4. Andernfalls nennt es die Server, die den Sitzungen fehlen würden:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. Bei `y` sichert es `.claude.json`, nimmt die Sperre von Claude Code, ergänzt die Server der Gruppe in `mcpServers` und behält die übrigen Server des Kontos. Bei `N`, ohne Terminal oder während das Konto läuft, ändert es nichts und sagt das:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Es vermerkt die Sitzung, setzt `CLAUDE_CONFIG_DIR` auf den Kontoordner und ersetzt sich per `exec` durch `claude rc`.

Remote Control läuft danach als Server weiter. Jede Sitzung, die sie für dich startet, ist ein Kindprozess davon und läuft deshalb als `work`: mit dem Login, den Einstellungen und dem Gesprächsverlauf von `work`, und sie zählt gegen das Nutzungslimit von `work`. Solange der Server läuft, betrachtet ccam das Konto als laufend.

Der Probelauf zeigt den Unterschied. Ein normaler Start übergibt die Server der Gruppe als Argument, Remote Control nicht:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## Warum es so funktioniert

**Remote Control nimmt kein `--mcp-config`.** Bei normalen Starts übergibt ccam die Server der Gruppe als `--mcp-config=<path>` an Claude Code und schreibt nie in die Dateien von Claude Code. `claude remote-control` lehnt diese Option vor seinem Befehlsnamen ab und nimmt danach keine an, deshalb schlug `ccam work rc` vor ccam 0.6.0 für jedes Konto fehl, das MCP-Server teilte. Jetzt lässt ccam die Option weg.

**Ihre Sitzungen lesen Server nur aus `.claude.json`.** Die Sitzungen, die Remote Control startet, erben keine Optionen von der Befehlszeile. Der einzige Ort, an dem sie Benutzer-MCP-Server finden, ist die `.claude.json` des Kontos. Diese Datei enthält auch die Anmeldedaten des Kontos und den Zustand pro Ordner, deshalb kann ccam sie nicht wie `settings.json` auf eine geteilte Kopie verlinken. Die Server hineinzukopieren ist der einzige Weg, den Sitzungen die geteilte Liste zu geben.

**ccam fragt, bevor es schreibt.** Mit der Standard-Flag-Strategie verspricht ccam, `.claude.json` in Ruhe zu lassen. Remote Control ist der eine Fall, in dem es dort schreiben muss, also fragt es zuerst, behält eine Sicherung und rührt nichts außer `mcpServers` an.

**Es schreibt nicht, während das Konto läuft.** Claude Code schreibt `.claude.json` während einer Sitzung neu, und eine Änderung von außen ginge verloren oder würde überschrieben. Darum überspringt ccam das Kopieren, solange irgendeine Sitzung des Kontos läuft, auch ein anderer Remote-Control-Server.

**Es fragt einmal pro Änderung.** Nach einer Kopie enthält `sync-stamp` die Änderungszeit der `mcp.json` der Gruppe. ccam fragt erst wieder, wenn sich diese Datei ändert, zum Beispiel nach `ccam mcp add`. Antwortest du `N`, gibt es noch keine Kopie, also fragt es bei jedem Start der Remote Control.

**`ccam mcp` hält die Kopien aktuell.** Ein Eintrag in `.claude.json` gewinnt bei gleichem Namen gegen `--mcp-config`, daher würde eine alte Kopie auch bei normalen Starts eine neuere Definition verdecken. `ccam mcp add`, `add-json`, `remove` und `import` aktualisieren deshalb die Kopien jedes Kontos, das gerade nicht läuft:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Ohne Terminal lautet die Antwort Nein.** Remote Control läuft oft aus einem Skript, einem Dienst oder `nohup`, wo niemand antworten kann. ccam startet dann ohne die Server, statt zu warten.

## Probleme und Lösungen

### Die geteilten MCP-Server fehlen in Sitzungen der Remote Control

Die Frage wurde mit `N` beantwortet, es gab kein Terminal, oder das Konto lief. Beende Remote Control, schließ die anderen Sitzungen des Kontos und führe `ccam work rc` erneut in einem Terminal aus; antworte `y`. Damit nie gefragt wird, setze `mcp_strategy = "sync"` in `~/.config/ccam/config.toml`. Die Server einer laufenden Sitzung ändern sich nicht, starte nach dem Kopieren also eine neue Sitzung.

### "Account work is running, so its .claude.json is left alone"

Eine Sitzung von `work` ist noch offen: ein normales `ccam work`, ein Kurzbefehl oder ein anderer Remote-Control-Server. Beende sie und führe `ccam work rc` erneut aus. Sitzungen, die ein anderes Konto ausführt, zählen nicht.

### ccam fragt jedes Mal

Du hast die Server noch nicht kopiert. Antworte einmal `y` oder setze `mcp_strategy = "sync"`. Wenn du die Server in Remote Control nie haben willst, beende das Teilen von MCP mit dem Konto: Nimm es aus dieser Gruppe (`ccam group leave default work`), oder halte MCP in einer eigenen Gruppe, in der das Konto nicht ist.

### `ccam work rc` schlägt sofort mit einem Fehler zu `--mcp-config` fehl

Dein ccam ist älter als 0.6.0 und übergibt `--mcp-config` an Remote Control. Führe `ccam update` aus. Derselbe Fehler mit einem aktuellen ccam bedeutet, dass vor `rc` eine Option mit Wert stand, wie in `ccam work --model opus rc`, sodass ccam `opus` für den Befehl hielt; schreib `rc` direkt hinter den Kontonamen.

### Eine Änderung mit `ccam mcp` kommt nicht bei Remote Control an

`ccam mcp add` und `remove` aktualisieren Kopien nur in Konten, die nicht laufen, und sagen das: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Beende Remote Control und starte sie neu; ccam erkennt die geänderte `mcp.json` und fragt erneut.

### Die Kopie hat einen eigenen Server ersetzt

Hatte das Konto einen eigenen Server mit demselben Namen wie ein Server der Gruppe, hat die Kopie ihn ersetzt. Die vorherige `.claude.json` liegt in `~/.local/state/ccam/mcp/work/backups/`, wo die fünf neuesten aufbewahrt werden. Um deine Version für alle Konten zu behalten, verschiebe sie vor dem nächsten Start der Remote Control mit `ccam mcp import work <server> --force` in die Gruppe.

### Ich will die kopierten Server entfernen

`ccam mcp remove <server>` entfernt einen Server aus der Gruppe und aus jeder Kopie. Um zum Zustand vor der Kopie zurückzukehren, stelle die neueste Sicherung aus `~/.local/state/ccam/mcp/work/backups/` wieder her, während keine Sitzung des Kontos läuft.

### Remote Control fragt, ob dem Ordner vertraut werden soll, oder beendet sich ohne Terminal

Das Vertrauen in Ordner wird pro Konto in `.claude.json` gespeichert, ein Ordner, dem du als `personal` vertraut hast, ist für `work` also neu. Führe in diesem Ordner einmal `ccam work` aus und bestätige, dann starte Remote Control dort.

### "You must be logged in with a Claude account that has a subscription"

Remote Control braucht eine Anmeldung mit einem Claude-Abonnement. Prüfe die E-Mail-Adresse des Kontos in `ccam list` und melde dich bei Bedarf mit `ccam login work` erneut an.

### Die Sitzung erscheint nicht auf dem Handy

Die Claude-App oder claude.ai/code ist mit einem anderen Claude-Konto angemeldet. Vergleiche es mit der E-Mail-Adresse, die `ccam list` für `work` anzeigt.

### ccam warnt, dass ein anderes Claude Code in diesem Ordner den Gesprächsverlauf teilt

Ein anderes Konto derselben Verlaufsgruppe läuft im selben Ordner. Sitzungen, die von beiden Konten aus dasselbe Gespräch fortsetzen, können den Verlauf durcheinanderbringen. Antworte nur dann `y`, wenn sie an verschiedenen Gesprächen arbeiten, oder starte Remote Control in einem anderen Ordner.

## Verwandte Themen

- [Claude Code starten](/features/launch)
- [Geteilte MCP-Server](/features/mcp)
- [Teilen zwischen Konten](/features/share)
- [Konten verwalten](/features/accounts)
