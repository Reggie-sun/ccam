---
title: ccam installieren
description: ccam unter Linux oder macOS mit einem Befehl, aus einem .deb-Paket oder aus dem Quellcode installieren. Claude Code muss vorher installiert sein.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: fc63978e8f9e
---

Führe dies in einem Terminal unter Linux oder macOS aus:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Das Skript lädt das Release für dein System herunter, prüft seine SHA-256-Prüfsumme und legt `ccam` in `~/.local/bin` ab, wo `ccam add` später die Kurzbefehle `claude1`, `claude2` usw. anlegt. Wenn dieser Ordner nicht in deinem `PATH` liegt, gibt das Skript die Zeile aus, die du zu `~/.bashrc` oder `~/.zshrc` hinzufügen musst.

ccam startet Claude Code für dich, also installiere zuerst [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) und stelle sicher, dass `claude --version` funktioniert.

## Unterstützte Systeme

- Linux auf x86-64 und ARM64: Ubuntu, Debian, Fedora, Arch, Alpine und andere Distributionen. Die Binärdatei ist statisch und hat keine Abhängigkeiten.
- macOS auf Apple Silicon und Intel.
- Windows wird nicht direkt unterstützt. Innerhalb von WSL funktioniert ccam wie auf jedem Linux-System.

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

## Shell-Vervollständigung

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Die Vervollständigung deckt alle Befehle und Optionen von `ccam` ab, einschließlich der Kontonamen.

## Installation prüfen

```sh
ccam version
ccam doctor
```

`ccam doctor` prüft, ob Claude Code gefunden wird und ob jeder Kontoordner, jeder Link und jede Anmeldung in Ordnung ist. Weiter geht es mit dem [Schnellstart](/docs/quick-start).

## Deinstallieren

Entferne den Befehl und die Konto-Kurzbefehle und danach die Datenordner, wenn du die Konten nicht mehr brauchst:

```sh
rm ~/.local/bin/ccam ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Wenn du `~/.local/share/ccam` löschst, werden auch die dort gespeicherten Anmeldungen und Gesprächsverläufe gelöscht. Ein `~/.claude`-Ordner, den du [importiert](/docs/commands) hast, bleibt, wo er ist.
