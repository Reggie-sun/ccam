---
title: Instalacja ccam
metaTitle: "Instalacja ccam w systemach Linux, macOS i Windows"
description: "Zainstaluj ccam na Linuksie, macOS lub Windows jednym poleceniem, z pakietu .deb lub ze źródeł; aktualizuj przez ccam update. Najpierw zainstaluj Claude Code."
updated: 2026-10-21
order: 1
translation: machine
sourceHash: "1e3a53623de0"
---

Uruchom to w terminalu na Linux lub macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Skrypt pobiera wydanie dla twojego systemu, sprawdza jego sumę kontrolną SHA-256 i umieszcza `ccam` w `~/.local/bin`, gdzie ccam tworzy też później wybrane przez ciebie skróty, na przykład `ccwork`. Jeśli tego folderu nie ma w `PATH`, skrypt wypisuje wiersz, który trzeba dodać do `~/.bashrc` lub `~/.zshrc`.

W Windows uruchom to w PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Skrypt pobiera `ccam_windows_<arch>.zip`, sprawdza jego sumę kontrolną SHA-256, umieszcza `ccam.exe` w `%USERPROFILE%\.local\bin`, czyli w folderze, do którego instalator Claude Code wstawia `claude.exe`, i dodaje ten folder do `PATH` twojego użytkownika, jeśli jeszcze go tam nie ma. Potem otwórz nowy terminal.

ccam uruchamia za ciebie Claude Code, więc najpierw zainstaluj [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) i upewnij się, że `claude --version` działa.

## Obsługiwane systemy

- Linux na x86-64 i ARM64: Ubuntu, Debian, Fedora, Arch, Alpine i inne dystrybucje. Plik binarny jest statyczny i nie ma zależności.
- macOS na Apple silicon i Intel.
- Windows 10 i 11 na x86-64 i ARM64, w PowerShell, Windows Terminal lub Git Bash. Współdzielone foldery działają dla każdego użytkownika. Współdzielone pliki (`settings.json`, `CLAUDE.md`, `history.jsonl`) wymagają Trybu dewelopera (Ustawienia > System > Dla deweloperów), bo Windows pozwala tworzyć linki do plików tylko w tym trybie albo administratorowi; zobacz [współdzielenie w Windows](/features/share#on-windows).
- W WSL ccam działa jak w każdym systemie Linux.

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

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Uzupełnianie w powłoce

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Uzupełnianie obejmuje wszystkie polecenia i opcje `ccam`, łącznie z nazwami kont.

## Plugin Claude Code

Plugin ccam uczy sam Claude Code korzystać z ccam: sprawdzać, które konto ma zapas, uruchamiać równoległe workery `claude -p` na różnych kontach przez `ccam pick` i wyjaśniać brakującą skill lub serwer MCP przez `ccam ext`. Zainstaluj go w Claude Code:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

Dodaje skill, polecenia `/ccam:usage`, `/ccam:ext` i `/ccam:switch` oraz dwa hooki: gdy sesja zatrzyma się na limicie użycia, dostajesz powiadomienie w terminalu, a `/ccam:switch [account]` przekazuje rozmowę innemu kontu ([szczegóły](/features/next#inside-claude-code)). Hooki uruchamiają `ccam`, więc musi on być w `PATH`, który widzi Claude Code. Skill nigdy nie uruchamia poleceń, które logują lub przenoszą pliki; zamiast tego je proponuje.

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

Panel (`ccam` w terminalu) także szuka nowego wydania przy każdym otwarciu i, jeśli jest, proponuje je: naciśnij `U` i potwierdź. Żadne inne polecenie go nie szuka, a ccam nigdy nie aktualizuje się bez pytania. Aby zainstalować konkretne wydanie, także starsze, podaj `--version`. `ccam --dry-run update` pokazuje, co zostałoby pobrane i podmienione, niczego nie zmieniając.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Ponowne uruchomienie skryptu instalacyjnego też aktualizuje kopię w `~/.local/bin`. Jeśli zbudowano ccam przez `go install`, uruchom to samo polecenie jeszcze raz.

**Windows.** `ccam update` pobiera plik `.zip` dla Windows. W Windows nie da się podmienić działającego programu, więc stary plik binarny dostaje nazwę `ccam.exe.old`, a następna aktualizacja go usuwa.

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

W Windows, w PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Usunięcie `~/.local/share/ccam` usuwa też zapisane tam loginy i historię rozmów. Folder `~/.claude`, który został [zaimportowany](/docs/commands), zostaje tam, gdzie jest.
