---
title: Instalacja ccam
description: Zainstaluj ccam na Linux lub macOS jednym poleceniem, z pakietu .deb lub ze źródeł. Najpierw musi być zainstalowany Claude Code.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 79dc1b55e774
---

Uruchom to w terminalu na Linux lub macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Skrypt pobiera wydanie dla twojego systemu, sprawdza jego sumę kontrolną SHA-256 i umieszcza `ccam` w `~/.local/bin`, gdzie ccam tworzy też później wybrane przez ciebie skróty, na przykład `ccwork`. Jeśli tego folderu nie ma w `PATH`, skrypt wypisuje wiersz, który trzeba dodać do `~/.bashrc` lub `~/.zshrc`.

ccam uruchamia za ciebie Claude Code, więc najpierw zainstaluj [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) i upewnij się, że `claude --version` działa.

## Obsługiwane systemy

- Linux na x86-64 i ARM64: Ubuntu, Debian, Fedora, Arch, Alpine i inne dystrybucje. Plik binarny jest statyczny i nie ma zależności.
- macOS na Apple silicon i Intel.
- Windows nie jest obsługiwany bezpośrednio. W WSL ccam działa jak w każdym systemie Linux.

## Inne sposoby instalacji

**Pakiet dla Debiana i Ubuntu.** Pobierz `ccam_<version>_linux_amd64.deb` ze [strony wydań](https://github.com/jhyunwoo/ccam/releases) i zainstaluj go:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Pakiet instaluje `/usr/bin/ccam`.

**Ze źródeł.** Z Go 1.26 lub nowszym:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**Konkretna wersja.** Ustaw `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Uzupełnianie w powłoce

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Uzupełnianie obejmuje wszystkie polecenia i opcje `ccam`, łącznie z nazwami kont.

## Sprawdzenie instalacji

```sh
ccam version
ccam doctor
```

`ccam doctor` sprawdza, czy Claude Code został znaleziony i czy każdy folder konta, link i login są w porządku. Przejdź dalej do [szybkiego startu](/docs/quick-start).

## Odinstalowanie

Najpierw usuń skróty kont, dopóki ccam jest jeszcze zainstalowany. `ccam list` pokazuje je w kolumnie Shortcut; usuń każdy poleceniem `ccam shortcut <account> --remove` albo sam usuń te linki z `~/.local/bin`. Potem usuń polecenie, a także foldery z danymi, jeśli konta nie są ci już potrzebne:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Usunięcie `~/.local/share/ccam` usuwa też zapisane tam loginy i historię rozmów. Folder `~/.claude`, który został [zaimportowany](/docs/commands), zostaje tam, gdzie jest.
