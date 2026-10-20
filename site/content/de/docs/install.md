---
title: ccam installieren
metaTitle: "ccam unter Linux, macOS und Windows installieren"
description: "ccam unter Linux, macOS oder Windows per Befehl, .deb-Paket oder Quellcode installieren und mit ccam update aktualisieren. Claude Code vorher installieren."
updated: 2026-10-20
order: 1
translation: machine
sourceHash: 07eadc0769c8
---

Führe dies in einem Terminal unter Linux oder macOS aus:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Das Skript lädt das Release für dein System herunter, prüft seine SHA-256-Prüfsumme und legt `ccam` in `~/.local/bin` ab, wo ccam später auch die Kurzbefehle anlegt, die du wählst, etwa `ccwork`. Wenn dieser Ordner nicht in deinem `PATH` liegt, gibt das Skript die Zeile aus, die du zu `~/.bashrc` oder `~/.zshrc` hinzufügen musst.

Unter Windows führst du dies in PowerShell aus:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Es lädt `ccam_windows_<arch>.zip` herunter, prüft seine SHA-256-Prüfsumme, legt `ccam.exe` in `%USERPROFILE%\.local\bin` ab, den Ordner, in den das Installationsprogramm von Claude Code `claude.exe` legt, und fügt diesen Ordner deinem Benutzer-`PATH` hinzu, falls er dort noch fehlt. Öffne danach ein neues Terminal.

ccam startet Claude Code für dich, also installiere zuerst [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) und stelle sicher, dass `claude --version` funktioniert.

## Unterstützte Systeme

- Linux auf x86-64 und ARM64: Ubuntu, Debian, Fedora, Arch, Alpine und andere Distributionen. Die Binärdatei ist statisch und hat keine Abhängigkeiten.
- macOS auf Apple Silicon und Intel.
- Windows 10 und 11 auf x86-64 und ARM64, in PowerShell, Windows Terminal oder Git Bash. Geteilte Ordner funktionieren für jeden Benutzer. Geteilte Dateien (`settings.json`, `CLAUDE.md`, `history.jsonl`) brauchen den Entwicklermodus (Einstellungen > System > Für Entwickler), weil Windows nur in diesem Modus oder einem Administrator erlaubt, Links auf Dateien anzulegen; siehe [Teilen unter Windows](/features/share#on-windows).
- Innerhalb von WSL funktioniert ccam wie auf jedem Linux-System.

## Andere Installationswege

**Paket für Debian und Ubuntu.** Lade `ccam_<version>_linux_amd64.deb` von der [Release-Seite](https://github.com/jhyunwoo/ccam/releases) herunter und installiere es:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Das Paket installiert `/usr/bin/ccam`.

**Aus dem Quellcode.** Mit Go 1.26 oder neuer:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**Eine bestimmte Version.** Setze `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Shell-Vervollständigung

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Die Vervollständigung deckt alle Befehle und Optionen von `ccam` ab, einschließlich der Kontonamen.

## Installation prüfen

```sh
ccam version
ccam doctor
```

`ccam doctor` prüft, ob Claude Code gefunden wird und ob jeder Kontoordner, jeder Link und jede Anmeldung in Ordnung ist. Weiter geht es mit dem [Schnellstart](/docs/quick-start).

## Aktualisieren

```sh
ccam update --check
ccam update
```

`ccam update --check` zeigt nur an, ob es ein neueres Release gibt. `ccam update` lädt das Release für dein System herunter, prüft seine SHA-256-Prüfsumme anhand der `checksums.txt` des Releases und ersetzt dann die Datei `ccam` in einem Schritt. Stimmt die Prüfsumme nicht, wird nichts verändert. Kurzbefehle wie `ccwork` sind Links auf ccam und starten deshalb sofort die neue Version.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam sucht nie von selbst nach Updates, sondern erst, wenn du diesen Befehl ausführst. Um ein bestimmtes Release zu installieren, auch ein älteres, gib `--version` an. `ccam --dry-run update` zeigt, was heruntergeladen und ersetzt würde, ohne etwas zu verändern.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Wenn du das Installationsskript erneut ausführst, wird die Kopie in `~/.local/bin` ebenfalls aktualisiert. Hast du ccam mit `go install` gebaut, führe denselben Befehl noch einmal aus.

**Windows.** `ccam update` lädt das `.zip` für Windows herunter. Ein laufendes Programm kann unter Windows nicht ersetzt werden, deshalb wird die alte Binärdatei in `ccam.exe.old` umbenannt, die das nächste Update löscht.

**Paket für Debian und Ubuntu.** `ccam update` ersetzt kein `/usr/bin/ccam`, das über apt installiert wurde, und gibt stattdessen diese Schritte aus. Lade das neue `ccam_<version>_linux_amd64.deb` von der [Release-Seite](https://github.com/jhyunwoo/ccam/releases) herunter und installiere es über das alte:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Deinstallieren

Entferne zuerst die Kurzbefehle der Konten, solange ccam noch installiert ist. `ccam list` zeigt sie in der Spalte Shortcut; entferne jeden mit `ccam shortcut <account> --remove` oder lösche diese Links selbst aus `~/.local/bin`. Entferne dann den Befehl und die Datenordner, wenn du die Konten nicht mehr brauchst:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Unter Windows, in PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Wenn du `~/.local/share/ccam` löschst, werden auch die dort gespeicherten Anmeldungen und Gesprächsverläufe gelöscht. Ein `~/.claude`-Ordner, den du [importiert](/docs/commands) hast, bleibt, wo er ist.
