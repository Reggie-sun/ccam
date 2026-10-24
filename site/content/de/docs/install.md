---
title: ccam installieren
metaTitle: "ccam unter Linux, macOS und Windows installieren"
description: "ccam unter Linux, macOS oder Windows per Befehl, .deb-Paket oder Quellcode installieren und mit ccam update aktualisieren. Claude Code vorher installieren."
updated: 2026-10-23
order: 1
translation: machine
sourceHash: af21dcc24169
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

Es lädt `ccam_windows_<arch>.zip` herunter, prüft seine SHA-256-Prüfsumme, legt `ccam.exe` in `%USERPROFILE%\.local\bin` ab, den Ordner, in den das Installationsprogramm von Claude Code `claude.exe` legt, und fügt diesen Ordner deinem Benutzer-`PATH` hinzu, falls er dort noch fehlt, wobei Einträge wie `%USERPROFILE%\bin` unverändert bleiben. Öffne danach ein neues Terminal. Führst du es erneut aus, ersetzt es ccam, auch während noch eine Claude Code-Sitzung mit der alten Binärdatei läuft (diese wird in `ccam.exe.old` umbenannt, oder in `ccam.exe.old-<number>`, wenn eine ältere noch in Gebrauch ist), und erneuert Kurzbefehle, die Hardlinks sind: `OK 2 shortcuts now run the new ccam`.

ccam startet Claude Code für dich, also installiere zuerst [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) und stelle sicher, dass `claude --version` funktioniert.

## Unterstützte Systeme

- Linux auf x86-64 und ARM64: Ubuntu, Debian, Fedora, Arch, Alpine und andere Distributionen. Die Binärdatei ist statisch und hat keine Abhängigkeiten.
- macOS auf Apple Silicon und Intel.
- Windows 10 und 11 auf x86-64 und ARM64, in PowerShell, Windows Terminal oder Git Bash. Geteilte Ordner funktionieren für jeden Benutzer. Geteilte Dateien (`settings.json`, `CLAUDE.md`, `history.jsonl`) brauchen den Entwicklermodus (Einstellungen > System > Für Entwickler), weil Windows nur in diesem Modus oder einem Administrator erlaubt, Links auf Dateien anzulegen; siehe [Teilen unter Windows](/features/share#on-windows). Ohne ihn legt das erste `ccam add` die Gruppe `default` nur mit MCP-Servern und Erweiterungen an und sagt, wie du den Rest hinzufügst, sobald der Entwicklermodus eingeschaltet ist.
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

Ein Build eines Release-Tags (`@latest` oder `@vX.Y.Z`) meldet diese Version, sodass `ccam update` dafür funktioniert; ein Build von `@main` oder eines lokalen Checkouts nach einem Tag meldet `dev`.

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

## Claude Code-Plugin

Das ccam-Plugin bringt Claude Code selbst bei, ccam zu verwenden: zu prüfen, welches Konto Spielraum hat, parallele `claude -p`-Worker mit `ccam pick` auf verschiedenen Konten zu starten und einen fehlenden Skill oder MCP-Server mit `ccam ext` zu erklären. Installiere es in Claude Code:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

Es fügt einen Skill hinzu, die Befehle `/ccam:usage`, `/ccam:ext` und `/ccam:switch` sowie zwei Hooks: Wenn eine Sitzung an einem Nutzungslimit anhält, bekommst du eine Terminal-Benachrichtigung, und `/ccam:switch [account]` übergibt das Gespräch an ein anderes Konto ([Details](/features/next#inside-claude-code)). Die Hooks führen `ccam` aus, daher muss es im `PATH` liegen, den Claude Code sieht. Der Skill führt nie Befehle aus, die sich anmelden oder Dateien verschieben; er schlägt sie stattdessen vor.

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

Auch das Dashboard (`ccam` im Terminal) sucht bei jedem Öffnen nach einem neuen Release und bietet es an, wenn es eines gibt: Drücke `U` und bestätige. Kein anderer Befehl sucht danach, und ccam aktualisiert sich nie, ohne zu fragen. Um ein bestimmtes Release zu installieren, auch ein älteres, gib `--version` an. `ccam --dry-run update` zeigt, was heruntergeladen und ersetzt würde, ohne etwas zu verändern.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Wenn du das Installationsskript erneut ausführst, wird die Kopie in `~/.local/bin` ebenfalls aktualisiert. Hast du ccam mit `go install` aus einem Release-Tag gebaut, funktioniert auch `ccam update`, oder du führst denselben Befehl noch einmal aus.

**Windows.** `ccam update` lädt das `.zip` für Windows herunter. Ein laufendes Programm kann unter Windows nicht ersetzt werden, deshalb wird die alte Binärdatei in `ccam.exe.old` umbenannt, die das nächste Update löscht; solange noch eine Sitzung mit einer früheren `ccam.exe.old` läuft, lautet der Name `ccam.exe.old-<time>`. Kurzbefehle, die Hardlinks sind, werden gleich danach erneuert.

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

Wenn du `~/.local/share/ccam` löschst, werden auch die dort gespeicherten Anmeldungen und Gesprächsverläufe gelöscht. Ein `~/.claude`-Ordner, den du [importiert](/docs/commands) hast, bleibt, wo er ist, doch solange er in einer Gruppe ist, liegen die Dateien, die er teilt, in `~/.local/share/ccam/groups/`. Führe zuerst `ccam remove <name>` für ihn aus (oder `ccam group leave`): Er verlässt seine Gruppen mit einer Kopie von allem, was er geteilt hat.
