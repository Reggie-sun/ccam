---
title: Claude Code starten
description: Claude Code mit ccam <name|number>, ccam run oder der ccam-Auswahl als beliebiges Konto starten; Argumente gehen unverändert an claude, und --dry-run zeigt den Befehl.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d033da86f44b
---

`ccam work` startet Claude Code als das Konto `work`, und alles nach dem Namen geht unverändert an `claude`. Verwende es überall dort, wo du `claude` tippen würdest: in einem Terminal, in einem Skript oder mit weitergeleiteter Eingabe.

## Verwendung

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

Gib den Namen eines Kontos oder seine Nummer aus `ccam list` an. ccam liest die Argumente nach dem Namen nicht, daher zeigt `ccam work --help` die Hilfe von Claude Code. Ein `--` nach dem Namen ist optional. Standardeingabe, Exit-Code und Signale verhalten sich wie bei `claude`.

Beginnt das erste Argument mit `-` und ist keine von ccams eigenen Optionen, startet ccam das Standardkonto mit allen Argumenten. Das Standardkonto ist das mit `ccam use` festgelegte Konto, oder das Konto mit der niedrigsten Nummer, wenn du keines festgelegt hast.

### Die Auswahl

`ccam` ohne Argumente öffnet die Auswahl, wenn Standardeingabe und -ausgabe ein Terminal sind, und gibt sonst `ccam list` aus. Sie zeigt für jedes Konto Nummer, E-Mail-Adresse, 5-Stunden-Nutzung, letzte Verwendung und geteilte Elemente; die Nutzung erscheint, sobald das jeweilige Konto antwortet.

| Taste | Wirkung |
| --- | --- |
| `Enter` | Claude Code als das ausgewählte Konto starten |
| `1` bis `9` | Das Konto mit dieser Nummer starten |
| `Up`/`Down`, `k`/`j` | Die Auswahl bewegen |
| `a` | Die Auswahl schließen und den Assistenten von `ccam add` ausführen |
| `s` | Die geteilten Elemente des ausgewählten Kontos bearbeiten (`Space` oder `x` schaltet um, `Enter` übernimmt, `Esc` geht zurück) |
| `d` | Das ausgewählte Konto zum Standardkonto machen |
| `r` | Die Nutzung erneut abrufen, ohne den kurzen Cache |
| `q`, `Esc`, `Ctrl+C` | Beenden |

