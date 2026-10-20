---
title: Kurzbefehle
description: Mit ccam shortcut einem Konto einen eigenen Befehl wie ccwork geben und mit ccam shims install die Links in ~/.local/bin neu erstellen.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 939256e7a67e
---

Ein Kurzbefehl ist ein Befehl, den du selbst benennst und der Claude Code mit einem bestimmten Konto startet: `ccwork --resume` macht dasselbe wie `ccam work --resume`. Er ist ein Link auf ccam in `~/.local/bin` und funktioniert daher in bash, zsh, Skripten und überall sonst, wo Programme aus dem `PATH` ausgeführt werden. Kurzbefehle sind optional; lege einen für die Konten an, die du oft startest.

## Verwendung

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam bestätigt jede Änderung:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

Ein Konto hat höchstens einen Kurzbefehl, gespeichert als `shortcut` in `~/.config/ccam/config.toml`. Der Assistent von `ccam add` fragt `Create a shortcut command for <name>? [y/N]` und dann `Shortcut name [cc<name>]:`. Außerhalb des Assistenten legt `ccam add` nur mit `--shortcut` einen an. `ccam list` zeigt den Kurzbefehl jedes Kontos oder `-`.

Ein Kurzbefehl gibt jedes Argument an Claude Code weiter, genauso wie `ccam <name>` nach dem Namen. Das gilt auch für `--dry-run` und `--help`: Für eine Vorschau eines Starts verwende `ccam --dry-run work`.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--remove` | Den Kurzbefehl des Kontos und seinen Link löschen | aus |
| `-h`, `--help` | Hilfe zum Befehl | |
| `--dry-run` | `ccam shortcut` beachtet die Option: Es prüft den Namen und gibt `Would set the shortcut of work to ccwork.` oder `Would remove ccwork, the shortcut of work.` aus, ohne zu speichern oder zu verlinken | aus |

`ccam shims install` nimmt keine Argumente und lehnt `--dry-run` mit `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed `` ab.

Regeln für Namen, in dieser Reihenfolge geprüft:

- 1 bis 64 Zeichen aus `a-z`, `0-9`, `.`, `-` und `_`, beginnend mit einem Buchstaben oder einer Ziffer;
- nicht `ccam`, `claude` oder ein ccam-Befehl wie `run` oder `usage`;
- nicht der Kurzbefehl eines anderen Kontos;
- keine Datei in `~/.local/bin`, die nicht bereits ein Link auf ccam ist;
- nicht der Name eines anderen Programms in deinem `PATH`.

Wer einen Befehl und `--remove` zusammen angibt, erhält `give a command or --remove, not both`. Setzt du den Namen, den ein Konto bereits hat, erscheint `work: ccwork (unchanged)`, und der Link wird neu erstellt, falls er fehlt.

## So funktioniert es

Das Setzen eines Namens führt `checkShortcut` aus, das `shim.Check` mit dem bin-Ordner und dem echten Pfad der ccam-Binärdatei aufruft:

<!-- motion:shortcuts-check -->

Der Name durchläuft die Prüfungen der Reihe nach. Neben jeder Prüfung steht ein Name, den sie ablehnt: `Work` passt nicht zum Muster, `run` ist ein ccam-Befehl, `ccside` gehört einem anderen Konto, `ccfile` ist eine vorhandene Datei in `~/.local/bin`, und `git` würde `/usr/bin/git` verdecken. Nur ein Name, der jede Prüfung besteht, wird gespeichert und verlinkt.

1. `shim.Check` prüft den Namen gegen `^[a-z0-9][a-z0-9._-]{0,63}$`, dann gegen `config.Reserved`, dann mit `cfg.FindShortcut`. Es ruft `os.Lstat` für `~/.local/bin/<name>` auf, ruft für jeden anderen Ordner im `PATH` `os.Stat` auf und lehnt eine ausführbare Datei ab, die kein Link auf ccam ist.
2. ccam speichert `config.toml` mit `fsutil.AtomicWrite` (temporäre Datei, `fsync`, rename).
3. `shim.Remove` löscht den Link des alten Kurzbefehls, aber nur, wenn er ein Link auf ccam ist.
4. `shim.Install` legt für jeden konfigurierten Kurzbefehl, dessen Link fehlt, den Link an, mit `os.Symlink` auf die ccam-Binärdatei, ermittelt mit `os.Executable` und `filepath.EvalSymlinks`. Es entfernt nie etwas.

Ein Link gilt als ccams eigener, wenn `os.Readlink` den ccam-Pfad zurückgibt oder wenn Link und ccam-Pfad auf dieselbe Datei führen. Jedes Entfernen und jede Prüfung vor dem Überschreiben läuft über diesen Test, daher fasst ccam keine Dateien an, die es nicht selbst angelegt hat.

