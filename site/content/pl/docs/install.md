---
title: Instalacja ccam
metaTitle: "Instalacja ccam w systemach Linux, macOS i Windows"
description: "Zainstaluj ccam na Linuksie, macOS lub Windows jednym poleceniem, z pakietu .deb lub przez go install i aktualizuj przez ccam update. Wymaga Claude Code."
updated: 2026-10-26
order: 1
translation: machine
sourceHash: "2ada7b2650d6"
---

ccam uruchamia za ciebie Claude Code, więc najpierw zainstaluj [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) i upewnij się, że `claude --version` działa.

Uruchom to w terminalu na Linux lub macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Skrypt pobiera wydanie dla twojego systemu, sprawdza jego sumę kontrolną SHA-256 i umieszcza `ccam` w `~/.local/bin`, gdzie ccam tworzy też później wybrane przez ciebie skróty, na przykład `ccwork`. Jeśli tego folderu nie ma w `PATH`, skrypt wypisuje wiersz, który trzeba dodać do `~/.bashrc` lub `~/.zshrc`.

W Windows uruchom to w PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Skrypt pobiera `ccam_windows_<arch>.zip`, sprawdza jego sumę kontrolną SHA-256 i umieszcza `ccam.exe` w `%USERPROFILE%\.local\bin`, czyli w folderze, do którego instalator Claude Code wstawia `claude.exe`. Jeśli tego folderu nie ma jeszcze w `PATH` twojego użytkownika, skrypt go dodaje i zachowuje bez zmian wpisy takie jak `%USERPROFILE%\bin`. Potem otwórz nowy terminal.

## Obsługiwane systemy

- Linux na x86-64 i ARM64: Ubuntu, Debian, Fedora, Arch, Alpine i inne dystrybucje. Plik binarny jest statyczny i nie ma zależności.
- macOS na Apple silicon i Intel.
- Windows 10 i 11 na x86-64 i ARM64, w PowerShell, Windows Terminal lub Git Bash:
  - Współdzielone foldery działają bez Trybu dewelopera.
  - Współdzielone pliki (`settings.json`, `CLAUDE.md`, `history.jsonl`) wymagają Trybu dewelopera (Ustawienia > System > Dla deweloperów), bo Windows pozwala tworzyć linki do plików tylko w tym trybie albo administratorowi; zobacz [współdzielenie w Windows](/features/share#on-windows).
  - Bez Trybu dewelopera pierwsze `ccam add` tworzy grupę `default` tylko z serwerami MCP i rozszerzeniami i podaje, jak dodać resztę po włączeniu Trybu dewelopera.
- W WSL ccam działa tak samo jak w każdym innym systemie Linux.

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

Jeśli zainstalujesz tag wydania (`@latest` lub `@vX.Y.Z`), ccam podaje tę wersję, a `ccam update` i propozycja aktualizacji w panelu działają. Kompilacja z `@main` albo z lokalnej kopii z commitami po tagu podaje `dev`, a `ccam update` odmawia jej aktualizacji, chyba że podasz `--version`.

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

Dodaje:

- skill, który proponuje polecenia logujące lub przenoszące pliki, ale nigdy ich nie uruchamia;
- polecenia `/ccam:usage`, `/ccam:ext` i `/ccam:switch`;
- dwa hooki ([szczegóły](/features/next#inside-claude-code)). Jeden pokazuje powiadomienie w terminalu, gdy sesja zatrzyma się na limicie użycia; na Linux i macOS, gdy włączony jest [`auto_switch`](/features/next#auto-switch), ccam przenosi potem sesję na inne konto, które ma jeszcze zapas. Drugi pozwala `/ccam:switch [account]` przekazać rozmowę innemu kontu.

Hooki uruchamiają `ccam`, więc musi on być w `PATH`, który widzi Claude Code.

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

Panel (`ccam` bez argumentów, w terminalu) także sprawdza przy każdym otwarciu, czy jest nowe wydanie, i, jeśli jest, proponuje je: naciśnij `U` i potwierdź. Żadne inne polecenie nie sprawdza aktualizacji, a ccam nigdy nie aktualizuje się bez pytania. Aby zainstalować konkretne wydanie, także starsze, podaj `--version`. `ccam --dry-run update` pokazuje, co zostałoby pobrane i podmienione, niczego nie zmieniając.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Ponowne uruchomienie skryptu instalacyjnego też aktualizuje kopię w `~/.local/bin`. Jeśli zbudowano ccam przez `go install` z tagu wydania, możesz uruchomić `ccam update` albo jeszcze raz to samo polecenie `go install`.

**Windows.** `ccam update` pobiera plik `.zip` dla Windows. W Windows nie da się podmienić działającego programu, więc stary plik binarny dostaje nazwę `ccam.exe.old`, a następna aktualizacja go usuwa; gdy sesja wciąż działa z wcześniejszego `ccam.exe.old`, nazwą jest `ccam.exe.old-<time>`. Skróty będące twardymi linkami są odnawiane zaraz potem. Ponowne uruchomienie `install.ps1` działa tak samo, nawet gdy sesja Claude Code wciąż działa ze starego pliku binarnego, i wypisuje `OK 2 shortcuts now run the new ccam`, gdy odnawia skróty będące twardymi linkami.

**Pakiet dla Debiana i Ubuntu.** `ccam update` nie podmienia `/usr/bin/ccam` zainstalowanego przez apt, tylko wypisuje poniższe kroki. Pobierz nowy `ccam_<version>_linux_amd64.deb` ze [strony wydań](https://github.com/jhyunwoo/ccam/releases) i zainstaluj go na miejsce starego:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Odinstalowanie

Usunięcie `~/.local/share/ccam` usuwa też zapisane tam loginy i historię rozmów. Folder `~/.claude`, który został [zaimportowany](/docs/commands), zostaje tam, gdzie jest, ale dopóki jest w grupie, pliki, które współdzieli, leżą w `~/.local/share/ccam/groups/`. Zanim cokolwiek usuniesz, uruchom dla niego `ccam remove <name>` (albo `ccam group leave`): wyłącza go to z jego grup i umieszcza z powrotem w folderze kopię wszystkiego, co współdzielił.

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
