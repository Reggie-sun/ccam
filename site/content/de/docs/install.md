---
title: ccam installieren
metaTitle: "ccam unter Linux, macOS und Windows installieren"
description: "ccam unter Linux, macOS oder Windows per Befehl, .deb-Paket oder go install installieren und mit ccam update aktualisieren. Claude Code vorher installieren."
updated: 2026-10-28
order: 1
translation: machine
sourceHash: "193df30c9d9a"
---

ccam startet Claude Code für dich, also installiere zuerst [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) und stelle sicher, dass `claude --version` funktioniert.

Führe dies in einem Terminal unter Linux oder macOS aus:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Das Skript lädt das Release für dein System herunter, prüft seine SHA-256-Prüfsumme und legt `ccam` in `~/.local/bin` ab, wo ccam später auch die Kurzbefehle anlegt, die du wählst, etwa `ccwork`. Wenn dieser Ordner nicht in deinem `PATH` liegt, gibt das Skript die Zeile aus, die du zu `~/.bashrc` oder `~/.zshrc` hinzufügen musst.

Unter Windows führst du dies in PowerShell aus:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Es lädt `ccam_windows_<arch>.zip` herunter, prüft seine SHA-256-Prüfsumme und legt `ccam.exe` in `%USERPROFILE%\.local\bin` ab, den Ordner, in den das Installationsprogramm von Claude Code `claude.exe` legt. Liegt dieser Ordner noch nicht in deinem Benutzer-`PATH`, fügt das Skript ihn hinzu und lässt Einträge wie `%USERPROFILE%\bin` unverändert. Öffne danach ein neues Terminal.

## Unterstützte Systeme

- Linux auf x86-64 und ARM64: Ubuntu, Debian, Fedora, Arch, Alpine und andere Distributionen. Die Binärdatei ist statisch und hat keine Abhängigkeiten.
- macOS auf Apple Silicon und Intel.
- Windows 10 und 11 auf x86-64 und ARM64, in PowerShell, Windows Terminal oder Git Bash:
  - Geteilte Ordner funktionieren ohne Entwicklermodus.
  - Geteilte Dateien (`settings.json`, `CLAUDE.md`, `history.jsonl`) brauchen den Entwicklermodus (Einstellungen > System > Für Entwickler), weil Windows nur in diesem Modus oder einem Administrator erlaubt, Links auf Dateien anzulegen; siehe [Teilen unter Windows](/features/share#on-windows).
  - Ohne Entwicklermodus legt das erste `ccam add` die Gruppe `default` nur mit MCP-Servern und Erweiterungen an und sagt, wie du den Rest hinzufügst, sobald der Entwicklermodus eingeschaltet ist.
- Innerhalb von WSL funktioniert ccam so wie auf jedem anderen Linux-System.

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

Installierst du ein Release-Tag (`@latest` oder `@vX.Y.Z`), meldet ccam diese Version, und `ccam update` sowie das Update-Angebot des Dashboards funktionieren. Ein Build von `@main` oder eines lokalen Checkouts mit Commits nach einem Tag meldet `dev`, und `ccam update` lehnt ihn ab, außer du gibst `--version` an.

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

Es fügt Folgendes hinzu:

- einen Skill, der Befehle vorschlägt, die anmelden oder Dateien verschieben, sie aber nie selbst ausführt;
- die Befehle `/ccam:usage`, `/ccam:ext` und `/ccam:switch`;
- zwei Hooks ([Details](/features/next#inside-claude-code)). Einer zeigt eine Terminal-Benachrichtigung, wenn eine Sitzung an einem Nutzungslimit anhält; unter Linux und macOS verschiebt ccam die Sitzung dann bei eingeschaltetem [`auto_switch`](/features/next#auto-switch) zu einem anderen Konto mit freiem Kontingent. Der andere lässt `/ccam:switch [account]` das Gespräch an ein anderes Konto übergeben.

Die Hooks führen `ccam` aus, daher muss es im `PATH` liegen, den Claude Code sieht.

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

Auch das Dashboard (`ccam` ohne Argumente, in einem Terminal) sucht bei jedem Öffnen nach einem neuen Release und bietet es an, wenn es eines gibt: Drücke `U` und bestätige. Kein anderer Befehl sucht nach Updates, und ccam aktualisiert sich nie, ohne zu fragen. Um ein bestimmtes Release zu installieren, auch ein älteres, gib `--version` an. `ccam --dry-run update` zeigt, was heruntergeladen und ersetzt würde, ohne etwas zu verändern.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Wenn du das Installationsskript erneut ausführst, wird die Kopie in `~/.local/bin` ebenfalls aktualisiert. Hast du ccam mit `go install` aus einem Release-Tag gebaut, kannst du `ccam update` oder denselben `go install`-Befehl noch einmal ausführen.

**Windows.** `ccam update` lädt das `.zip` für Windows herunter. Ein laufendes Programm kann unter Windows nicht ersetzt werden, deshalb wird die alte Binärdatei in `ccam.exe.old` umbenannt, die das nächste Update löscht; solange noch eine Sitzung mit einer früheren `ccam.exe.old` läuft, lautet der Name `ccam.exe.old-<time>`. Kurzbefehle, die Hardlinks sind, werden gleich danach erneuert. Ein erneutes Ausführen von `install.ps1` funktioniert genauso, auch während noch eine Claude Code-Sitzung mit der alten Binärdatei läuft, und gibt `OK 2 shortcuts now run the new ccam` aus, wenn es Kurzbefehle erneuert, die Hardlinks sind.

**Paket für Debian und Ubuntu.** `ccam update` ersetzt kein `/usr/bin/ccam`, das über apt installiert wurde, und gibt stattdessen diese Schritte aus. Lade das neue `ccam_<version>_linux_amd64.deb` von der [Release-Seite](https://github.com/jhyunwoo/ccam/releases) herunter und installiere es über das alte:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Deinstallieren

`ccam uninstall` entfernt ccam und alles, was es angelegt hat. Es listet zuerst auf, was wegkommt, und fragt nach:

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- Jedes Konto wird so entfernt, wie `ccam remove` es tut. Ein Konto in den Ordnern von ccam wird zuerst abgemeldet und mit seinem Ordner gelöscht, samt Gesprächsverlauf und Einstellungen.
- Ein `~/.claude`-Ordner, den du [importiert](/docs/commands) hast, bleibt für das normale `claude`, wo er ist. Ist er in einer Gruppe, bekommt er eine Kopie von allem zurück, was er geteilt hat. Die ccam-Statuszeile in seiner `settings.json` weicht der Statuszeile, die ccam umschlossen hat, falls es eine gab. Hat er das ccam-Plugin, gibt ccam den Befehl `/plugin uninstall` aus, den du in Claude Code ausführst.
- Die Kurzbefehle, `~/.config/ccam`, `~/.local/share/ccam` (samt dem Papierkorb von ccam) und `~/.local/state/ccam` werden gelöscht, danach der Befehl `ccam` selbst.

ccam verweigert das, solange ein Konto läuft; beende zuerst seine Claude Code-Sitzungen. Ein `/usr/bin/ccam`, das apt installiert hat, bleibt; entferne es danach mit `sudo apt remove ccam`. Unter Windows kann sich ein laufendes Programm nicht selbst löschen, daher gibt `ccam uninstall` zum Schluss den Befehl `Remove-Item` für `ccam.exe` aus.