Die Auswahl schließt sich, bevor sie handelt, sodass das Terminal wiederhergestellt ist, bevor Claude Code startet. Nach `d` oder einer Änderung am Teilen öffnet sie sich erneut. Eine Änderung am Teilen zeigt zuerst den Plan und fragt nach, wie `ccam share`.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--dry-run` | Umgebung und Befehl ausgeben, die ccam ausführen würde, dann beenden | aus |
| `-v`, `--version` | `ccam <version>` ausgeben (nur als erstes Argument) | |
| `-h`, `--help` | Die Hilfe von ccam, als `ccam --help` oder `ccam run --help` | |

Regeln für Argumente:

- ccam liest seine eigenen Optionen nur vor dem Kontonamen. Nach dem Namen gehört jedes Argument Claude Code, auch `--dry-run` und `--help`.
- `--dry-run` funktioniert als `ccam --dry-run work`, `ccam run --dry-run work` und `ccam --dry-run` (das Standardkonto). Die Option wird beachtet, nicht abgelehnt.
- `ccam run <name|number> [claude args]` ist die Langform von `ccam <name>`. `ccam run` ohne Namen verwendet das Standardkonto.
- Kontonamen dürfen keine ccam-Befehlsnamen wie `add`, `list`, `run` oder `usage` sein, daher ist `ccam <word>` nie mehrdeutig.

Der Probelauf gibt eine Zeile aus, die du in eine Shell einfügen könntest:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## So funktioniert es

Jeder Start durchläuft dieselben Schritte, von der Befehlszeile bis `exec`:

<!-- motion:launch -->

Jedes Feld leuchtet auf, während sein Schritt läuft, und die Zeilen darunter sind die Umgebung und die Argumente, mit denen Claude Code startet.

1. `cmd/ccam/main.go` ruft `app.Main` auf. Wurde ccam unter einem anderen Namen gestartet, ist es ein [Kurzbefehl](/features/shortcuts) und geht direkt an den Starter.
2. `app.parse` wertet `os.Args` aus, bevor cobra beteiligt ist. Nur `ccam` allein (die Auswahl), ein ccam-Befehlsname (`config.Reserved`), `-h`/`--help` und `-v`/`--version` gehen an `cli.Execute`. Alles andere ist ein Start.
3. `launcher.Launch` liest `~/.config/ccam/config.toml`, die einzige Datei, die es parst.
4. `launcher.Prepare` findet das Konto über Namen oder Nummer, entfernt ein führendes `--` und ermittelt die ausführbare Datei `claude` mit `ResolveClaude`. Die Suche im PATH wird zusammen mit der Änderungszeit der Datei in `~/.local/state/ccam/claude-path` zwischengespeichert, sodass ccam nicht bei jedem Start `PATH` durchsucht. Eine geänderte oder fehlende Datei löst eine neue Suche aus. Ein `claude_bin`-Pfad in `config.toml` wird unverändert verwendet. ccam lehnt ein `claude` ab, das ccam selbst ist (Vergleich über den Inode), weil das eine Endlosschleife wäre.
5. `checkLinks` vergleicht den Kontoordner `~/.local/share/ccam/profiles/<name>/` mit seinen Freigabeeinstellungen über `share.Check`, das nur `lstat` und `readlink` aufruft, ein oder zwei Systemaufrufe pro Element. Ein fehlender Link wird sofort angelegt. Alles, was Daten kosten könnte, hält den Start an.
6. `applyMCP` fügt die geteilten MCP-Server hinzu (siehe unten).
7. `session.Register` erfasst die Sitzung. ccam schreibt eine nach seiner PID benannte Markierung mit der Startzeit des Prozesses nach `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (oder nach `~/.local/state/ccam/locks/`, wenn `XDG_RUNTIME_DIR` fehlt). Bei einem Konto, das den Gesprächsverlauf teilt, sucht `historyCheck` zuerst in `history-<group>-<hash>` nach einer laufenden Sitzung, wobei der Hash aus dem aktuellen Ordner gebildet wird, und registriert sich auch dort. `exec` behält die PID, daher bezeichnet die Markierung danach den Claude Code-Prozess selbst. Eine Markierung statt einer gehaltenen Sperre, weil Prozesse, die Claude Code startet, eine Sperre erben und sie nach dem Ende der Sitzung behalten würden.
8. `Plan.Exec` setzt den Terminaltitel auf `claude · <name>` (wenn die Standardfehlerausgabe ein Terminal ist), schreibt die Startzeit ohne `fsync` nach `~/.local/state/ccam/last-used/<name>` und ruft `syscall.Exec` auf.

Claude Code läuft mit `CLAUDE_CONFIG_DIR` auf den Kontoordner und `CCAM_PROFILE` auf den Kontonamen gesetzt. ccam entfernt geerbte Werte dieser Variablen. Für einen mit `ccam import` registrierten Ordner `~/.claude` unter macOS setzt es außerdem `CLAUDE_SECURESTORAGE_CONFIG_DIR` auf einen leeren Wert, damit das Konto und einfaches `claude` dieselbe Anmeldung im Schlüsselbund verwenden.

### Warum exec

`syscall.Exec` ersetzt den ccam-Prozess durch `claude`. ccam bleibt also nicht im Speicher, Tastendrücke, Signale und Exit-Code gehören Claude Code selbst, und Argumente bleiben ein Array, nie ein Shell-String. Arbeit nach einer Sitzung, etwa das Bemerken von Servern, die mit `claude mcp add` hinzugefügt wurden, geschieht beim nächsten Start: ccam vergleicht die Änderungszeit von `.claude.json` mit `~/.local/state/ccam/mcp/<name>/seen` und parst die Datei nur, wenn sie sich geändert hat.