Wenn du einen Kurzbefehl ausführst, ist der eingetippte Name das Einzige, was ccam braucht:

<!-- motion:shortcuts-link -->

Die Shell findet `ccwork` im `PATH`, der Kernel folgt dem Link und startet die ccam-Binärdatei, und `argv[0]` ist weiterhin `ccwork`. In `app.Main` lädt jeder Name außer `ccam` einmal `config.toml` und ruft `shim.Lookup` auf, das das Konto findet, dessen `shortcut` gleich `filepath.Base(argv[0])` ist. `launcher.LaunchWith` durchläuft dann dieselben Schritte wie [`ccam <name>`](/features/launch) mit der bereits geladenen Konfiguration und gibt `argv[1:]` unverändert an `claude` weiter.

Warum ein symbolischer Link und kein Shell-Alias oder Wrapper-Skript: Ein symbolischer Link funktioniert in jeder Shell und in Skripten, er kostet keinen zusätzlichen Shell-Start (etwa 1 bis 2 ms), und der Start liest weiterhin nur eine Datei. `scripts/perf-check.sh` misst einen Kurzbefehl mit `hyperfine` gegen dasselbe Limit von 10 ms (p95) wie `ccam <name>`.

### Unter Windows

Ein Kurzbefehl ist `ccwork.exe` in `%USERPROFILE%\.local\bin`, und `shim.Name` entfernt vor der Suche das `.exe` (in beliebiger Groß- und Kleinschreibung), sodass `ccwork`, `ccwork.exe` und `CCWORK.EXE` alle `work` starten. ccam legt zuerst einen symbolischen Link an, den Windows im Entwicklermodus oder für einen Administrator erlaubt. Ohne dieses Recht legt es einen Hardlink an: einen zweiten Namen für dieselbe Datei `ccam.exe`, der keine besonderen Rechte braucht. Ein Hardlink zählt als ccams eigener, wenn er dieselbe Datei wie `ccam.exe` ist. Weil `ccam update` eine neue Datei an die Stelle von `ccam.exe` setzt, erneuert es danach jeden per Hardlink angelegten Kurzbefehl mit `shim.Relink` und gibt `✓ shortcut ccwork renewed` aus. Die Namensprüfung sucht in den anderen `PATH`-Ordnern nach `<name>` mit jeder Endung aus `PATHEXT` (`.exe`, `.cmd` und so weiter).

## Sicherheit und Sonderfälle

- **Ein Link, den kein Konto beansprucht.** Hat ein Link auf ccam in `~/.local/bin` einen Namen, den kein Konto verwendet, etwa nachdem du `shortcut` in `config.toml` von Hand bearbeitet hast, greift ccam nicht auf das Standardkonto zurück. Es hält mit `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` an. Setze den Namen mit `ccam shortcut` erneut oder lösche den Link.
- **Vorhandene Dateien.** ccam ersetzt nie eine Datei, die es nicht selbst angelegt hat. `ccam shims install` hält mit `<path> already exists and is not a ccam link, so it was left alone` an.
- **Eine verschobene ccam-Binärdatei.** Alte Links zeigen dann auf einen Pfad, den es nicht mehr gibt, und ccam betrachtet sie nicht als seine eigenen. Lösche sie und führe dann `ccam shims install` aus.
- **Umbenennen und entfernen.** `ccam rename` behält den Namen des Kurzbefehls, und der Link funktioniert weiter, weil er auf ccam zeigt und ccam den Namen in `config.toml` nachschlägt. `ccam remove` löscht den Link des Kontos. `ccam shortcut <account> <new>` löscht den alten Link, bevor es den neuen anlegt.
- **PATH.** Wenn ccam einen Link anlegt und `~/.local/bin` nicht in deinem `PATH` liegt, gibt es `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` aus, gefolgt von der `export PATH=...`-Zeile, die du hinzufügen musst. Unter Windows gibt es stattdessen die PowerShell-Zeile aus, die den Ordner deinem Benutzer-`PATH` hinzufügt.
- **Später prüfen.** `ccam doctor` meldet `missing shortcuts in /home/me/.local/bin: ccwork` und warnt, wenn der Ordner nicht im `PATH` liegt. `ccam doctor --fix` legt die fehlenden Links an.

## Verwandte Themen

- [Claude Code starten](/features/launch)
- [Konten verwalten](/features/accounts)
- [Befehlsreferenz](/docs/commands#shortcuts)
- [ccam installieren](/docs/install)
