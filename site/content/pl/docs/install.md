---
title: Instalacja ccam
metaTitle: "Instalacja ccam w systemach Linux i macOS"
description: "Zainstaluj ccam na Linux lub macOS jednym poleceniem, z pakietu .deb albo ze źródeł, i aktualizuj przez ccam update. Najpierw zainstaluj Claude Code."
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 7ecd1ed00e09
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

## Aktualizacja

```sh
ccam update --check
ccam update
```

`ccam update --check` tylko informuje, czy jest dostępne nowsze wydanie. `ccam update` pobiera wydanie dla twojego systemu, porównuje jego sumę kontrolną SHA-256 z plikiem `checksums.txt` z tego wydania, a potem w jednym kroku podmienia plik `ccam`. Jeśli suma kontrolna się nie zgadza, nic nie zostaje zmienione. Skróty, na przykład `ccwork`, są linkami do ccam, więc od razu uruchamiają nową wersję.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam nigdy sam nie sprawdza aktualizacji. Szuka nowego wydania tylko wtedy, gdy uruchomisz to polecenie. Aby zainstalować konkretne wydanie, także starsze, podaj `--version`. `ccam --dry-run update` pokazuje, co zostałoby pobrane i podmienione, niczego nie zmieniając.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Ponowne uruchomienie skryptu instalacyjnego też aktualizuje kopię w `~/.local/bin`. Jeśli zbudowano ccam przez `go install`, uruchom to samo polecenie jeszcze raz.

**Pakiet dla Debiana i Ubuntu.** `ccam update` nie podmienia `/usr/bin/ccam` zainstalowanego przez apt, tylko wypisuje poniższe kroki. Pobierz nowy `ccam_<version>_linux_amd64.deb` ze [strony wydań](https://github.com/jhyunwoo/ccam/releases) i zainstaluj go na miejsce starego:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Odinstalowanie

Najpierw usuń skróty kont, dopóki ccam jest jeszcze zainstalowany. `ccam list` pokazuje je w kolumnie Shortcut; usuń każdy poleceniem `ccam shortcut <account> --remove` albo sam usuń te linki z `~/.local/bin`. Potem usuń polecenie, a także foldery z danymi, jeśli konta nie są ci już potrzebne:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Usunięcie `~/.local/share/ccam` usuwa też zapisane tam loginy i historię rozmów. Folder `~/.claude`, który został [zaimportowany](/docs/commands), zostaje tam, gdzie jest.