### Wie eine Befehlszeile weitergeleitet wird

Drei Befehlszeilen, drei Wege durch `internal/app`:

<!-- motion:launch-route -->

Ein Kurzbefehl wird über `argv[0]` erkannt, `ccam <name>` über das erste Argument, und beide enden in `syscall.Exec`. Nur ein ccam-Befehl wie `ccam usage` erreicht cobra.

### Geteilte MCP-Server

Mit der Standardstrategie `flag` übergibt ccam die MCP-Server der Gruppe als Argument an Claude Code, vor deinen eigenen:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` gibt nur dann ein Argument zurück, wenn die `mcp.json` der Gruppe existiert, und `Plan.InsertArgs` setzt es direkt hinter `claude`. Dein eigenes `--mcp-config` oder `--strict-mcp-config` gilt dann zusätzlich weiterhin. Es ist ein einziges Argument, `--mcp-config=<path>`, weil `--mcp-config` von Claude Code mehrere Werte nimmt: Als zwei Argumente geschrieben, würde es auch einen folgenden Prompt als Konfigurationsdatei lesen.

Mit `mcp_strategy = "sync"` in `config.toml` schreibt ccam die Server stattdessen in die `.claude.json` des Kontos und lässt das aus, wenn das Konto bereits läuft.

### Das Budget von 10 ms

ccam fügt höchstens 10 ms (p95) hinzu, bevor Claude Code startet. Der Startpfad baut nie den cobra-Befehlsbaum auf, startet nie Bubble Tea und öffnet keine Netzwerkverbindung, und `TestPackageInitIsCheap` schlägt fehl, wenn die `init()`-Funktionen der Pakete zusammen mehr als 5 ms brauchen. `scripts/perf-check.sh` misst `ccam --dry-run work`, `ccam work` und einen Kurzbefehl mit `hyperfine` und schlägt oberhalb von 10 ms fehl, und `BenchmarkLaunchPath` in `internal/launcher` führt die Schritte in einer Schleife aus. Um die Zeit pro Schritt zu sehen, setze `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Sicherheit und Sonderfälle

- Ein unbekannter Name oder eine unbekannte Nummer hält mit `no account "<name>"` an und verweist auf `ccam list`. Gibt es gar keine Konten, bittet ccam dich, mit `ccam add <name>` eines hinzuzufügen.
- Ist ein geteilter Pfad zu einer gewöhnlichen Datei geworden oder zeigt er an die falsche Stelle, hält ccam mit `stopped because of share link problems:` an, listet die Probleme auf und bittet dich, `ccam doctor --fix` auszuführen. Auf dem Startpfad verschiebt es keine Daten.
- Läuft im selben Ordner ein anderes Claude Code, das den Gesprächsverlauf mit diesem Konto teilt (dieselbe Gruppe), warnt ccam, dass sich der Verlauf vermischen kann, und fragt `Run anyway? [y/N]`. Ohne Terminal zeigt es die Warnung und macht weiter.
- `--dry-run` überspringt die MCP-Synchronisierung, den Hinweis auf neue Server, die Verlaufsprüfung, die Sitzungsmarkierung, den Titel und die Zeit der letzten Verwendung und führt kein exec aus. Die Freigabelinks prüft es trotzdem und legt einen fehlenden Link an wie ein normaler Start.
- Mit `mcp_strategy = "sync"` gibt eine fehlgeschlagene Synchronisierung `ccam: MCP sync failed (launching anyway): …` aus, und Claude Code startet mit den Servern, die es bereits hatte.

## Verwandte Themen

- [Kurzbefehle](/features/shortcuts)
- [Konten verwalten](/features/accounts)
- [Befehlsreferenz](/docs/commands)
- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Wie CLAUDE_CONFIG_DIR funktioniert](/guides/claude-config-dir)
